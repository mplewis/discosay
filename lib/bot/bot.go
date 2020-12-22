package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/mplewis/discosay/lib/responder"
)

type Bot struct {
	name *string
	sess *discordgo.Session
}

func New(name *string, authToken string, responders []responder.Responder) (*Bot, error) {
	sess, err := discordgo.New("Bot " + authToken)
	if err != nil {
		return nil, err
	}
	handler := buildMessageHandler(name, responders)
	sess.AddHandler(handler)

	b := Bot{name, sess}
	return &b, err
}

func (b *Bot) Close() {
	b.sess.Close()
}

func buildMessageHandler(botName *string, responders []responder.Responder) func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		for _, r := range responders {
			log.Printf("<- %s: %s", m.Author.Username, m.Message.Content)
			if out := r.Respond(m.Message.Content); out != nil {
				log.Printf("-> %s(%s): %s", *botName, r.Name, *out)
				s.ChannelMessageSend(m.ChannelID, *out)
				return
			}
		}
	}
}
