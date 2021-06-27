package cardsagainstdiscord

func init() {
	pack := &CardPack{
		Name:        "geek",
		Description: "Geek Pack - 30 cards about video games, D&D, Game of Thrones, and all the other bullshit you like",
		Prompts: []*PromptCard{
			&PromptCard{Prompt: `%s is OP. Please nerf.`},
			&PromptCard{Prompt: `%s is way better in %s mode. (Pax East 2014)`},
			&PromptCard{Prompt: `%s: Achievement unlocked. (Pax Prime 2013)`},
			&PromptCard{Prompt: `(Heavy breathing) Luke, I am %s. (Pax East 2014)`},
			&PromptCard{Prompt: `Press ↓ ↓ ← → B to unleash %s. (Pax East 2013 Promo Pack C)`},
			&PromptCard{Prompt: `What made Spock cry? (Pax Prime 2013)`},
			&PromptCard{Prompt: `What's the latest bullshit that's troubling this quaint fantasy town? (Pax Prime 2013)`},
		},
		Responses: []ResponseCard{
			`A fully-dressed female videogame character. (Pax Prime 2013)`,
			`A grumpy old Harrison Ford who'd rather be doing anything else. (Pax East 2014)`,
			`A homemade, cum-stained Star Trek Uniform. (Pax Prime 2013)`,
			`Achieving 500 actions per minute. (Pax East 2013 Promo Pack C)`,
			`Charging up all the way. (Pax East 2013 Promo Pack C)`,
			`Eating a pizza that's lying in the street to gain health. (Pax Prime 2013)`,
			`Endless ninjas. (Pax East 2014)`,
			`Forgetting to eat, and consequently dying. (Pax East 2013 Promo Pack C)`,
			`Getting bitch slapped by Dhalsim. (Pax East 2013 Promo Pack A)`,
			`Getting bitten by a radioactive spider and then battling leukimia for 30 years. (Pax East 2014)`,
			`KHAAAAAAAAN! (Pax East 2014)`,
			`Loading from a previous save. (Pax East 2013 Promo Pack B)`,
			`Offering sexual favors for an ore and a sheep. (Pax Prime 2013)`,
			`Running out of stamina. (Pax East 2013 Promo Pack A)`,
			`Separate drinking fountains for dark elves. (Pax East 2014)`,
			`Ser Jorah Mormont's cerulean-blue balls. (Pax East 2014)`,
			`Sharpening a foam broadsword on a foam whetstone. (Pax East 2013 Promo Pack B)`,
			`Stuffing my balls into a Sega Genesis and pressing the power button. (Pax East 2014)`,
			`Taking 2d6 emotional damage. (Pax East 2014)`,
			`Tapping Serra Angel. (Pax Prime 2013)`,
			`The Cock Ring of Alacrity. (Pax Prime 2013)`,
			`The collective wail of every Magic player suddenly realizing that they've spent hundreds of dollars on pieces of cardboard. (Pax Prime 2013)`,
			`The depression that ensues after catching 'em all. (Pax East 2013 Promo Pack B)`,
			`Yoshi's huge egg-laying cloaca. (Pax Prime 2013)`,
		},
	}

	AddPack(pack)
}
