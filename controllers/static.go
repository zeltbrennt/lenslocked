package controllers

import (
	"html/template"
	"net/http"
)

func StaticHandler(tpl Executer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r, nil)
	}
}

func FAQ(tpl Executer) http.HandlerFunc {
	questions := []struct {
		Question string
		Answer   template.HTML
	}{
		{
			Question: "Is there a free version?",
			Answer:   "Yes! We offer a free trial of 30 days on any paid plans",
		},
		{
			Question: "What are your support hours?",
			Answer:   "We have support staff answering emails 24/7",
		},
		{
			Question: "How do I contact support?",
			Answer:   `Email us - <a class="underline" href="mailto:support@lenslocked.com">support@lenslocked.com</a>`,
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r, questions)
	}
}
