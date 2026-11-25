package config

import "os"

func IsDevServer(guildID string) bool {
	devServerId := os.Getenv("DEV_SERVER")
	if guildID == devServerId {
		return true
	}

	return false
}
