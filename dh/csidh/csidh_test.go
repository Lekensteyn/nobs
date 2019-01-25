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
	for i := 0; i < kNumIter; i++ {
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

func TestEphemeralKeyExchange(t *testing.T) {
	var buf [64]uint8
	var ss [64]uint8

	prv_bytes1 := []byte{0xaa, 0x54, 0xe4, 0xd4, 0xd0, 0xbd, 0xee, 0xcb, 0xf4, 0xd0, 0xc2, 0xbc, 0x52, 0x44, 0x11, 0xee, 0xe1, 0x14, 0xd2, 0x24, 0xe5, 0x0, 0xcc, 0xf5, 0xc0, 0xe1, 0x1e, 0xb3, 0x43, 0x52, 0x45, 0xbe, 0xfb, 0x54, 0xc0, 0x55, 0xb2}
	prv_bytes2 := []byte{0xbb, 0x54, 0xe4, 0xd4, 0xd0, 0xbd, 0xee, 0xcb, 0xf4, 0xd0, 0xc2, 0xbc, 0x52, 0x44, 0x11, 0xee, 0xe1, 0x14, 0xd2, 0x24, 0xe5, 0x0, 0xcc, 0xf5, 0xc0, 0xe1, 0x1e, 0xb3, 0x43, 0x52, 0x45, 0xbe, 0xfb, 0x54, 0xc0, 0x55, 0xb2}
	var prv1, prv2 PrivateKey
	var pub1, pub2 PublicKey
	prv1.Import(prv_bytes1)
	pub1.Generate(&prv1)
	pub1.Export(buf[:])

	prv2.Import(prv_bytes2)
	pub2.Generate(&prv2)
	pub2.Export(buf[:])

	pub1.DeriveSecret(ss[:], &pub2, &prv1)
}

func TestPrivateKeyExportImport(t *testing.T) {
	var buf [37]uint8
	for i := 0; i < 100; i++ {
		var prv1, prv2 PrivateKey
		prv1.Generate(crand.Reader)
		prv1.Export(buf[:])
		prv2.Import(buf[:])

		for i := 0; i < len(prv1.e); i++ {
			if prv1.e[i] != prv2.e[i] {
				t.Error("Error occured when public key export/import")
			}
		}
	}
}

func TestPublicKeyExportImport(t *testing.T) {
	var buf [64]uint8
	for i := 0; i < 100; i++ {
		var prv PrivateKey
		var pub1, pub2 PublicKey
		prv.Generate(crand.Reader)
		pub1.Generate(&prv)

		pub1.Export(buf[:])
		pub2.Import(buf[:])

		if eq64(pub1.A[:], pub2.A[:]) != 1 {
			t.Error("Error occured when public key export/import")
		}
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
	var prv PrivateKey
	for n := 0; n < b.N; n++ {
		prv.Generate(crand.Reader)
		pub.Generate(&prv)
		pub.validate()
	}
}

func BenchmarkEphemeralKeyExchange(b *testing.B) {
	var ss [64]uint8
	var prv1, prv2 PrivateKey
	var pub1, pub2 PublicKey
	for n := 0; n < b.N; n++ {
		prv1.Generate(crand.Reader)
		pub1.Generate(&prv1)

		prv2.Generate(crand.Reader)
		pub2.Generate(&prv2)

		pub1.DeriveSecret(ss[:], &pub2, &prv1)
	}
}
