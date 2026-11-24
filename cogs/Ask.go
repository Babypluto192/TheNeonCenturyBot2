package cogs

import (
	"TheNeonCenturyBot2/config"
	"TheNeonCenturyBot2/storage"
	"log"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
	"google.golang.org/genai"
	"gorm.io/gorm"
)

func Ask(db *gorm.DB, ai *genai.Client) func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(d *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == d.State.User.ID {
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
		if !storage.IsChannelAllowed(db, m.ChannelID) {
			d.ChannelMessageSend(m.ChannelID, "You are not allowed to use this command here")
			return
		}
		if storage.IsUserBanned(db, m.Author.ID) {
			d.ChannelMessageSend(m.ChannelID, "You are not allowed to use commands")
			return
		}

		if !storage.IsUserExist(db, m.Author.ID) {
			isAdmin := config.IsAdmin(d, m.Author.ID, m.ChannelID)
			storage.AddUser(db, m.Author.ID, false, isAdmin)
		}

		if storage.IsUserExist(db, m.Author.ID) {
			requests, err := storage.GetUserRequests(db, m.Author.ID)
			if err != nil {
				d.ChannelMessageSend(m.ChannelID, "There was an error getting the user requests")
				return
			}
			isJailBreak := config.IsTryingToJailBreak(db, requests, m.Author.ID)
			if isJailBreak && !config.IsAdmin(d, m.Author.ID, m.ChannelID) {
				d.ChannelMessageSend(m.ChannelID, "You appear to be attempting to jailbreak the bot, you have been banned [bye byee](https://klipy.com/gifs/mai-sakurajima-rascal-does-not-dream-of-bunny-girl-senpai-3) ")
			}
		}

		userFilteredQuestion := config.GeminiFitler(ai, config.FilterPrompt, question)

		if userFilteredQuestion == "" {
			log.Println("Something went wrong with filter")
			d.ChannelMessageSend(m.ChannelID, "Something went wrong with filter, try again later")
			return
		}

		if strings.Contains(userFilteredQuestion, "OK") {
			err := storage.AddRequestToUser(db, m.Author.ID, userFilteredQuestion)
			d.ChannelMessageSend(m.ChannelID, "Question Passed Filter")
			if err != nil {
				log.Println("AddRequestToUser error:", err)
				return
			}
			return
		}
		if strings.Contains(userFilteredQuestion, "429") {
			err := storage.AddRequestToUser(db, m.Author.ID, userFilteredQuestion)
			d.ChannelMessageSend(m.ChannelID, "Your question is not related to lore")
			if err != nil {
				log.Println("AddRequestToUser error:", err)
				return
			}
			return
		}

		return

	}
}
