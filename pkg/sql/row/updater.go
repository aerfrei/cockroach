// Copyright 2019 The Cockroach Authors.
//
// Use of this software is governed by the CockroachDB Software License
// included in the /LICENSE file.

package row

import (
	"bytes"
	"context"
	"sort"

	"github.com/cockroachdb/cockroach/pkg/keys"
	"github.com/cockroachdb/cockroach/pkg/kv"
	"github.com/cockroachdb/cockroach/pkg/roachpb"
	"github.com/cockroachdb/cockroach/pkg/settings"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
	"github.com/cockroachdb/cockroach/pkg/sql/rowenc"
	"github.com/cockroachdb/cockroach/pkg/sql/rowinfra"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/idxtype"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/sessiondata"
	"github.com/cockroachdb/cockroach/pkg/util/buildutil"
	"github.com/cockroachdb/cockroach/pkg/util/deduplicate"
	"github.com/cockroachdb/cockroach/pkg/util/intsets"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/errors"
)

// Updater abstracts the key/value operations for updating table rows.
type Updater struct {
	Helper       RowHelper
	DeleteHelper *RowHelper
	FetchCols    []catalog.Column
	// FetchColIDtoRowIndex must be kept in sync with FetchCols.
	FetchColIDtoRowIndex  catalog.TableColMap
	UpdateCols            []catalog.Column
	UpdateColIDtoRowIndex catalog.TableColMap
	primaryKeyColChange   bool

	// primaryLocked, if true, indicates that no lock is needed when modifying
	// old KVs in the primary index because the caller already acquired it.
	primaryLocked bool
	// secondaryLocked, if set, indicates that no lock is needed when modifying
	// old KVs in this secondary index because the caller already acquired it.
	secondaryLocked catalog.Index

	// rd and ri are used when the update this Updater is created for modifies
	// the primary key of the table. In that case, rows must be deleted and
	// re-added instead of merely updated, since the keys are changing.
	rd Deleter
	ri Inserter

	// For allocation avoidance.
	newValues       []tree.Datum
	key             roachpb.Key
	valueBuf        []byte
	value           roachpb.Value
	oldIndexEntries [][]rowenc.IndexEntry
	newIndexEntries [][]rowenc.IndexEntry
}

type rowUpdaterType int

const (
	// UpdaterDefault indicates that an Updater should update everything
	// about a row, including secondary indexes.
	UpdaterDefault rowUpdaterType = 0
	// UpdaterOnlyColumns indicates that an Updater should only update the
	// columns of a row and not the secondary indexes.
	UpdaterOnlyColumns rowUpdaterType = 1
)

