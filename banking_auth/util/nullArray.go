package util

import (
	"encoding/json"
)

// NullArray handles nullable arrays in MongoDB
type NullStringArray struct {
	Accounts []NullString
}

// MarshalBSON for MongoDB serialization
// UnmarshalJSON to handle array of strings or null values
func (nsa *NullStringArray) UnmarshalJSON(data []byte) error {
	var temp []interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	for _, v := range temp {
		ns := NullString{}
		switch val := v.(type) {
		case string:
			ns.String = val
			ns.Valid = true
		case nil:
			ns.Valid = false
		default:
			ns.Valid = false
		}
		nsa.Accounts = append(nsa.Accounts, ns)
	}

	return nil
}
