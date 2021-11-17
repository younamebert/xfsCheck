package common

import "encoding/json"

func MarshalIndent(val interface{}) ([]byte, error) {
	return json.MarshalIndent(val, "", "    ")
}
