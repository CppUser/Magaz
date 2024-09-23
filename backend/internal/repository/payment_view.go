package repository

type CardView struct {
	BankName   string `json:"bank_name"`
	BankUrl    string `json:"bank_url"`
	CardNumber string `json:"card_number"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	UserName   string `json:"user_name"`
	Password   string `json:"password"`
	QuickPay   string `json:"quick_pay"`
}

type CryptoView struct {
	WalletName    string `json:"wallet_name"`
	WalletAddress string `json:"wallet_address"`
	WalletURL     string `json:"wallet_url"`
	UserName      string `json:"user_name"`
	Password      string `json:"password"`
}

type PaymentView struct {
	PaymentCategory string     `json:"payment_category"`
	CardPayment     CardView   `json:"card_payment"`
	CryptoPayment   CryptoView `json:"crypto_payment"`
}
