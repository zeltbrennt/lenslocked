package controllers

import (
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/zeltbrennt/lenslocked/context"
	"github.com/zeltbrennt/lenslocked/cookie"
	"github.com/zeltbrennt/lenslocked/models"
)

type Users struct {
	Templates struct {
		Signup         Executer
		Signin         Executer
		CurrentUser    Executer
		ForgotPassword Executer
		CheckYourMail  Executer
		ResetPassword  Executer
	}
	UserService          *models.UserService
	SessionService       *models.SessionService
	PasswordResetService *models.PasswordResetService
	EmailService         *models.EmailService
}

func (u Users) SignupPage(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.Signup.Execute(w, r, data)
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
	http.Redirect(w, r, "/user/me", http.StatusFound)
}

func (u Users) SigninPage(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.Signin.Execute(w, r, data)
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
	u.Templates.CurrentUser.Execute(w, r, *user)
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

func (u Users) ForgotPasswordPage(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.ForgotPassword.Execute(w, r, data)
}

func (u Users) HandleForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	pwReset, err := u.PasswordResetService.Create(data.Email)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	vals := url.Values{
		"token": {pwReset.Token},
	}
	// TODO: use config type
	resetURL := os.Getenv("DOMAIN") + "/reset-pw?" + vals.Encode()
	err = u.EmailService.ForgotPassword(data.Email, resetURL)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	u.Templates.CheckYourMail.Execute(w, r, data)
}

func (u Users) ResetPasswordPage(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token string
	}
	data.Token = r.FormValue("token")
	u.Templates.ResetPassword.Execute(w, r, data)
}

func (u Users) HandleResetPassword(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Token    string
		Password string
	}
	data.Token = r.FormValue("token")
	data.Password = r.FormValue("password")
	user, err := u.PasswordResetService.Consume(data.Token)
	if err != nil {
		log.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	cookie.SetCookie(w, cookie.CookieSession, session.NewToken)
	http.Redirect(w, r, "/user/me", http.StatusFound)
}
