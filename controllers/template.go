package controllers

import "net/http"

type Executer interface {
	Execute(http.ResponseWriter, any)
}
