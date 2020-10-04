package main

import (
	"encoding/json"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

func (p *Photo) save(db *bolt.DB) error {
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

func setupDB() (*bolt.DB, error) {
	db, err := bolt.Open("gallery.db", 0600, nil)
	if err != nil {
		return nil, errors.Wrap(err, "could not open db")
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("photos"))
		if err != nil {
			return errors.Wrap(err, "could not create photos bucket")
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not set up bucket")
	}

	return db, nil
}

func findPhotoByID(db *bolt.DB, id string) (*Photo, error) {
	var photo *Photo
	err := db.View(func(tx *bolt.Tx) error {
		var err error
		photos := tx.Bucket([]byte("photos"))
		photo, err = decode(photos.Get([]byte(id)))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "Could not find Photo")
	}

	return photo, nil
}

func seedDB(db *bolt.DB) error {
	p, _ := listPhotos(db, "photos")
	if len(p) > 1 {
		return nil
	}

	var photos []Photo

	res, err := http.Get("https://api.unsplash.com/photos?client_id=niSDwlj_iyHtD5u0e5UJvy55XrVvuaV6MH7NSUHupT4")
	if err != nil {
		return err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&photos)
	if err != nil {
		return err
	}

	for _, photo := range photos {
		photo.save(db)
	}

	return nil
}

func listPhotos(db *bolt.DB, bucket string) ([]*Photo, error) {
	var photos []*Photo
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(bucket)).Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			photo, _ := decode(v)
			photos = append(photos, photo)
		}
		return nil
	})
	return photos, nil
}
