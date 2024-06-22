package data

import "gorm.io/gorm"

type Order struct {
	gorm.Model

	AccountId uint
	Account

	PlanId uint
	Plan   Plan

	Amount      float64
	CurrencyISO string
}
