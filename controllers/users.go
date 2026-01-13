package controllers

import (
	"fmt"
	"net/http"

	"github.com/zeltbrennt/lenslocked/models"
)

type Users struct {
	Templates struct {
		Signup Executer
		Signin Executer
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
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		// TODO: Warning here
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	setCookie(w, cookieSession, session.NewToken)
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
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	setCookie(w, cookieSession, session.NewToken)
	fmt.Fprintf(w, "User authenticated: %+v", user)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	token, err := readCookie(r, "session")
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	user, err := u.SessionService.User(token)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	fmt.Fprintf(w, "current User: %s\n", user.Email)
}
