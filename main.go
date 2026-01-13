package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/zeltbrennt/lenslocked/controllers"
	mw "github.com/zeltbrennt/lenslocked/middleware"
	"github.com/zeltbrennt/lenslocked/models"
	"github.com/zeltbrennt/lenslocked/templates"
	"github.com/zeltbrennt/lenslocked/views"
)

func main() {
	// setup controllers
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	userService := models.UserService{
		DB: db,
	}
	var userController controllers.Users
	userController.Templates.Signup = views.Must(views.ParseFS(templates.FS, "signup.html", "layout.html"))
	userController.Templates.Signin = views.Must(views.ParseFS(templates.FS, "signin.html", "layout.html"))
	userController.UserService = &userService

	// setup router
	r := chi.NewRouter()
	protection := http.NewCrossOriginProtection()
	r.Use(protection.Handler)
	r.Use(middleware.Logger)

	r.Get("/", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "home.html", "layout.html"))))
	r.Get("/contact", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "contact.html", "layout.html"))))
	r.Get("/faq", controllers.FAQ(
		views.Must(views.ParseFS(templates.FS, "faq.html", "layout.html"))))
	r.Get("/signin", userController.SigninPage)
	r.Post("/signin", userController.HandleSignin)
	r.Get("/signup", userController.SignupPage)
	r.Post("/signup", userController.HandleSignup)
	r.Get("/cookie", userController.CurrentUser)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Println("Starting Server on :3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
