// Copyright 2025 The Cockroach Authors.
//
// Use of this software is governed by the CockroachDB Software License
// included in the /LICENSE file.

package changefeedccl

import (
	"github.com/cockroachdb/cockroach/pkg/ccl/changefeedccl/avro"
	"github.com/cockroachdb/cockroach/pkg/ccl/changefeedccl/cdcevent"
	"github.com/cockroachdb/cockroach/pkg/ccl/changefeedccl/changefeedbase"
	"github.com/cockroachdb/cockroach/pkg/util/cache"
	"github.com/cockroachdb/cockroach/pkg/util/json"
	"github.com/linkedin/goavro/v2"
)

type enrichedSourceProviderOpts struct {
	updated, mvccTimestamp bool
}
type enrichedSourceData struct {
	jobId string
	// TODO(#139692): Add schema info support.
	// TODO(#139691): Add job info support.
	// TODO(#139690): Add node/cluster info support.
}
type enrichedSourceProvider struct {
	opts       enrichedSourceProviderOpts
	sourceData enrichedSourceData

	// Caches for encoding. The key is the table id since the only values
	// dependent on the row are table info.
	// TODO: that's ok right?
	jsonCache *cache.UnorderedCache
	avroCache *cache.UnorderedCache
}

func newEnrichedSourceProvider(
	opts changefeedbase.EncodingOptions, sourceData enrichedSourceData,
) *enrichedSourceProvider {
	return &enrichedSourceProvider{
		sourceData: sourceData,
		opts: enrichedSourceProviderOpts{
			mvccTimestamp: opts.MVCCTimestamps,
			updated:       opts.UpdatedTimestamps,
		},
		jsonCache: cache.NewUnorderedCache(encoderCacheConfig),
		avroCache: cache.NewUnorderedCache(encoderCacheConfig),
	}
}

func (p *enrichedSourceProvider) Schema() (*avro.FunctionalRecord, error) {
	fields := []*avro.SchemaField{
		{Name: "job_id", SchemaType: []avro.SchemaType{avro.SchemaTypeNull, avro.SchemaTypeString}},
	}
	// TODO: move this to something callable by GetAvro(), since it needs it too
	sourceFn := func(row cdcevent.Row) (map[string]any, error) {
		return map[string]any{
			"job_id": goavro.Union(avro.SchemaTypeString, p.sourceData.jobId),
		}, nil
	}
	rec, err := avro.NewFunctionalRecord("source", "" /* namespace */, fields, sourceFn)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

func (p *enrichedSourceProvider) GetJSON(row cdcevent.Row) (json.JSON, error) {
	ck := row.TableID
	if j, ok := p.jsonCache.Get(ck); ok {
		return j.(json.JSON), nil
	}

	// TODO(various): Add fields here.
	keys := []string{"job_id"}
	b, err := json.NewFixedKeysObjectBuilder(keys)
	if err != nil {
		return nil, err
	}

	if err := b.Set("job_id", json.FromString(p.sourceData.jobId)); err != nil {
		return nil, err
	}

	j, err := b.Build()
	if err != nil {
		return nil, err
	}

	p.jsonCache.Add(ck, j)
	return j, nil
}

func (p *enrichedSourceProvider) GetAvro(row cdcevent.Row) ([]byte, error) {
	ck := row.TableID
	if bs, ok := p.avroCache.Get(ck); ok {
		return bs.([]byte), nil
	}

	// TODO(#139655): Implement this.

	p.avroCache.Add(ck, nil)
	return nil, nil
}
