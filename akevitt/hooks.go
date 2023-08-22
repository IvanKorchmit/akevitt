package akevitt

import "errors"

// Send the message to other current sessions
func (engine *Akevitt) Message(channel, message, username string, session ActiveSession) error {
	if engine.onMessage == nil {
		return errors.New("onMessage func is nil")
	}
	purgeDeadSessions(&engine.sessions, engine, engine.onDeadSession)

	for _, v := range engine.sessions {

		err := engine.onMessage(engine, v, channel, message, username)

		if session != v {
			v.GetApplication().Draw()
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// Invokes dialogue event.
// Make sure you have installed the hook during initalisation.
func (engine *Akevitt) Dialogue(dialogue *Dialogue, session ActiveSession) error {
	if engine.onDialogue == nil {
		return errors.New("dialogue callback is not installed")
	}

	return engine.onDialogue(engine, session, dialogue)
}
