package routes

import (
	"github.com/martini-contrib/render"
	"go-blog/session"
	"net/http"
)

func GetLoginHandler(rnd render.Render) {
	rnd.HTML(200, "login", nil)
}

func PostLoginHandler(rnd render.Render, r *http.Request, s *session.Session) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	s.Username = username
	s.Password = password

	rnd.Redirect("/")
}
