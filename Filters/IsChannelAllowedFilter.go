package Filters

import (
	"TheNeonCenturyBot2/storage"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

func IsChannelAllowedFilter(db *gorm.DB, channelId string, d *discordgo.Session) bool {
	if !storage.IsChannelAllowed(db, channelId) {
		d.ChannelMessageSend(channelId, "You are not allowed to use this command here")
		return false
	}
	return true
}
