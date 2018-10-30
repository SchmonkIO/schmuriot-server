package lobbyactions

import (
	"encoding/json"
	"fmt"

	"github.com/schmonk.io/schmuriot-server/actions/ingame"
	"github.com/schmonk.io/schmuriot-server/constants"
	"github.com/schmonk.io/schmuriot-server/models"
	"github.com/schmonk.io/schmuriot-server/utils"
)

type StartGameAction struct {
	models.BaseAction
	Rounds    int `json:"rounds"`
	Countdown int `json:"countdown"`
}

func StartGame(player *models.Player, message []byte, mt int) {
	if player.State != constants.StateLobby {
		models.SendJsonResponse(false, constants.ActionToggleReady, constants.ErrActionNotPossible.Error(), mt, player)
		return
	}
	data := StartGameAction{}
	err := json.Unmarshal(message, &data)
	if err != nil {
		utils.LogToConsole(err.Error())
		models.SendJsonResponse(false, constants.ActionStartGame, constants.ErrInvalidJSON.Error(), mt, player)
		return
	}
	fmt.Print("Game Rounds: ")
	fmt.Println(data.Rounds)
	fmt.Print("Game Countdown: ")
	fmt.Println(data.Countdown)
	r := models.Rooms.GetRoom(player.GetRoomID())
	if r != nil {
		ready := r.CheckAllReady()
		if !ready {
			models.SendJsonResponse(false, constants.ActionStartGame, constants.ErrNotReady.Error(), mt, player)
			return
		}
		for element := range r.Players {
			p := r.Players[element]
			p.SetState(constants.StateInGame)
		}
		data.Countdown = 15
		data.Rounds = 5
		game, _ := models.CreateCoinHunter(data.Rounds, data.Countdown)
		fmt.Print("Game info ")
		fmt.Println(game.Rounds)
		r.Game = &game
		ingameactions.RoundStart(player, mt)
		// r.SendToAllPlayers(true, constants.ActionStartGame, "", nil)
		return
	}
	models.SendJsonResponse(false, constants.ActionStartGame, constants.ErrRoomNotFound.Error(), mt, player)
	return
}
