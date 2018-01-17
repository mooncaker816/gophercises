package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

const tpl = `<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8" />
	<title>Choose Your Own Advernture</title>	
</head>
<body>
	<h1>{{ .Title }}</h1>
	{{ range .Story }}
	<p>{{ . }}</p>
	{{ end }}
	<ul>
	{{ range .Options }}
	<li><a href="/{{ .Arc }}">{{ .Text }}</a></li>
	{{ end }}
	</ul>
</body>
</html>
`

var t *template.Template

func init() {
	t, _ = template.New("").Parse(tpl)
}

func main() {
	file, err := os.Open("gopher.json")
	if err != nil {
		log.Println("open input error!")
	}
	defer file.Close()
	var adv Adventure
	err = json.NewDecoder(file).Decode(&adv)
	if err != nil {
		log.Println("json decode error!", err)
	}
	http.Handle("/", http.HandlerFunc(adv.hello))
	log.Fatal(http.ListenAndServe(":8000", nil))
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

func (adv Adventure) hello(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	//fmt.Printf("%+v\n%s", adv, path)
	if chapter, ok := adv[path[1:]]; ok {
		err := t.Execute(w, chapter)
		if err != nil {
			log.Printf("execute template error: %v", err)
			http.Error(w, "sorry the server is down now...", http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "page not found!", http.StatusNotFound)
	}
	return
}
