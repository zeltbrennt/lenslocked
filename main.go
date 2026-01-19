package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/zeltbrennt/lenslocked/controllers"
	mw "github.com/zeltbrennt/lenslocked/middleware"
	"github.com/zeltbrennt/lenslocked/migrations"
	"github.com/zeltbrennt/lenslocked/models"
	"github.com/zeltbrennt/lenslocked/templates"
	"github.com/zeltbrennt/lenslocked/views"
)

func main() {
	// env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// setup DB
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// setup services
	userService := models.UserService{
		DB: db,
	}
	sessionService := models.SessionService{
		DB: db,
		TM: models.TokenManager{BytesPerToken: 32},
	}
	mailService := models.NewMailService(models.SMTPConfigFromEnv())

	// setup middleware
	protection := http.NewCrossOriginProtection()
	protection.AddTrustedOrigin("http://localhost:5173")
	umw := mw.UserMiddleware{
		SessionService: &sessionService,
	}

	// setup controllers
	userController := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}
	userController.Templates.Signup = views.Must(views.ParseFS(templates.FS, "signup.html", "layout.html"))
	userController.Templates.Signin = views.Must(views.ParseFS(templates.FS, "signin.html", "layout.html"))
	userController.Templates.CurrentUser = views.Must(views.ParseFS(templates.FS, "currentUser.html", "layout.html"))
	// TODO: implement a test to make sure, all Services and Templates are set

	// setup router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(protection.Handler)
	r.Use(mw.SetHeaders)
	r.Use(umw.SetUser)

	// setup routes
	// static
	r.Get("/", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "home.html", "layout.html"))))
	r.Get("/contact", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "contact.html", "layout.html"))))
	r.Get("/faq", controllers.FAQ(
		views.Must(views.ParseFS(templates.FS, "faq.html", "layout.html"))))
	// dynamic
	r.Get("/signin", userController.SigninPage)
	r.Post("/signin", userController.HandleSignin)
	r.Get("/signup", userController.SignupPage)
	r.Post("/signup", userController.HandleSignup)
	r.Route("/user/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", userController.CurrentUser)
	})
	r.Post("/signout", userController.HandleSignOut)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	log.Println("Starting Server on :3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
