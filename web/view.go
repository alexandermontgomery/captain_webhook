package main

import (
	"html/template"
	"net/http"
	"bytes"
	"log"
)

type Page struct{
	Content template.HTML
}

type View struct{
	Templates *template.Template
}

func NewView() *View {
	v := new(View)
	v.Templates = template.Must(template.ParseFiles(homeDir + "/templates/home.html",homeDir + "/templates/page.html"))
	return v
}

func (v *View) RenderPage(w http.ResponseWriter, name string, data interface{}) {
	var content bytes.Buffer
	err := v.Templates.ExecuteTemplate(&content, name + ".html", data)

	var page = Page{
		Content : template.HTML(content.Bytes()),
	}

	if err != nil{
		log.Printf("%+v\n", err)
	}
	v.Templates.ExecuteTemplate(w, "page.html", page)
}