package fsm

import "fmt"

type State string
type Event string

type ConditionFunc func(context map[string]interface{}) bool
type ActionFunc func(context map[string]interface{}) error

type Rule struct {
	Event      Event
	Conditions []ConditionFunc
	Actions    []ActionFunc
}

// RuleBasedFSM is a finite state machine that is based on rules
type RuleBasedFSM struct {
	Rules   []Rule
	Context map[string]interface{}
}

func NewRuleBasedFSM(rules []Rule) *RuleBasedFSM {
	return &RuleBasedFSM{
		Rules:   rules,
		Context: make(map[string]interface{}),
	}
}

func (fsm *RuleBasedFSM) Trigger(event Event) error {
	fmt.Printf("Triggering event: %s\n", event)

	// Iterate over all rules to find matching events
	for _, rule := range fsm.Rules {
		if rule.Event == event {
			// Check if all conditions are met
			conditionsMet := true
			for _, condition := range rule.Conditions {
				if !condition(fsm.Context) {
					conditionsMet = false
					break
				}
			}

			// If conditions are met, execute the actions
			if conditionsMet {
				fmt.Printf("Executing actions for event: %s\n", event)
				for _, action := range rule.Actions {
					if err := action(fsm.Context); err != nil {
						return err
					}
				}
				return nil
			}
		}
	}

	fmt.Printf("No valid actions found for event: %s\n", event)
	return fmt.Errorf("no valid actions found for event %s", event)
}
