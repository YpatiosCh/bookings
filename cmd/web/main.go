package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/YpatiosCh/bookings/pkg/config"
	"github.com/YpatiosCh/bookings/pkg/handlers"
	"github.com/YpatiosCh/bookings/pkg/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig

var session *scs.SessionManager

func main() {
	//change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateChache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	fmt.Printf("Starting application on port %s\n", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: Routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
