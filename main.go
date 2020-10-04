package main

import (
	"fmt"
	"html/template"

	"github.com/BurntSushi/toml"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

var config struct {
	Port string `toml:"port"`
	Host string `toml:"host"`
}

func main() {
	e := echo.New()

	_, err := toml.DecodeFile("config.toml", &config)
	if err != nil {
		e.Logger.Fatal(err)
	}

	db, err := setupDB()
	if err != nil {
		e.Logger.Fatal(errors.Wrap(err, "could not set up db"))
	}
	defer db.Close()
	seedDB(db)


	e.Renderer = &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e.GET("/", root(db))
	e.GET("/like", like(db))
	e.Static("/assets", "public/assets")
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", config.Host, config.Port)))

}
