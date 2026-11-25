package cogs

import (
	"TheNeonCenturyBot2/Filters"
	"TheNeonCenturyBot2/config"
	"TheNeonCenturyBot2/storage"
	"log"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
	"github.com/revrost/go-openrouter"
	"google.golang.org/genai"
	"gorm.io/gorm"
)

func Ask(db *gorm.DB, ai *genai.Client, clients []*openrouter.Client) func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(d *discordgo.Session, m *discordgo.MessageCreate) {
		if !Filters.AntiBotFilter(d, m.Author.ID) {
			return
		}

		cmd := config.CommandPrefix + config.CommandAsk
		if !strings.HasPrefix(m.Content, cmd) {
			return
		}
		rest, ok := strings.CutPrefix(strings.TrimSpace(m.Content), cmd)
		if !ok {
			d.ChannelMessageSend(m.ChannelID, "Wrong use of command >ask question")
			return
		}

		if rest != "" {
			r, _ := utf8.DecodeRuneInString(rest)
			if !unicode.IsSpace(r) {
				return
			}
		}
		question := strings.TrimSpace(rest)
		if question == "" {
			d.ChannelMessageSend(m.ChannelID, "Wrong number of arguments, use "+cmd+" question")
			return
		}
		if !Filters.IsChannelAllowedFilter(db, m.ChannelID, d) {
			return
		}

		if !Filters.IsUserBannedFilter(db, m.Author.ID, d, m.ChannelID) {
			return
		}

		if !Filters.CreateUserFilter(db, m.ChannelID, d, m.Author.ID) {
			return
		}

		if !Filters.AutoBanFilter(db, m.ChannelID, d, m.Author.ID) {
			return
		}

		if config.IsDevServer(m.GuildID) {
			log.Println(m.Author.Username, "asks", question)
			stopTyping := startTyping(d, m.ChannelID)
			response := config.GenerateResponse(m.GuildID, question, clients)
			stopTyping()

			const limit = 1900
			responseMessage := []rune(response)
			for i := 0; i < len(responseMessage); i += limit {
				end := min(i+limit, len(responseMessage))
				if _, err := d.ChannelMessageSend(m.ChannelID, string(responseMessage[i:end])); err != nil {
					log.Println("ChannelMessageSend error:", err)
					return
				}
			}

			return
		}
		userFilteredQuestion := config.GeminiFitler(ai, config.FilterPrompt, question)

		if !Filters.GeminiFilterError(m.ChannelID, d, userFilteredQuestion) {
			return
		}

		if !Filters.FilterJailBreakAttempt(db, m.ChannelID, d, userFilteredQuestion, m.Author.ID) {
			return
		}

		if strings.Contains(userFilteredQuestion, "OK") {
			log.Println(m.Author.Username, "asks", question)
			err := storage.AddRequestToUser(db, m.Author.ID, userFilteredQuestion)
			stopTyping := startTyping(d, m.ChannelID)
			response := config.GenerateResponse(m.GuildID, question, clients)
			stopTyping()

			const limit = 1900
			responseMessage := []rune(response)
			for i := 0; i < len(responseMessage); i += limit {
				end := min(i+limit, len(responseMessage))
				if _, err := d.ChannelMessageSend(m.ChannelID, string(responseMessage[i:end])); err != nil {
					log.Println("ChannelMessageSend error:", err)
					return
				}
			}

			// d.ChannelMessageSend(m.ChannelID, "Question Passed Filter")an
			if err != nil {
				log.Println("AddRequestToUser error:", err)
				return
			}
			return
		}

		return

	}
}

func startTyping(d *discordgo.Session, channelID string) func() {
	done := make(chan struct{})

	go func() {
		_ = d.ChannelTyping(channelID)
		ticker := time.NewTicker(8 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				_ = d.ChannelTyping(channelID)
			}
		}
	}()

	return func() {
		close(done)
	}
}
