/*
Program written by Ivan Korchmit (c) 2023
Licensed under European Union Public Licence 1.2.
For more information, view LICENCE or README
*/

package akevitt

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/gliderlabs/ssh"
	"github.com/rivo/tview"
)

type Akevitt struct {
	sessions    Sessions
	root        UIFunc
	bind        string
	mouse       bool
	dbPath      string
	commands    map[string]CommandFunc
	db          *bolt.DB
	onMessage   MessageFunc
	defaultRoom Room
}

// Engine default constructor
func NewEngine() *Akevitt {
	engine := &Akevitt{}
	engine.bind = ":2222"
	engine.sessions = make(Sessions)
	engine.dbPath = "data/database.db"
	engine.mouse = false
	return engine
}

func (engine *Akevitt) UseBind(bindAddress string) *Akevitt {
	engine.bind = bindAddress

	return engine
}

func (engine *Akevitt) UseRootUI(uiFunc UIFunc) *Akevitt {
	engine.root = uiFunc

	return engine
}

func (engine *Akevitt) UseDBPath(path string) *Akevitt {
	engine.dbPath = path

	return engine
}

func (engine *Akevitt) UseMouse() *Akevitt {
	engine.mouse = true

	return engine
}

func (engine *Akevitt) RegisterCommand(command string, function CommandFunc) *Akevitt {
	command = strings.TrimSpace(command)
	engine.commands[command] = function
	return engine
}

func (engine *Akevitt) Login(username, password string, session ActiveSession) error {
	account, err := login(username, password, engine.db)
	if err != nil {
		return err
	}
	if isSessionAlreadyActive(*account, &engine.sessions) {
		return errors.New("the session is already active")
	}
	session.SetAccount(account)

	return nil
}

func (engine *Akevitt) Register(username, password string, session ActiveSession) error {
	exists := isAccountExists(username, engine.db)

	if exists {
		return errors.New("account already exists")
	}
	account, err := createAccount(engine.db, username, password)
	session.SetAccount(account)

	return err
}

func (engine *Akevitt) ProcessCommand(command string, session ActiveSession) error {
	zeroArg := strings.Fields(command)[0]
	noZeroArgArray := strings.Fields(command)[1:]
	noZeroArg := strings.Join(noZeroArgArray, " ")
	commandFunc, ok := engine.commands[zeroArg]
	if !ok {
		return errors.New("command not found")
	}

	return commandFunc(engine, session, noZeroArg)
}

func (engine *Akevitt) UseOnMessage(f MessageFunc) *Akevitt {
	engine.onMessage = f

	return engine
}

func (engine *Akevitt) UseSpawnRoom(r Room) *Akevitt {
	engine.defaultRoom = r

	return engine
}

func (engine *Akevitt) GetSpawnRoom() Room {
	return engine.defaultRoom
}

func (engine *Akevitt) SaveGameObject(gameObject GameObject, key uint64, account *Account) error {
	return overwriteObject(engine.db, key, account.Username, gameObject)
}

func (engine *Akevitt) OnMessage(channel, message string, session ActiveSession) error {
	if engine.onMessage == nil {
		return errors.New("onMessage func is nil")
	}
	return engine.onMessage(engine, session, channel, message)
}

func (engine *Akevitt) Run(sessionTemplate ActiveSession) error {
	fmt.Println("Running Akevitt")

	err := createDatabase(engine)

	if err != nil {
		log.Fatal(err)
	}

	defer engine.db.Close()

	gob.Register(Account{})

	if engine.root == nil {
		return errors.New("base screen is not provided")
	}

	ssh.Handle(func(sesh ssh.Session) {
		screen, err := newSessionScreen(sesh)
		if err != nil {
			fmt.Fprintln(sesh.Stderr(), "unable to create screen:", err)
			return
		}
		purgeDeadSessions(&engine.sessions)
		app := tview.NewApplication().SetScreen(screen).EnableMouse(engine.mouse)
		engine.sessions[sesh] = sessionTemplate
		engine.sessions[sesh].SetApplication(app)
		engine.sessions[sesh].GetApplication().SetRoot(engine.root(engine, engine.sessions[sesh]), true)
		if err := app.Run(); err != nil {
			fmt.Fprintln(sesh.Stderr(), err)
			return
		}
		sesh.Exit(0)
	})
	return ssh.ListenAndServe(engine.bind, nil)
}
