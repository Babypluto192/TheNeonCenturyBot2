package cogs

import (
	"TheNeonCenturyBot2/config"
	"TheNeonCenturyBot2/storage"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

func DisAllowChannel(db *gorm.DB) func(d *discordgo.Session, m *discordgo.MessageCreate) {
	return func(d *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == d.State.User.ID {
			return
		}

		if strings.HasPrefix(m.Content, config.CommandPrefix+config.CommandDisallowChannelNew) {
			if !config.IsAdmin(d, m.Author.ID, m.ChannelID) {
				d.ChannelMessageSend(m.ChannelID, "You are not allowed to use this command.")
				return
			}
			if !storage.IsChannelAllowed(db, m.ChannelID) {
				d.ChannelMessageSend(m.ChannelID, "This channel haven't been allowed")
			}
			err := storage.DisAllowChannel(db, m.ChannelID)
			if err != nil {
				return
			}
			d.ChannelMessageSend(m.ChannelID, "This channel is now disallowed")
		}
	}
}
