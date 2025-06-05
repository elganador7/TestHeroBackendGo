package models

import (
	"encoding/json"
)

// StripeCustomer represents a customer in the Stripe customers table
type StripeCustomer struct {
	ID                  string          `gorm:"primaryKey"`
	Object              string          `gorm:"column:object"`
	Address             json.RawMessage `gorm:"type:jsonb"`
	Description         string
	Email               string          `gorm:"index"`
	Metadata            json.RawMessage `gorm:"type:jsonb"`
	Name                string
	Phone               string
	Shipping            json.RawMessage `gorm:"type:jsonb"`
	Balance             int64
	Created             int64
	Currency            string
	DefaultSource       string `gorm:"column:default_source"`
	Delinquent          bool
	Discount            json.RawMessage `gorm:"type:jsonb"`
	InvoicePrefix       string          `gorm:"column:invoice_prefix"`
	InvoiceSettings     json.RawMessage `gorm:"type:jsonb"`
	Livemode            bool
	NextInvoiceSequence int64           `gorm:"column:next_invoice_sequence"`
	PreferredLocales    json.RawMessage `gorm:"type:jsonb"`
	TaxExempt           string          `gorm:"column:tax_exempt"`
}

// TableName specifies the table name for the StripeCustomer model
func (StripeCustomer) TableName() string {
	return "customer"
}
