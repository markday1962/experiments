package main

import (
	"log"
	"os"
	"text/template"
)

var tpl *template.Template

type chart struct {
	Chart string
	Version string
	AppVersion string
	Type string
}

func init() {
	tpl = template.Must(template.ParseFiles("Chart.tpl"))
}

func main() {

	name := "backend-service"

	c := chart{
		Chart: name,
		Type: "application",
		Version: "0.1.0",
		AppVersion: "0.1.0",
	}

	err := tpl.Execute(os.Stdout, c)
	if err != nil {
		log.Fatalln(err)
	}
}
