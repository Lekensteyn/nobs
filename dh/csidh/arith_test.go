package csidh

import (
	"fmt"
	"math/big"
	mrand "math/rand"
	"testing"
)

// Commonly used variables
var (
	// Number of interations
	kNumIter = 10000
	// Modulus
	kModulus, _ = new(big.Int).SetString(u512toS(p), 16)
	// Zero in Fp512
	ZeroFp512 = Fp{}
	// One in Fp512
	OneFp512 = Fp{v: u512{1, 0, 0, 0, 0, 0, 0, 0}}
)

func u512toS(v u512) string {
	var str string
	for i := 0; i < 8; i++ {
		str = fmt.Sprintf("%016x", v[i]) + str
	}
	return str
}

// zeroize u512
func zero(v *u512) {
	for i, _ := range *v {
		v[i] = 0
	}
}

// returns random value in a range (0,p)
func randomU512() u512 {
	var u u512
	for i := 0; i < 8; i++ {
		u[i] = mrand.Uint64()
	}
	return u
}

// return true if x==y, otherwise false
func cmp512(x, y *u512) bool {
	for i, _ := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return len(*x) == len(*y)
}

// Check if mul512 produces result
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

		mul512(&fV, &fV, multiplier64)
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

		add512(&ret, &a, &b)
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
	if add512(&a, &b, &c) != 1 {
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

	if add512(&a, &b, &c) != 0 {
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

		sub512(&ret, &a, &b)
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

	if sub512(&a, &b, &c) != 1 {
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

	if sub512(&a, &b, &c) != 0 {
		t.Error("Carry not returned")
	}
}

func TestCswap(t *testing.T) {
	arg1 := randomU512()
	arg2 := randomU512()

	arg1cpy := arg1
	cswap512(&arg1, &arg2, 0)
	if !cmp512(&arg1, &arg1cpy) {
		t.Error("cswap swapped")
	}

	arg1cpy = arg1
	cswap512(&arg1, &arg2, 1)
	if cmp512(&arg1, &arg1cpy) {
		t.Error("cswap didn't swapped")
	}

	arg1cpy = arg1
	cswap512(&arg1, &arg2, 0xF2)
	if cmp512(&arg1, &arg1cpy) {
		t.Error("cswap didn't swapped")
	}
}

func TestAddRdc(t *testing.T) {
	var res Fp
	var Fp512 = Fp{v: p}

	tmp := OneFp512
	addRdc(&res, &tmp, &Fp512)
	if !cmp512(&res.v, &tmp.v) {
		t.Errorf("Wrong value\n%X", res.v)
	}

	tmp = ZeroFp512
	addRdc(&res, &Fp512, &Fp512)
	if !cmp512(&res.v, &Fp512.v) {
		t.Errorf("Wrong value\n%X", res.v)
	}

	tmp = Fp{v: u512{1, 1, 1, 1, 1, 1, 1, 1}}
	addRdc(&res, &Fp512, &tmp)
	if !cmp512(&res.v, &tmp.v) {
		t.Errorf("Wrong value\n%X", res.v)
	}

	tmp = Fp{v: u512{1, 1, 1, 1, 1, 1, 1, 1}}
	exp := Fp{v: u512{2, 2, 2, 2, 2, 2, 2, 2}}
	addRdc(&res, &tmp, &tmp)
	if !cmp512(&res.v, &exp.v) {
		t.Errorf("Wrong value\n%X", res.v)
	}
}

func TestSubRdc(t *testing.T) {
	var res Fp
	var Fp512 = Fp{v: p}

	// 1 - 1 mod P
	tmp := OneFp512
	subRdc(&res, &tmp, &tmp)
	if !cmp512(&res.v, &ZeroFp512.v) {
		t.Errorf("Wrong value\n%X", res.v)
	}
	zero(&res.v)

	// 0 - 1 mod P
	exp := Fp512
	exp.v[0]--

	subRdc(&res, &ZeroFp512, &OneFp512)
	if !cmp512(&res.v, &exp.v) {
		t.Errorf("Wrong value\n%X\n%X", res.v, exp.v)
	}
	zero(&res.v)

	// P - (P-1)
	pMinusOne := Fp512
	pMinusOne.v[0]--
	subRdc(&res, &Fp512, &pMinusOne)
	if !cmp512(&res.v, &OneFp512.v) {
		t.Errorf("Wrong value\n[%X != %X]", res.v, OneFp512.v)
	}
	zero(&res.v)

	subRdc(&res, &Fp512, &OneFp512)
	if !cmp512(&res.v, &pMinusOne.v) {
		t.Errorf("Wrong value\n[%X != %X]", res.v, pMinusOne.v)
	}
}

func BenchmarkFp512Add(b *testing.B) {
	var arg1 u512
	arg2 := randomU512()
	arg3 := randomU512()
	for n := 0; n < b.N; n++ {
		add512(&arg1, &arg2, &arg3)
	}
}

func BenchmarkFp512Sub(b *testing.B) {
	var arg1 u512
	arg2, arg3 := randomU512(), randomU512()
	for n := 0; n < b.N; n++ {
		sub512(&arg1, &arg2, &arg3)
	}
}

func BenchmarkFp512Mul(b *testing.B) {
	var arg1 = mrand.Uint64()
	arg2, arg3 := randomU512(), randomU512()
	for n := 0; n < b.N; n++ {
		mul512(&arg2, &arg3, arg1)
	}
}

func BenchmarkCSwap(b *testing.B) {
	arg1 := randomU512()
	arg2 := randomU512()
	for n := 0; n < b.N; n++ {
		cswap512(&arg1, &arg2, uint8(n%2))
	}
}

func BenchmarkAddRdc(b *testing.B) {
	arg1 := Fp{v: randomU512()}
	arg2 := Fp{v: randomU512()}
	var res Fp
	for n := 0; n < b.N; n++ {
		addRdc(&res, &arg1, &arg2)
	}
}

func BenchmarkSubRdc(b *testing.B) {
	arg1 := Fp{v: randomU512()}
	arg2 := Fp{v: randomU512()}
	var res Fp
	for n := 0; n < b.N; n++ {
		subRdc(&res, &arg1, &arg2)
	}
}
