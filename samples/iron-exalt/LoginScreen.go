package main

import (
	"akevitt/akevitt"
	"fmt"

	"github.com/rivo/tview"
)

func loginScreen(engine *akevitt.Akevitt, session *ActiveSession) tview.Primitive {
	var username string
	var password string
	loginScreen := tview.NewForm().
		AddInputField("Username: ", "", 32, nil, func(text string) {
			username = text
		}).
		AddPasswordField("Password: ", "", 32, '*', func(text string) {
			password = text
		})
	loginScreen.
		AddButton("Login", func() {
			err := engine.Login(username, password, session)
			if err != nil {
				ErrorBox(err.Error(), session.app, session.previousUI)
				return
			}
			character, err := akevitt.FindObject[*Character](engine, session, CharacterKey)

			if err != nil {
				ErrorBox(err.Error(), session.app, session.previousUI)
				return
			}
			session.character = character
			fmt.Printf("session.character.CurrentRoomKey: %v\n", session.character.CurrentRoomKey)
			room, err := engine.GetRoom(session.character.CurrentRoomKey)

			if err != nil {
				ErrorBox(err.Error(), session.app, session.previousUI)
				return
			}
			session.character.account = session.account
			session.character.currentRoom = room
			session.SetRoot(gameScreen(engine, session))
		}).
		AddButton("Back", func() {
			session.app.SetRoot(rootScreen(engine, session), true)
		})
	return loginScreen
}
