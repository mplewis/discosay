package config

import (
	"crypto/sha256"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/mplewis/discosay/lib/bot"
	"github.com/mplewis/discosay/lib/responder"
	"gopkg.in/yaml.v3"
)

// Source indicates the source location for a Discosay config.
type Source struct {
	Path *string
	URL  *string
}

// Load loads and parses a config blob into Bot Specs, including a hash of the loaded config.
func Load(s Source) ([]bot.Spec, [32]byte, error) {
	if s.Path == nil && s.URL == nil {
		return nil, [32]byte{}, errors.New("config path and URL both unset")
	}

	var rawYaml []byte
	if s.Path != nil {
		y, err := ioutil.ReadFile(*s.Path)
		if err != nil {
			return nil, [32]byte{}, err
		}
		rawYaml = y
	} else {
		resp, err := http.Get(*s.URL)
		if err != nil {
			return nil, [32]byte{}, err
		}
		rawYaml, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, [32]byte{}, err
		}
	}

	configBlob := make(map[string]interface{})
	err := yaml.Unmarshal(rawYaml, configBlob)
	if err != nil {
		return nil, [32]byte{}, err
	}

	hash := sha256.Sum256(rawYaml)
	specs, err := Parse(configBlob)
	return specs, hash, err
}

// Parse parses a config blob into Bot Specs.
func Parse(config map[string]interface{}) ([]bot.Spec, error) {
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

	botSpecs := []bot.Spec{}
	for name, respNames := range config["bots"].(map[string]interface{}) {

		s := bot.Spec{Name: name}
		for _, rawName := range respNames.([]interface{}) {
			responder := responders[rawName.(string)]
			s.Responders = append(s.Responders, &responder)
		}
		botSpecs = append(botSpecs, s)
	}
	return botSpecs, nil
}
