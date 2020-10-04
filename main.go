package main

import (
	"fmt"
	"html/template"

	"github.com/alirezaly/photo-gallery/pkg/config"
	"github.com/alirezaly/photo-gallery/pkg/http"
	"github.com/alirezaly/photo-gallery/pkg/storage"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func newTemplate() *Template {
	return &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
}
func main() {
	e := echo.New()

	config, err := config.SetupConfig()
	if err != nil {
		e.Logger.Fatal(errors.Wrap(err, "could not set up db"))
	}

	db, err := storage.SetupDB(config)
	if err != nil {
		e.Logger.Fatal(errors.Wrap(err, "could not set up db"))
	}
	defer db.Close()
	storage.Seed(db)

	e.Renderer = newTemplate()
	e.GET("/", http.Root(db))
	e.GET("/like", http.Like(db))
	e.Static("/assets", "public/assets")
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", config.Host, config.Port)))

}
