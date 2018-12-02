package csidh

import (
	// "fmt"
	"math/big"
	"testing"
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
	var P, Q, PdQ Point
	var PaQ Point
	var expPaQ big.Int

	// Points from a Elliptic Curve defined in sage as follows:
	// A = 0x6055947AAFEBF773CE912680A6A32656073233D2FD6FDF4A143BE82D25B44ECC0431DE564C0F0D6591ACC62D6876E86F5D06B68C9EAF20D0DB0A6B99ED558512
	// E = EllipticCurve(GF(p), [0, A, 0, 1, 0])
	// where p is CSIDH's 511-bit prime

	expPaQ.SetString("0x41C98C5D7FF118B1A3987733581FD69C0CC27D7B63BCCA525106B9945869C6DAEDAA3D5D9D2679237EF0D013BE68EF12731DBFB26E12576BAD1E824C67ABD125", 0)
	P.x = toFp("0x5840FD8E0165F7F474260F99337461AF195233F791FABE735EC2634B74A95559568B4CEB23959C8A01C5C57E215D22639868ED840D74FE2BAC04830CF75047AD")
	P.z = toFp("1")
	Q.x = toFp("0x3C1A003C71436698B4A181CEB12BA4B4D1FF7BB14AAAF6FBDA6957C4EBA20AD8E3893DF6F64E67E81163E024C19C7E975F3EC61862F75502C3ED802370E75A3F")
	Q.z = toFp("1")
	PdQ.x = toFp("0x519B1928F752B0B2143C1C23EB247B370DBB5B9C29B9A3A064D7FBC1B67FAC34B6D3DDA0F3CB87C387B425B36F31B93A8E73252BA701927B767A9DE89D5A92AE")
	PdQ.z = toFp("1")

	xAdd(&PaQ, &P, &Q, &PdQ)
	ret := toNormX(&PaQ)
	if ret.Cmp(&expPaQ) != 0 {
		t.Errorf("\nExp: %s\nGot: %s", expPaQ.Text(16), ret.Text(16))
	}
}

func TestXDbl(t *testing.T) {
	var P, A Point
	var PaP Point
	var expPaP big.Int

	// Points from a Elliptic Curve defined in sage as follows:
	// A = 0x599841D7D1FCD92A85759B7A3D2D5E4C56EFB17F19F86EB70E121EA16305EDE45A55868BE069313F821F7D94069EC220A4AC3B85500376710538246E9B3BC138
	// E = EllipticCurve(GF(p), [0, A, 0, 1, 0])
	// where p is CSIDH's 511-bit prime

	expPaP.SetString("0x6115B5D8BB613D11BDFEA70D436D87C1515553F6A15061727B4001E0AF745AAA9F39EB9464982829D931F77DAB9D71B24FF0D1D34C347F2A51FD45821F2EA06F", 0)
	P.x = toFp("0x6C5B4D4AB0765AAB23C10F8455BE522D3A5363324D7AD641CC67C0A52FC1FFE9F3F8EDFE641478CA93D4D0016D83F21487FD4AF4E02F8A2C237CF27C5604BCC")
	P.z = toFp("1")
	A.x = toFp("0x599841D7D1FCD92A85759B7A3D2D5E4C56EFB17F19F86EB70E121EA16305EDE45A55868BE069313F821F7D94069EC220A4AC3B85500376710538246E9B3BC138")
	A.z = toFp("1")

	xDbl(&PaP, &P, &A)
	ret := toNormX(&PaP)
	if ret.Cmp(&expPaP) != 0 {
		t.Errorf("\nExp: %s\nGot: %s", expPaP.Text(16), ret.Text(16))
	}
}

