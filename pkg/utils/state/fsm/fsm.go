package fsm

import (
	"errors"
)

// State represents a state of a finite state machine.
type State string

// Event represents an event that triggers a state transition.
type Event struct {
	Name    string
	Payload interface{}
}

// StateHandler represents a handler for a state.
type StateHandler func(event Event, f *FSM) error

// FSM represents a finite state machine.
type FSM struct {
	CurrentState State
	Transitions  map[State]map[Event]State
	Handlers     map[State]StateHandler
}

// NewFSM creates a new finite state machine.
func NewFSM(initialState State) *FSM {
	return &FSM{
		CurrentState: initialState,
		Transitions:  make(map[State]map[Event]State),
	}
}

// AddState adds a new state to the finite state machine.
func (f *FSM) AddState(state State) {
	if _, exists := f.Transitions[state]; !exists {
		f.Transitions[state] = make(map[Event]State)
	}
}

// AddTransition adds a new transition to the finite state machine.
func (f *FSM) AddTransition(from State, event Event, to State) {
	if _, exists := f.Transitions[from]; !exists {
		f.AddState(from)
	}
	f.Transitions[from][event] = to
}

// Trigger triggers a state transition.
func (f *FSM) Trigger(event Event) error {
	if nextState, ok := f.Transitions[f.CurrentState][event]; ok {
		f.CurrentState = nextState
		if handler, exists := f.Handlers[nextState]; exists {
			return handler(event, f)
		}
		return nil
	}
	return errors.New("invalid event for current state")
}

// AddStateHandler adds a handler function for a specific state
func (f *FSM) AddStateHandler(state State, handler StateHandler) {
	if f.Handlers == nil {
		f.Handlers = make(map[State]StateHandler)
	}
	f.Handlers[state] = handler
}

// HandleState executes the handler associated with the current state
func (f *FSM) HandleState(event Event) error {
	if handler, exists := f.Handlers[f.CurrentState]; exists {
		return handler(event, f)
	}
	return errors.New("no handler for current state")
}
