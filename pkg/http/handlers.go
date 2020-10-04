package http

import (
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/labstack/echo/v4"
	"github.com/alirezaly/photo-gallery/pkg/storage"
)

func Root(db *bolt.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		photos, _ := storage.Photos(db)
		return c.Render(http.StatusOK, "main", photos)
	}
}

func Like(db *bolt.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		photo, _ := storage.FindPhoto(db, c.QueryParam("id"))
		photo.Liked = !photo.Liked
		photo.Save(db)

		photos, _ := storage.Photos(db)
		return c.Render(http.StatusOK, "main", photos)
	}
}
