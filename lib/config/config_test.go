package config_test

import (
	"testing"

	"github.com/mplewis/discosay/lib/bot"
	"github.com/mplewis/discosay/lib/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v3"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}

var raw = `
templates:
  bracket: "<<$MSG>>"
  warning: "!!! $MSG !!!"

bots:
  bracketeer:
    - bracksay
  warner:
    - dangerwill

responders:
  bracksay:
    match: "^!bracksay (.+)$"
    template: bracket
  dangerwill:
    match: \bdanger\b
    template: warning
    probability: 0.1
    responses:
      - It sounds like old morse code.
      - Danger, Will Robinson!
`

func load(raw string) map[string]interface{} {
	blob := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(raw), blob)
	if err != nil {
		panic(err)
	}
	return blob
}

var _ = Describe("config", func() {
	Describe("Parse", func() {
		specs, err := config.Parse(load(raw))
		It("parses the expected bot specs", func() {
			Expect(err).To(BeNil())
			names := []string{}
			botsByName := make(map[string]bot.Spec)
			for _, spec := range specs {
				botsByName[spec.Name] = spec
				names = append(names, spec.Name)
			}
			Expect(names).To(ContainElements("bracketeer", "warner"))

			br := botsByName["bracketeer"]
			Expect(br.AuthToken).To(BeNil())
			Expect(br.Responders).To(HaveLen(1))

			brr := br.Responders[0]
			Expect(*brr.Name).To(Equal("bracksay"))
			Expect(brr.Match.String()).To(Equal("^!bracksay (.+)$"))
			Expect(*brr.Template).To(Equal("<<$MSG>>"))

			wa := botsByName["warner"]
			Expect(wa.AuthToken).To(BeNil())
			Expect(wa.Responders).To(HaveLen(1))

			dw := wa.Responders[0]
			Expect(*dw.Name).To(Equal("dangerwill"))
			Expect(dw.Match.String()).To(Equal(`\bdanger\b`))
			Expect(*dw.Template).To(Equal("!!! $MSG !!!"))
			Expect(*dw.Probability).To(Equal(0.1))
			Expect(*dw.Responses).To(Equal([]string{
				"It sounds like old morse code.",
				"Danger, Will Robinson!",
			}))
		})
	})
})
