package cogs

import (
	"TheNeonCenturyBot2/config"
	"TheNeonCenturyBot2/storage"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

func AllowChannel(db *gorm.DB) func(d *discordgo.Session, m *discordgo.MessageCreate) {
	return func(d *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == d.State.User.ID {
			return
		}

		if strings.HasPrefix(m.Content, config.CommandPrefix+config.CommandAllowChannelNew) {
			if !config.IsAdmin(d, m.Author.ID, m.ChannelID) {
				d.ChannelMessageSend(m.ChannelID, "You are not allowed to use this command.")
				return
			}
			if storage.IsChannelAllowed(db, m.ChannelID) {
				d.ChannelMessageSend(m.ChannelID, "Messages in this channel are already been allowed.")
				return
			}
			d.ChannelMessageSend(m.ChannelID, "Channel has been allowed!")
			err := storage.AllowChannel(db, m.ChannelID)
			if err != nil {
				log.Fatal("AllowChannel Error: " + err.Error())
				return
			}
		}
	}

}
