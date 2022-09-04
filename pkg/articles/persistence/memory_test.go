package persistence_test

import (
	"github/com/stereoit/e2etests/pkg/articles/domain"
	"github/com/stereoit/e2etests/pkg/articles/persistence"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	assert := assert.New(t)
	repo := persistence.New()
	assert.NotNil(repo)
}

func Test_Save(t *testing.T) {
	assert := assert.New(t)
	article := &domain.Article{
		ID: "2342",
	}
	repo := persistence.New()
	id, err := repo.Save(article)
	assert.Nil(err)
	assert.NotEmpty(id)
}

func Test_FindbyID(t *testing.T) {
	assert := assert.New(t)
	id := "123"

	repo := persistence.New()
	article, err := repo.FindByID(id)
	assert.Nil(err)
	assert.Nil(article)
	// assert.Equal(article.ID, id)

	testUser := &domain.Article{
		ID:    "",
		Title: "Our test title",
		Slug:  "our-test-slug",
	}

	id, err = repo.Save(testUser)
	assert.Nil(err)
	got, err := repo.FindByID(id)
	assert.Nil(err)
	assert.Equal(got.Title, testUser.Title)
}

func Test_Delete(t *testing.T) {
	assert := assert.New(t)
	repo := persistence.New()

	testArticle := &domain.Article{
		ID:    "",
		Title: "Our test title",
		Slug:  "our-test-slug",
	}
	id, err := repo.Save(testArticle)
	assert.Nil(err)
	err = repo.Delete(id)
	assert.Nil(err, "repository should delete existing user")

	got, err := repo.FindByID(id)
	assert.Nil(got)
	assert.Nil(err)

	err = repo.Delete("missing-id")
	assert.NotNil(err, "missing user should throw error")
}

func Test_Update(t *testing.T) {
	assert := assert.New(t)
	repo := persistence.New()
	testArticle := &domain.Article{
		ID:    "",
		Title: "Our test title",
		Slug:  "our-test-slug",
	}
	id, err := repo.Save(testArticle)
	assert.Nil(err)

	testArticle.ID = id
	testArticle.Title = "Test title"

	err = repo.Update(testArticle)
	assert.Nil(err)

}

func Test_List(t *testing.T) {
	assert := assert.New(t)
	repo := persistence.New()

	articles, err := repo.List()
	assert.Nil(err)
	assert.NotEmpty(articles)
}
