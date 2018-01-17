package cyoa

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.New("tpl0.html").ParseFiles("tpl0.html"))
	// [Min] 用ParseFiles的时候，注意New的名字和ParseFiles中的名字保持一致
	// First template becomes return value if not already defined,
	// and we use that one for subsequent New calls to associate
	// all the templates together. Also, if this file has the same name
	// as t, this file becomes the contents of t, so
	//  t, err := New(name).Funcs(xxx).ParseFiles(name)
	// works. Otherwise we create a new template associated with t.
}

type Adventure map[string]Chapter

type Chapter struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"`
}

type Customize func(h *handler)

func NewHandler(adv Adventure, customizes ...Customize) http.Handler {
	var h = handler{adv, tpl}
	for _, cust := range customizes {
		cust(&h)
	}
	return h
}

type handler struct {
	adv Adventure
	t   *template.Template
}

func CustomizeTemplate(newTpl *template.Template) Customize {
	return func(h *handler) {
		h.t = newTpl
	}
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/" || path == "" {
		path = "/intro"
	}
	if chapter, ok := h.adv[path[1:]]; ok {
		err := h.t.Execute(w, chapter)
		if err != nil {
			log.Printf("%v\n", err)
			http.Error(w, "sorry, the server is not working now...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "page not found.", http.StatusFound)
}

// func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	for k, v := range h.adv {
// 		fmt.Fprintf(w, "Key:%s\nTitle:%s\n", k, v.Title)
// 		for _, v2 := range v.Story {
// 			fmt.Fprintf(w, "Story:%s\n", v2)
// 		}
// 		for _, v3 := range v.Options {
// 			fmt.Fprintf(w, "Options:%s\n%s\n", v3.Text, v3.Arc)
// 		}
// 	}
// }
