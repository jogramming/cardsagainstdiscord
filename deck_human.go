package cardsagainstdiscord

func init() {
	pack := &CardPack{
		Name:        "human",
		Description: "Human Pack - 30 beautiful cards about the human condition straight from the hearts of our human writers",
		Prompts: []*PromptCard{
			&PromptCard{Prompt: `Do not go gentle into that good night. Rage, rage against %s?`},
			&PromptCard{Prompt: `Sure, sex is great, but have you tried %s?`},
		},
		Responses: []ResponseCard{
			`A pointless job at a soulless corporation.`,
			`Being terrified of a single bee.`,
			`Being young and in love in New York City.`,
			`Chemotherapy.`,
			`Delighting in the pain of others.`,
			`Forming a monogamous pairbond.`,
			`Getting my butt trapped hilariously in a tuba.`,
			`Getting sad for no reason.`,
			`Good ol’ fashioned face-to-face conversation.`,
			`Harnessing the power of steam!!!`,
			`How good this cantaloupe feels on my penis.`,
			`Just hangin’ out, ya know?`,
			`Like 50 mosquito bites.`,
			`Losing a loved one to Fox News.`,
			`Really tall people.`,
			`Saying “what’s that?” And then stealing a french fry.`,
			`Sucking all the oil out of the planet and fucking off to Mars.`,
			`The bomb strapped to my chest.`,
			`The joy of song.`,
			`The salmon.`,
			`The unspeakable horrors of factory farming.`,
			`The whole Jeffrey Epstein thing.`,
			`This horrible thing called twitter.`,
			`Three hours of pimple-popping videos.`,
			`Unwelcome sexual attention from grown men.`,
			`Using words to communicate.`,
			`War with China.`,
			`Watching the life drain from the eyes of my son’s killer.`,
		},
	}

	AddPack(pack)
}
