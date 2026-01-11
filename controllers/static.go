package controllers

import (
	"html/template"
	"net/http"

	"github.com/zeltbrennt/lenslocked/views"
)

type Static struct {
	Template views.Template
}

func StaticHandler(tpl views.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, nil)
	}
}

func FAQ(tpl views.Template) http.HandlerFunc {
	questions := []struct {
		Question string
		Answer   template.HTML
	}{
		{
			Question: "Is there a free version?",
			Answer:   "Yes! W offer a free trial of 30 days on any paid plans",
		},
		{
			Question: "What are your support hours?",
			Answer:   "We have support staff answering emails 24/7",
		},
		{
			Question: "How do I contact support?",
			Answer:   `Email us - <a href="mailto:support@lenslocked.com">support@lenslocked.com</a>`,
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, questions)
	}
}
