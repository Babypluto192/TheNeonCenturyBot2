package cogs

import (
	"TheNeonCenturyBot2/Filters"
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
			if !Filters.IsAdminFilter(m.ChannelID, d, m.Author.ID) {
				return
			}
			if !Filters.IsChannelAllowedFilter(db, m.ChannelID, d) {
				return
			}
			err := storage.DisAllowChannel(db, m.ChannelID)
			if err != nil {
				return
			}
			d.ChannelMessageSend(m.ChannelID, "This channel is now disallowed")
		}
	}
}
