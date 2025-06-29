// Copyright 2022 The Cockroach Authors.
//
// Use of this software is governed by the CockroachDB Software License
// included in the /LICENSE file.

// Code generated by sctestgen, DO NOT EDIT.

package sctestbackupccl

import (
	"testing"

	"github.com/cockroachdb/cockroach/pkg/sql/schemachanger/sctest"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
)

func TestBackupRollbacks_multiregion_add_column_multiple_regional_by_row(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/add_column_multiple_regional_by_row"
	sctest.BackupRollbacks(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupRollbacks_multiregion_add_column_regional_by_row(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/add_column_regional_by_row"
	sctest.BackupRollbacks(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupRollbacks_multiregion_add_column_subzones(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/add_column_subzones"
	sctest.BackupRollbacks(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupRollbacks_multiregion_alter_index_configure_zone(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/alter_index_configure_zone"
	sctest.BackupRollbacks(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupRollbacks_multiregion_alter_index_configure_zone_multiple(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/alter_index_configure_zone_multiple"
	sctest.BackupRollbacks(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupRollbacks_multiregion_alter_partition_configure_zone(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/alter_partition_configure_zone"
	sctest.BackupRollbacks(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupRollbacks_multiregion_alter_partition_configure_zone_discard(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/alter_partition_configure_zone_discard"
	sctest.BackupRollbacks(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupRollbacks_multiregion_alter_partition_configure_zone_multiple(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/alter_partition_configure_zone_multiple"
	sctest.BackupRollbacks(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupRollbacks_multiregion_alter_partition_configure_zone_subpartitions(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/alter_partition_configure_zone_subpartitions"
	sctest.BackupRollbacks(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupRollbacks_multiregion_alter_table_alter_primary_key_rbr(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/alter_table_alter_primary_key_rbr"
	sctest.BackupRollbacks(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupRollbacks_multiregion_create_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/create_index"
	sctest.BackupRollbacks(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupRollbacks_multiregion_create_index_regional_by_row(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/create_index_regional_by_row"
	sctest.BackupRollbacks(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupRollbacks_multiregion_create_trigger(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/create_trigger"
	sctest.BackupRollbacks(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupRollbacks_multiregion_drop_column_regional_by_row(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/drop_column_regional_by_row"
	sctest.BackupRollbacks(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupRollbacks_multiregion_drop_database_multiregion_primary_region(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/drop_database_multiregion_primary_region"
	sctest.BackupRollbacks(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupRollbacks_multiregion_drop_table_multiregion(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/drop_table_multiregion"
	sctest.BackupRollbacks(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupRollbacks_multiregion_drop_table_multiregion_primary_region(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/drop_table_multiregion_primary_region"
	sctest.BackupRollbacks(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupRollbacks_multiregion_drop_table_trigger(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/drop_table_trigger"
	sctest.BackupRollbacks(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupRollbacks_multiregion_drop_trigger(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/drop_trigger"
	sctest.BackupRollbacks(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupRollbacksMixedVersion_multiregion_add_column_multiple_regional_by_row(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/add_column_multiple_regional_by_row"
	sctest.BackupRollbacksMixedVersion(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupRollbacksMixedVersion_multiregion_add_column_regional_by_row(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/add_column_regional_by_row"
	sctest.BackupRollbacksMixedVersion(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupRollbacksMixedVersion_multiregion_add_column_subzones(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/add_column_subzones"
	sctest.BackupRollbacksMixedVersion(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupRollbacksMixedVersion_multiregion_alter_index_configure_zone(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/alter_index_configure_zone"
	sctest.BackupRollbacksMixedVersion(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupRollbacksMixedVersion_multiregion_alter_index_configure_zone_multiple(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/alter_index_configure_zone_multiple"
	sctest.BackupRollbacksMixedVersion(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupRollbacksMixedVersion_multiregion_alter_partition_configure_zone(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/alter_partition_configure_zone"
	sctest.BackupRollbacksMixedVersion(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupRollbacksMixedVersion_multiregion_alter_partition_configure_zone_discard(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/alter_partition_configure_zone_discard"
	sctest.BackupRollbacksMixedVersion(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupRollbacksMixedVersion_multiregion_alter_partition_configure_zone_multiple(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/alter_partition_configure_zone_multiple"
	sctest.BackupRollbacksMixedVersion(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupRollbacksMixedVersion_multiregion_alter_partition_configure_zone_subpartitions(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/alter_partition_configure_zone_subpartitions"
	sctest.BackupRollbacksMixedVersion(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupRollbacksMixedVersion_multiregion_alter_table_alter_primary_key_rbr(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/alter_table_alter_primary_key_rbr"
	sctest.BackupRollbacksMixedVersion(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupRollbacksMixedVersion_multiregion_create_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/create_index"
	sctest.BackupRollbacksMixedVersion(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupRollbacksMixedVersion_multiregion_create_index_regional_by_row(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/create_index_regional_by_row"
	sctest.BackupRollbacksMixedVersion(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupRollbacksMixedVersion_multiregion_create_trigger(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/create_trigger"
	sctest.BackupRollbacksMixedVersion(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupRollbacksMixedVersion_multiregion_drop_column_regional_by_row(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/drop_column_regional_by_row"
	sctest.BackupRollbacksMixedVersion(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupRollbacksMixedVersion_multiregion_drop_database_multiregion_primary_region(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/drop_database_multiregion_primary_region"
	sctest.BackupRollbacksMixedVersion(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupRollbacksMixedVersion_multiregion_drop_table_multiregion(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/drop_table_multiregion"
	sctest.BackupRollbacksMixedVersion(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupRollbacksMixedVersion_multiregion_drop_table_multiregion_primary_region(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/drop_table_multiregion_primary_region"
	sctest.BackupRollbacksMixedVersion(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupRollbacksMixedVersion_multiregion_drop_table_trigger(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/drop_table_trigger"
	sctest.BackupRollbacksMixedVersion(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupRollbacksMixedVersion_multiregion_drop_trigger(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/drop_trigger"
	sctest.BackupRollbacksMixedVersion(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccess_multiregion_add_column_multiple_regional_by_row(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/add_column_multiple_regional_by_row"
	sctest.BackupSuccess(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccess_multiregion_add_column_regional_by_row(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/add_column_regional_by_row"
	sctest.BackupSuccess(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccess_multiregion_add_column_subzones(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/add_column_subzones"
	sctest.BackupSuccess(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccess_multiregion_alter_index_configure_zone(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/alter_index_configure_zone"
	sctest.BackupSuccess(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccess_multiregion_alter_index_configure_zone_multiple(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/alter_index_configure_zone_multiple"
	sctest.BackupSuccess(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccess_multiregion_alter_partition_configure_zone(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/alter_partition_configure_zone"
	sctest.BackupSuccess(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccess_multiregion_alter_partition_configure_zone_discard(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/alter_partition_configure_zone_discard"
	sctest.BackupSuccess(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccess_multiregion_alter_partition_configure_zone_multiple(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/alter_partition_configure_zone_multiple"
	sctest.BackupSuccess(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccess_multiregion_alter_partition_configure_zone_subpartitions(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/alter_partition_configure_zone_subpartitions"
	sctest.BackupSuccess(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccess_multiregion_alter_table_alter_primary_key_rbr(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/alter_table_alter_primary_key_rbr"
	sctest.BackupSuccess(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccess_multiregion_create_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/create_index"
	sctest.BackupSuccess(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccess_multiregion_create_index_regional_by_row(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/create_index_regional_by_row"
	sctest.BackupSuccess(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccess_multiregion_create_trigger(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/create_trigger"
	sctest.BackupSuccess(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccess_multiregion_drop_column_regional_by_row(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/drop_column_regional_by_row"
	sctest.BackupSuccess(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccess_multiregion_drop_database_multiregion_primary_region(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/drop_database_multiregion_primary_region"
	sctest.BackupSuccess(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccess_multiregion_drop_table_multiregion(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/drop_table_multiregion"
	sctest.BackupSuccess(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccess_multiregion_drop_table_multiregion_primary_region(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/drop_table_multiregion_primary_region"
	sctest.BackupSuccess(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccess_multiregion_drop_table_trigger(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/drop_table_trigger"
	sctest.BackupSuccess(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccess_multiregion_drop_trigger(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/drop_trigger"
	sctest.BackupSuccess(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccessMixedVersion_multiregion_add_column_multiple_regional_by_row(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/add_column_multiple_regional_by_row"
	sctest.BackupSuccessMixedVersion(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccessMixedVersion_multiregion_add_column_regional_by_row(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/add_column_regional_by_row"
	sctest.BackupSuccessMixedVersion(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccessMixedVersion_multiregion_add_column_subzones(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/add_column_subzones"
	sctest.BackupSuccessMixedVersion(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccessMixedVersion_multiregion_alter_index_configure_zone(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/alter_index_configure_zone"
	sctest.BackupSuccessMixedVersion(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccessMixedVersion_multiregion_alter_index_configure_zone_multiple(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/alter_index_configure_zone_multiple"
	sctest.BackupSuccessMixedVersion(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccessMixedVersion_multiregion_alter_partition_configure_zone(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/alter_partition_configure_zone"
	sctest.BackupSuccessMixedVersion(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccessMixedVersion_multiregion_alter_partition_configure_zone_discard(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/alter_partition_configure_zone_discard"
	sctest.BackupSuccessMixedVersion(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccessMixedVersion_multiregion_alter_partition_configure_zone_multiple(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/alter_partition_configure_zone_multiple"
	sctest.BackupSuccessMixedVersion(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccessMixedVersion_multiregion_alter_partition_configure_zone_subpartitions(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/alter_partition_configure_zone_subpartitions"
	sctest.BackupSuccessMixedVersion(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccessMixedVersion_multiregion_alter_table_alter_primary_key_rbr(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/alter_table_alter_primary_key_rbr"
	sctest.BackupSuccessMixedVersion(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccessMixedVersion_multiregion_create_index(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/create_index"
	sctest.BackupSuccessMixedVersion(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccessMixedVersion_multiregion_create_index_regional_by_row(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/create_index_regional_by_row"
	sctest.BackupSuccessMixedVersion(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccessMixedVersion_multiregion_create_trigger(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/create_trigger"
	sctest.BackupSuccessMixedVersion(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccessMixedVersion_multiregion_drop_column_regional_by_row(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/drop_column_regional_by_row"
	sctest.BackupSuccessMixedVersion(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccessMixedVersion_multiregion_drop_database_multiregion_primary_region(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/drop_database_multiregion_primary_region"
	sctest.BackupSuccessMixedVersion(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccessMixedVersion_multiregion_drop_table_multiregion(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/drop_table_multiregion"
	sctest.BackupSuccessMixedVersion(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccessMixedVersion_multiregion_drop_table_multiregion_primary_region(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/drop_table_multiregion_primary_region"
	sctest.BackupSuccessMixedVersion(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccessMixedVersion_multiregion_drop_table_trigger(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/drop_table_trigger"
	sctest.BackupSuccessMixedVersion(t, path, MultiRegionTestClusterFactory{})
}

func TestBackupSuccessMixedVersion_multiregion_drop_trigger(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)
	const path = "pkg/ccl/schemachangerccl/testdata/end_to_end/drop_trigger"
	sctest.BackupSuccessMixedVersion(t, path, MultiRegionTestClusterFactory{})
}
