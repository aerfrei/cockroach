// Copyright 2021 The Cockroach Authors.
//
// Use of this software is governed by the CockroachDB Software License
// included in the /LICENSE file.

package closedts

import (
	"testing"
	"time"

	"github.com/cockroachdb/cockroach/pkg/kv/kvserver/closedts/ctpb"
	"github.com/cockroachdb/cockroach/pkg/util/hlc"
	"github.com/cockroachdb/cockroach/pkg/util/leaktest"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/stretchr/testify/require"
)

func TestTargetForPolicy(t *testing.T) {
	defer leaktest.AfterTest(t)()
	defer log.Scope(t).Close(t)

	cast := func(i int, unit time.Duration) time.Duration { return time.Duration(i) * unit }
	secs := func(i int) time.Duration { return cast(i, time.Second) }
	millis := func(i int) time.Duration { return cast(i, time.Millisecond) }

	now := hlc.Timestamp{WallTime: millis(100).Nanoseconds()}
	maxClockOffset := millis(500)

	for _, tc := range []struct {
		lagTargetNanos             time.Duration
		leadTargetOverride         time.Duration
		sideTransportCloseInterval time.Duration
		rangePolicy                ctpb.RangeClosedTimestampPolicy
		expClosedTSTarget          hlc.Timestamp
	}{
		{
			lagTargetNanos:    secs(3),
			rangePolicy:       ctpb.LAG_BY_CLUSTER_SETTING,
			expClosedTSTarget: now.Add(-secs(3).Nanoseconds(), 0),
		},
		{
			lagTargetNanos:    secs(1),
			rangePolicy:       ctpb.LAG_BY_CLUSTER_SETTING,
			expClosedTSTarget: now.Add(-secs(1).Nanoseconds(), 0),
		},
		{
			sideTransportCloseInterval: millis(200),
			rangePolicy:                ctpb.LEAD_FOR_GLOBAL_READS_WITH_NO_LATENCY_INFO,
			expClosedTSTarget: now.
				Add((maxClockOffset +
					millis(275) /* sideTransportPropTime */ +
					millis(25) /* bufferTime */).Nanoseconds(), 0),
		},
		{
			sideTransportCloseInterval: millis(50),
			rangePolicy:                ctpb.LEAD_FOR_GLOBAL_READS_WITH_NO_LATENCY_INFO,
			expClosedTSTarget: now.
				Add((maxClockOffset +
					millis(245) /* raftTransportPropTime */ +
					millis(25) /* bufferTime */).Nanoseconds(), 0),
		},
		{
			leadTargetOverride:         millis(1234),
			sideTransportCloseInterval: millis(200),
			rangePolicy:                ctpb.LEAD_FOR_GLOBAL_READS_WITH_NO_LATENCY_INFO,
			expClosedTSTarget:          now.Add(millis(1234).Nanoseconds(), 0),
		},
	} {
		t.Run("", func(t *testing.T) {
			target := TargetForPolicy(
				now.UnsafeToClockTimestamp(),
				maxClockOffset,
				tc.lagTargetNanos,
				tc.leadTargetOverride,
				tc.sideTransportCloseInterval,
				tc.rangePolicy,
			)
			require.Equal(t, tc.expClosedTSTarget, target)
		})
	}
}
