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
	base64Service := models.Base64Service{
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
			Index: views.Must(views.ParseFS(templates.FS, "blog/index.gohtml", "base.gohtml")),
			Post:  views.Must(views.ParseFS(templates.FS, "blog/post.gohtml", "base.gohtml")),
		},
	}

	utilsController := controllers.Utils{
		Base64Service: base64Service,
		Templates: controllers.UtilsTemplates{
			Index: views.Must(views.ParseFS(templates.FS, "utils/index.gohtml", "base.gohtml")),
			Base64: controllers.Base64Templates{
				Base64Response: views.Must(views.ParseFS(templates.FS, "utils/base64/base64_response.gohtml")),
			},
		},
	}

	// Setup our router and routes
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(csrfMw)
	r.Get("/", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "home/index.gohtml", "home/infocard.gohtml", "base.gohtml"))))

	// Utils
	r.Get("/utils", utilsController.Index)

	// Base64 Utils
	r.Get("/utils/base64/encode", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "utils/base64/encode.gohtml", "base.gohtml"))))
	r.Post("/utils/base64/encode", utilsController.Encode)
	r.Get("/utils/base64/decode", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "utils/base64/decode.gohtml", "base.gohtml"))))
	r.Post("/utils/base64/decode", utilsController.Decode)

	// Blog
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
