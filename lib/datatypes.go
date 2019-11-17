package lib

import "encoding/json"

type AccessKey struct {
	Token string `json:"accessKey"`
}

func (a *AccessKey) ToJson() ([]byte, error){

	data, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return data, nil
}