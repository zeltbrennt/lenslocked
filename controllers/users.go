package controllers

import (
	"log"
	"net/http"

	"github.com/zeltbrennt/lenslocked/context"
	"github.com/zeltbrennt/lenslocked/cookie"
	"github.com/zeltbrennt/lenslocked/models"
)

type Users struct {
	Templates struct {
		Signup      Executer
		Signin      Executer
		CurrentUser Executer
	}
	UserService    *models.UserService
	SessionService *models.SessionService
}

func (u Users) SignupPage(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.Signup.Execute(w, data)
}

func (u Users) HandleSignup(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := u.UserService.Create(email, password)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		log.Println(err)
		// TODO: Warning here
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	cookie.SetCookie(w, cookie.CookieSession, session.NewToken)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) SigninPage(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.Signin.Execute(w, data)
}

func (u Users) HandleSignin(w http.ResponseWriter, r *http.Request) {
	var data struct {
		email    string
		password string
	}
	data.email = r.FormValue("email")
	data.password = r.FormValue("password")
	user, err := u.UserService.Authenticate(data.email, data.password)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	cookie.SetCookie(w, cookie.CookieSession, session.NewToken)
	http.Redirect(w, r, "/user/me", http.StatusFound)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	u.Templates.CurrentUser.Execute(w, *user)
}

func (u Users) HandleSignOut(w http.ResponseWriter, r *http.Request) {
	token, err := cookie.ReadCookie(r, cookie.CookieSession)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	err = u.SessionService.Delete(token)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	cookie.DeleteCookie(w, cookie.CookieSession)
	http.Redirect(w, r, "/signin", http.StatusFound)
}
