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
