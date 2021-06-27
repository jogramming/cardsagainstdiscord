package cardsagainstdiscord

func init() {
	pack := &CardPack{
		Name:        "fantasy",
		Description: "Fantasy Pack - Dragons, wizards, Orlando Bloom, etc",
		Prompts: []*PromptCard{
			&PromptCard{Prompt: `And in the end, the dragon was not evil; he just wanted %s.`},
			&PromptCard{Prompt: `Critics are raving about HBO's new Game of Thrones spin-off, "%s of %s."`},
			&PromptCard{Prompt: `Having tired of poetry and music, the immortal elves now fill their days with %s.`},
			&PromptCard{Prompt: `Legend tells of a princess who has been asleep for a thousand years and can only be awoken by %s.`},
			&PromptCard{Prompt: `Who blasphemes and bubbles at the center of all infinity, whose name no lips dare speak aloud, and who gnaws hungrily in inconceivable, unlighted chambers beyond time?`},
			&PromptCard{Prompt: `Your father was a powerful wizard, Harry. Before he died, he left you something very precious: %s.`},
		},
		Responses: []ResponseCard{
			`A CGI dragon.`,
			`A dwarf who won't leave you alone until you compare penis sizes.`,
			`A gay sorcerer who turns everyone gay.`,
			`A ghoul.`,
			`A Hitachi Magic Wand.`,
			`A magical kingdom with dragons and elves and no black people.`,
			`A mysterious, floating orb.`,
			`A weed elemental who gets everyone high.`,
			`Accidentally conjuring a legless horse that can't stop ejaculating.`,
			`Astral-projecting into the Goat Realm.`,
			`Bathing naked in a moonlit grove.`,
			`Dinosaurs who wear armor and you ride them and they kick ass.`,
			`Eternal darkness.`,
			`Foraging for berries in yonder wood.`,
			`Freaky, pan-dimensional sex with a demigod.`,
			`Gender equality.`,
			`Going on an epic adventure and learning a valuable lesson about friendship.`,
			`Handcuffing a wizard to a radiator and dousing him with kerosene.`,
			`Hodor.`,
			`How hot Orlando Bloom was in Lord of the Rings.`,
			`Kneeing a wizard in the balls.`,
			`Make-believe stories for autistic white men.`,
			`Reading The Hobbit under the covers while mom and dad scream at each other downstairs.`,
			`Shitting in a wizard's spell book and jizzing in his hat.`,
			`Shooting a wizard with a gun.`,
			`The all-seeing Eye of Sauron.`,
			`The card Neil Gaiman wrote: "Three elves at a time."`,
			`True love's kiss.`,
		},
	}

	AddPack(pack)
}
