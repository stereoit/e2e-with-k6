package articles

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github/com/stereoit/e2etests/pkg/articles/domain"
	"github/com/stereoit/e2etests/pkg/rest"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type ArticleSVC interface {
	ArticleCtx(next http.Handler) http.Handler
	CreateArticle(http.ResponseWriter, *http.Request)
	DeleteArticle(http.ResponseWriter, *http.Request)
	UpdateArticle(http.ResponseWriter, *http.Request)
	ListArticles(http.ResponseWriter, *http.Request)
	GetArticle(http.ResponseWriter, *http.Request)
}

type articleSVC struct {
	repo ArticleRepo
}

func New(repo ArticleRepo) ArticleSVC {
	return &articleSVC{
		repo: repo,
	}
}

func (s *articleSVC) ListArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := s.repo.List()
	if err != nil {
		render.Render(w, r, rest.ErrRender(err))
		return
	}

	if err := render.RenderList(w, r, NewArticleListResponse(articles)); err != nil {
		render.Render(w, r, rest.ErrRender(err))
		return
	}
}

type ArticleKey string

var articleKey ArticleKey = "article"

// ArticleCtx middleware is used to load an Article object from
// the URL parameters passed through as the request. In case
// the Article could not be found, we stop here and return a 404.
func (s *articleSVC) ArticleCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var article *domain.Article
		var err error

		if articleID := chi.URLParam(r, "articleID"); articleID != "" {
			article, err = s.dbGetArticle(articleID)
		} else {
			render.Render(w, r, rest.ErrNotFound)
			return
		}
		if err != nil {
			render.Render(w, r, rest.ErrNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), articleKey, article)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// CreateArticle persists the posted Article and returns it
// back to the client as an acknowledgement.
func (s *articleSVC) CreateArticle(w http.ResponseWriter, r *http.Request) {
	data := &ArticleRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, rest.ErrInvalidRequest(err))
		return
	}

	article := data.Article
	s.dbNewArticle(article)

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewArticleResponse(article))
}

// GetArticle returns the specific Article. You'll notice it just
// fetches the Article right off the context, as its understood that
// if we made it this far, the Article must be on the context. In case
// its not due to a bug, then it will panic, and our Recoverer will save us.
func (s *articleSVC) GetArticle(w http.ResponseWriter, r *http.Request) {
	// Assume if we've reach this far, we can access the article
	// context because this handler is a child of the ArticleCtx
	// middleware. The worst case, the recoverer middleware will save us.
	article := r.Context().Value(articleKey).(*domain.Article)

	if err := render.Render(w, r, NewArticleResponse(article)); err != nil {
		render.Render(w, r, rest.ErrRender(err))
		return
	}
}

// UpdateArticle updates an existing Article in our persistent store.
func (s *articleSVC) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	article := r.Context().Value(articleKey).(*domain.Article)

	data := &ArticleRequest{Article: article}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, rest.ErrInvalidRequest(err))
		return
	}
	article = data.Article
	err := s.repo.Update(article)
	if err != nil {
		render.Render(w, r, rest.ErrRender(err))
		return
	}

	render.Render(w, r, NewArticleResponse(article))
}

// DeleteArticle removes an existing Article from our persistent store.
func (s *articleSVC) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	var err error

	// Assume if we've reach this far, we can access the article
	// context because this handler is a child of the ArticleCtx
	// middleware. The worst case, the recoverer middleware will save us.
	article := r.Context().Value(articleKey).(*domain.Article)

	article, err = s.dbRemoveArticle(article.ID)
	if err != nil {
		render.Render(w, r, rest.ErrInvalidRequest(err))
		return
	}

	render.Render(w, r, NewArticleResponse(article))
}

// ArticleRequest is the request payload for Article data model.
//
// NOTE: It's good practice to have well defined request and response payloads
// so you can manage the specific inputs and outputs for clients, and also gives
// you the opportunity to transform data on input or output, for example
// on request, we'd like to protect certain fields and on output perhaps
// we'd like to include a computed field based on other values that aren't
// in the data model. Also, check out this awesome blog post on struct composition:
// http://attilaolah.eu/2014/09/10/json-and-struct-composition-in-go/
type ArticleRequest struct {
	*domain.Article

	ProtectedID string `json:"id"` // override 'id' json to have more control
}

func (a *ArticleRequest) Bind(r *http.Request) error {
	// a.Article is nil if no Article fields are sent in the request. Return an
	// error to avoid a nil pointer dereference.
	if a.Article == nil {
		return errors.New("missing required Article fields")
	}

	// just a post-process after a decode..
	a.ProtectedID = ""                                 // unset the protected ID
	a.Article.Title = strings.ToLower(a.Article.Title) // as an example, we down-case
	return nil
}

// ArticleResponse is the response payload for the Article data model.
// See NOTE above in ArticleRequest as well.
//
// In the ArticleResponse object, first a Render() is called on itself,
// then the next field, and so on, all the way down the tree.
// Render is called in top-down order, like a http handler middleware chain.
type ArticleResponse struct {
	*domain.Article
}

func NewArticleResponse(article *domain.Article) *ArticleResponse {
	resp := &ArticleResponse{Article: article}

	return resp
}

func (rd *ArticleResponse) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

func NewArticleListResponse(articles []*domain.Article) []render.Renderer {
	list := []render.Renderer{}
	for _, article := range articles {
		list = append(list, NewArticleResponse(article))
	}
	return list
}

type ArticleRepo interface {
	Save(*domain.Article) (string, error)
	FindByID(string) (*domain.Article, error)
	Delete(string) error
	Update(*domain.Article) error
	List() ([]*domain.Article, error)

	// Populate
	Populate() error
}

func (s *articleSVC) dbNewArticle(article *domain.Article) (string, error) {
	id, err := s.repo.Save(article)
	if err != nil {
		return "", err
	}

	// articles = append(articles, article)
	return id, nil
}

func (s *articleSVC) dbGetArticle(id string) (*domain.Article, error) {
	article, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("article not found")
	}
	return article, nil
}

func (s *articleSVC) dbRemoveArticle(id string) (*domain.Article, error) {
	article, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("article not found")
	}

	err = s.repo.Delete(id)
	if err != nil {
		return nil, err
	}

	return article, nil

}
