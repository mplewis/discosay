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

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	configPath := env("CONFIG_PATH")
	rawYaml, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}
	config := make(map[string]interface{})
	err = yaml.Unmarshal(rawYaml, config)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%#v\n", config)

	// templates := config["templates"]
	// log.Printf("%#v\n", templates)

	// bots := config["bots"]
	// log.Printf("%#v\n", bots)

	responders := config["responders"].(map[string]interface{})
	// log.Printf("%#v\n", responders)

	for name, rspec := range responders {
		spec := rspec.(map[string]interface{})
		log.Printf("%#v: %#v\n", name, spec)
		r, err := responder.Parse(&name, spec)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%#v\n", r)

		response := r.Respond("!riir")
		if response != nil {
			log.Printf("-> %s: %s", *r.Name, *response)
		}

		response = r.Respond("!retf")
		if response != nil {
			log.Printf("-> %s: %s", *r.Name, *response)
		}

		response = r.Respond("!retf Rust rules!")
		if response != nil {
			log.Printf("-> %s: %s", *r.Name, *response)
		}

		response = r.Respond("!gotime")
		if response != nil {
			log.Printf("-> %s: %s", *r.Name, *response)
		}

		response = r.Respond("!gotime yay for go!")
		if response != nil {
			log.Printf("-> %s: %s", *r.Name, *response)
		}
	}

	// x, err := regexp.Compile("^!cheer$")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// r := responder.Responder{
	// 	Name:      "cheer",
	// 	Match:     *x,
	// 	Responses: []string{"yay!", "hooray!"},
	// }
	// log.Println(*r.Respond("!cheer"))

	// x, err = regexp.Compile("^!say (.+)$")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// r = responder.Responder{
	// 	Name:     "say",
	// 	Match:    *x,
	// 	Template: "$MSG!",
	// }
	// log.Println(*r.Respond("!say Hi"))
}
