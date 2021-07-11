package responder

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
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

	// Optional. If set, delete the parent message which triggered the responder
	DeleteParent bool

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
func New(name string, config map[string]interface{}) (*Responder, error) {
	reTmpl := config["match"].(string)
	if ci := config["case_sensitive"]; ci != nil {
		if !ci.(bool) {
			reTmpl = fmt.Sprintf("(?i)%s", reTmpl)
		}
	}
	re, err := regexp.Compile(reTmpl)
	if err != nil {
		return nil, err
	}

	responses := []string{}
	if rr := config["responses"]; rr != nil {
		for _, resp := range rr.([]interface{}) {
			responses = append(responses, resp.(string))
		}
	}

	var template *string = nil
	if tp := config["template"]; tp != nil {
		s := tp.(string)
		template = &s
	}

	deleteParent := false
	if dp := config["delete_parent"]; dp != nil {
		deleteParent = dp.(bool)
	}

	var probability *float64 = nil
	if pb := config["probability"]; pb != nil {
		f := pb.(float64)
		probability = &f
	}

	return &Responder{
		Name:         &name,
		Match:        re,
		Responses:    &responses,
		TemplateName: template,
		DeleteParent: deleteParent,
		Probability:  probability,
	}, err
}

// Respond returns what a responder chose to say in reply to a message, if anything.
func (r *Responder) Respond(in string) *string {
	l := log.Debug().Str("responder", *r.Name)
	matched, msg := r.match(in)
	if !matched {
		l.Msg("not matched")
		return nil
	}
	if msg == nil {
		msg = r.response()
	}
	if msg == nil {
		l.Msg("no capture, no response")
		return nil
	}
	if !r.roll() {
		l.Msg("not responding due to probability")
		return nil
	}
	if r.Template != nil {
		l.Msg("using template")
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
