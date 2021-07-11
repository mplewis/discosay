package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/mplewis/discosay/lib/responder"
	"github.com/rs/zerolog/log"
)

// Bot is a set of Responders connected to Discord.
type Bot struct {
	Name *string
	sess *discordgo.Session
}

// Spec is the data needed to build and connect a bot.
type Spec struct {
	Name       string
	AuthToken  *string
	Responders []*responder.Responder
}

// SetAuthToken sets the auth token used to connect a bot to Discord.
// You must set this before creating a bot with bot.New.
// You can get this from the Discord Developer Portal (https://discord.com/developers/applications)
// under [your app] -> Bot -> Token.
func (spec *Spec) SetAuthToken(authToken string) {
	spec.AuthToken = &authToken
}

// New builds a bot and connects it to Discord.
func New(spec Spec) (*Bot, error) {
	if spec.AuthToken == nil {
		return nil, fmt.Errorf("must set auth token in bot.Spec before attempting to connect for %s", spec.Name)
	}
	sess, err := discordgo.New("Bot " + *spec.AuthToken)
	if err != nil {
		return nil, err
	}
	handler := buildMessageHandler(&spec.Name, spec.Responders)
	sess.AddHandler(handler)

	err = sess.Open()
	if err != nil {
		return nil, err
	}

	b := Bot{&spec.Name, sess}
	return &b, nil
}

// Close disconnects the bot from Discord cleanly.
func (b *Bot) Close() error {
	return b.sess.Close()
}

func buildMessageHandler(botName *string, responders []*responder.Responder) func(*discordgo.Session, *discordgo.MessageCreate) {
	l := log.Info().Str("bot", *botName)
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}
		for _, r := range responders {
			if out := r.Respond(m.Message.Content); out != nil {
				l.Str("inUser", m.Author.Username).
					Str("inMsg", m.Message.Content).
					Str("responder", *r.Name).
					Str("msg", *out).
					Bool("deleteParent", r.DeleteParent).
					Send()

				if r.DeleteParent {
					fmt.Println(m.ChannelID, m.Message.ID)
					s.ChannelMessageDelete(m.ChannelID, m.Message.ID)
				}
				s.ChannelMessageSend(m.ChannelID, *out)

				return
			}
		}
	}
}
