package persistence

import (
	"errors"
	"fmt"
	"github/com/stereoit/e2etests/pkg/articles"
	"github/com/stereoit/e2etests/pkg/articles/domain"
	"math/rand"
	"sync"
)

type repo struct {
	mu       *sync.Mutex
	articles map[string]*domain.Article
}

func New() articles.ArticleRepo {
	return &repo{
		mu:       &sync.Mutex{},
		articles: map[string]*domain.Article{},
	}
}

func (r *repo) Save(article *domain.Article) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if article.ID == "" {
		article.ID = fmt.Sprintf("%d", rand.Intn(100)+10)
	}
	r.articles[article.ID] = &domain.Article{
		ID:    article.ID,
		Title: article.Title,
		Slug:  article.Slug,
	}

	return article.ID, nil
}

func (r *repo) FindByID(id string) (*domain.Article, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, article := range r.articles {
		if article.ID == id {
			return article, nil
		}
	}
	return nil, errors.New("resource not found")
}

func (r *repo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.articles[id]; !ok {
		return errors.New("article does not exist")
	}

	delete(r.articles, id)
	return nil
}

func (r *repo) Update(article *domain.Article) error {
	if _, err := r.Save(article); err != nil {
		return err
	}
	return nil
}

func (r *repo) List() ([]*domain.Article, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	articles := make([]*domain.Article, len(r.articles))
	i := 0
	for _, article := range r.articles {
		articles[i] = domain.NewArticle(article.ID, article.Title, article.Slug)
		i = i + 1
	}

	return articles, nil
}

func (r *repo) Populate() error {
	// Article fixture data
	var articles = []*domain.Article{
		domain.NewArticle("1", "Zaporizhzhia nuclear plant", "zaporizhzhia-nuclear-plant"),
		domain.NewArticle("2", "Foo Fighters pay tribute", "foo-fighters-pay-tribute"),
		domain.NewArticle("3", "22 of the USA's most underrated destinations", "22-of-the-USA-s-most-underrated-destinations"),
	}
	for _, article := range articles {
		if _, err := r.Save(article); err != nil {
			return err
		}

	}
	return nil
}
