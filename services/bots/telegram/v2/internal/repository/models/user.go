package models

type CachedUserState struct {
	ID             int64  `json:"id"`
	PlacingOrder   bool   `json:"placing_order"`
	CancelingOrder bool   `json:"canceling_order"`
	OrderStatus    string `json:"status"`
	State          string `json:"state"`
}
