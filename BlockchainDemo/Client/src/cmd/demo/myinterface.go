package demo

import (
	"fmt"
)

// Checkouter checkouts order
type Payment interface {
	// Pay from email to email this amount
	Pay(fromEmail, toEmail string, amount float64) error
}

type BankAdapter struct {
	name string
}

func (b *BankAdapter) Pay(fromEmail, toEmail string, amount float64) error {
	fmt.Println(b.name, fromEmail, toEmail, amount)
	return nil
}

type PaypalAdapter struct {
	name string
}

func NewBankAdapter() *BankAdapter {
	return &BankAdapter{name: "Bank"}
}

func NewPaypalAdapter() *PaypalAdapter {
	return &PaypalAdapter{name: "pay"}
}

func (p *PaypalAdapter) Pay(fromEmail, toEmail string, amount float64) error {
	fmt.Println(p.name, fromEmail, toEmail, amount)
	return nil
}

func PrintAdapter(p Payment) {
	p.Pay("from", "to", 1666.6)
}

/*
type I interface {
	M()
}

type T1 struct{}

func (T1) M() { fmt.Println("T1.M") }

type T2 struct{}

func (T2) M() { fmt.Println("T2.M") }

func F(i I) { i.M() }
*/
