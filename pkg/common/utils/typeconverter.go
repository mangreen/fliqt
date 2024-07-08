package utils

import "encoding/json"

// https://stackoverflow.com/a/72050244
func TypeConverter[R any](data any) (*R, error) {
	var result R

	b, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, err
	}

	return &result, err
}