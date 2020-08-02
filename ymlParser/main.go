package main
//https://stackoverflow.com/questions/26290485/golang-yaml-reading-with-map-of-maps
import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

type Config struct {
	Description string
	Services map[string] string
}

func main() {
	var updatedService = "frontend-service"
	var newTag = "2.10.9"

	filename, _ := filepath.Abs("./tags.yaml")
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	var config Config

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	for k, _ := range config.Services{
		if k == updatedService{
			config.Services[k] = newTag
		}
	}

	for k, v := range config.Services{
		fmt.Printf("%#v,%v\n",k,v)
	}

	//write back to file


}