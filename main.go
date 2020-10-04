package main

import (
	"html/template"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func main() {
	e := echo.New()

	db, err := setupDB()
	if err != nil {
		e.Logger.Fatal(errors.Wrap(err, "could not set up db"))
	}
	defer db.Close()
	seedDB(db)

	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}

	e.Renderer = t
	e.GET("/", root(db))
	e.GET("/like", like(db))
	e.Static("/src", "public")
	e.Logger.Fatal(e.Start(":1323"))

}
