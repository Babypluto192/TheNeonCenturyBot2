package cogs

import (
	"TheNeonCenturyBot2/Filters"
	"TheNeonCenturyBot2/config"
	"TheNeonCenturyBot2/storage"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

func AllowChannel(db *gorm.DB) func(d *discordgo.Session, m *discordgo.MessageCreate) {
	return func(d *discordgo.Session, m *discordgo.MessageCreate) {
		if !Filters.AntiBotFilter(d, m.Author.ID) {
			return
		}

		if strings.HasPrefix(m.Content, config.CommandPrefix+config.CommandAllowChannelNew) {
			if !Filters.IsAdminFilter(m.ChannelID, d, m.Author.ID) {
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
