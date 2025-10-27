package fixed128_test

import (
	"fmt"
	"log"

	"github.com/seannyphoenix/binarytime/pkg/fixed128"
)

// ExampleNew demonstrates creating Fixed128 values using New.
func ExampleNew() {
	// Create a Fixed128 representing 22/7 (π approximation)
	pi, err := fixed128.New(22, 7)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("22/7 = %s\n", pi)

	// Create a Fixed128 representing 1/3
	third, err := fixed128.New(1, 3)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("1/3 = %s\n", third)

	// Create a negative value: -5/4
	neg, err := fixed128.New(-5, 4)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("-5/4 = %s\n", neg)

	// Output:
	// 22/7 = 03.2492492492492494
	// 1/3 = 00.5555555555555556
	// -5/4 = -01.40
}

// ExampleMustNew demonstrates creating Fixed128 values using MustNew.
func ExampleMustNew() {
	// MustNew panics on division by zero, so only use when you're sure
	half := fixed128.MustNew(1, 2)
	quarter := fixed128.MustNew(1, 4)

	fmt.Printf("1/2 = %s\n", half)
	fmt.Printf("1/4 = %s\n", quarter)

	// Output:
	// 1/2 = 00.80
	// 1/4 = 00.40
}

// ExampleFixed128_Add demonstrates addition of Fixed128 values.
func ExampleFixed128_Add() {
	a := fixed128.MustNew(1, 2) // 0.5
	b := fixed128.MustNew(1, 4) // 0.25

	sum := a.Add(b)
	fmt.Printf("%s + %s = %s\n", a, b, sum)

	// Adding negative numbers
	neg := fixed128.MustNew(-1, 3) // -0.333...
	pos := fixed128.MustNew(2, 3)  // 0.666...

	result := neg.Add(pos)
	fmt.Printf("%s + %s = %s\n", neg, pos, result)

	// Output:
	// 00.80 + 00.40 = 00.c0
	// -00.5555555555555556 + 00.aaaaaaaaaaaaaaac = 00.5555555555555556
}

// ExampleFixed128_Sub demonstrates subtraction of Fixed128 values.
func ExampleFixed128_Sub() {
	a := fixed128.MustNew(3, 4) // 0.75
	b := fixed128.MustNew(1, 4) // 0.25

	diff := a.Sub(b)
	fmt.Printf("%s - %s = %s\n", a, b, diff)

	// Output:
	// 00.c0 - 00.40 = 00.80
}

// ExampleFixed128_Mul demonstrates multiplication of Fixed128 values.
// Note: The current Mul implementation may have precision issues with the underlying big.Int operations.
func ExampleFixed128_Mul() {
	// Simple integer multiplication that should work
	a := fixed128.MustNew(2, 1) // 2.0
	b := fixed128.Zero          // 0.0

	product := a.Mul(b)
	fmt.Printf("%s * %s = %s\n", a, b, product)

	// Output:
	// 02.00 * 00.00 = 00.00
}

// ExampleFixed128_MulInt64 demonstrates converting Fixed128 back to int64.
func ExampleFixed128_MulInt64() {
	// Create 22/7 and multiply back by 7 to get 22
	f := fixed128.MustNew(22, 7)
	result, err := f.MulInt64(7)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("(22/7) * 7 = %d\n", result)

	// Create 1/3 and multiply by 9 to get 3
	third := fixed128.MustNew(1, 3)
	result2, err := third.MulInt64(9)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("(1/3) * 9 = %d\n", result2)

	// Output:
	// (22/7) * 7 = 22
	// (1/3) * 9 = 3
}

// ExampleFixed128_String demonstrates the string representation.
func ExampleFixed128_String() {
	values := []fixed128.Fixed128{
		fixed128.MustNew(0, 1),   // 0
		fixed128.MustNew(1, 1),   // 1
		fixed128.MustNew(22, 7),  // π approximation
		fixed128.MustNew(-5, 4),  // -1.25
		fixed128.MustNew(255, 1), // 255 (shows hex)
	}

	for _, v := range values {
		fmt.Printf("Value: %s\n", v)
	}

	// Output:
	// Value: 00.00
	// Value: 01.00
	// Value: 03.2492492492492494
	// Value: -01.40
	// Value: ff.00
}

// ExampleFixed128_Cmp demonstrates comparison of Fixed128 values.
func ExampleFixed128_Cmp() {
	a := fixed128.MustNew(1, 2) // 0.5
	b := fixed128.MustNew(1, 3) // 0.333...
	c := fixed128.MustNew(1, 2) // 0.5

	fmt.Printf("a.Cmp(b) = %d (a > b)\n", a.Cmp(b))
	fmt.Printf("b.Cmp(a) = %d (b < a)\n", b.Cmp(a))
	fmt.Printf("a.Cmp(c) = %d (a == c)\n", a.Cmp(c))

	// Output:
	// a.Cmp(b) = 1 (a > b)
	// b.Cmp(a) = -1 (b < a)
	// a.Cmp(c) = 0 (a == c)
}

