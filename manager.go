package cardsagainstdiscord

import (
	"github.com/jonas747/discordgo"
	"sync"
)

type GameManager struct {
	sync.RWMutex
	SessionProvider SessionProvider
	ActiveGames     map[int64]*Game
	NumActiveGames  int
}

func NewGameManager(sessionProvider SessionProvider) *GameManager {
	return &GameManager{
		ActiveGames:     make(map[int64]*Game),
		SessionProvider: sessionProvider,
	}
}

func (gm *GameManager) CreateGame(guildID int64, channelID int64, userID int64, username string, packs ...string) (*Game, error) {
	for _, v := range packs {
		_, ok := Packs[v]
		if !ok {
			return nil, &ErrUnknownPack{
				PassedPack: v,
			}
		}
	}

	if len(packs) < 1 {
		return nil, ErrNoPacks
	}

	gm.Lock()
	defer gm.Unlock()

	if _, ok := gm.ActiveGames[channelID]; ok {
		return nil, ErrGameAlreadyInChannel
	}

	if _, ok := gm.ActiveGames[userID]; ok {
		return nil, ErrPlayerAlreadyInGame
	}

	game := &Game{
		MasterChannel: channelID,
		Manager:       gm,
		GuildID:       guildID,
		Packs:         packs,
		GameMaster:    userID,
		PlayerLimit:   10,
		WinLimit:      10,
		Session:       gm.SessionProvider.SessionForGuild(guildID),
	}

	game.Created()

	game.AddPlayer(userID, username)

	gm.ActiveGames[channelID] = game
	gm.ActiveGames[userID] = game
	gm.NumActiveGames++

	return game, nil
}

func (gm *GameManager) FindGameFromChannelOrUser(id int64) *Game {
	gm.RLock()
	defer gm.RUnlock()

	if g, ok := gm.ActiveGames[id]; ok {
		return g
	}

	return nil
}

func (gm *GameManager) PlayerTryJoinGame(gameID, playerID int64, username string) error {
	gm.Lock()
	defer gm.Unlock()

	if _, ok := gm.ActiveGames[playerID]; ok {
		return ErrPlayerAlreadyInGame
	}

	if g, ok := gm.ActiveGames[gameID]; ok {
		if g.AddPlayer(playerID, username) {
			gm.ActiveGames[playerID] = g
			return nil
		}

		return ErrGameFull
	}

	return ErrGameNotFound
}

func (gm *GameManager) PlayerTryLeaveGame(playerID int64) error {
	gm.Lock()
	defer gm.Unlock()

	if g, ok := gm.ActiveGames[playerID]; ok {
		delete(gm.ActiveGames, playerID)
		g.RemovePlayer(playerID)
		return nil
	}

	return ErrGameNotFound
}

func (gm *GameManager) RemoveGame(gameID int64) error {
	gm.Lock()
	defer gm.Unlock()

	g, ok := gm.ActiveGames[gameID]
	if !ok {
		return ErrGameNotFound
	}

	g.Stop()

	// Remove all references to the game
	g.RLock()
	defer g.RUnlock()

	delete(gm.ActiveGames, g.MasterChannel)
	delete(gm.ActiveGames, g.GameMaster)
	for _, v := range g.Players {
		delete(gm.ActiveGames, v.ID)
	}

	gm.NumActiveGames--

	return nil
}

func (gm *GameManager) HandleReactionAdd(ra *discordgo.MessageReactionAdd) {
	cid := ra.ChannelID
	userID := ra.UserID

	gm.RLock()
	if game, ok := gm.ActiveGames[cid]; ok {
		gm.RUnlock()
		game.HandleRectionAdd(ra)
	} else if game, ok := gm.ActiveGames[userID]; ok {
		gm.RUnlock()
		game.HandleRectionAdd(ra)
	} else {
		gm.RUnlock()
	}
}

func (gm *GameManager) HandleMessageCreate(msgCreate *discordgo.MessageCreate) {
	userID := msgCreate.Author.ID

	gm.RLock()
	if game, ok := gm.ActiveGames[userID]; ok {
		gm.RUnlock()
		game.HandleMessageCreate(msgCreate)
	} else {
		gm.RUnlock()
	}
}
