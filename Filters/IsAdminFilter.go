package Filters

import (
	"TheNeonCenturyBot2/config"

	"github.com/bwmarrin/discordgo"
)

func IsAdminFilter(channelId string, d *discordgo.Session, userId string) bool {
	if !config.IsAdmin(d, userId, channelId) {
		d.ChannelMessageSend(channelId, "You are not allowed to use this command.")
		return false
	}

	return true
}
