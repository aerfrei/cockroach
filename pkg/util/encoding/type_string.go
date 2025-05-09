// Copyright 2023 The Cockroach Authors.
//
// Use of this software is governed by the CockroachDB Software License
// included in the /LICENSE file.

// Code generated by "stringer"; DO NOT EDIT.

package encoding

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Unknown-0]
	_ = x[Null-1]
	_ = x[NotNull-2]
	_ = x[Int-3]
	_ = x[Float-4]
	_ = x[Decimal-5]
	_ = x[Bytes-6]
	_ = x[BytesDesc-7]
	_ = x[Time-8]
	_ = x[Duration-9]
	_ = x[True-10]
	_ = x[False-11]
	_ = x[UUID-12]
	_ = x[Array-13]
	_ = x[IPAddr-14]
	_ = x[SentinelType-15]
	_ = x[JSON-15]
	_ = x[Tuple-16]
	_ = x[BitArray-17]
	_ = x[BitArrayDesc-18]
	_ = x[TimeTZ-19]
	_ = x[Geo-20]
	_ = x[GeoDesc-21]
	_ = x[ArrayKeyAsc-22]
	_ = x[ArrayKeyDesc-23]
	_ = x[Box2D-24]
	_ = x[Void-25]
	_ = x[TSQuery-26]
	_ = x[TSVector-27]
	_ = x[JSONNull-28]
	_ = x[JSONNullDesc-29]
	_ = x[JSONString-30]
	_ = x[JSONStringDesc-31]
	_ = x[JSONNumber-32]
	_ = x[JSONNumberDesc-33]
	_ = x[JSONFalse-34]
	_ = x[JSONFalseDesc-35]
	_ = x[JSONTrue-36]
	_ = x[JSONTrueDesc-37]
	_ = x[JSONArray-38]
	_ = x[JSONArrayDesc-39]
	_ = x[JSONObject-40]
	_ = x[JSONObjectDesc-41]
	_ = x[JsonEmptyArray-42]
	_ = x[JsonEmptyArrayDesc-43]
	_ = x[PGVector-44]
}

func (i Type) String() string {
	switch i {
	case Unknown:
		return "Unknown"
	case Null:
		return "Null"
	case NotNull:
		return "NotNull"
	case Int:
		return "Int"
	case Float:
		return "Float"
	case Decimal:
		return "Decimal"
	case Bytes:
		return "Bytes"
	case BytesDesc:
		return "BytesDesc"
	case Time:
		return "Time"
	case Duration:
		return "Duration"
	case True:
		return "True"
	case False:
		return "False"
	case UUID:
		return "UUID"
	case Array:
		return "Array"
	case IPAddr:
		return "IPAddr"
	case SentinelType:
		return "SentinelType"
	case Tuple:
		return "Tuple"
	case BitArray:
		return "BitArray"
	case BitArrayDesc:
		return "BitArrayDesc"
	case TimeTZ:
		return "TimeTZ"
	case Geo:
		return "Geo"
	case GeoDesc:
		return "GeoDesc"
	case ArrayKeyAsc:
		return "ArrayKeyAsc"
	case ArrayKeyDesc:
		return "ArrayKeyDesc"
	case Box2D:
		return "Box2D"
	case Void:
		return "Void"
	case TSQuery:
		return "TSQuery"
	case TSVector:
		return "TSVector"
	case JSONNull:
		return "JSONNull"
	case JSONNullDesc:
		return "JSONNullDesc"
	case JSONString:
		return "JSONString"
	case JSONStringDesc:
		return "JSONStringDesc"
	case JSONNumber:
		return "JSONNumber"
	case JSONNumberDesc:
		return "JSONNumberDesc"
	case JSONFalse:
		return "JSONFalse"
	case JSONFalseDesc:
		return "JSONFalseDesc"
	case JSONTrue:
		return "JSONTrue"
	case JSONTrueDesc:
		return "JSONTrueDesc"
	case JSONArray:
		return "JSONArray"
	case JSONArrayDesc:
		return "JSONArrayDesc"
	case JSONObject:
		return "JSONObject"
	case JSONObjectDesc:
		return "JSONObjectDesc"
	case JsonEmptyArray:
		return "JsonEmptyArray"
	case JsonEmptyArrayDesc:
		return "JsonEmptyArrayDesc"
	case PGVector:
		return "PGVector"
	default:
		return "Type(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