// MakeUpdater creates an Updater for the given table.
//
// lockedIndexes, if non-empty, specify the indexes in which old keys (those
// that are read as part of the initial scan) have already been locked.
//
// UpdateCols are the columns being updated and correspond to the updateValues
// that will be passed to UpdateRow.
//
// The returned Updater contains a FetchCols field that defines the
// expectation of which values are passed as oldValues to UpdateRow.
// requestedCols must be non-nil and define the schema that determines
// FetchCols.
func MakeUpdater(
	codec keys.SQLCodec,
	tableDesc catalog.TableDescriptor,
	uniqueWithTombstoneIndexes []catalog.Index,
	lockedIndexes []catalog.Index,
	updateCols []catalog.Column,
	requestedCols []catalog.Column,
	updateType rowUpdaterType,
	sd *sessiondata.SessionData,
	sv *settings.Values,
	metrics *rowinfra.Metrics,
) (Updater, error) {
	if requestedCols == nil {
		return Updater{}, errors.AssertionFailedf("requestedCols is nil in MakeUpdater")
	}

	var primaryIndexCols catalog.TableColSet
	for i := 0; i < tableDesc.GetPrimaryIndex().NumKeyColumns(); i++ {
		colID := tableDesc.GetPrimaryIndex().GetKeyColumnID(i)
		primaryIndexCols.Add(colID)
	}

	var primaryKeyColChange bool
	for _, c := range updateCols {
		if primaryIndexCols.Contains(c.GetID()) {
			primaryKeyColChange = true
			break
		}
	}

	updateColIDToRowIndex := ColIDtoRowIndexFromCols(updateCols)
	var includeIndexes []catalog.Index
	var deleteOnlyIndexes []catalog.Index
	// If the UPDATE is set to only update columns, do not collect secondary
	// indexes to update.
	if updateType != UpdaterOnlyColumns {
		for _, index := range tableDesc.DeletableNonPrimaryIndexes() {
			// If the primary key changed, we need to update all secondary
			// indexes, regardless of what other columns are being updated.
			if !primaryKeyColChange && !indexNeedsUpdate(index, updateColIDToRowIndex) {
				continue
			}
			if !index.DeleteOnly() {
				if includeIndexes == nil {
					includeIndexes = make([]catalog.Index, 0, len(tableDesc.WritableNonPrimaryIndexes()))
				}
				includeIndexes = append(includeIndexes, index)
			} else {
				if deleteOnlyIndexes == nil {
					// Allocate at most once.
					deleteOnlyIndexes = make([]catalog.Index, 0, len(tableDesc.DeleteOnlyNonPrimaryIndexes()))
				}
				deleteOnlyIndexes = append(deleteOnlyIndexes, index)
			}
		}
	}

	var deleteOnlyHelper *RowHelper
	if len(deleteOnlyIndexes) > 0 {
		rh := NewRowHelper(codec, tableDesc, deleteOnlyIndexes, nil /* uniqueWithTombstoneIndexes */, sd, sv, metrics)
		deleteOnlyHelper = &rh
	}

	var primaryLocked bool
	var secondaryLocked catalog.Index
	for _, index := range lockedIndexes {
		if index.Primary() {
			primaryLocked = true
		} else {
			secondaryLocked = index
		}
	}
	if buildutil.CrdbTestBuild && len(lockedIndexes) > 1 && !primaryLocked {
		// We don't expect multiple secondary indexes to be locked, yet if that
		// happens in prod, we'll just not use the already acquired locks on all
		// but the last secondary index, which means a possible performance hit
		// but no correctness issues.
		panic(errors.AssertionFailedf("locked at least two secondary indexes in the initial scan: %v", lockedIndexes))
	}
	numEntries := len(includeIndexes)
	indexEntries := make([][]rowenc.IndexEntry, numEntries*2)
	ru := Updater{
		Helper:                NewRowHelper(codec, tableDesc, includeIndexes, uniqueWithTombstoneIndexes, sd, sv, metrics),
		DeleteHelper:          deleteOnlyHelper,
		FetchCols:             requestedCols,
		FetchColIDtoRowIndex:  ColIDtoRowIndexFromCols(requestedCols),
		UpdateCols:            updateCols,
		UpdateColIDtoRowIndex: updateColIDToRowIndex,
		primaryKeyColChange:   primaryKeyColChange,
		primaryLocked:         primaryLocked,
		secondaryLocked:       secondaryLocked,
		oldIndexEntries:       indexEntries[:numEntries:numEntries],
		newIndexEntries:       indexEntries[numEntries:],
	}

	if primaryKeyColChange {
		// These fields are only used when the primary key is changing.
		var err error
		// All indexes that were locked during the initial scan don't require
		// locking when deleting from them - we only would delete KVs that we've
		// scanned (and locked) already.
		deleteLockedIndexes := lockedIndexes
		ru.rd = MakeDeleter(codec, tableDesc, deleteLockedIndexes, requestedCols, sd, sv, metrics)
		if ru.ri, err = MakeInserter(
			codec, tableDesc, uniqueWithTombstoneIndexes, requestedCols, sd, sv, metrics,
		); err != nil {
			return Updater{}, err
		}
	}

	// If we are fetching from specific families, we might get
	// less columns than in the table. So we cannot assign this to
	// have length len(tableCols).
	ru.newValues = make(tree.Datums, len(ru.FetchCols))

	return ru, nil
}

