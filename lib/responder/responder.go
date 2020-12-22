package responder

import (
	"log"
	"math/rand"
	"regexp"
	"strings"
)

type Responder struct {
	Name  string        // The name of the responder
	Match regexp.Regexp // Only respond to messages that match this regex

	// Optional. A list of responses, from which one is randomly selected
	Responses []string

	// Optional. Insert the response into this template, replacing $MSG. If omitted, just send the response as-is
	Template string

	// Optional. If provided, the probability from 0.0 to 1.0 that this
	// responder will fire on any given matched message
	Probability float32
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
	n := rand.Intn(len(r.Responses))
	return r.Responses[n]
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
	if r.Template != "" {
		log.Println("using template")
		msg = strings.Replace(r.Template, "$MSG", msg, 1)
	}
	return &msg
}
