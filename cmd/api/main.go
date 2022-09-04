package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github/com/stereoit/e2etests/pkg/articles"
	"github/com/stereoit/e2etests/pkg/fly"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	r := chi.NewRouter()
	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", fly.RenderIndex)

	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all your base belongs to us!"))
	})

	// RESTy routes for "articles" resource
	r.Route("/articles", func(r chi.Router) {
		r.Get("/", articles.ListArticles)   // GET /articles
		r.Post("/", articles.CreateArticle) // POST /articles

		// Subrouters:
		r.Route("/{articleID}", func(r chi.Router) {
			r.Use(articles.ArticleCtx)
			r.Get("/", articles.GetArticle)       // GET /articles/123
			r.Put("/", articles.UpdateArticle)    // PUT /articles/123
			r.Delete("/", articles.DeleteArticle) // DELETE /articles/123
		})
	})

	log.Println("listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
