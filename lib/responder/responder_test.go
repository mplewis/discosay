package responder_test

import (
	"math/rand"
	"testing"

	"github.com/mplewis/discosay/lib/responder"
	"gopkg.in/yaml.v3"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestResponder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Responder Suite")
}

var responses = `
match: '^!pick$'
responses:
  - '1'
  - '2'
  - '3'
`

var template = `
match: '^!say (.+)$'
template: whisper
`

var templateResponses = `
match: '^!say$'
template: whisper
responses:
  - Can I get a cigarette?
  - Do you know the answer for number six?
`

var probability = `
match: spaghetti
probability: 0.05
responses:
  - SPAGETT!!
`

func load(raw string) map[string]interface{} {
	blob := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(raw), blob)
	if err != nil {
		panic(err)
	}
	return blob
}

var _ = Describe("responder", func() {
	BeforeEach(func() {
		rand.Seed(1)
	})

	Context("responses", func() {
		resp, _ := responder.New("responses", load(responses))
		It("picks a random message", func() {
			Expect(*resp.Name).To(Equal("responses"))
			Expect(resp.Respond("hello")).To(BeNil())
			Expect(*resp.Respond("!pick")).To(Equal("3"))
			Expect(*resp.Respond("!pick")).To(Equal("1"))
			Expect(*resp.Respond("!pick")).To(Equal("3"))
			Expect(*resp.Respond("!pick")).To(Equal("3"))
			Expect(*resp.Respond("!pick")).To(Equal("2"))
		})
	})

	Context("template", func() {
		resp, _ := responder.New("template", load(template))
		It("fills in the template with the capture", func() {
			Expect(*resp.TemplateName).To(Equal("whisper"))
			tmpl := `You hear a faint whisper: "$MSG"`
			resp.Template = &tmpl
			Expect(*resp.Respond("!say Psst. You wanna buy a boat?")).
				To(Equal(`You hear a faint whisper: "Psst. You wanna buy a boat?"`))
		})
	})

	Context("with responses", func() {
		resp, _ := responder.New("templateResponses", load(templateResponses))
		tmpl := `You hear a faint whisper: "$MSG"`
		resp.Template = &tmpl
		It("fills in the template with a random response", func() {
			Expect(*resp.Respond("!say")).To(Equal(`You hear a faint whisper: "Do you know the answer for number six?"`))
			Expect(*resp.Respond("!say")).To(Equal(`You hear a faint whisper: "Do you know the answer for number six?"`))
			Expect(*resp.Respond("!say")).To(Equal(`You hear a faint whisper: "Do you know the answer for number six?"`))
			Expect(*resp.Respond("!say")).To(Equal(`You hear a faint whisper: "Do you know the answer for number six?"`))
			Expect(*resp.Respond("!say")).To(Equal(`You hear a faint whisper: "Do you know the answer for number six?"`))
			Expect(*resp.Respond("!say")).To(Equal(`You hear a faint whisper: "Can I get a cigarette?"`))
		})
	})

	Context("probability", func() {
		resp, _ := responder.New("probability", load(probability))
		It("fires with the given probability", func() {
			Expect(*resp.Probability).To(Equal(0.05))
			fired := 0
			for i := 0; i < 1000; i++ {
				if resp.Respond("spaghetti") != nil {
					fired++
				}
			}
			Expect(fired).To(Equal(47))
		})
	})
})
