package cardsagainstdiscord

import (
	"github.com/jonas747/discordgo"
)

var Packs = make(map[string]*CardPack)

type CardPack struct {
	Name      string
	Prompts   []PromptCard
	Responses []ResponseCard
}

type PromptCard struct {
	Prompt  string
	NumPick int
}

type ResponseCard string

type SessionProvider interface {
	SessionForGuild(guildID int64) *discordgo.Session
}

type StaticSessionProvider struct {
	Session *discordgo.Session
}

func (sp *StaticSessionProvider) SessionForGuild(guildID int64) *discordgo.Session {
	return sp.Session
}
