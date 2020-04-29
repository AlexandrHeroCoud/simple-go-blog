package main

import (
	"fmt"
	"go-blog/db/documents"
	"html/template"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"go-blog/models"
	"labix.org/v2/mgo"
)

var postsCollection *mgo.Collection

func indexHandler(rnd render.Render) {
	PostDocument := []documents.PostDocument{}
	postsCollection.Find(nil).All(&PostDocument)
	posts := []models.Post{}
	for _, doc := range PostDocument {
		post := models.Post{doc.Id, doc.Title, doc.ContentHtml, doc.ContentMarkdown}
		posts = append(posts, post)
	}
	rnd.HTML(200, "index", posts)
}

func writeHandler(rnd render.Render) {
	post := models.Post{}
	rnd.HTML(200, "write", post)
}

func editHandler(rnd render.Render, params martini.Params) {
	id := params["id"]
	postDocument := documents.PostDocument{}
	err := postsCollection.FindId(id).One(&postDocument)
	if err != nil {
		rnd.Redirect("/")
		return
	}
	post := models.Post{postDocument.Id, postDocument.Title, postDocument.ContentHtml, postDocument.ContentMarkdown}
	rnd.HTML(200, "write", post)
}

func savePostHandler(rnd render.Render, r *http.Request) {
	id := r.FormValue("id")
	title := r.FormValue("title")
	contentMarkdown := r.FormValue("content")
	contentHtml := ConvertMarkdownToHtml(contentMarkdown)

	postDocument := documents.PostDocument{id, title, contentHtml, contentMarkdown}

	if id != "" {
		postsCollection.UpdateId(id, postDocument)
	} else {
		id = GenerateId()
		postDocument.Id = id
		postsCollection.Insert(postDocument)
	}

	rnd.Redirect("/", 302)
}

func deletePostHandler(rnd render.Render, params martini.Params) {
	id := params["id"]
	if id == "" {
		rnd.Error(404)
	}
	postsCollection.RemoveId(id)
	rnd.Redirect("/", 302)
}

func getHtmlHandler(rnd render.Render, r *http.Request) {
	md := r.FormValue("md")
	html := ConvertMarkdownToHtml(md)

	rnd.JSON(200, map[string]interface{}{"html": html})
}

func unescape(x string) interface{} {
	return template.HTML(x)
}

func main() {
	m := martini.Classic()

	fmt.Println("Listen port: 3000")

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}

	postsCollection = session.DB("go-blog").C("posts")

	unescapeFuncMap := template.FuncMap{"unescape": unescape}

	// ...
	m.Use(render.Renderer(render.Options{
		Directory:  "templates",                // Specify what path to load the templates from.
		Layout:     "layout",                   // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions: []string{".tmpl", ".html"}, // Specify extensions to load for templates.
		Charset:    "UTF-8",                    // Sets encoding for json and html content-types. Default is "UTF-8".
		Funcs:      []template.FuncMap{unescapeFuncMap},
		IndentJSON: true, // Output human readable JSON
		IndentXML:  true, // Output human readable XML
	}))
	// ...
	StaticOptions := martini.StaticOptions{Prefix: "assets"}
	m.Use(martini.Static("assets", StaticOptions))
	m.Get("/", indexHandler)
	m.Get("/write", writeHandler)
	m.Get("/edit/:id", editHandler)
	m.Get("/delete/:id", deletePostHandler)
	m.Post("/SavePost", savePostHandler)
	m.Post("/gethtml", getHtmlHandler)
	m.Run()
}
