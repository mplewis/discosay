package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/mplewis/discosay/lib/responder"
	"gopkg.in/yaml.v3"
)

func env(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Missing mandatory environment variable: %s", key)
	}
	return val
}

func parseConfig(rawYaml []byte) ([]responder.Responder, error) {
	config := make(map[string]interface{})
	err := yaml.Unmarshal(rawYaml, config)
	if err != nil {
		return nil, err
	}

	templates := map[string]string{}
	for name, template := range config["templates"].(map[string]interface{}) {
		templates[name] = template.(string)
	}

	responders := []responder.Responder{}
	for name, rspec := range config["responders"].(map[string]interface{}) {
		spec := rspec.(map[string]interface{})
		r, err := responder.Parse(name, spec)
		if err != nil {
			return nil, err
		}

		if r.TemplateName != nil {
			template := templates[*r.TemplateName]
			r.Template = &template
		}
		responders = append(responders, *r)
	}
	return responders, err
}

func test(resp responder.Responder, msg string) {
	response := resp.Respond(msg)
	if response != nil {
		log.Printf("-> %s: %s", *resp.Name, *response)
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	configPath := env("CONFIG_PATH")
	rawYaml, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	responders, err := parseConfig(rawYaml)
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range responders {
		test(r, "!riir")
		test(r, "!retf")
		test(r, "!retf Rust rules!")
		test(r, "!gotime")
		test(r, "!gotime yay for go!")
	}
}
