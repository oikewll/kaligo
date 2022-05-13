package controller

import (
	"net/http"

	"github.com/owner888/kaligo"
)

func AddRoutes(r kaligo.Router) {
	r.AddRoute("user/login", map[string]string{
		http.MethodPost: "Login",
	}, &User{})

	r.AddRoute("user/logout", map[string]string{
		http.MethodPost: "Logout",
	}, &User{})
}
