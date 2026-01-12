package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/zeltbrennt/lenslocked/controllers"
	"github.com/zeltbrennt/lenslocked/templates"
	"github.com/zeltbrennt/lenslocked/views"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "home.html", "layout.html"))))
	r.Get("/contact", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "contact.html", "layout.html"))))
	r.Get("/faq", controllers.FAQ(
		views.Must(views.ParseFS(templates.FS, "faq.html", "layout.html"))))

	var userController controllers.Users
	userController.Templates.Signup = views.Must(views.ParseFS(templates.FS, "signup.html", "layout.html"))

	r.Get("/signup", userController.Signup)
	r.Post("/signup", userController.Create)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Println("Starting Server on :3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
