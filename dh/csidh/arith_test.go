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

// x<y: <0
// x>y: >0
// x==y: 0
func cmp512(x, y *u512) int {
	if len(*x) == len(*y) {
		for i := len(*x) - 1; i >= 0; i-- {
			if x[i] < y[i] {
				return -1
			} else if x[i] > y[i] {
				return 1
			}
		}
		return 0
	}
	return len(*x) - len(*y)
}

// return x==y
func ceq512(x, y *u512) bool {
	return cmp512(x, y) == 0
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
	if !ceq512(&arg1, &arg1cpy) {
		t.Error("cswap swapped")
	}

	arg1cpy = arg1
	cswap512(&arg1, &arg2, 1)
	if ceq512(&arg1, &arg1cpy) {
		t.Error("cswap didn't swapped")
	}

	arg1cpy = arg1
	cswap512(&arg1, &arg2, 0xF2)
	if ceq512(&arg1, &arg1cpy) {
		t.Error("cswap didn't swapped")
	}
}

func TestAddRdc(t *testing.T) {
	var res Fp
	var Fp512 = Fp{v: p}

	tmp := OneFp512
	addRdc(&res, &tmp, &Fp512)
	if !ceq512(&res.v, &tmp.v) {
		t.Errorf("Wrong value\n%X", res.v)
	}

	tmp = ZeroFp512
	addRdc(&res, &Fp512, &Fp512)
	if !ceq512(&res.v, &Fp512.v) {
		t.Errorf("Wrong value\n%X", res.v)
	}

	tmp = Fp{v: u512{1, 1, 1, 1, 1, 1, 1, 1}}
	addRdc(&res, &Fp512, &tmp)
	if !ceq512(&res.v, &tmp.v) {
		t.Errorf("Wrong value\n%X", res.v)
	}

	tmp = Fp{v: u512{1, 1, 1, 1, 1, 1, 1, 1}}
	exp := Fp{v: u512{2, 2, 2, 2, 2, 2, 2, 2}}
	addRdc(&res, &tmp, &tmp)
	if !ceq512(&res.v, &exp.v) {
		t.Errorf("Wrong value\n%X", res.v)
	}
}

func TestSubRdc(t *testing.T) {
	var res Fp
	var Fp512 = Fp{v: p}

	// 1 - 1 mod P
	tmp := OneFp512
	subRdc(&res, &tmp, &tmp)
	if !ceq512(&res.v, &ZeroFp512.v) {
		t.Errorf("Wrong value\n%X", res.v)
	}
	zero(&res.v)

	// 0 - 1 mod P
	exp := Fp512
	exp.v[0]--

	subRdc(&res, &ZeroFp512, &OneFp512)
	if !ceq512(&res.v, &exp.v) {
		t.Errorf("Wrong value\n%X\n%X", res.v, exp.v)
	}
	zero(&res.v)

	// P - (P-1)
	pMinusOne := Fp512
	pMinusOne.v[0]--
	subRdc(&res, &Fp512, &pMinusOne)
	if !ceq512(&res.v, &OneFp512.v) {
		t.Errorf("Wrong value\n[%X != %X]", res.v, OneFp512.v)
	}
	zero(&res.v)

	subRdc(&res, &Fp512, &OneFp512)
	if !ceq512(&res.v, &pMinusOne.v) {
		t.Errorf("Wrong value\n[%X != %X]", res.v, pMinusOne.v)
	}
}

func TestMulRdc(t *testing.T) {
	var res Fp
	var fp1 = Fp{v: fp_1}
	var m1 = Fp{v: u512{
		0x85E2579C786882D0, 0x4E3433657E18DA95,
		0x850AE5507965A0B3, 0xA15BC4E676475964}}
	var m2 = Fp{v: u512{
		0x85E2579C786882CF, 0x4E3433657E18DA95,
		0x850AE5507965A0B3, 0xA15BC4E676475964}}

	// Expected
	var m1m1 = u512{
		0xAEBF46E92C88A4B4, 0xCFE857977B946347,
		0xD3B264FF08493901, 0x6EEB3D23746B6C7C,
		0xC0CA874A349D64B4, 0x7AD4A38B406F8504,
		0x38B6B6CEB82472FB, 0x1587015FD7DDFC7D}
	var m1m2 = u512{
		0x51534771258C4624, 0x2BFEDE86504E2160,
		0xE8127D5E9329670B, 0x0C84DBD584491D75,
		0x656C73C68B16E38C, 0x01C0DA470B30B8DE,
		0x2532E3903EAA950B, 0x3F2C28EA97FE6FEC}

	// 0*0
	tmp := ZeroFp512
	mulRdc(&res, &tmp, &tmp)
	if !ceq512(&res.v, &tmp.v) {
		t.Errorf("Wrong value\n%X", res.v)
	}

	// 1*m1 == m1
	zero(&res.v)
	mulRdc(&res, &m1, &fp1)
	if !ceq512(&res.v, &m1.v) {
		t.Errorf("Wrong value\n%X", res.v)
	}

	// OZAPTF: I don't understand those results. but they are correct
	// m1*m2 < p
	zero(&res.v)
	mulRdc(&res, &m1, &m2)
	if !ceq512(&res.v, &m1m2) {
		t.Errorf("Wrong value\n%X", res.v)
	}

	// m1*m1 > p
	zero(&res.v)
	mulRdc(&res, &m1, &m1)
	if !ceq512(&res.v, &m1m1) {
		t.Errorf("Wrong value\n%X", res.v)
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
