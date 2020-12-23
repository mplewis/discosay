package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/mplewis/discosay/lib/bot"
	"github.com/mplewis/discosay/lib/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

func env(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatal().Str("key", key).Msg("Missing mandatory environment variable")
	}
	return val
}

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if os.Getenv("DEBUG") != "" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	if os.Getenv("DEVELOPMENT") != "" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}
	rand.Seed(time.Now().UTC().UnixNano())

	configPath := env("CONFIG_PATH")
	rawYaml, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal().Str("config_path", configPath).Err(err).Msg("Could not read config file")
	}

	configBlob := make(map[string]interface{})
	err = yaml.Unmarshal(rawYaml, configBlob)
	if err != nil {
		log.Fatal().Str("config_path", configPath).Err(err).Msg("Could not unmarshal config file")
	}

	botSpecs, err := config.Parse(configBlob)
	if err != nil {
		log.Fatal().Str("config_path", configPath).Err(err).Msg("Could not parse config file")
	}

	log.Info().Msg("Registered bots")
	for _, spec := range botSpecs {
		for _, resp := range spec.Responders {
			log.Info().Str("bot", spec.Name).Str("responder", resp.String()).Send()
		}
	}

	bots := []*bot.Bot{}
	for _, spec := range botSpecs {
		spec.AuthToken = env(fmt.Sprintf("%s_AUTH_TOKEN", strings.ToUpper(spec.Name)))
		log.Info().Str("bot", spec.Name).Msg("Connecting...")

		b, err := bot.New(spec)
		if err != nil {
			log.Fatal().Str("bot", *b.Name).Err(err).Msg("Failed to connect")
		}

		log.Info().Str("bot", *b.Name).Msg("Connected")
		bots = append(bots, b)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	log.Info().Msg("Shutting down")
	for _, b := range bots {
		if err := b.Close(); err != nil {
			log.Err(err).Msg("Failed to close connection")
		}
	}
}
