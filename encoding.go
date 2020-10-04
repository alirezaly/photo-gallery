package main

import "encoding/json"

type Photo struct {
	ID             string `json:"id"`
	AltDescription string `json:"alt_description"`
	Liked          bool   `json:"liked"`
	User           struct {
		Name string `json:"name"`
	} `json:"user"`
}

func (p *Photo) encode() ([]byte, error) {
	enc, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

func decode(data []byte) (*Photo, error) {
	var p *Photo
	err := json.Unmarshal(data, &p)
	if err != nil {
		return nil, err
	}
	return p, nil
}
