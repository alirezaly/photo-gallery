package storage

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

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

func (p *Photo) Save(db *bolt.DB) error {
	err := db.Update(func(tx *bolt.Tx) error {
		photos, err := tx.CreateBucketIfNotExists([]byte("photos"))
		if err != nil {
			return errors.Wrap(err, "could not create photos bucket")
		}
		enc, err := p.encode()
		if err != nil {
			return errors.Wrap(err, "could not encode photo")
		}

		err = photos.Put([]byte(p.ID), enc)
		return err
	})
	return err
}