// indexNeedsUpdate returns true if the given index may need to be updated based
// on the columns in the map.
func indexNeedsUpdate(index catalog.Index, updateCols catalog.TableColMap) bool {
	// If the index is a partial index, an update may be required even if
	// the indexed columns aren't changing. For example, an index entry must
	// be added when an update to a non-indexed column causes a row to
	// satisfy the partial index predicate when it did not before.
	// TODO(mgartner): This function does not need to return true for every
	// partial index. A partial index will never require updating if neither
	// its indexed columns nor the columns referenced in its predicate
	// expression are changing.
	if index.IsPartial() {
		return true
	}
	colIDs := index.CollectKeyColumnIDs()
	colIDs.UnionWith(index.CollectSecondaryStoredColumnIDs())
	colIDs.UnionWith(index.CollectKeySuffixColumnIDs())
	for colID, ok := colIDs.Next(0); ok; colID, ok = colIDs.Next(colID + 1) {
		if _, ok := updateCols.Get(colID); ok {
			return true
		}
	}
	return false
}

// UpdateRow adds to the batch the kv operations necessary to update a table row
// with the given values.
//
// The row corresponding to oldValues is updated with the ones in updateValues.
// Note that updateValues only contains the ones that are changing.
//
// The return value is only good until the next call to UpdateRow.
func (ru *Updater) UpdateRow(
	ctx context.Context,
	batch *kv.Batch,
	oldValues []tree.Datum,
	updateValues []tree.Datum,
	pm PartialIndexUpdateHelper,
	vh VectorIndexUpdateHelper,
	oth OriginTimestampCPutHelper,
	mustValidateOldPKValues bool,
	traceKV bool,
) ([]tree.Datum, error) {
	if len(oldValues) != len(ru.FetchCols) {
		return nil, errors.Errorf("got %d values but expected %d", len(oldValues), len(ru.FetchCols))
	}
	if len(updateValues) != len(ru.UpdateCols) {
		return nil, errors.Errorf("got %d values but expected %d", len(updateValues), len(ru.UpdateCols))
	}

	primaryIndexKey, err := ru.Helper.encodePrimaryIndexKey(ru.FetchColIDtoRowIndex, oldValues)
	if err != nil {
		return nil, err
	}
	var deleteOldSecondaryIndexEntries map[catalog.Index][]rowenc.IndexEntry
	if ru.DeleteHelper != nil {
		// We want to include empty k/v pairs because we want
		// to delete all k/v's for this row. By setting includeEmpty
		// to true, we will get a k/v pair for each family in the row,
		// which will guarantee that we delete all the k/v's in this row.
		// N.B. that setting includeEmpty to true will sometimes cause
		// deletes of keys that aren't present. We choose to make this
		// compromise in order to avoid having to read all values of
		// the row that is being updated.
		_, deleteOldSecondaryIndexEntries, err = ru.DeleteHelper.encodeIndexes(
			ctx, ru.FetchColIDtoRowIndex, oldValues, vh.GetDel(), pm.IgnoreForDel, true, /* includeEmpty */
		)
		if err != nil {
			return nil, err
		}
	}

	// Update the row values.
	copy(ru.newValues, oldValues)
	for i, updateCol := range ru.UpdateCols {
		idx, ok := ru.FetchColIDtoRowIndex.Get(updateCol.GetID())
		if !ok {
			return nil, errors.AssertionFailedf("update column without a corresponding fetch column")
		}
		ru.newValues[idx] = updateValues[i]
	}

	rowPrimaryKeyChanged := false
	if ru.primaryKeyColChange {
		var newPrimaryIndexKey []byte
		newPrimaryIndexKey, err =
			ru.Helper.encodePrimaryIndexKey(ru.FetchColIDtoRowIndex, ru.newValues)
		if err != nil {
			return nil, err
		}
		rowPrimaryKeyChanged = !bytes.Equal(primaryIndexKey, newPrimaryIndexKey)
	}

	for i, index := range ru.Helper.Indexes {
		// We don't want to insert any empty k/v's, so set includeEmpty to false.
		// Consider the following case:
		// TABLE t (
		//   x INT PRIMARY KEY, y INT, z INT, w INT,
		//   INDEX (y) STORING (z, w),
		//   FAMILY (x), FAMILY (y), FAMILY (z), FAMILY (w)
		//)
		// If we are to perform an update on row (1, 2, 3, NULL), the k/v pair
		// for index i that encodes column w would have an empty value because w
		// is null and the sole resident of that family. We want to ensure that
		// we don't insert empty k/v pairs during the process of the update, so
		// set includeEmpty to false while generating the old and new index
		// entries.
		//
		// Also, we don't build entries for old and new values if the index
		// exists in ignoreIndexesForDel and ignoreIndexesForPut, respectively.
		// Index IDs in these sets indicate that old and new values for the row
		// do not satisfy a partial index's predicate expression.
		if pm.IgnoreForDel.Contains(int(index.GetID())) {
			ru.oldIndexEntries[i] = nil
		} else {
			ru.oldIndexEntries[i], err = rowenc.EncodeSecondaryIndex(
				ctx,
				ru.Helper.Codec,
				ru.Helper.TableDesc,
				index,
				ru.FetchColIDtoRowIndex,
				oldValues,
				vh.GetDel(),
				false, /* includeEmpty */
			)
			if err != nil {
				return nil, err
			}
		}
		if pm.IgnoreForPut.Contains(int(index.GetID())) {
			ru.newIndexEntries[i] = nil
		} else {
			ru.newIndexEntries[i], err = rowenc.EncodeSecondaryIndex(
				ctx,
				ru.Helper.Codec,
				ru.Helper.TableDesc,
				index,
				ru.FetchColIDtoRowIndex,
				ru.newValues,
				vh.GetPut(),
				false, /* includeEmpty */
			)
			if err != nil {
				return nil, err
			}
		}
		if ru.Helper.Indexes[i].GetType() == idxtype.INVERTED && !ru.Helper.Indexes[i].IsTemporaryIndexForBackfill() && !ru.Helper.Indexes[i].Merging() {
			// Deduplicate the keys we're adding and removing if we're updating an
			// inverted index. For example, imagine a table with an inverted index on j:
			//
			// a | j
			// --+----------------
			// 1 | {"foo": "bar"}
			//
			// If we update the json value to be {"foo": "bar", "baz": "qux"}, we don't
			// want to delete the /foo/bar key and re-add it, that would be wasted work.
			// So, we are going to remove keys from both the new and old index entry
			// array if they're identical.
			//
			// We don't do this deduplication on temporary indexes used during the
			// backfill because any deletes that are elided here are not elided on the
			// newly added index when it is in DELETE_ONLY.
			newIndexEntries := ru.newIndexEntries[i]
			oldIndexEntries := ru.oldIndexEntries[i]
			sort.Slice(oldIndexEntries, func(i, j int) bool {
				return compareIndexEntries(oldIndexEntries[i], oldIndexEntries[j]) < 0
			})
			sort.Slice(newIndexEntries, func(i, j int) bool {
				return compareIndexEntries(newIndexEntries[i], newIndexEntries[j]) < 0
			})
			oldLen, newLen := deduplicate.AcrossSlices(
				oldIndexEntries, newIndexEntries,
				func(l, r int) int {
					return compareIndexEntries(oldIndexEntries[l], newIndexEntries[r])
				},
				func(i, j int) {
					oldIndexEntries[i] = oldIndexEntries[j]
				},
				func(i, j int) {
					newIndexEntries[i] = newIndexEntries[j]
				})
			ru.oldIndexEntries[i] = oldIndexEntries[:oldLen]
			ru.newIndexEntries[i] = newIndexEntries[:newLen]
		}
	}

	b := &KVBatchAdapter{Batch: batch}
	if rowPrimaryKeyChanged {
		// TODO(#143175): the current pattern of deleting the full row and then
		// inserting the full row is suboptimal in how it handles unique
		// secondary indexes when the row doesn't contain NULL values in the
		// indexed columns. In such a scenario, the key in the unique secondary
		// index doesn't necessarily change, so rather than performing a Del
		// followed by a CPut we could skip the Del altogether. Furthermore, if
		// we acquired the lock on this index during the initial scan we could
		// replace a CPut with a Put.
		if err := ru.rd.DeleteRow(
			ctx, batch, oldValues, pm, vh, oth, mustValidateOldPKValues, traceKV,
		); err != nil {
			return nil, err
		}
		if err := ru.ri.InsertRow(
			ctx, b, ru.newValues, pm, vh, oth, CPutOp, traceKV,
		); err != nil {
			return nil, err
		}

		return ru.newValues, nil
	}

	// Add the new values to the primary index.
	kvOp := PutMustAcquireExclusiveLockOp
	if ru.primaryLocked {
		// Since the row PK doesn't change, and we've already locked it, we can
		// skip the lock acquisition.
		kvOp = PutOp
	}
	ru.valueBuf, err = prepareInsertOrUpdateBatch(
		ctx, b, &ru.Helper, primaryIndexKey, ru.FetchCols, ru.newValues, ru.FetchColIDtoRowIndex,
		ru.UpdateColIDtoRowIndex, &ru.key, &ru.value, ru.valueBuf, oth, oldValues,
		kvOp, mustValidateOldPKValues, traceKV,
	)
	if err != nil {
		return nil, err
	}

	// Update secondary indexes.
	// We're iterating through all of the indexes, which should have corresponding entries
	// in the new and old values.
	var writtenIndexes intsets.Fast
	for i, index := range ru.Helper.Indexes {
		alreadyLocked := ru.secondaryLocked != nil && ru.secondaryLocked.GetID() == index.GetID()
		// putFn and sameKeyPutFn are the functions that should be invoked in
		// order to write the new k/v entry. If the key doesn't change,
		// sameKeyPutFn will be used, otherwise putFn will be used.
		var putFn, sameKeyPutFn func(context.Context, Putter, *roachpb.Key, *roachpb.Value, bool, *RowHelper, lazyIndexDirs)
		if index.ForcePut() {
			// See the comment on (catalog.Index).ForcePut() for more details.
			// TODO(#140695): re-evaluate the lock need when we enable buffered
			// writes with DDLs.
			putFn = insertPutFn
			sameKeyPutFn = insertPutFn
		} else if index.IsUnique() {
			// For unique indexes we need to ensure that key doesn't exist
			// already.
			putFn = insertCPutFn
			// However, when updating an existing key, we must use a locking Put
			// (unless we already acquired a lock on the key in which case we
			// can elide the lock).
			sameKeyPutFn = insertPutMustAcquireExclusiveLockFn
			if alreadyLocked {
				sameKeyPutFn = insertPutFn
			}
		} else {
			// For non-unique indexes we don't care whether there exists an
			// entry already, so we can just use the Put. (In fact, since we
			// always include the PK columns into the non-unique secondary index
			// key, the current key should never already exist (unless we have a
			// duplicate PK which will be detected when modifying the primary
			// index).)
			//
			// We also don't need the lock.
			putFn = insertPutFn
			sameKeyPutFn = insertPutFn
			if ru.Helper.sd.BufferedWritesUseLockingOnNonUniqueIndexes {
				// When dealing with contention on this non-unique index, we
				// might benefit from locking the keys (when the corresponding
				// session var is enabled).
				putFn = insertPutMustAcquireExclusiveLockFn
				sameKeyPutFn = insertPutMustAcquireExclusiveLockFn
			}
			if ru.Helper.sd.UseCPutsOnNonUniqueIndexes {
				// We'll use CPuts for new keys if the session variable dictates
				// that.
				putFn = insertCPutFn
			}
		}
		if index.GetType() == idxtype.FORWARD {
			oldIdx, newIdx := 0, 0
			oldEntries, newEntries := ru.oldIndexEntries[i], ru.newIndexEntries[i]
			// The index entries for a particular index are stored in
			// family sorted order. We use this fact to update rows.
			// The algorithm to update a row using the old k/v pairs
			// for the row and the new k/v pairs for the row is very
			// similar to the algorithm to merge two sorted lists.
			// We move in lock step through the entries, and potentially
			// update k/v's that belong to the same family.
			// If we are in the case where there exists a family's k/v
			// in the old entries but not the new entries, we need to
			// delete that k/v. If we are in the case where a family's
			// k/v exists in the new index entries, then we need to just
			// insert that new k/v.
			for oldIdx < len(oldEntries) && newIdx < len(newEntries) {
				oldEntry, newEntry := &oldEntries[oldIdx], &newEntries[newIdx]
				if oldEntry.Family == newEntry.Family {
					// If the families are equal, then check if the keys have
					// changed. If so, delete the old key. Then, perform the
					// write for the new k/v entry.
					oldIdx++
					newIdx++
					var sameKey bool
					if !bytes.Equal(oldEntry.Key, newEntry.Key) {
						if err = ru.Helper.deleteIndexEntry(
							ctx, b, index, &oldEntry.Key, alreadyLocked,
							ru.Helper.sd.BufferedWritesUseLockingOnNonUniqueIndexes,
							traceKV, secondaryIndexDirs(i),
						); err != nil {
							return nil, err
						}
					} else if !newEntry.Value.EqualTagAndData(oldEntry.Value) {
						sameKey = true
					} else if !index.IsTemporaryIndexForBackfill() && !index.Merging() {
						// If this is a temporary index for backfill, we want to make sure we write out all
						// index values even in the case where they should be the same. We do this because the
						// temporary index is eventually merged into a newly added index that might be in a
						// DELETE_ONLY state at the time of this update and thus the temporary index needs to
						// have all of the entries.
						//
						// For merging indexes we will compare timestamps during the merge process so all
						// updates should always be captured into the final secondary index that is being
						// merged.  Otherwise, skip this put since the key and value are the same.
						continue
					}

					if sameKey {
						sameKeyPutFn(ctx, b, &newEntry.Key, &newEntry.Value, traceKV, &ru.Helper, secondaryIndexDirs(i))
					} else {
						putFn(ctx, b, &newEntry.Key, &newEntry.Value, traceKV, &ru.Helper, secondaryIndexDirs(i))
					}
					writtenIndexes.Add(i)
				} else if oldEntry.Family < newEntry.Family {
					if oldEntry.Family == descpb.FamilyID(0) {
						return nil, errors.AssertionFailedf(
							"index entry for family 0 for table %s, index %s was not generated",
							ru.Helper.TableDesc.GetName(), index.GetName(),
						)
					}
					// In this case, the index has a k/v for a family that does not exist in
					// the new set of k/v's for the row. So, we need to delete the old k/v.
					if err = ru.Helper.deleteIndexEntry(
						ctx, b, index, &oldEntry.Key, alreadyLocked,
						ru.Helper.sd.BufferedWritesUseLockingOnNonUniqueIndexes,
						traceKV, secondaryIndexDirs(i),
					); err != nil {
						return nil, err
					}
					oldIdx++
				} else {
					if newEntry.Family == descpb.FamilyID(0) {
						return nil, errors.AssertionFailedf(
							"index entry for family 0 for table %s, index %s was not generated",
							ru.Helper.TableDesc.GetName(), index.GetName(),
						)
					}

					// In this case, the index now has a k/v that did not exist
					// in the old row, so we put the new key in place.
					putFn(ctx, b, &newEntry.Key, &newEntry.Value, traceKV, &ru.Helper, secondaryIndexDirs(i))
					writtenIndexes.Add(i)
					newIdx++
				}
			}
			for oldIdx < len(oldEntries) {
				// Delete any remaining old entries that are not matched by new
				// entries in this row because 1) the family does not exist in
				// the new set of k/v's or 2) the index is a partial index and
				// the new row values do not match the partial index predicate.
				oldEntry := &oldEntries[oldIdx]
				if err = ru.Helper.deleteIndexEntry(
					ctx, b, index, &oldEntry.Key, alreadyLocked,
					ru.Helper.sd.BufferedWritesUseLockingOnNonUniqueIndexes,
					traceKV, secondaryIndexDirs(i),
				); err != nil {
					return nil, err
				}
				oldIdx++
			}
			for newIdx < len(newEntries) {
				// Insert any remaining new entries that are not present in the
				// old row. Insert any remaining new entries that are not
				// present in the old row because 1) the family does not exist
				// in the old set of k/v's or 2) the index is a partial index
				// and the old row values do not match the partial index
				// predicate.
				newEntry := &newEntries[newIdx]
				putFn(ctx, b, &newEntry.Key, &newEntry.Value, traceKV, &ru.Helper, secondaryIndexDirs(i))
				writtenIndexes.Add(i)
				newIdx++
			}
		} else {
			// Remove all inverted index entries, and re-add them.
			for j := range ru.oldIndexEntries[i] {
				if err = ru.Helper.deleteIndexEntry(
					ctx, b, index, &ru.oldIndexEntries[i][j].Key, alreadyLocked,
					ru.Helper.sd.BufferedWritesUseLockingOnNonUniqueIndexes, traceKV, emptyIndexDirs,
				); err != nil {
					return nil, err
				}
			}
			// We're adding all of the inverted index entries from the row being updated.
			for j := range ru.newIndexEntries[i] {
				putFn(ctx, b, &ru.newIndexEntries[i][j].Key, &ru.newIndexEntries[i][j].Value, traceKV,
					&ru.Helper, secondaryIndexDirs(i))
			}
		}
	}

	writtenIndexes.ForEach(func(idx int) {
		if err == nil {
			err = writeTombstones(ctx, &ru.Helper, ru.Helper.Indexes[idx], b, ru.FetchColIDtoRowIndex, ru.newValues, traceKV)
		}
	})
	if err != nil {
		return nil, err
	}

	// We're deleting indexes in a delete only state. We're bounding this by the number of indexes because inverted
	// indexed will be handled separately.
	if ru.DeleteHelper != nil {
		// For determinism, add the entries for the secondary indexes in the same
		// order as they appear in the helper.
		for _, index := range ru.DeleteHelper.Indexes {
			alreadyLocked := ru.secondaryLocked != nil && ru.secondaryLocked.GetID() == index.GetID()
			deletedSecondaryIndexEntries, ok := deleteOldSecondaryIndexEntries[index]

			if ok {
				for _, deletedSecondaryIndexEntry := range deletedSecondaryIndexEntries {
					if err = ru.DeleteHelper.deleteIndexEntry(
						ctx, b, index, &deletedSecondaryIndexEntry.Key, alreadyLocked,
						ru.Helper.sd.BufferedWritesUseLockingOnNonUniqueIndexes, traceKV, emptyIndexDirs,
					); err != nil {
						return nil, err
					}
				}
			}
		}
	}

	return ru.newValues, nil
}

