// Copyright 2015 The Cockroach Authors.
//
// Use of this software is governed by the CockroachDB Software License
// included in the /LICENSE file.

package types

import (
	"bytes"

	"github.com/gogo/protobuf/jsonpb"
)

// This file contains logic to allow the *types.T to properly marshal to json.
// It is a separate file to make it straightforward to defeat the linter that
// refuses to allow one to call a method Marshal unless it's protoutil.Marshal.

// MarshalJSONPB marshals the T to json. This is necessary as otherwise
// this field will be lost to the crdb_internal.pb_to_json and the likes.
func (t *T) MarshalJSONPB(marshaler *jsonpb.Marshaler) ([]byte, error) {
	// Map empty locale to nil so empty string does not appear in the JSON result.
	// TODO(rafi): When we upgrade to go1.24, we can modify the proto definition
	//  of the locale field to use `[(gogoproto.jsontag) = ",omitzero"]` instead of
	//  this workaround.
	temp := *t
	if temp.InternalType.Locale != nil && len(*temp.InternalType.Locale) == 0 {
		temp.InternalType.Locale = nil
	}
	var buf bytes.Buffer
	if err := marshaler.Marshal(&buf, &temp.InternalType); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalJSONPB unmarshals the T to json. This is necessary as otherwise
// this field will be lost to the crdb_internal.json_to_pb and the likes.
func (t *T) UnmarshalJSONPB(unmarshaler *jsonpb.Unmarshaler, data []byte) error {
	if err := unmarshaler.Unmarshal(bytes.NewReader(data), &t.InternalType); err != nil {
		return err
	}
	return t.upgradeType()
}
