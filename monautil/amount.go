// Copyright (c) 2013, 2014 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package monautil

import (
	"errors"
	"math"
	"strconv"

	"github.com/shopspring/decimal"
)

// AmountUnit describes a method of converting an Amount to something
// other than the base unit of a monacoin.  The value of the AmountUnit
// is the exponent component of the decadic multiple to convert from
// an amount in monacoin to an amount counted in units.
type AmountUnit int

// These constants define various units used when describing a monacoin
// monetary amount.
const (
	AmountMegaBTC  AmountUnit = 6
	AmountKiloBTC  AmountUnit = 3
	AmountBTC      AmountUnit = 0
	AmountMilliBTC AmountUnit = -3
	AmountMicroBTC AmountUnit = -6
	AmountSatoshi  AmountUnit = -8
)

// String returns the unit as a string.  For recognized units, the SI
// prefix is used, or "Satoshi" for the base unit.  For all unrecognized
// units, "1eN BTC" is returned, where N is the AmountUnit.
func (u AmountUnit) String() string {
	switch u {
	case AmountMegaBTC:
		return "MMONA"
	case AmountKiloBTC:
		return "kMONA"
	case AmountBTC:
		return "MONA"
	case AmountMilliBTC:
		return "mMONA"
	case AmountMicroBTC:
		return "μMONA"
	case AmountSatoshi:
		return "Watanabe"
	default:
		return "1e" + strconv.FormatInt(int64(u), 10) + " MONA"
	}
}

// Amount represents the base monacoin monetary unit (colloquially referred
// to as a `Satoshi').  A single Amount is equal to 1e-8 of a monacoin.
type Amount int64

// round converts a floating point number, which may or may not be representable
// as an integer, to the Amount integer type by rounding to the nearest integer.
// This is performed by adding or subtracting 0.5 depending on the sign, and
// relying on integer truncation to round the value to the nearest Amount.
//func round(f float64) Amount {
//	if f < 0 {
//		return Amount(f - 0.5)
//	}
//	return Amount(f + 0.5)
//}

// NewAmount creates an Amount from a floating point value representing
// some value in monacoin.  NewAmount errors if f is NaN or +-Infinity, but
// does not check that the amount is within the total amount of monacoin
// producible as f may not refer to an amount at a single moment in time.
//
// NewAmount is for specifically for converting BTC to Satoshi.
// For creating a new Amount with an int64 value which denotes a quantity of Satoshi,
// do a simple type conversion from type int64 to Amount.
// See GoDoc for example: http://godoc.org/github.com/monasuite/monautil#example-Amount
func NewAmount(d decimal.Decimal) (Amount, error) {
	// The amount is only considered invalid if it cannot be represented
	// as an integer type.  This may happen if f is NaN or +-Infinity.
	//switch {
	//case math.IsNaN(f):
	//	fallthrough
	//case math.IsInf(f, 1):
	//	fallthrough
	//case math.IsInf(f, -1):
	//	return 0, errors.New("invalid monacoin amount")
	//}

	// for mwat(Inf+-)
	if d.Cmp(decimal.NewFromInt(105.12e6*1e3)) == 1 {
		return Amount(0), errors.New("invalid monacoin amount")
	}
	if d.Cmp(decimal.NewFromInt(-105.12e6*1e3)) == -1 {
		return Amount(0), errors.New("invalid monacoin amount")
	}
	d2 := d.Round(8)
	decimalSatoshi := decimal.NewFromInt(SatoshiPerBitcoin)
	d3 := d2.Mul(decimalSatoshi)

	return Amount(d3.IntPart()), nil
}

// ToDecimalUnit converts a monetary amount counted in monacoin base units to a
// Decimal value representing an amount of monacoin.
func (a *Amount) ToDecimalUnit(u AmountUnit) decimal.Decimal {
	return decimal.NewFromInt(int64(*a)).Shift(int32(-u - 8))
}

// ToDecimalBTC is the equivalent of calling ToDecimalUnit with AmountBTC.
func (a *Amount) ToDecimalBTC() decimal.Decimal {
	return a.ToDecimalUnit(AmountBTC)
}

// ToUnit converts a monetary amount counted in monacoin base units to a
// floating point value representing an amount of monacoin.
func (a Amount) ToUnit(u AmountUnit) float64 {
	return float64(a) / math.Pow10(int(u+8))
}

// ToBTC is the equivalent of calling ToUnit with AmountBTC.
func (a Amount) ToBTC() float64 {
	return a.ToUnit(AmountBTC)
}

// Format formats a monetary amount counted in monacoin base units as a
// string for a given unit.  The conversion will succeed for any unit,
// however, known units will be formated with an appended label describing
// the units with SI notation, or "Satoshi" for the base unit.
func (a Amount) Format(u AmountUnit) string {
	units := " " + u.String()
	return a.ToDecimalUnit(u).String() + units
}

// String is the equivalent of calling Format with AmountBTC.
func (a Amount) String() string {
	return a.Format(AmountBTC)
}

// MulF64 multiplies an Amount by a floating point value.  While this is not
// an operation that must typically be done by a full node or wallet, it is
// useful for services that build on top of monacoin (for example, calculating
// a fee by multiplying by a percentage).
func (a Amount) MulF64(f float64) Amount {
	return Amount(decimal.NewFromInt(int64(a)).Mul(decimal.NewFromFloat(f)).Round(0).IntPart())
	//return round(float64(a) * f)
}
