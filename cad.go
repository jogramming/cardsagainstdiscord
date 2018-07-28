package cardsagainstdiscord

import (
	"github.com/jonas747/discordgo"
	"strings"
)

var Packs = make(map[string]*CardPack)

func AddPack(name string, pack *CardPack) {
	// Count picks
	for _, v := range pack.Prompts {
		numPicks := strings.Count(v.Prompt, "%s")
		if numPicks == 0 {
			v.Prompt += " %s"
			v.NumPick = 1
		} else {
			v.NumPick = numPicks
		}
	}

	Packs[name] = pack
}

type CardPack struct {
	Name      string
	Prompts   []*PromptCard
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
