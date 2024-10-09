package fsm

//
//import "fmt"
//
//type State string
//type Event string
//
//type ConditionFunc func(context map[string]interface{}) bool
//type ActionFunc func(context map[string]interface{}) error
//
//type Rule struct {
//	Event      Event
//	Conditions []ConditionFunc
//	Actions    []ActionFunc
//}
//
//// RuleBasedFSM is a finite state machine that is based on rules
//type RuleBasedFSM struct {
//	Rules   []Rule
//	Context map[string]interface{}
//}
//
//func NewRuleBasedFSM(rules []Rule) *RuleBasedFSM {
//	return &RuleBasedFSM{
//		Rules:   rules,
//		Context: make(map[string]interface{}),
//	}
//}
//
//func (fsm *RuleBasedFSM) Trigger(event Event) error {
//	fmt.Printf("Triggering event: %s\n", event)
//
//	// Iterate over all rules to find matching events
//	for _, rule := range fsm.Rules {
//		if rule.Event == event {
//			// Check if all conditions are met
//			conditionsMet := true
//			for _, condition := range rule.Conditions {
//				if !condition(fsm.Context) {
//					conditionsMet = false
//					break
//				}
//			}
//
//			// If conditions are met, execute the actions
//			if conditionsMet {
//				fmt.Printf("Executing actions for event: %s\n", event)
//				for _, action := range rule.Actions {
//					if err := action(fsm.Context); err != nil {
//						return err
//					}
//				}
//				return nil
//			}
//		}
//	}
//
//	fmt.Printf("No valid actions found for event: %s\n", event)
//	return fmt.Errorf("no valid actions found for event %s", event)
//}
//
////
////type State string
////type Event string
////
////type EventHandler func(payload interface{})
////
////type FSM struct {
////	CurrentState  State
////	PreviousState State
////	SubStates     map[State]*FSM
////	Transitions   map[State]map[Event]State
////	Handlers      map[Event]EventHandler
////}
////
////func NewFSM(initialState State) *FSM {
////	return &FSM{
////		CurrentState: initialState,
////		SubStates:    make(map[State]*FSM),
////		Transitions:  make(map[State]map[Event]State),
////		Handlers:     make(map[Event]EventHandler),
////	}
////}
////
////// AddSubState adds a substate FSM to a particular state
////func (f *FSM) AddSubState(state State, subFSM *FSM) {
////	f.SubStates[state] = subFSM
////}
////
////// AddTransition adds a valid transition from one state to another based on an event
////func (f *FSM) AddTransition(from State, event Event, to State) {
////	if f.Transitions[from] == nil {
////		f.Transitions[from] = make(map[Event]State)
////	}
////	f.Transitions[from][event] = to
////}
////
////// AddHandler assigns a handler function to an event
////func (f *FSM) AddHandler(event Event, handler EventHandler) {
////	f.Handlers[event] = handler
////}
////
////// TriggerEvent triggers an event and passes a payload to the handler
////func (f *FSM) TriggerEvent(event Event, payload interface{}) error {
////	// If the current state has a sub-FSM, delegate the event to it
////	if subFSM, exists := f.SubStates[f.CurrentState]; exists {
////		return subFSM.TriggerEvent(event, payload)
////	}
////
////	// If the current state has a transition for the event, handle it
////	if nextState, ok := f.Transitions[f.CurrentState][event]; ok {
////		// Execute the handler if it exists
////		if handler, handlerExists := f.Handlers[event]; handlerExists {
////			handler(payload)
////		}
////
////		// Track the previous state for going back
////		f.PreviousState = f.CurrentState
////
////		// Transition to the next state
////		fmt.Printf("Transitioning from %s to %s\n", f.CurrentState, nextState)
////		f.CurrentState = nextState
////		return nil
////	}
////
////	return errors.New(fmt.Sprintf("Invalid transition from %s on event %s", f.CurrentState, event))
////}
////
////// GoBack transitions back to the previous state
////func (f *FSM) GoBack(payload interface{}) error {
////	if f.PreviousState == "" {
////		return errors.New("no previous state to go back to")
////	}
////
////	fmt.Printf("Going back from %s to %s\n", f.CurrentState, f.PreviousState)
////	f.CurrentState, f.PreviousState = f.PreviousState, "" // Clear the previous state after going back
////	return nil
////}
////
////// GetCurrentState returns the current state of the FSM
////func (f *FSM) GetCurrentState() State {
////	return f.CurrentState
////}
////
////// GetSubState returns the current substate FSM if it exists
////func (f *FSM) GetSubState() *FSM {
////	if subFSM, exists := f.SubStates[f.CurrentState]; exists {
////		return subFSM
////	}
////	return nil
////}
//
////
////// State represents a state of a finite state machine.
////type State string
////
////// Event represents an event that triggers a state transition.
////type Event struct {
////	Name    string
////	Payload interface{}
////}
////
////// StateHandler represents a handler for a state.
////type StateHandler func(event Event, f *FSM) error
////
////// FSM represents a finite state machine.
////type FSM struct {
////	CurrentState State
////	Transitions  map[State]map[Event]State
////	Handlers     map[State]StateHandler
////}
////
////// NewFSM creates a new finite state machine.
////func NewFSM(initialState State) *FSM {
////	return &FSM{
////		CurrentState: initialState,
////		Transitions:  make(map[State]map[Event]State),
////	}
////}
////
////// AddState adds a new state to the finite state machine.
////func (f *FSM) AddState(state State) {
////	if _, exists := f.Transitions[state]; !exists {
////		f.Transitions[state] = make(map[Event]State)
////	}
////}
////
////// AddTransition adds a new transition to the finite state machine.
////func (f *FSM) AddTransition(from State, event Event, to State) {
////	if _, exists := f.Transitions[from]; !exists {
////		f.AddState(from)
////	}
////	f.Transitions[from][event] = to
////}
////
////// Trigger triggers a state transition.
////func (f *FSM) Trigger(event Event) error {
////	if nextState, ok := f.Transitions[f.CurrentState][event]; ok {
////		f.CurrentState = nextState
////		if handler, exists := f.Handlers[nextState]; exists {
////			return handler(event, f)
////		}
////		return nil
////	}
////	return errors.New("invalid event for current state")
////}
////
////// AddStateHandler adds a handler function for a specific state
////func (f *FSM) AddStateHandler(state State, handler StateHandler) {
////	if f.Handlers == nil {
////		f.Handlers = make(map[State]StateHandler)
////	}
////	f.Handlers[state] = handler
////}
////
////// HandleState executes the handler associated with the current state
////func (f *FSM) HandleState(event Event) error {
////	if handler, exists := f.Handlers[f.CurrentState]; exists {
////		return handler(event, f)
////	}
////	return errors.New("no handler for current state")
////}
