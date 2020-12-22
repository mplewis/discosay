package responder

import (
	"log"
	"math/rand"
	"regexp"
	"strings"
)

type Responder struct {
	Name  string         // The name of the responder
	Match *regexp.Regexp // Only respond to messages that match this regex

	// Optional. A list of responses, from which one is randomly selected
	Responses []string

	// Optional. Insert the response into this template, replacing $MSG. If omitted, just send the response as-is
	Template string

	// Optional. If provided, the probability from 0.0 to 1.0 that this
	// responder will fire on any given matched message
	Probability float32
}

func Parse(name string, from map[string]interface{}) (*Responder, error) {
	re, err := regexp.Compile(from["match"].(string))
	if err != nil {
		return nil, err
	}

	log.Println(from["responses"])
	responses := []string{}
	rawResps, found := from["responses"]
	if found {
		for _, resp := range rawResps.([]interface{}) {
			responses = append(responses, resp.(string))
		}
	}
	log.Println(responses)

	probability, found := from["probability"].(float32)
	if !found {
		probability = 0.0
	}

	template, found := from["template"].(string)
	if !found {
		template = ""
	}

	return &Responder{
		Name:        name,
		Match:       re,
		Responses:   responses,
		Template:    template,
		Probability: probability,
	}, err
}

func (r *Responder) roll() bool {
	if r.Probability > 0.0 && rand.Float32() > r.Probability {
		return false
	}
	return true
}

func (r *Responder) match(msg string) (bool, string) {
	m := r.Match.FindStringSubmatch(msg)
	if m == nil {
		return false, ""
	}
	if len(m) < 2 {
		return true, ""
	}
	return true, m[1]
}

func (r *Responder) response() string {
	rc := len(r.Responses)
	if rc == 0 {
		return ""
	}
	return r.Responses[rand.Intn(rc)]
}

func (r *Responder) Respond(in string) *string {
	matched, substr := r.match(in)
	if !matched {
		log.Println("not matched")
		return nil
	}
	msg := substr
	if msg == "" {
		log.Println("no captured msg")
		msg = r.response()
	}
	if msg == "" {
		log.Println("no response msg")
		return nil
	}
	if r.Template != "" {
		log.Println("using template")
		msg = strings.Replace(r.Template, "$MSG", msg, 1)
	}
	return &msg
}