// TODO: test C!=1
func TestXDblAdd(t *testing.T) {
	var P, Q, PdQ Point
	var PaP, PaQ Point
	var expPaP, expPaQ big.Int
	var A, A24 Coeff

	// 2*P
	expPaP.SetString("0x38F5B37271A3D8FA50107F88045D6F6B08355DD026C02E0306CE5875F47422736AD841B4122B2BD7DE6166BB6498F6A283378FF8250948E834F15CEA2D59A57B", 0)
	// P+Q
	expPaQ.SetString("0x53D9B44C5F61651612243CF7987F619FE6ACB5CF29538F96A63E7278E131F41A17D64388E31B028A5183EF9096AE82724BC34D8DDFD67AD68BD552A33C345B8C", 0)
	P.x = toFp("0x4FE17B4CC66E85960F57033CD45996C99248DA09DF2E36F8840657B52F74ED8173E0D322FA57D7B4D0EE7F12967BBD59140B42F2626E29167D6419E851E5A4C9")
	P.z = toFp("1")
	Q.x = toFp("0x465047949CD6574FDBE00EA365CAF7A95DC9DEBE96A188823CA8C9DD9F527CF81290D49864F61DF0C08C1D6052139230735CA6CFDBDC1A8820610CCD71861176")
	Q.z = toFp("1")
	PdQ.x = toFp("0x49D3B999A0A020B34473568A8F75B5405F2D3BE5A006595015FC6DDC6BED8AB2A51A887B6DC62C64354466865FFD69E50AD37F6F4FBD74119EB65EBC9367B556")
	PdQ.z = toFp("1")
	A.a = toFp("0x118F955D498D902FD42E5B2926F297CC814CD7649EC5B070295622F97C4A0D9BD34058A7E0E00CB73ED32FCC237F9F6B7D2A15F5CC7C4EC61ECEF80ACBB0EFA4")
	A.c = toFp("1")

	// A24.a = 2*A.z + A.a
	addRdc(&A24.a, &A.c, &A.c)
	addRdc(&A24.a, &A24.a, &A.a)
	// A24.z = 4*A.z
	mulRdc(&A24.c, &A.c, &four)

	xDblAdd(&PaP, &PaQ, &P, &Q, &PdQ, &A24)
	retPaP := toNormX(&PaP)
	retPaQ := toNormX(&PaQ)
	if retPaP.Cmp(&expPaP) != 0 {
		t.Errorf("\nExp: %s\nGot: %s", expPaP.Text(16), retPaP.Text(16))
	}

	if retPaQ.Cmp(&expPaQ) != 0 {
		t.Errorf("\nExp: %s\nGot: %s", expPaQ.Text(16), retPaQ.Text(16))
	}
}

// TODO: test C!=1
func TestXMul(t *testing.T) {
	var kP, P Point
	var co Coeff
	var expKP big.Int
	var k Fp

	// Case C=1
	expKP.SetString("0x582B866603E6FBEBD21FE660FB34EF9466FDEC55FFBCE1073134CC557071147821BBAD225E30F7B2B6790B00ED9C39A29AA043F58AF995E440AFB13DA8E6D788", 0)
	P.x = toFp("0x1C5CA539C1D5B52DE4750C390C24C05251E8B1D33E48971FA86F5ADDED2D06C8CD31E94887541468BB2925EBD693C9DDFF5BD9508430F25FE28EE30C0760C0FE")
	P.z = toFp("1")
	co.a = toFp("0x538F785D52996919C8D5C73D842A0249669B5B6BB05338B74EAE8094AE5009A3BA2D73730F527D7403E8184D9B1FA11C0C4C40E7B328A84874A6DBCE99E1DF92")
	co.c = toFp("1")
	k = Fp{0x7A36C930A83EFBD5, 0xD0E80041ED0DDF9F, 0x5AA17134F1B8F877, 0x975711EC94168E51, 0xB3CAD962BED4BAC5, 0x3026DFDD7E4F5687, 0xE67F91AB8EC9C3AF, 0x34671D3FD8C317E7}

	xMul512(&kP, &P, &co, &k)
	retKP := toNormX(&kP)
	if expKP.Cmp(&retKP) != 0 {
		t.Errorf("\nExp: %s\nGot: %s", expKP.Text(16), retKP.Text(16))
	}
}

func BenchmarkXMul(b *testing.B) {
	var kP, P Point
	var co Coeff
	var expKP big.Int
	var k Fp

	// Case C=1
	expKP.SetString("0x582B866603E6FBEBD21FE660FB34EF9466FDEC55FFBCE1073134CC557071147821BBAD225E30F7B2B6790B00ED9C39A29AA043F58AF995E440AFB13DA8E6D788", 0)
	P.x = toFp("0x1C5CA539C1D5B52DE4750C390C24C05251E8B1D33E48971FA86F5ADDED2D06C8CD31E94887541468BB2925EBD693C9DDFF5BD9508430F25FE28EE30C0760C0FE")
	P.z = toFp("1")
	co.a = toFp("0x538F785D52996919C8D5C73D842A0249669B5B6BB05338B74EAE8094AE5009A3BA2D73730F527D7403E8184D9B1FA11C0C4C40E7B328A84874A6DBCE99E1DF92")
	co.c = toFp("1")
	k = Fp{0x7A36C930A83EFBD5, 0xD0E80041ED0DDF9F, 0x5AA17134F1B8F877, 0x975711EC94168E51, 0xB3CAD962BED4BAC5, 0x3026DFDD7E4F5687, 0xE67F91AB8EC9C3AF, 0x34671D3FD8C317E7}

	for n := 0; n < b.N; n++ {
		xMul512(&kP, &P, &co, &k)
	}
}
