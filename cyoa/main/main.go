package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/mooncaker816/gophercises/cyoa"
)

func main() {
	file, err := os.Open("gopher.json")
	if err != nil {
		log.Println("open input error!")
		return
	}
	defer file.Close()
	var adv cyoa.Adventure
	err = json.NewDecoder(file).Decode(&adv)
	if err != nil {
		log.Println("json decode error!", err)
		return
	}
	mytpl := template.Must(template.New("tpl1.html").ParseFiles("tpl1.html"))
	myhandler := cyoa.NewHandler(adv, cyoa.CustomizeTemplate(mytpl))
	log.Fatal(http.ListenAndServe(":8000", myhandler))
}
