package config

import (
	"github.com/mplewis/discosay/lib/bot"
	"github.com/mplewis/discosay/lib/responder"
)

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
