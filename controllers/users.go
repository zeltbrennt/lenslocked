package controllers

import (
	"fmt"
	"net/http"

	"github.com/zeltbrennt/lenslocked/templates/models"
)

type Users struct {
	Templates struct {
		Signup Executer
		Signin Executer
	}
	UserService *models.UserService
}

// GET
func (u Users) Signup(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.Signup.Execute(w, data)
}

// POST
func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	user, err := u.UserService.Create(email, password)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "user created: %+v", user)
}

// Get
func (u Users) Signin(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.Signin.Execute(w, data)
}
