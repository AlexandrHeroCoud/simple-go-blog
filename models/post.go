package models

type Post struct {
	Id              string
	Title           string
	ContentHtml     string
	ContentMarkdown string
}

func NewPost(id, title, ContentHtml, ContentMarkdown string) *Post {
	return &Post{id, title, ContentHtml, ContentMarkdown}
}
