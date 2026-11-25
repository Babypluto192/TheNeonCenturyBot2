package Filters

import (
	"TheNeonCenturyBot2/config"
	"TheNeonCenturyBot2/storage"

	"github.com/bwmarrin/discordgo"
	"gorm.io/gorm"
)

func AutoBanFilter(db *gorm.DB, channelId string, d *discordgo.Session, userId string) bool {
	if storage.IsUserExist(db, userId) {
		requests, err := storage.GetUserRequests(db, userId)
		if err != nil {
			d.ChannelMessageSend(channelId, "There was an error getting the user requests")
			return false
		}
		isJailBreak := config.IsTryingToJailBreak(db, requests, userId)
		if isJailBreak && !config.IsAdmin(d, userId, channelId) {
			d.ChannelMessageSend(channelId, "You appear to be attempting to jailbreak the bot, you have been banned [bye byee](https://klipy.com/gifs/mai-sakurajima-rascal-does-not-dream-of-bunny-girl-senpai-3) ")
			return false
		}
		return true
	}
	return false
}
