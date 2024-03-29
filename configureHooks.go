package akevitt

// Accepts the function which gets invoked when someone sends the message (engine.Message)
func (builder *akevittBuilder) UseOnMessage(f MessageFunc) *akevittBuilder {
	builder.engine.onMessage = f

	return builder
}

// Called when engine.Dialogue is called
func (builder *akevittBuilder) UseOnDialogue(f DialogueFunc) *akevittBuilder {
	builder.engine.onDialogue = f

	return builder
}

// Accepts function which gets called when the user lefts the game.
// Note: use with caution, because calling methods from the engine like Message
// will cause an infinite recursion
// and in result: the application will crash.
func (builder *akevittBuilder) UseOnSessionEnd(f DeadSessionFunc) *akevittBuilder {
	builder.engine.onDeadSession = f

	return builder
}
