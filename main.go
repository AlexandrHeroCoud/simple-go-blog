package main

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/russross/blackfriday"
	"html/template"
	"net/http"

	"go-blog/models"
)

var posts map[string]*models.Post
var counter int

func indexHandler(rnd render.Render) {
	rnd.HTML(200, "index", posts)
}

func writeHandler(rnd render.Render) {
	post := models.Post{}
	rnd.HTML(200, "write", post)
}

func editHandler(rnd render.Render, r *http.Request, params martini.Params) {
	id := params["id"]

	post, found := posts[id]
	if !found {
		rnd.Error(404)
	}
	rnd.HTML(200, "write", post)
}

func savePostHandler(rnd render.Render, r *http.Request) {
	id := r.FormValue("id")
	title := r.FormValue("title")
	contentMarkdown := r.FormValue("content")
	contentHtml := string(blackfriday.MarkdownBasic([]byte(contentMarkdown)))
	r.FormValue("content")
	var post *models.Post

	if id != "" {
		post = posts[id]
		post.Title = title
		post.ContentHtml = contentHtml
		post.ContentMarkdown = contentMarkdown
	} else {
		id = GenerateId()
		post := models.NewPost(id, title, contentHtml, contentMarkdown)
		posts[post.Id] = post
	}

	rnd.Redirect("/", 302)
}

func deletePostHandler(rnd render.Render, params martini.Params) {
	id := params["id"]
	if id == "" {
		rnd.Error(404)
	} else {
		delete(posts, id)
	}
	rnd.Redirect("/", 302)
}

func getHtmlHandler(rnd render.Render, r *http.Request) {
	md := r.FormValue("md")
	htmlBytes := blackfriday.MarkdownBasic([]byte(md))

	rnd.JSON(200, map[string]interface{}{"html": string(htmlBytes)})
}

func unescape(x string) interface{} {
	return template.HTML(x)
}

func main() {
	m := martini.Classic()

	fmt.Println("Listen port: 3000")

	posts = make(map[string]*models.Post, 0)

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
	m.Use(func(r *http.Request) {
		if r.URL.Path == "/write" {
			counter++
		}
	})
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
