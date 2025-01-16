// Copyright 2019 The Cockroach Authors.
//
// Use of this software is governed by the CockroachDB Software License
// included in the /LICENSE file.

package cdctest

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
)

type ChangefeedOption struct {
	Format           string
	PubsubSinkConfig string
	BooleanOptions   map[string]bool
	KafkaSinkConfig  string
}

func newChangefeedOption(testName string) ChangefeedOption {
	isCloudstorage := strings.Contains(testName, "cloudstorage")
	isWebhook := strings.Contains(testName, "webhook")
	isPubsub := strings.Contains(testName, "pubsub")
	isKafka := strings.Contains(testName, "kafka")

	cfo := ChangefeedOption{}

	cfo.BooleanOptions = make(map[string]bool)

	booleanOptionEligibility := map[string]bool{
		"full_table_name": true,
		// Because key_in_value is on by default for cloudstorage and webhook sinks,
		// the key in the value is extracted and removed from the test feed
		// messages (see extractKeyFromJSONValue function).
		"key_in_value":   !isCloudstorage && !isWebhook,
		"diff":           true,
		"mvcc_timestamp": true,
	}

	for option, eligible := range booleanOptionEligibility {
		cfo.BooleanOptions[option] = eligible && rand.Intn(2) < 1
	}

	if isPubsub {
		cfo.PubsubSinkConfig = newSinkConfig(isKafka).OptionString()
	}

	if isKafka {
		cfo.KafkaSinkConfig = newSinkConfig(isKafka).OptionString()
	}

	if isCloudstorage && rand.Intn(2) < 1 {
		cfo.Format = "parquet"
	} else {
		cfo.Format = "json"
	}

	return cfo
}

func (cfo ChangefeedOption) OptionString() string {
	var options []string
	for option, value := range cfo.BooleanOptions {
		if value {
			options = append(options, option)
		}
	}
	if cfo.Format != "" {
		option := fmt.Sprintf("format=%s", cfo.Format)
		options = append(options, option)
	}
	if cfo.PubsubSinkConfig != "" {
		option := fmt.Sprintf("pubsub_sink_config='%s'", cfo.PubsubSinkConfig)
		options = append(options, option)
	}
	if cfo.KafkaSinkConfig != "" {
		option := fmt.Sprintf("kafka_sink_config='%s'", cfo.KafkaSinkConfig)
		options = append(options, option)
	}
	return fmt.Sprintf("WITH updated, resolved, %s", strings.Join(options, ","))
}

type SinkConfig struct {
	Flush map[string]any
	Retry map[string]any
	Base  map[string]any
}

func (sk SinkConfig) OptionString() string {
	nonEmptyConfig := make(map[string]any)

	if len(sk.Flush) > 0 {
		nonEmptyConfig["Flush"] = sk.Flush
	}
	if len(sk.Retry) > 0 {
		nonEmptyConfig["Retry"] = sk.Retry
	}
	for k, v := range sk.Base {
		nonEmptyConfig[k] = v
	}

	jsonData, err := json.Marshal(nonEmptyConfig)
	if err != nil {
		return "{}"
	}

	return string(jsonData)
}

func newSinkConfig(isKafka bool) SinkConfig {
	flush := make(map[string]any)
	nonZeroInterval := "500ms"
	if rand.Intn(2) < 1 {
		// Setting either Messages or Bytes with a non-zero value without setting
		// Frequency is an invalid configuration. We set Frequency to a non-zero
		// interval here but can reset it later.
		flush["Messages"] = rand.Intn(10) + 1
		flush["Frequency"] = nonZeroInterval
	}
	if rand.Intn(2) < 1 {
		flush["Bytes"] = rand.Intn(1000) + 1
		flush["Frequency"] = nonZeroInterval
	}
	if isKafka && rand.Intn(2) < 1 {
		flush["MaxMessages"] = rand.Intn(10000) + 1
	}
	if rand.Intn(2) < 1 {
		intervals := []string{"100ms", "500ms", "1s", "5s"}
		interval := intervals[rand.Intn(len(intervals))]
		flush["Frequency"] = interval
	}

	retry := make(map[string]any)
	if !isKafka && rand.Intn(2) < 1 {
		if rand.Intn(2) < 1 {
			retry["Max"] = "inf"
		} else {
			retry["Max"] = rand.Intn(5) + 1
		}
	}
	if !isKafka && rand.Intn(2) < 1 {
		intervals := []string{"100ms", "500ms", "1s", "5s"}
		interval := intervals[rand.Intn(len(intervals))]
		retry["Backoff"] = interval
	}

	base := make(map[string]any)
	if isKafka && rand.Intn(2) < 1 {
		clientIds := []string{"ABCabc123._-", "FooBar", "2002-02-02.1_1"}
		clientId := clientIds[rand.Intn(len(clientIds))]
		base["ClientID"] = clientId
	}

	if isKafka && rand.Intn(2) < 1 {
		versions := []string{"2.7.2", "0.8.2.0"}
		version := versions[rand.Intn(len(versions))]
		base["Version"] = version
	}

	if isKafka && rand.Intn(2) < 1 {
		compressions := []string{"NONE", "GZIP", "SNAPPY"}
		// lz4 compression requires Version >= V0_10_0_0
		if base["Version"] != "0.8.2.0" {
			compressions = append(compressions, "LZ4")
		}
		// zstd compression requires Version >= V2_1_0_0
		if base["Version"] == "2.7.2" {
			compressions = append(compressions, "ZSTD")
		}
		compression := compressions[rand.Intn(len(compressions))]
		base["Compression"] = compression
		if compression == "GZIP" {
			level := rand.Intn(10)
			base["CompressionLevel"] = level

		}
		if base["Version"] == "2.7.2" && compression == "ZSTD" {
			level := rand.Intn(4) + 1
			base["CompressionLevel"] = level
		}

		if base["Version"] != "0.8.2.0" && compression == "LZ4" {
			levels := []int{0, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536, 131072}
			level := levels[rand.Intn(len(levels))]
			base["CompressionLevel"] = level
		}
	}

	if isKafka && rand.Intn(2) < 1 {
		levels := []string{"ONE", "NONE", "ALL"}
		level := levels[rand.Intn(len(levels))]
		base["RequiredAcks"] = level
	}

	return SinkConfig{
		Flush: flush,
		Retry: retry,
		Base:  base,
	}
}
