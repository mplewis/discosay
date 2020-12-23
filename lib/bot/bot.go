package bot

import (
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
	AuthToken  string
	Responders []*responder.Responder
}

// New builds a bot and connects it to Discord.
func New(spec Spec) (*Bot, error) {
	sess, err := discordgo.New("Bot " + spec.AuthToken)
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
					Send()
				s.ChannelMessageSend(m.ChannelID, *out)
				return
			}
		}
	}
}
