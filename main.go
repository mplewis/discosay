package main

import (
	"log"
	"math/rand"
	"os"
	"time"
)

func env(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Missing mandatory environment variable: %s", val)
	}
	return val
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	// configPath := env("CONFIG_PATH")
	// rawYaml, err := ioutil.ReadFile(configPath)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// config := make(map[string]interface{})
	// err = yaml.Unmarshal(rawYaml, config)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Printf("%#v\n", config)

	// templates := config["templates"]
	// log.Printf("%#v\n", templates)

	// bots := config["bots"]
	// log.Printf("%#v\n", bots)

	// responders := config["responders"]
	// log.Printf("%#v\n", responders)

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
