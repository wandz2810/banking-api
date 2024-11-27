package util

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

type NullString struct {
	String string
	Valid  bool // Valid is true if the string is not NULL
}

// MarshalJSON handles marshalling NullString into JSON format.
func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON handles unmarshalling NullString from JSON format.
func (ns *NullString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		ns.Valid = false
		ns.String = ""
		return nil
	}
	ns.Valid = true
	return json.Unmarshal(data, &ns.String)
}

// MarshalBSONValue handles marshalling NullString into BSON format for MongoDB.
func (ns NullString) MarshalBSONValue() (bsontype.Type, []byte, error) {
	if ns.Valid {
		return bson.MarshalValue(ns.String)
	}
	return bson.MarshalValue(nil)
}

// UnmarshalBSONValue handles unmarshalling NullString from BSON format in MongoDB.
func (ns *NullString) UnmarshalBSONValue(t bsontype.Type, data []byte) error {
	if t == bsontype.Null {
		ns.Valid = false
		ns.String = ""
		return nil
	}
	ns.Valid = true
	return bson.Unmarshal(data, &ns.String)
}
