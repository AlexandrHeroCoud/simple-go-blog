package routes

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"go-blog/db/documents"
	"go-blog/models"
	"go-blog/session"
	"go-blog/utils"
	"labix.org/v2/mgo"
	"net/http"
)

func WriteHandler(rnd render.Render, s *session.Session) {
	if !s.IsAuthorized {
		rnd.Redirect("/")
	}
	model := models.EditPostModel{}
	model.IsAutorized = s.IsAuthorized
	model.Post = models.Post{}
	rnd.HTML(200, "write", model)
}

func EditHandler(rnd render.Render, params martini.Params, db *mgo.Database, s *session.Session) {
	if !s.IsAuthorized {
		rnd.Redirect("/")
	}
	postsCollection := db.C("Posts")
	id := params["id"]
	postDocument := documents.PostDocument{}
	err := postsCollection.FindId(id).One(&postDocument)
	if err != nil {
		rnd.Redirect("/")
		return
	}
	post := models.Post{postDocument.Id, postDocument.Title, postDocument.ContentHtml, postDocument.ContentMarkdown}
	model := models.EditPostModel{}
	model.IsAutorized = s.IsAuthorized
	model.Post = post
	rnd.HTML(200, "write", model)
}

func ViewHandler(rnd render.Render, params martini.Params, db *mgo.Database, s *session.Session) {
	postsCollection := db.C("Posts")
	id := params["id"]
	postDocument := documents.PostDocument{}
	err := postsCollection.FindId(id).One(&postDocument)
	if err != nil {
		rnd.Redirect("/")
		return
	}
	post := models.Post{postDocument.Id, postDocument.Title, postDocument.ContentHtml, postDocument.ContentMarkdown}
	model := models.EditPostModel{}
	model.IsAutorized = s.IsAuthorized
	model.Post = post
	rnd.HTML(200, "view", model)
}

func SavePostHandler(rnd render.Render, r *http.Request, db *mgo.Database) {

	id := r.FormValue("id")
	title := r.FormValue("title")
	contentMarkdown := r.FormValue("content")
	contentHtml := utils.ConvertMarkdownToHtml(contentMarkdown)

	postDocument := documents.PostDocument{id, title, contentHtml, contentMarkdown}

	postsCollection := db.C("Posts")
	if id != "" {
		postsCollection.UpdateId(id, postDocument)
	} else {
		id = utils.GenerateId()
		postDocument.Id = id
		postsCollection.Insert(postDocument)
	}

	rnd.Redirect("/")
}

func DeletePostHandler(rnd render.Render, params martini.Params, db *mgo.Database, s *session.Session) {
	if !s.IsAuthorized {
		rnd.Redirect("/")
	}
	id := params["id"]
	if id == "" {
		rnd.Error(404)
	}
	postsCollection := db.C("Posts")
	postsCollection.RemoveId(id)
	rnd.Redirect("/")
}

func GetHtmlHandler(rnd render.Render, r *http.Request) {
	md := r.FormValue("md")
	html := utils.ConvertMarkdownToHtml(md)

	rnd.JSON(200, map[string]interface{}{"html": html})
}
