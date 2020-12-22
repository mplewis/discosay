package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/mplewis/discosay/lib/bot"
	"github.com/mplewis/discosay/lib/responder"
	"gopkg.in/yaml.v3"
)

type botSpec struct {
	name       string
	responders []*responder.Responder
}

func env(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Missing mandatory environment variable: %s", key)
	}
	return val
}

func parseConfig(rawYaml []byte) ([]botSpec, error) {
	config := make(map[string]interface{})
	err := yaml.Unmarshal(rawYaml, config)
	if err != nil {
		return nil, err
	}

	templates := map[string]string{}
	for name, template := range config["templates"].(map[string]interface{}) {
		templates[name] = template.(string)
	}

	responders := map[string]responder.Responder{}
	for name, rspec := range config["responders"].(map[string]interface{}) {
		spec := rspec.(map[string]interface{})
		r, err := responder.New(name, spec)
		if err != nil {
			return nil, err
		}

		if r.TemplateName != nil {
			template := templates[*r.TemplateName]
			r.Template = &template
		}
		responders[*r.Name] = *r
	}

	botSpecs := []botSpec{}
	for name, respNames := range config["bots"].(map[string]interface{}) {

		s := botSpec{name: name}
		for _, rawName := range respNames.([]interface{}) {
			responder := responders[rawName.(string)]
			s.responders = append(s.responders, &responder)
		}
		botSpecs = append(botSpecs, s)
	}
	return botSpecs, nil
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

	botSpecs, err := parseConfig(rawYaml)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Registered bots:")
	for _, botSpec := range botSpecs {
		fmt.Println(botSpec.name)
		for _, resp := range botSpec.responders {
			fmt.Printf("    %s\n", resp)
		}
	}

	fmt.Println("Connecting...")
	bots := []*bot.Bot{}
	for _, botSpec := range botSpecs {
		authToken := env(fmt.Sprintf("%s_AUTH_TOKEN", strings.ToUpper(botSpec.name)))
		b, err := bot.New(botSpec.name, authToken, botSpec.responders)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s connected.", *b.Name)
		bots = append(bots, b)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	fmt.Println("Shutting down...")
	for _, b := range bots {
		if err := b.Close(); err != nil {
			log.Println(err)
		}
	}
}
