package main

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/boltdb/bolt"
)

type Photo struct {
	ID             string      `json:"id"`
	CreatedAt      string      `json:"created_at"`
	UpdatedAt      string      `json:"updated_at"`
	PromotedAt     interface{} `json:"promoted_at"`
	Width          int         `json:"width"`
	Height         int         `json:"height"`
	Color          string      `json:"color"`
	BlurHash       string      `json:"blur_hash"`
	Description    interface{} `json:"description"`
	AltDescription string      `json:"alt_description"`
	Urls           struct {
		Raw     string `json:"raw"`
		Full    string `json:"full"`
		Regular string `json:"regular"`
		Small   string `json:"small"`
		Thumb   string `json:"thumb"`
	} `json:"urls"`
	Links struct {
		Self             string `json:"self"`
		HTML             string `json:"html"`
		Download         string `json:"download"`
		DownloadLocation string `json:"download_location"`
	} `json:"links"`
	Categories             []interface{} `json:"categories"`
	Likes                  int           `json:"likes"`
	LikedByUser            bool          `json:"liked_by_user"`
	Liked                  bool          `json:"liked"`
	CurrentUserCollections []interface{} `json:"current_user_collections"`
	Sponsorship            struct {
		ImpressionUrls []string `json:"impression_urls"`
		Tagline        string   `json:"tagline"`
		TaglineURL     string   `json:"tagline_url"`
		Sponsor        struct {
			ID              string      `json:"id"`
			UpdatedAt       string      `json:"updated_at"`
			Username        string      `json:"username"`
			Name            string      `json:"name"`
			FirstName       string      `json:"first_name"`
			LastName        interface{} `json:"last_name"`
			TwitterUsername string      `json:"twitter_username"`
			PortfolioURL    string      `json:"portfolio_url"`
			Bio             string      `json:"bio"`
			Location        interface{} `json:"location"`
			Links           struct {
				Self      string `json:"self"`
				HTML      string `json:"html"`
				Photos    string `json:"photos"`
				Likes     string `json:"likes"`
				Portfolio string `json:"portfolio"`
				Following string `json:"following"`
				Followers string `json:"followers"`
			} `json:"links"`
			ProfileImage struct {
				Small  string `json:"small"`
				Medium string `json:"medium"`
				Large  string `json:"large"`
			} `json:"profile_image"`
			InstagramUsername string `json:"instagram_username"`
			TotalCollections  int    `json:"total_collections"`
			TotalLikes        int    `json:"total_likes"`
			TotalPhotos       int    `json:"total_photos"`
			AcceptedTos       bool   `json:"accepted_tos"`
		} `json:"sponsor"`
	} `json:"sponsorship"`
	User struct {
		ID              string      `json:"id"`
		UpdatedAt       string      `json:"updated_at"`
		Username        string      `json:"username"`
		Name            string      `json:"name"`
		FirstName       string      `json:"first_name"`
		LastName        interface{} `json:"last_name"`
		TwitterUsername string      `json:"twitter_username"`
		PortfolioURL    string      `json:"portfolio_url"`
		Bio             string      `json:"bio"`
		Location        interface{} `json:"location"`
		Links           struct {
			Self      string `json:"self"`
			HTML      string `json:"html"`
			Photos    string `json:"photos"`
			Likes     string `json:"likes"`
			Portfolio string `json:"portfolio"`
			Following string `json:"following"`
			Followers string `json:"followers"`
		} `json:"links"`
		ProfileImage struct {
			Small  string `json:"small"`
			Medium string `json:"medium"`
			Large  string `json:"large"`
		} `json:"profile_image"`
		InstagramUsername string `json:"instagram_username"`
		TotalCollections  int    `json:"total_collections"`
		TotalLikes        int    `json:"total_likes"`
		TotalPhotos       int    `json:"total_photos"`
		AcceptedTos       bool   `json:"accepted_tos"`
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

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
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

func FindPhotoByID(db *bolt.DB, id string) (*Photo, error) {
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

func root(db *bolt.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		photos, _ := List(db, "photos")
		return c.Render(http.StatusOK, "main", photos)
	}
}

func like(db *bolt.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		photo, _ := FindPhotoByID(db, c.QueryParam("id"))
		photo.Liked = !photo.Liked
		photo.save(db)

		photos, _ := List(db, "photos")
		return c.Render(http.StatusOK, "main", photos)
	}
}

func List(db *bolt.DB, bucket string) ([]*Photo, error) {
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

func main() {
	e := echo.New()

	db, err := setupDB()
	if err != nil {
		e.Logger.Fatal(errors.Wrap(err, "could not set up db"))
	}
	defer db.Close()

	// seed
	photos, _ := List(db, "photos")
	if len(photos) < 1 {
		seedDB(db)
	}

	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}

	e.Renderer = t
	e.GET("/", root(db))
	e.GET("/like", like(db))
	e.Static("/src", "public")
	e.Logger.Fatal(e.Start(":1323"))

}
