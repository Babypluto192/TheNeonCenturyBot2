package Filters

import (
	"TheNeonCenturyBot2/storage"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

func IsUserBannedFilter(db *gorm.DB, userId string, d *discordgo.Session, channelId string) bool {
	if storage.IsUserBanned(db, userId) {
		d.ChannelMessageSend(channelId, "You are not allowed to use commands")
		return false
	}

	return true
}