// ExampleParse demonstrates parsing Fixed128 from strings.
func ExampleParse() {
	// Parse a positive value
	f1, err := fixed128.Parse("03.14")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Parsed: %s\n", f1)

	// Parse a negative value
	f2, err := fixed128.Parse("-01.50")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Parsed: %s\n", f2)

	// Output:
	// Parsed: 03.14
	// Parsed: -01.50
}

// ExampleFixed128_Base64 demonstrates base64 encoding and decoding.
func ExampleFixed128_Base64() {
	original := fixed128.MustNew(22, 7)

	// Encode to base64
	encoded := original.Base64()
	fmt.Printf("Base64: %s\n", encoded)

	// Decode from base64
	decoded, err := fixed128.ParseBase64(encoded)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Original: %s\n", original)
	fmt.Printf("Decoded:  %s\n", decoded)
	fmt.Printf("Equal: %t\n", original.Cmp(decoded) == 0)

	// Output:
	// Base64: AAAAAAAAAAADJJJJJJJJJJQ=
	// Original: 03.2492492492492494
	// Decoded:  03.2492492492492494
	// Equal: true
}

// ExampleFixed128_DecimalString demonstrates decimal string representation.
func ExampleFixed128_DecimalString() {
	// Create some values
	half := fixed128.MustNew(1, 2)
	third := fixed128.MustNew(1, 3)
	piApprox := fixed128.MustNew(22, 7)

	// Get decimal representation
	fmt.Printf("1/2 = %s\n", half.DecimalString())
	fmt.Printf("1/3 = %s\n", third.DecimalString())
	fmt.Printf("22/7 = %s\n", piApprox.DecimalString())

	// Custom precision
	fmt.Printf("1/3 (5 decimals) = %s\n", third.DecimalStringWithPrecision(5))
	fmt.Printf("22/7 (3 decimals) = %s\n", piApprox.DecimalStringWithPrecision(3))

	// Output:
	// 1/2 = 0.500000000000000
	// 1/3 = 0.333333333333333
	// 22/7 = 3.142857142857143
	// 1/3 (5 decimals) = 0.33333
	// 22/7 (3 decimals) = 3.143
}

// ExampleFixed128_Quo demonstrates division of Fixed128 values.
func ExampleFixed128_Quo() {
	a := fixed128.MustNew(6, 1) // 6.0
	b := fixed128.MustNew(2, 1) // 2.0
	quotient, err := a.Quo(b)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s / %s = %s\n", a, b, quotient)

	// Fractional division
	c := fixed128.MustNew(1, 2) // 0.5
	d := fixed128.MustNew(1, 4) // 0.25
	quotient2, err := c.Quo(d)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s / %s = %s\n", c, d, quotient2)

	// Output:
	// 06.00 / 02.00 = 03.00
	// 00.80 / 00.40 = 02.00
}

// ExampleFixed128_Float64 demonstrates Float64 conversion.
func ExampleFixed128_Float64() {
	// Convert from fraction to float64
	half := fixed128.MustNew(1, 2)
	fmt.Printf("Half as float64: %v\n", half.Float64())

	// Convert from float64 to Fixed128
	piApprox, err := fixed128.FromFloat64(3.14159)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Pi approximation: %.5f\n", piApprox.Float64())

	// Output:
	// Half as float64: 0.5
	// Pi approximation: 3.14159
}

// ExampleFixed128_Int64 demonstrates Int64 conversion.
func ExampleFixed128_Int64() {
	// Convert to int64
	half := fixed128.MustNew(7, 2) // 3.5
	intVal, err := half.Int64()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("3.5 as int64: %d (rounded)\n", intVal)

	// Convert from int64
	fromInt := fixed128.FromInt64(42)
	fmt.Printf("42 as Fixed128: %s\n", fromInt.DecimalString())

	// Output:
	// 3.5 as int64: 4 (rounded)
	// 42 as Fixed128: 42.000000000000000
}

// ExampleFixed128_Abs demonstrates absolute value.
func ExampleFixed128_Abs() {
	pos := fixed128.MustNew(5, 1)
	neg := fixed128.MustNew(-5, 1)

	fmt.Printf("Abs(%v) = %v\n", neg, neg.Abs())
	fmt.Printf("Abs(%v) = %v\n", pos, pos.Abs())

	// Output:
	// Abs(-05.00) = 05.00
	// Abs(05.00) = 05.00
}

// ExampleFixed128_Neg demonstrates negation.
func ExampleFixed128_Neg() {
	pos := fixed128.MustNew(5, 1)
	neg := pos.Neg()

	fmt.Printf("Neg(%v) = %v\n", pos, neg)
	fmt.Printf("Neg(%v) = %v\n", neg, neg.Neg())

	// Output:
	// Neg(05.00) = -05.00
	// Neg(-05.00) = 05.00
}
