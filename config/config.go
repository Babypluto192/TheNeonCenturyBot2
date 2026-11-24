package config

const CommandPrefix = ">"
const CommandHelp = "help"
const CommandAllowChannelNew = "allow_channel"
const CommandDisallowChannelNew = "disallow_channel"
const CommandBanUser = "ban"
const CommandUnbanUser = "unban"
const CommandAsk = "ask"
const CommandCleanHistory = "clean_history"
const CommandUpdateLoreDoc = "update_lore_doc"
const FilterPrompt = `You are strict topic filter for a Discord Bot dedicated to Hearts of Iron 4 mod "The Neon Century"
	
	Decide whether user's' message is a question related to mod's lore, characters, world, flag, map, events or question like what happens to person or country or any other variation of this questions.
	- If it is related to the mod's lore reply with exactly: OK
	- If it is NOT related to the mod's lore reply with exactly: 429BadQuestion

Do not add explanations, punctuation, or markdown.
`
