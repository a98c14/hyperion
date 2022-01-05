package data

import "encoding/json"

type ModulePartValue struct {
	Id         int             `json:"id"`
	ValueType  int             `json:"valueType"`
	ArrayIndex int             `json:"arrayIndex"`
	Value      json.RawMessage `json:"value"`
}
