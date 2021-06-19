package cardsagainstdiscord

func init() {
	pack := &CardPack{
		Name:		 "science",
		Description: "Science Pack - The Science Pack is a pack about the hit system of knowledge known worldwide as “science.” Written with Phil Plait (Bad Astronomy) and Zach Weinersmith (SMBC). 30 cards about captivating theories like evolution and global warming, with a special guest appearance by Uranus. All profits donated to the Cards Against Humanity & SMBC Science Ambassador Scholarship for women in STEM. Released March 30, 2015. Icon: An Erlenmeyer flask.",
		Prompts: []*PromptCard{
			&PromptCard{Prompt: `A study published in Nature this week found that %s is good for you in small doses.`},
			&PromptCard{Prompt: `Hey there, Young Scientists! Put on your labcoats and strap on your safety goggles, because today we're learning about %s!`},
			&PromptCard{Prompt: `In an attempt to recreate conditions just after the Big Bang, physicists at the LHC are observing collisions between %s and %s.`},
			&PromptCard{Prompt: `In line with our predictions, we find a robust correlation between %s and %s (p<.05).`},
			&PromptCard{Prompt: `In what's being hailed as a major breakthrough, scientists have synthesized %s in the lab.`},
			&PromptCard{Prompt: `Today on Mythbusters, we found out how long %s can withstand %s.`},
			&PromptCard{Prompt: `What really killed the dinosaurs?`},
		},
		Responses: []ResponseCard{
			`3.7 billion years of evolution.`,
			`A 0.7 waist-to-hip ratio.`,
			`A supermassive black hole.`,
			`Achieving reproductive success.`,
			`Being knowledgeable in a narrow domain that nobody understands or cares about.`,
			`David Attenborough watching us mate.`,
			`Developing secondary sex characteristics.`,
			`Driving into a tornado to learn about tornadoes.`,
			`Electroejaculating a capuchin monkey.`,
			`Evolving a labyrinthine vagina.`,
			`Explosive decompression.`,
			`Failing the Turing test.`,
			`Fun and interesting facts about rocks.`,
			`Getting really worried about global warming for a few seconds.`,
			`Infinity.`,
			`Insufficient serotonin.`,
			`Oxytocin release via manual stimulation of the nipples.`,
			`Photosynthesis.`,
			`Reconciling quantum theory with general relativity.`,
			`Slowly evaporating.`,
			`The quiet majesty of the sea turtle.`,
			`The Sun engulfing the Earth.`,
			`Uranus.`,
		},
	}

	AddPack(pack)
}
