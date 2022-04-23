package monautil_test

import (
	"fmt"
	//"math"

	"github.com/monasuite/monad/monautil"
	"github.com/shopspring/decimal"
)

func ExampleAmount() {

	a := monautil.Amount(0)
	fmt.Println("Zero Satoshi:", a)

	a = monautil.Amount(1e8)
	fmt.Println("100,000,000 Satoshis:", a)

	a = monautil.Amount(1e5)
	fmt.Println("100,000 Satoshis:", a)
	// Output:
	// Zero Satoshi: 0 MONA
	// 100,000,000 Satoshis: 1 MONA
	// 100,000 Satoshis: 0.001 MONA
}

func ExampleNewAmount() {
	amountOne, err := monautil.NewAmount(decimal.NewFromFloat(1))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountOne) //Output 1

	amountFraction, err := monautil.NewAmount(decimal.NewFromFloat(0.01234567))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountFraction) //Output 2

	amountZero, err := monautil.NewAmount(decimal.NewFromFloat(0))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountZero) //Output 3

	amountNaN, err := monautil.NewAmount(decimal.NewFromFloat(1051200000000.00000001))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(amountNaN) //Output 4

	// Output: 1 MONA
	// 0.01234567 MONA
	// 0 MONA
	// invalid monacoin amount
}

func ExampleAmount_unitConversions() {
	amount := monautil.Amount(44433322211100)

	fmt.Println("Satoshi to kBTC:", amount.Format(monautil.AmountKiloBTC))
	fmt.Println("Satoshi to BTC:", amount)
	fmt.Println("Satoshi to MilliBTC:", amount.Format(monautil.AmountMilliBTC))
	fmt.Println("Satoshi to MicroBTC:", amount.Format(monautil.AmountMicroBTC))
	fmt.Println("Satoshi to Satoshi:", amount.Format(monautil.AmountSatoshi))

	// Output:
	// Satoshi to kBTC: 444.333222111 kMONA
	// Satoshi to BTC: 444333.222111 MONA
	// Satoshi to MilliBTC: 444333222.111 mMONA
	// Satoshi to MicroBTC: 444333222111 μMONA
	// Satoshi to Satoshi: 44433322211100 Watanabe
}
