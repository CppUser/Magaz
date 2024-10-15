package client

type RuleConfig struct {
	Rules []RuleEntry `mapstructure:"rules"`
}

type RuleEntry struct {
	Event      string       `mapstructure:"event"`
	Conditions []string     `mapstructure:"conditions"`
	Actions    []ActionItem `mapstructure:"actions"`
}

type ActionItem struct {
	Type   string            `mapstructure:"type"`
	Text   string            `mapstructure:"text,omitempty"`
	Markup map[string]string `mapstructure:"markup,omitempty"`
}

func LoadRulesUsingViper() (*RuleConfig, error) {
	return nil, nil
}

//
//func mapActions(actions []ActionItem) []fsm.ActionFunc {
//	var fsmActions []fsm.ActionFunc
//
//	for _, action := range actions {
//		switch action.Type {
//		case "send_message":
//			fsmActions = append(fsmActions, handlers.SendMessageWithMarkup(bot, action.Text, action.Markup))
//		case "edit_message":
//			fsmActions = append(fsmActions, handlers.EditMessage(bot, action.Text))
//		case "fetch_cities_and_prompt_selection":
//			fsmActions = append(fsmActions, fetchCitiesAndPromptSelection(bot, b.DB))
//		case "fetch_products_for_city":
//			fsmActions = append(fsmActions, fetchProductsForCity(bot, b.DB, b.Cache))
//		case "fetch_quantity_for_product":
//			fsmActions = append(fsmActions, fetchQuantityForProduct(bot, b.DB, b.Cache))
//		case "fetch_available_addresses_and_payment":
//			fsmActions = append(fsmActions, fetchAvailableAddressesAndPayment(bot, b.DB, b.Cache))
//		}
//	}
//
//	return fsmActions
//}
//
//func mapConditions(conditions []string) []fsm.ConditionFunc {
//	// Placeholder for mapping conditions, depending on your condition handling logic
//	return []fsm.ConditionFunc{}
//}
//
//
//func mapRules(ruleConfig *RuleConfig) []fsm.Rule {
//	var fsmRules []fsm.Rule
//
//	for _, rule := range ruleConfig.Rules {
//		actions := mapActions(rule.Actions)
//		conditions := mapConditions(rule.Conditions)
//
//		// Create a rule for each event
//		for _, event := range rule.Events {
//			fsmRule := fsm.Rule{
//				Event:      event,
//				Conditions: conditions,
//				Actions:    actions,
//			}
//			fsmRules = append(fsmRules, fsmRule)
//		}
//	}
//
//	return fsmRules
//}

////////////////////////////////////////////////////////////////////
//  TODO: Refactor rules
////////////////////////////////////////////////////////////////////
