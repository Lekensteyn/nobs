package csidh

import (
	crand "crypto/rand"
	mrand "math/rand"
	"testing"
)

func eq64(x, y []uint64) uint {
	for i, _ := range x {
		if x[i] != y[i] {
			return 0
		}
	}
	return 1
}

func TestCtEq64(t *testing.T) {
	var t1, t2 [8]uint64
	for i := 0; i < 100000; i++ {
		for i, _ := range t1 {
			t1[i] = mrand.Uint64()
			t2[i] = mrand.Uint64()
		}

		if ctEq64(t1[:], t2[:]) != eq64(t1[:], t2[:]) {
			t.FailNow()
		}
	}

	var t3 = [8]uint64{1, 2, 3, 4, 5, 6, 7, 8}
	var t4 = [8]uint64{1, 2, 3, 4, 5, 6, 7, 8}
	if ctEq64(t3[:], t4[:]) != eq64(t3[:], t4[:]) {
		t.FailNow()
	}
}

func BenchmarkGeneratePrivate(b *testing.B) {
	var prv PrivateKey
	for n := 0; n < b.N; n++ {
		prv.Generate(crand.Reader)
	}
}

func BenchmarkValidate(b *testing.B) {
	var pub PublicKey
	for n := 0; n < b.N; n++ {
		pub.Validate()
	}
}
