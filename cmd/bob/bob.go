package main

import (
	"fmt"
	"github.com/DryWaters/bitofbytes/controllers"
	"github.com/DryWaters/bitofbytes/controllers/middleware"
	"github.com/DryWaters/bitofbytes/models"
	"github.com/DryWaters/bitofbytes/templates"
	"github.com/DryWaters/bitofbytes/views"
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
	csrfMw := middleware.CSRF(cfg.CSRF.Key, cfg.CSRF.Secure)
	cspMw := middleware.SecureHeaders

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
	r := http.NewServeMux()
	csrfRouter := csrfMw(r)
	secureRouter := cspMw(csrfRouter)
	r.HandleFunc("GET /", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "home/index.gohtml", "home/infocard.gohtml", "base.gohtml"))))

	// Healthcheck
	r.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// Utils
	r.HandleFunc("GET /utils", utilsController.Index)

	// Base64 Utils
	r.HandleFunc("GET /utils/base64/encode", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "utils/base64/encode.gohtml", "base.gohtml"))))
	r.HandleFunc("POST /utils/base64/encode", utilsController.Encode)
	r.HandleFunc("GET /utils/base64/decode", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "utils/base64/decode.gohtml", "base.gohtml"))))
	r.HandleFunc("POST /utils/base64/decode", utilsController.Decode)

	// Blog
	r.HandleFunc("GET /blog", blogController.Index)
	r.HandleFunc("GET /posts/{slug}", blogController.Blog)

	staticHandler := http.FileServer(http.Dir("static"))
	r.Handle("GET /static/", http.StripPrefix("/static/", staticHandler))

	// Start the server
	fmt.Println("Starting the server on ", cfg.Server.Address)
	return http.ListenAndServe(cfg.Server.Address, secureRouter)
}
