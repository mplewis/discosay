package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/mplewis/discosay/lib/bot"
	"github.com/mplewis/discosay/lib/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// mustEnv fetches the value of a mandatory environment variable.
func mustEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatal().Str("key", key).Msg("Missing mandatory environment variable")
	}
	return val
}

// maybeEnv returns a value if an environment variable is found, and nil if not.
func maybeEnv(key string) *string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return nil
	}
	return &val
}

// checkEnv returns true if an environment variable is set.
func checkEnv(key string) bool {
	_, ok := os.LookupEnv(key)
	return ok
}

// tty returns true if the current terminal is interactive.
func tty() bool {
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// openBots instantiates and connects all bots to Discord.
func openBots(botSpecs []bot.Spec) ([]*bot.Bot, error) {
	bots := []*bot.Bot{}
	for _, spec := range botSpecs {
		spec.SetAuthToken(mustEnv(fmt.Sprintf("%s_AUTH_TOKEN", strings.ToUpper(spec.Name))))
		log.Info().Str("bot", spec.Name).Msg("Connecting...")

		b, err := bot.New(spec)
		log := log.With().Str("name", *b.Name).Logger()
		if err != nil {
			log.Error().Err(err).Msg("Could not connect bot")
			return nil, err
		}

		log.Info().Msg("Connected")
		bots = append(bots, b)
	}
	return bots, nil
}

// closeBots disconnects all bots from Discord.
func closeBots(bots []*bot.Bot) error {
	var e error
	for _, b := range bots {
		if err := b.Close(); err != nil {
			e = err
			log.Err(err).Msg("Failed to close connection")
		}
	}
	return e
}

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerolog.DurationFieldUnit = time.Second
	if checkEnv("DEBUG") {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	if tty() {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	}

	rand.Seed(time.Now().UTC().UnixNano())

	csrc := config.Source{Path: maybeEnv("CONFIG_PATH"), URL: maybeEnv("CONFIG_URL")}
	log.Info().Interface("config_source", csrc).Msg("Loading Discosay config")
	botSpecs, hash, err := config.Load(csrc)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not load config")
	}

	log.Info().Msg("Registered bots")
	for _, spec := range botSpecs {
		for _, resp := range spec.Responders {
			log.Info().Str("bot", spec.Name).Str("responder", resp.String()).Send()
		}
	}

	bots, err := openBots(botSpecs)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not connect bot")
	}

	intRaw := maybeEnv("CONFIG_UPDATE_INTERVAL")
	if intRaw != nil {
		log := log.With().Interface("config_source", csrc).Logger()

		intSecs, err := strconv.Atoi(*intRaw)
		if err != nil {
			log.Fatal().Err(err).Msg("Could not parse CONFIG_UPDATE_INTERVAL")
		}

		interval := time.Duration(intSecs) * time.Second
		t := time.Tick(interval)
		for {
			<-t
			botSpecs, newHash, err := config.Load(csrc)
			if err != nil {
				log.Error().Err(err).Interface("source", csrc).Msg("Failed to load config, using last good config")
			}
			if hash == newHash {
				log.Debug().Msg("Config unchanged")
				continue
			}

			log.Info().Msg("Config changed, reloading bots")
			hash = newHash
			err = closeBots(bots)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to shutdown bots after config update")
			}
			newBots, err := openBots(botSpecs)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to start bots after config update")
			}
			bots = newBots
		}
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Info().Msg("Shutting down")
	closeBots(bots)
}
