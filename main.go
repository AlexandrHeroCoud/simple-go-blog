package main

import (
	"go-blog/routes"
	"go-blog/session"
	"html/template"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"labix.org/v2/mgo"
)

func unescape(x string) interface{} {
	return template.HTML(x)
}
func main() {
	m := martini.Classic()

	sessionDB, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}

	db := sessionDB.DB("go-blog")
	m.Map(db)
	m.Use(session.Middleware)
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
	m.Get("/", routes.IndexHandler)
	m.Get("/login", routes.GetLoginHandler)
	m.Post("/login", routes.PostLoginHandler)
	m.Get("/write", routes.WriteHandler)
	m.Get("/edit/:id", routes.EditHandler)
	m.Get("/delete/:id", routes.DeletePostHandler)
	m.Post("/SavePost", routes.SavePostHandler)
	m.Post("/gethtml", routes.GetHtmlHandler)
	m.Run()
}
