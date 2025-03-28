// Copyright 2018 The Cockroach Authors.
//
// Use of this software is governed by the CockroachDB Software License
// included in the /LICENSE file.

package ordering

import (
	"github.com/cockroachdb/cockroach/pkg/sql/opt"
	"github.com/cockroachdb/cockroach/pkg/sql/opt/memo"
	"github.com/cockroachdb/cockroach/pkg/sql/opt/props"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/errors"
)

func scanCanProvideOrdering(
	mem *memo.Memo, expr memo.RelExpr, required *props.OrderingChoice,
) bool {
	ok, _ := ScanPrivateCanProvide(
		mem.Metadata(),
		&expr.(*memo.ScanExpr).ScanPrivate,
		required,
	)
	return ok
}

// ScanIsReverse returns true if the scan must be performed in reverse order
// in order to satisfy the required ordering. If either direction is ok (e.g. no
// required ordering), reutrns false. The scan must be able to satisfy the
// required ordering, according to ScanCanProvideOrdering.
func ScanIsReverse(mem *memo.Memo, scan *memo.ScanExpr, required *props.OrderingChoice) bool {
	ok, reverse := ScanPrivateCanProvide(
		mem.Metadata(),
		&scan.ScanPrivate,
		required,
	)
	if !ok {
		panic(errors.AssertionFailedf("scan can't provide required ordering"))
	}
	return reverse
}

// ScanDirection represents the direction of a scan, either for a Scan operator
// or a LookupJoin.
type ScanDirection uint8

const (
	// EitherDirection indicates that the scan can be in either direction.
	EitherDirection ScanDirection = iota
	// ForwardDirection indicates a forward scan.
	ForwardDirection
	// ReverseDirection indicates a reverse scan.
	ReverseDirection
)

// String implements the fmt.Stringer interface.
func (d ScanDirection) String() string {
	switch d {
	case EitherDirection:
		return "either"
	case ForwardDirection:
		return "forward"
	case ReverseDirection:
		return "reverse"
	default:
		return "unknown"
	}
}

// ScanPrivateCanProvide returns true if the scan operator returns rows
// that satisfy the given required ordering; it also returns whether the scan
// needs to be in reverse order to match the required ordering.
func ScanPrivateCanProvide(
	md *opt.Metadata, s *memo.ScanPrivate, required *props.OrderingChoice,
) (ok bool, reverse bool) {
	// Scan naturally orders according to scanned index's key columns. A scan can
	// be executed either as a forward or as a reverse scan (unless it has a row
	// limit, in which case the direction is fixed).
	//
	// The code below follows the structure of OrderingChoice.Implies. We go
	// through the columns and determine if the ordering matches with either scan
	// direction.

	// We start off as accepting either a forward or a reverse scan. Until then,
	// the reverse variable is unset. Once the direction is known, reverseSet is
	// true and reverse indicates whether we need to do a reverse scan.
	var direction ScanDirection
	if s.HardLimit.IsSet() {
		// When we have a limit, the limit forces a certain scan direction (because
		// it affects the results, not just their ordering).
		direction = ForwardDirection
		if s.HardLimit.Reverse() {
			direction = ReverseDirection
		}
	} else if s.Flags.Direction != 0 {
		direction = ForwardDirection
		if s.Flags.Direction == tree.Descending {
			direction = ReverseDirection
		}
	}
	index := md.Table(s.Table).Index(s.Index)
	for left, right := 0, 0; right < len(required.Columns); {
		if left >= index.KeyColumnCount() {
			return false, false
		}
		indexCol := index.Column(left)
		indexColID := s.Table.ColumnID(indexCol.Ordinal())
		if required.Optional.Contains(indexColID) {
			left++
			continue
		}
		reqCol := &required.Columns[right]
		if !reqCol.Group.Contains(indexColID) {
			if left < s.ExactPrefix {
				// All columns in the exact prefix are constant and can be ignored.
				left++
				continue
			}
			return false, false
		}
		// The directions of the index column and the required column impose either
		// a forward or a reverse scan.
		requiredDirection := ForwardDirection
		if indexCol.Descending != reqCol.Descending {
			requiredDirection = ReverseDirection
		}
		if direction == EitherDirection {
			direction = requiredDirection
		} else if direction != requiredDirection {
			// We already determined the direction, and according to it, this column
			// has the wrong direction.
			return false, false
		}
		left, right = left+1, right+1
	}
	// If direction is either, we prefer forward scan.
	return true, direction == ReverseDirection
}

func scanBuildProvided(
	mem *memo.Memo, expr memo.RelExpr, required *props.OrderingChoice,
) opt.Ordering {
	scan := expr.(*memo.ScanExpr)
	md := mem.Metadata()
	index := md.Table(scan.Table).Index(scan.Index)
	fds := &scan.Relational().FuncDeps

	// We need to know the direction of the scan.
	reverse := ScanIsReverse(mem, scan, required)

	// We generate the longest ordering that this scan can provide, then we trim
	// it. This is the longest prefix of index columns that are output by the scan
	// (ignoring constant columns, in the case of constrained scans).
	outCols := expr.Relational().OutputCols
	constCols := fds.ComputeClosure(opt.ColSet{})
	numCols := index.KeyColumnCount()
	provided := make(opt.Ordering, 0, numCols)
	for i := 0; i < numCols; i++ {
		indexCol := index.Column(i)
		colID := scan.Table.ColumnID(indexCol.Ordinal())
		if i < scan.ExactPrefix && !outCols.Contains(colID) {
			// All columns in the exact prefix are constant and can be ignored as long
			// as they are not in the output of the scan. If an exact-prefix column is
			// in the output, it may still be constant, but can only be ignored if the
			// FDs "know" it is constant. This latter case is handled below as normal.
			continue
		}
		if constCols.Contains(colID) {
			// Column constrained to a constant, ignore.
			continue
		}
		if !scan.Cols.Contains(colID) {
			// Column not in output; we are done.
			break
		}
		direction := indexCol.Descending != reverse // != is bool XOR
		provided = append(provided, opt.MakeOrderingColumn(colID, direction))
	}

	return trimProvided(provided, required, fds)
}

func init() {
	memo.ScanIsReverseFn = func(
		md *opt.Metadata, s *memo.ScanPrivate, required *props.OrderingChoice,
	) bool {
		ok, reverse := ScanPrivateCanProvide(md, s, required)
		// If the scan cannot provide the ordering, then the scan direction
		// cannot be determined. We return false in that case, even though it is
		// not confirmed that the scan is not in reverse. This is possible when
		// a memo is not fully optimized.
		return ok && reverse
	}
}
