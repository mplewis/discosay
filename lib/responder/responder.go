package responder

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strings"
)

// Responder is a set of stimuli and response behavior used by a Bot.
type Responder struct {
	Name  *string        // The name of the responder
	Match *regexp.Regexp // Only respond to messages that match this regex

	// Optional. A list of responses, from which one is randomly selected
	Responses *[]string

	// Optional. Metadata that specifies the template to be used
	TemplateName *string

	// Optional. Insert the response into this template, replacing $MSG. If omitted, just send the response as-is
	Template *string

	// Optional. If provided, the probability from 0.0 to 1.0 that this
	// responder will fire on any given matched message
	Probability *float64
}

func (r *Responder) String() string {
	tmpl := "no template"
	if r.Template != nil {
		tmpl = fmt.Sprintf("template %s (%d chars)", *r.TemplateName, len(*r.Template))
	}
	prob := "always"
	if r.Probability != nil {
		prob = fmt.Sprintf("with prob. %f", *r.Probability)
	}
	return fmt.Sprintf(
		"{Responder %s: match /%s/, %d response(s), %s, fires %s}",
		*r.Name, r.Match.String(), len(*r.Responses), tmpl, prob)
}

// New builds a Responder from a YAML config blob.
func New(name string, from map[string]interface{}) (*Responder, error) {
	re, err := regexp.Compile(from["match"].(string))
	if err != nil {
		return nil, err
	}

	responses := []string{}
	if rr := from["responses"]; rr != nil {
		for _, resp := range rr.([]interface{}) {
			responses = append(responses, resp.(string))
		}
	}

	var template *string = nil
	if tp := from["template"]; tp != nil {
		s := tp.(string)
		template = &s
	}

	var probability *float64 = nil
	if pb := from["probability"]; pb != nil {
		f := pb.(float64)
		probability = &f
	}

	return &Responder{
		Name:         &name,
		Match:        re,
		Responses:    &responses,
		TemplateName: template,
		Probability:  probability,
	}, err
}

// Respond returns what a responder chose to say in reply to a message, if anything.
func (r *Responder) Respond(in string) *string {
	matched, msg := r.match(in)
	if !matched {
		log.Println("not matched")
		return nil
	}
	if msg == nil {
		log.Println("no capture")
		msg = r.response()
	}
	if msg == nil {
		log.Println("no capture, no response")
		return nil
	}
	if !r.roll() {
		fmt.Println("not responding due to probability")
		return nil
	}
	if r.Template != nil {
		log.Println("using template")
		m := strings.Replace(*r.Template, "$MSG", *msg, 1)
		msg = &m
	}
	return msg
}

func (r *Responder) roll() bool {
	if r.Probability == nil {
		return true
	}
	return rand.Float64() <= *r.Probability
}

func (r *Responder) match(msg string) (bool, *string) {
	m := r.Match.FindStringSubmatch(msg)
	if m == nil {
		return false, nil
	}
	if len(m) < 2 {
		return true, nil
	}
	return true, &m[1]
}

func (r *Responder) response() *string {
	if r.Responses == nil {
		return nil
	}
	return &(*r.Responses)[rand.Intn(len(*r.Responses))]
}
