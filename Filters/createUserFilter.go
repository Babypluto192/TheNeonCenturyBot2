package Filters

import (
	"TheNeonCenturyBot2/config"
	"TheNeonCenturyBot2/storage"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

func CreateUserFilter(db *gorm.DB, channelId string, d *discordgo.Session, userId string) bool {
	if !storage.IsUserExist(db, userId) {
		isAdmin := config.IsAdmin(d, userId, channelId)
		err := storage.AddUser(db, userId, false, isAdmin)
		if err != nil {
			d.ChannelMessageSend(channelId, "Error while saving user to db, please contact devs")
			return false
		}
	}
	return true
}
