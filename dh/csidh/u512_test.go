package csidh

import (
	"fmt"
	"math/big"
	mrand "math/rand"
	"testing"
)

// Number of iterations
var (
	kNumIter    = 10000
	kModulus, _ = new(big.Int).SetString(u512toS(p), 16)
)

func u512toS(val u512) string {
	var str string
	for i := 0; i < 8; i++ {
		str = fmt.Sprintf("%016x", val[i]) + str
	}
	return str
}

// returns random value in a range (0,p)
func randomU512() u512 {
	var u u512
	for i := 0; i < 8; i++ {
		u[i] = mrand.Uint64()
	}
	return u
}

// Check if fp512Mul3 produces result
// z = x*y mod 2^512
func TestFp512Mul3_Nominal(t *testing.T) {
	var multiplier64 uint64
	var mod big.Int

	// modulus: 2^512
	mod.SetUint64(1).Lsh(&mod, 512)

	for i := 0; i < kNumIter; i++ {
		multiplier64 = mrand.Uint64()

		fV := randomU512()
		exp, _ := new(big.Int).SetString(u512toS(fV), 16)
		exp.Mul(exp, new(big.Int).SetUint64(multiplier64))
		// Truncate to 512 bits
		exp.Mod(exp, &mod)

		fp512Mul3(&fV, &fV, multiplier64)
		res, _ := new(big.Int).SetString(u512toS(fV), 16)

		if exp.Cmp(res) != 0 {
			t.Errorf("%X != %X", exp, res)
		}
	}
}

func TestFp512Add3_Nominal(t *testing.T) {
	var ret u512
	var mod big.Int
	// modulus: 2^512
	mod.SetUint64(1).Lsh(&mod, 512)

	for i := 0; i < kNumIter; i++ {
		a := randomU512()
		bigA, _ := new(big.Int).SetString(u512toS(a), 16)
		b := randomU512()
		bigB, _ := new(big.Int).SetString(u512toS(b), 16)

		fp512Add3(&ret, &a, &b)
		bigRet, _ := new(big.Int).SetString(u512toS(ret), 16)
		bigA.Add(bigA, bigB)
		// Truncate to 512 bits
		bigA.Mod(bigA, &mod)

		if bigRet.Cmp(bigA) != 0 {
			t.Errorf("%X != %X", bigRet, bigA)
		}
	}
}

func TestFp512Add3_ReturnsCarry(t *testing.T) {
	a := u512{}
	b := u512{
		0, 0,
		0, 0,
		0, 0,
		0, 1}
	c := u512{
		0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF,
		0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF,
		0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF,
		0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF}
	if fp512Add3(&a, &b, &c) != 1 {
		t.Error("Carry not returned")
	}
}

func TestFp512Add3_DoesntReturnCarry(t *testing.T) {
	a := u512{
		0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF,
		0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF,
		0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF,
		0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF}
	b := u512{
		0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF,
		0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF,
		0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF,
		0xFFFFFFFFFFFFFFFF, 0}
	c := u512{
		0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF,
		0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF,
		0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF,
		0xFFFFFFFFFFFFFFFA, 0}

	if fp512Add3(&a, &b, &c) != 0 {
		t.Error("Carry returned")
	}
}

func TestFp512Sub3_Nominal(t *testing.T) {
	var ret u512
	var mod big.Int
	// modulus: 2^512
	mod.SetUint64(1).Lsh(&mod, 512)

	for i := 0; i < kNumIter; i++ {
		a := randomU512()
		bigA, _ := new(big.Int).SetString(u512toS(a), 16)
		b := randomU512()
		bigB, _ := new(big.Int).SetString(u512toS(b), 16)

		fp512Sub3(&ret, &a, &b)
		bigRet, _ := new(big.Int).SetString(u512toS(ret), 16)
		bigA.Sub(bigA, bigB)
		// Truncate to 512 bits
		bigA.Mod(bigA, &mod)

		if bigRet.Cmp(bigA) != 0 {
			t.Errorf("%X != %X", bigRet, bigA)
		}
	}
}

func TestFp512Sub3_DoesntReturnCarry(t *testing.T) {
	a := u512{}
	b := u512{
		0xFFFFFFFFFFFFFFFF, 1,
		0, 0,
		0, 0,
		0, 0}
	c := u512{
		0xFFFFFFFFFFFFFFFF, 2,
		0, 0,
		0, 0,
		0, 0}

	if fp512Sub3(&a, &b, &c) != 1 {
		t.Error("Carry not returned")
	}
}

func TestFp512Sub3_ReturnsCarry(t *testing.T) {
	a := u512{}
	b := u512{
		0xFFFFFFFFFFFFFFFF, 2,
		0, 0,
		0, 0,
		0, 0}
	c := u512{
		0xFFFFFFFFFFFFFFFF, 1,
		0, 0,
		0, 0,
		0, 0}

	if fp512Sub3(&a, &b, &c) != 0 {
		t.Error("Carry not returned")
	}
}

func BenchmarkFp512Add(b *testing.B) {
	var arg1 u512
	arg2 := randomU512()
	arg3 := randomU512()
	for n := 0; n < b.N; n++ {
		fp512Add3(&arg1, &arg2, &arg3)
	}
}

func BenchmarkFp512Sub(b *testing.B) {
	var arg1 u512
	arg2, arg3 := randomU512(), randomU512()
	for n := 0; n < b.N; n++ {
		fp512Sub3(&arg1, &arg2, &arg3)
	}
}

func BenchmarkFp512Mul(b *testing.B) {
	var arg1 = mrand.Uint64()
	arg2, arg3 := randomU512(), randomU512()
	for n := 0; n < b.N; n++ {
		fp512Mul3(&arg2, &arg3, arg1)
	}
}
