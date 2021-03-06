package routes

import (
	"fmt"
	"github.com/martini-contrib/render"
	"go-blog/db/documents"
	"go-blog/models"
	"go-blog/session"
	"labix.org/v2/mgo"
)

func IndexHandler(rnd render.Render, s *session.Session, db *mgo.Database) {
	PostDocument := []documents.PostDocument{}
	postsCollection := db.C("Posts")
	postsCollection.Find(nil).All(&PostDocument)
	posts := []models.Post{}
	for _, doc := range PostDocument {
		post := models.Post{doc.Id, doc.Title, doc.ContentHtml, doc.ContentMarkdown}
		posts = append(posts, post)
	}

	model := models.PostsListModel{}
	model.IsAuthorized = s.IsAuthorized
	model.Posts = posts
	fmt.Println(model)
	rnd.HTML(200, "index", model)
}
