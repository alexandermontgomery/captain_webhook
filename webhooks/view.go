package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
)

type Page struct {
	Content template.HTML
}

type View struct {
	Templates *template.Template
}

func NewView() *View {
	v := new(View)
	if env != "dev" {
		v.Templates = template.Must(template.New("").Delims("[[", "]]").ParseFiles(homeDir+"/templates/home.html", homeDir+"/templates/page.html"))
	}
	return v
}

func (v *View) RenderPage(w http.ResponseWriter, name string, data interface{}) {
	var content bytes.Buffer

	// Parse dev files on demand
	if env == "dev" {
		v.Templates = template.Must(template.New("").Delims("[[", "]]").ParseFiles(homeDir+"/templates/page.html", homeDir+"/templates/"+name+".html"))
	}

	err := v.Templates.ExecuteTemplate(&content, name+".html", data)
	var page = Page{
		Content: template.HTML(content.Bytes()),
	}
	if err != nil {
		log.Printf("%+v\n", err)
	}
	v.Templates.ExecuteTemplate(w, "page.html", page)
}

func (v *View) ServeSingleTemplate(w http.ResponseWriter, name string) {
	v.Templates.ExecuteTemplate(w, name, nil)
}
