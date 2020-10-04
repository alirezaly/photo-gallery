package main

import (
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/labstack/echo/v4"
)

func root(db *bolt.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		photos, _ := listPhotos(db, "photos")
		return c.Render(http.StatusOK, "main", photos)
	}
}

func like(db *bolt.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		photo, _ := findPhotoByID(db, c.QueryParam("id"))
		photo.Liked = !photo.Liked
		photo.save(db)

		photos, _ := listPhotos(db, "photos")
		return c.Render(http.StatusOK, "main", photos)
	}
}
