package entity

type Transaction struct {
	ID        string
	Type      string // "topup" or "transfer"
	Amount    float64
	CreatedAt string
}
