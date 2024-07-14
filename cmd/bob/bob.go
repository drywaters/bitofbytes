package main

import (
	"fmt"
	"github.com/DryWaters/bitofbytes/controllers"
	"github.com/DryWaters/bitofbytes/models"
	"github.com/DryWaters/bitofbytes/templates"
	"github.com/DryWaters/bitofbytes/views"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"net/http"
)

func main() {
	// load config
	cfg, err := models.LoadEnvConfig()
	if err != nil {
		panic(err)
	}

	// run with config
	err = run(cfg)
}

func run(cfg models.Config) error {
	// TODO:  Setup database

	// setup services
	postService := models.PostService{
		// Add DB when needed
	}

	// setup CSRF protection
	csrfKey := []byte(cfg.CSRF.Key)
	csrfMw := csrf.Protect(
		csrfKey,
		csrf.Secure(cfg.CSRF.Secure),
		csrf.Path("/"),
	)

	// setup controllers
	blogController := controllers.Blog{
		PostService: postService,
		Templates: controllers.BlogTemplates{
			Index: views.Must(views.ParseFS(templates.FS, "blog/index.templ", "base.templ")),
			Post:  views.Must(views.ParseFS(templates.FS, "blog/post.templ", "base.templ")),
		},
	}

	// Setup our router and routes
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(csrfMw)
	r.Get("/", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "home/index.templ", "home/infocard.templ", "base.templ"))))
	r.Get("/utils", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "utils/index.templ", "base.templ"))))

	r.Get("/blog", blogController.Index)
	r.Get("/posts/{slug}", blogController.Blog)

	staticHandler := http.FileServer(http.Dir("static"))
	r.Get("/static/*", http.StripPrefix("/static", staticHandler).ServeHTTP)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	// Start the server
	fmt.Println("Starting the server on ", cfg.Server.Address)
	return http.ListenAndServe(cfg.Server.Address, r)
}
