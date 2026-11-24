package config

import "github.com/bwmarrin/discordgo"

func IsAdmin(d *discordgo.Session, userId, channelId string) bool {
	perms, err := d.UserChannelPermissions(userId, channelId)
	if err != nil {
		return false
	}

	return perms&discordgo.PermissionAdministrator != 0
}
