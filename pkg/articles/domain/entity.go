package domain

// Article data model. I suggest looking at https://upper.io for an easy
// and powerful data persistence adapter.
type Article struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Slug  string `json:"slug"`
}

func NewArticle(id, title, slug string) *Article {
	return &Article{
		ID:    id,
		Title: title,
		Slug:  slug,
	}
}