func compareIndexEntries(left, right rowenc.IndexEntry) int {
	cmp := bytes.Compare(left.Key, right.Key)
	if cmp != 0 {
		return cmp
	}

	return bytes.Compare(left.Value.RawBytes, right.Value.RawBytes)
}

// IsColumnOnlyUpdate returns true if this Updater is only updating column
// data (in contrast to updating the primary key or other indexes).
func (ru *Updater) IsColumnOnlyUpdate() bool {
	// TODO(dan): This is used in the schema change backfill to assert that it was
	// configured correctly and will not be doing things it shouldn't. This is an
	// unfortunate bleeding of responsibility and indicates the abstraction could
	// be improved. Specifically, Updater currently has two responsibilities
	// (computing which indexes need to be updated and mapping sql rows to k/v
	// operations) and these should be split.
	return !ru.primaryKeyColChange && ru.DeleteHelper == nil && len(ru.Helper.Indexes) == 0
}

func updateCPutFn(
	ctx context.Context,
	b Putter,
	key *roachpb.Key,
	value *roachpb.Value,
	expVal []byte,
	traceKV bool,
	rh *RowHelper,
	dirs lazyIndexDirs,
) {
	if traceKV {
		log.VEventfDepth(
			ctx, 1, 2, "CPut %s -> %s (swap)", keys.PrettyPrint(dirs.compute(rh), *key),
			value.PrettyPrint(),
		)
	}
	b.CPut(key, value, expVal)
}
