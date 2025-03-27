// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Mutation struct {
}

type Payment struct {
	ID            string  `json:"id"`
	StudentID     string  `json:"studentId"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	TransactionID string  `json:"transactionId"`
	Status        string  `json:"status"`
	CreatedAt     string  `json:"createdAt"`
	Refund        *Refund `json:"refund,omitempty"`
}

type Query struct {
}

type Refund struct {
	ID        string `json:"id"`
	PaymentID string `json:"paymentId"`
	Reason    string `json:"reason"`
	IssuedAt  string `json:"issuedAt"`
}
