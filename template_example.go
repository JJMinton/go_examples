package main

import (
"log"
"net/http"
"io/ioutil"
"encoding/json"
"html/template"
)

//Config
type Config struct {
    RootURL string "json:rootURL"
    Port string "json:port"
}
var config Config

func init() {
    file, err := ioutil.ReadFile("config.json")
    if err != nil {
        log.Fatal("Can't open config")
    }
    json.Unmarshal(file, &config)
}

//Router
func main() {
    http.HandleFunc("/", rootHandler)
    http.ListenAndServe(config.Port, nil)
}

//Templated page
func rootHandler(res http.ResponseWriter, req *http.Request) {
    title := req.URL.Path[1:]
    tpl, _ := template.ParseFiles("templates/home.html")
    tpl.Execute(res, title)
}
