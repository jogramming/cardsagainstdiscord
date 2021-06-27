package cardsagainstdiscord

func init() {
	pack := &CardPack{
		Name:        "2000s-nostalgia",
		Description: "2000s Nostalgia Pack - Butterfly clips, Bootcut jeans, Bush did 9/11, etc",
		Prompts: []*PromptCard{
			&PromptCard{Prompt: `%s? That's a no from me, dawg.`},
			&PromptCard{Prompt: `16 people. 39 days of %s. One Survivor.`},
			&PromptCard{Prompt: `I couldn't help but wonder. Was "having it all" an unattainable myth? Was the secret to a truly happy life just %s?`},
			&PromptCard{Prompt: `Oh my god! %s killed Kenny!`},
			&PromptCard{Prompt: `Oops! I did it again. I played with %s.`},
		},
		Responses: []ResponseCard{
			`All these boys in my yard asking for milkshakes.`,
			`Being a total Miranda.`,
			`Blowing Lucas on an inflatable chair.`,
			`Chunky highlights.`,
			`Collateralized debt obligations.`,
			`Competing against 20 other women for the love of Flava Flav.`,
			`Cybering with diselman89.`,
			`Defining marriage at the union of one man and one woman.`,
			`Doing 9/11.`,
			`Downloading a bunch of weird porn on Kazaa.`,
			`Getting crunk.`,
			`Getting drenched in Dave Matthew's shit and piss.`,
			`Getting pwned.`,
			`Going to prom with a 108-year-old vampire.`,
			`Hoobastank.`,
			`Killing Dumbledore.`,
			`Letting the terrorists win.`,
			`My son and business partnew H.W. Plainview.`,
			`Starting to be cool about gay people.`,
			`Suspicious white powder.`,
			`Taking a blurry photo of my penis on my Motorola Razr.`,
			`The conservative blogosphere.`,
			`The new Sweet Onion Chicken Teriyaki sandwich from Subway!`,
			`Trapping my Sim in a swimming pool until she shits herself.`,
			`Two girls sharing one cup of poop.`,
			`Vajazzling my vajayjay.`,
			`Weapons of mass destruction.`,
		},
	}

	AddPack(pack)
}
