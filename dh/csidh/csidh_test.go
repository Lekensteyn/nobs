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

func TestPublicKeyGen(t *testing.T) {
	prv_bytes := []byte{0xdb, 0x54, 0xe4, 0xd4, 0xd0, 0xbd, 0xee, 0xcb, 0xf4, 0xd0, 0xc2, 0xbc, 0x52, 0x44, 0x11, 0xee, 0xe1, 0x14, 0xd2, 0x24, 0xe5, 0x0, 0xcc, 0xf5, 0xc0, 0xe1, 0x1e, 0xb3, 0x43, 0x52, 0x45, 0xbe, 0xfb, 0x54, 0xc0, 0x55, 0xb2}
	//pub_bytes := []uint64{0x6BCAAD7EFD426976,0x743D780A06D2CDC5,0x841A2D76984849F7,0x1523EB45B3B78D5F,0xCF7A093C773EDF8D,0xFAB0FF04A7B4A54D,0x05DE322C864069D2,0x0C55DC69711DF47A}

	var prv PrivateKey
	var pub PublicKey
	prv.Import(prv_bytes)
	pub.Generate(&prv)
	fp_print(pub.A)
	//print(pub.Export())
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
