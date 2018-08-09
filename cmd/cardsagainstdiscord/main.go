package main

import (
	"github.com/jonas747/cardsagainstdiscord"
	"github.com/jonas747/dcmd"
	"github.com/jonas747/discordgo"
	"github.com/jonas747/dutil/dstate"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
)

var cahManager *cardsagainstdiscord.GameManager

func panicErr(err error, msg string) {
	if err != nil {
		panic(msg + ": " + err.Error())
	}
}

func main() {
	session, err := discordgo.New(os.Getenv("DG_TOKEN"))
	panicErr(err, "Failed initializing discordgo")

	cahManager = cardsagainstdiscord.NewGameManager(&cardsagainstdiscord.StaticSessionProvider{
		Session: session,
	})

	state := dstate.NewState()
	state.TrackMembers = false
	state.TrackPresences = false
	session.StateEnabled = false

	cmdSys := dcmd.NewStandardSystem("!cah")
	cmdSys.State = state
	cmdSys.Root.AddCommand(CreateGameCommand, dcmd.NewTrigger("create", "c").SetDisableInDM(true))
	cmdSys.Root.AddCommand(StopCommand, dcmd.NewTrigger("stop", "end", "s").SetDisableInDM(true))
	cmdSys.Root.AddCommand(KickCommand, dcmd.NewTrigger("kick").SetDisableInDM(true))
	cmdSys.Root.AddCommand(PacksCommand, dcmd.NewTrigger("packs").SetDisableInDM(true))

	session.AddHandler(state.HandleEvent)
	session.AddHandler(cmdSys.HandleMessageCreate)
	session.AddHandler(func(s *discordgo.Session, ra *discordgo.MessageReactionAdd) {
		go cahManager.HandleReactionAdd(ra)
	})

	session.AddHandler(func(s *discordgo.Session, msg *discordgo.MessageCreate) {
		go cahManager.HandleMessageCreate(msg)
	})

	err = session.Open()
	panicErr(err, "Failed opening gateway connection")
	log.Println("Running...")

	// We import http/pprof above to be ale to inspect shizz and do profiling
	go http.ListenAndServe(":7447", nil)
	select {}
}

var CreateGameCommand = &dcmd.SimpleCmd{
	ShortDesc: "Creates a cards against humanity game in this channel",
	CmdArgDefs: []*dcmd.ArgDef{
		&dcmd.ArgDef{Name: "packs", Type: dcmd.String, Default: "main", Help: "Packs seperated by space"},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		pStr := data.Args[0].Str()
		packs := strings.Fields(pStr)

		_, err := cahManager.CreateGame(data.GS.ID, data.CS.ID, data.Msg.Author.ID, data.Msg.Author.Username, packs...)
		if err == nil {
			log.Println("Created a new game in ", data.CS.ID)
			return "", nil
		}

		if cahErr := cardsagainstdiscord.HumanizeError(err); cahErr != "" {
			return cahErr, nil
		}

		return "Something went wrong", err

	},
}

var StopCommand = &dcmd.SimpleCmd{
	ShortDesc: "Stops a cards against humanity game in this channel",
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		err := cahManager.TryAdminRemoveGame(data.Msg.Author.ID)
		if err != nil {
			if cahErr := cardsagainstdiscord.HumanizeError(err); cahErr != "" {
				return cahErr, nil
			}

			return "Something went wrong", err
		}

		return "Stopped the game", nil
	},
}

var KickCommand = &dcmd.SimpleCmd{
	ShortDesc:       "Kicks a player from the card against humanity game in this channel, only the game master can do this",
	RequiredArgDefs: 1,
	CmdArgDefs: []*dcmd.ArgDef{
		&dcmd.ArgDef{Name: "user", Type: dcmd.UserID},
	},
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		userID := data.Args[0].Int64()

		err := cahManager.AdminKickUser(data.Msg.Author.ID, userID)
		if err != nil {
			if cahErr := cardsagainstdiscord.HumanizeError(err); cahErr != "" {
				return cahErr, nil
			}

			return "Something went wrong", err
		}

		return "User removed", nil
	},
}

var PacksCommand = &dcmd.SimpleCmd{
	ShortDesc: "Lists available packs",
	RunFunc: func(data *dcmd.Data) (interface{}, error) {
		resp := "Available packs: \n\n"
		for _, v := range cardsagainstdiscord.Packs {
			resp += "`" + v.Name + "` - " + v.Description + "\n"
		}

		return resp, nil
	},
}
