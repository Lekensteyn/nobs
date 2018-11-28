package csidh

import (
	"math/big"
	"testing"
	//        "fmt"
	//        mrand "math/rand"
)

// Converst src to big.Int. Function assumes that src is a slice of uint64
// values encoded in little-endian byte order.
func intSetU64(dst *big.Int, src []uint64) *big.Int {
	var tmp big.Int

	dst.SetUint64(0)
	for i, _ := range src {
		tmp.SetUint64(src[i])
		tmp.Lsh(&tmp, uint(i*64))
		dst.Add(dst, &tmp)
	}
	return dst
}

// Convers src to an array of uint64 values encoded in little-endian
// byte order.
func intGetU64(src *big.Int) []uint64 {
	var tmp, mod big.Int
	dst := make([]uint64, (src.BitLen()/64)+1)

	u64 := uint64(0)
	u64--
	mod.SetUint64(u64)
	for i := 0; i < (src.BitLen()/64)+1; i++ {
		tmp.Set(src)
		tmp.Rsh(&tmp, uint(i)*64)
		tmp.And(&tmp, &mod)
		dst[i] = tmp.Uint64()
	}
	return dst
}

// Converts dst to Montgomery domain of cSIDH-512
func toMont(dst *big.Int, toMont bool) *big.Int {
	var bigP, bigR big.Int

	intSetU64(&bigP, p[:])
	bigR.SetUint64(1)
	bigR.Lsh(&bigR, 512)

	if !toMont {
		bigR.ModInverse(&bigR, &bigP)
	}
	dst.Mul(dst, &bigR)
	dst.Mod(dst, &bigP)
	return dst
}

// Returns projective coordinate X of normalized EC 'point' (point.x / point.z).
func toNormX(point *Point) big.Int {
	var bigP, bigDnt, bigDor big.Int

	intSetU64(&bigP, p[:])
	intSetU64(&bigDnt, point.x[:])
	intSetU64(&bigDor, point.z[:])

	bigDor.ModInverse(&bigDor, &bigP)
	bigDnt.Mul(&bigDnt, &bigDor)
	bigDnt.Mod(&bigDnt, &bigP)
	return bigDnt
}

// Converts string to Fp element in Montgomery domain of cSIDH-512
func toFp(num string) Fp {
	var tmp big.Int
	var ok bool
	var ret Fp

	_, ok = tmp.SetString(num, 0)
	if !ok {
		panic("Can't parse a number")
	}
	toMont(&tmp, true)
	copy(ret[:], intGetU64(&tmp))
	return ret
}

// Actual test implementation

func TestXAdd(t *testing.T) {
	var P, Q, D Point
	var sumPQ Point
	var sumPQExp big.Int

	// Points from a Elliptic Curve defined in sage as follows:
	// A = 5045436521140567532715475549653557926015565274495421515721690830514805087260227651930593135277458179466002951380639083992282549121114024513660204218680594
	// E = EllipticCurve(GF(p), [0, A, 0, 1, 0])
	// where p is CSIDH's 511-bit prime

	sumPQExp.SetString("0x41C98C5D7FF118B1A3987733581FD69C0CC27D7B63BCCA525106B9945869C6DAEDAA3D5D9D2679237EF0D013BE68EF12731DBFB26E12576BAD1E824C67ABD125", 0)
	P.x = toFp("0x5840FD8E0165F7F474260F99337461AF195233F791FABE735EC2634B74A95559568B4CEB23959C8A01C5C57E215D22639868ED840D74FE2BAC04830CF75047AD")
	P.z = toFp("1")
	Q.x = toFp("0x3C1A003C71436698B4A181CEB12BA4B4D1FF7BB14AAAF6FBDA6957C4EBA20AD8E3893DF6F64E67E81163E024C19C7E975F3EC61862F75502C3ED802370E75A3F")
	Q.z = toFp("1")
	D.x = toFp("0x519B1928F752B0B2143C1C23EB247B370DBB5B9C29B9A3A064D7FBC1B67FAC34B6D3DDA0F3CB87C387B425B36F31B93A8E73252BA701927B767A9DE89D5A92AE")
	D.z = toFp("1")

	xAdd(&sumPQ, &P, &Q, &D)
	ret := toNormX(&sumPQ)
	if ret.Cmp(&sumPQExp) != 0 {
		t.Fail()
	}
}
