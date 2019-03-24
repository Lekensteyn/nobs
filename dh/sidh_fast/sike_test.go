package sike

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"
)

var tdata = map[uint8]struct {
	name     string
	PrB_sidh string
	PkB_sidh string
	PkB_sike string
	PrB_sike string
	PrA_sike string
	PkA_sike string
}{
	FP_503: {
		name:     "P-503",
		PkB_sike: "68460C22466E95864CFEA7B5D9077E768FF4F9ED69AE56D7CF3F236FB06B31020EEE34B5B572CEA5DDF20B531966AA8F5F3ACC0C6D1CE04EEDC30FD1F1233E2D96FE60C6D638FC646EAF2E2246F1AEC96859CE874A1F029A78F9C978CD6B22114A0D5AB20101191FD923E80C76908B1498B9D0200065CCA09159A0C65A1E346CC6470314FE78388DAA89DD08EC67DBE63C1F606674ACC49EBF9FDBB2B898B3CE733113AA6F942DB401A76D629CE6EE6C0FDAF4CFB1A5E366DB66C17B3923A1B7FB26A3FF25B9018869C674D3DEF4AF269901D686FE4647F9D2CDB2CEB3AFA305B27C885F037ED167F595066C21E7DD467D8332B934A5102DA5F13332DFA356B82156A0BB2E7E91C6B85B7D1E381BC9E3F0FC4DB9C36016D9ECEC415D7E977E9AC29910D934BA2FE4EE49D3B387607A4E1AFABF495FB86A77194626589E802FF5167C7A25C542C1EAD25A6E0AA931D94F2F9AFD3DBDF222E651F729A90E77B20974905F1E65E041CE6C95AAB3E1F22D332E0A5DE9C5DB3D9C7A38",
		PrB_sike: "80FC55DA74DEFE3113487B80841E678AF9ED4E0599CF07353A4AB93971C090A0" +
			"A9402C9DC98AC6DC8F5FDE5E970AE22BA48A400EFC72851C",
		PrB_sidh: "A885A8B889520A6DBAD9FB33365E5B77FDED629440A16A533F259A510F63A822",
		PrA_sike: "B0AD510708F4ABCF3E0D97DC2F2FF112D9D2AAE49D97FFD1E4267F21C6E71C03",
		PkA_sike: "A6BADBA04518A924B20046B59AC197DCDF0EA48014C9E228C4994CCA432F360E" +
			"2D527AFB06CA7C96EE5CEE19BAD53BF9218A3961CAD7EC092BD8D9EBB22A3D51" +
			"33008895A3F1F6A023F91E0FE06A00A622FD6335DAC107F8EC4283DC2632F080" +
			"4E64B390DAD8A2572F1947C67FDF4F8787D140CE2C6B24E752DA9A195040EDFA" +
			"C27333FAE97DBDEB41DA9EEB2DB067AE7DA8C58C0EF57AEFC18A3D6BD0576FF2" +
			"F1CFCAEC50C958331BF631F3D2E769790C7B6DF282B74BBC02998AD10F291D47" +
			"C5A762FF84253D3B3278BDF20C8D4D4AA317BE401B884E26A1F02C7308AADB68" +
			"20EBDB0D339F5A63346F3B40CACED72F544DAF51566C6E807D0E6E1E38514342" +
			"432661DC9564DA07548570E256688CD9E8060D8775F95D501886D958588CACA0" +
			"9F2D2AE1913F996E76AF63E31A179A7A7D2A46EDA03B2BCCF9020A5AA15F9A28" +
			"9340B33F3AE7F97360D45F8AE1B9DD48779A57E8C45B50A02C00349CD1C58C55" +
			"1D68BC2A75EAFED944E8C599C288037181E997471352E24C952B",
		PkB_sidh: "244AF1F367C2C33912750A98497CC8214BC195BD52BD76513D32ACE4B75E31F0" +
			"281755C265F5565C74E3C04182B9C244071859C8588CC7F09547CEFF8F7705D2" +
			"60CE87D6BFF914EE7DBE4B9AF051CA420062EEBDF043AF58184495026949B068" +
			"98A47046BFAE8DF3B447746184AF550553BB5D266D6E1967ACA33CAC5F399F90" +
			"360D70867F2C71EF6F94FF915C7DA8BC9549FB7656E691DAEFC93CF56876E482" +
			"CA2F8BE2D6CDCC374C31AD8833CABE997CC92305F38497BEC4DFD1821B004FEC" +
			"E16448F9A24F965EFE409A8939EEA671633D9FFCF961283E59B8834BDF7EDDB3" +
			"05D6275B61DA6692325432A0BAA074FC7C1F51E76208AB193A57520D40A76334" +
			"EE5712BDC3E1EFB6103966F2329EDFF63082C4DFCDF6BE1C5A048630B81871B8" +
			"83B735748A8FD4E2D9530C272163AB18105B10015CA7456202FE1C9B92CEB167" +
			"5EAE1132E582C88E47ED87B363D45F05BEA714D5E9933D7AF4071CBB5D49008F" +
			"3E3DAD7DFF935EE509D5DE561842B678CCEB133D62E270E9AC3E",
	},
}

/* -------------------------------------------------------------------------
   Helpers
   -------------------------------------------------------------------------*/
// Fail if err !=nil. Display msg as an error message
func checkErr(t testing.TB, err error, msg string) {
	if err != nil {
		t.Error(msg)
	}
}

// Utility used for running same test with all registered prime fields
type MultiIdTestingFunc func(testing.TB, uint8)

func Do(f MultiIdTestingFunc, t testing.TB) {
	for id, val := range tdata {
		fmt.Printf("\tTesting: %s\n", val.name)
		f(t, id)
	}
}

// Converts string to private key
func convToPrv(s string, v KeyVariant, id uint8) *PrivateKey {
	key := NewPrivateKey(id, v)
	hex, e := hex.DecodeString(s)
	if e != nil {
		panic("non-hex number provided")
	}
	e = key.Import(hex)
	if e != nil {
		panic("Can't import private key")
	}
	return key
}

// Converts string to public key
func convToPub(s string, v KeyVariant, id uint8) *PublicKey {
	key := NewPublicKey(id, v)
	hex, e := hex.DecodeString(s)
	if e != nil {
		panic("non-hex number provided")
	}
	e = key.Import(hex)
	if e != nil {
		panic("Can't import public key")
	}
	return key
}

/* -------------------------------------------------------------------------
   Unit tests
   -------------------------------------------------------------------------*/
func testKeygen(t testing.TB, id uint8) {
	alicePrivate := convToPrv(tdata[id].PrA_sike, KeyVariant_SIDH_A, id)
	bobPrivate := convToPrv(tdata[id].PrB_sidh, KeyVariant_SIDH_B, id)
	expPubA := convToPub(tdata[id].PkA_sike, KeyVariant_SIDH_A, id)
	expPubB := convToPub(tdata[id].PkB_sidh, KeyVariant_SIDH_B, id)

	pubA := alicePrivate.GeneratePublicKey()
	pubB := bobPrivate.GeneratePublicKey()

	if !bytes.Equal(pubA.Export(), expPubA.Export()) {
		t.Fatalf("unexpected value of public key A")
	}
	if !bytes.Equal(pubB.Export(), expPubB.Export()) {
		t.Fatalf("unexpected value of public key B")
	}
}

func testImportExport(t testing.TB, id uint8) {
	var err error
	a := NewPublicKey(id, KeyVariant_SIDH_A)
	b := NewPublicKey(id, KeyVariant_SIDH_B)

	// Import keys
	a_hex, err := hex.DecodeString(tdata[id].PkA_sike)
	checkErr(t, err, "invalid hex-number provided")

	err = a.Import(a_hex)
	checkErr(t, err, "import failed")

	b_hex, err := hex.DecodeString(tdata[id].PkB_sike)
	checkErr(t, err, "invalid hex-number provided")

	err = b.Import(b_hex)
	checkErr(t, err, "import failed")

	// Export and check if same
	if !bytes.Equal(b.Export(), b_hex) || !bytes.Equal(a.Export(), a_hex) {
		t.Fatalf("export/import failed")
	}

	if (len(b.Export()) != b.Size()) || (len(a.Export()) != a.Size()) {
		t.Fatalf("wrong size of exported keys")
	}
}

func testPrivateKeyBelowMax(t testing.TB, id uint8) {
	params := Params(id)
	for variant, keySz := range map[KeyVariant]*DomainParams{
		KeyVariant_SIDH_A: &params.A,
		KeyVariant_SIDH_B: &params.B} {

		func(v KeyVariant, dp *DomainParams) {
			var blen = int(dp.SecretByteLen)
			var prv = NewPrivateKey(id, v)

			// Calculate either (2^e2 - 1) or (2^s - 1); where s=ceil(log_2(3^e3)))
			maxSecertVal := big.NewInt(int64(dp.SecretBitLen))
			maxSecertVal.Exp(big.NewInt(int64(2)), maxSecertVal, nil)
			maxSecertVal.Sub(maxSecertVal, big.NewInt(1))

			// Do same test 1000 times
			for i := 0; i < 1000; i++ {
				err := prv.Generate(rand.Reader)
				checkErr(t, err, "Private key generation")

				// Convert to big-endian, as that's what expected by (*Int)SetBytes()
				secretBytes := prv.Export()
				for i := 0; i < int(blen/2); i++ {
					tmp := secretBytes[i] ^ secretBytes[blen-i-1]
					secretBytes[i] = tmp ^ secretBytes[i]
					secretBytes[blen-i-1] = tmp ^ secretBytes[blen-i-1]
				}
				prvBig := new(big.Int).SetBytes(secretBytes)
				// Check if generated key is bigger then acceptable
				if prvBig.Cmp(maxSecertVal) == 1 {
					t.Error("Generated private key is wrong")
				}
			}
		}(variant, keySz)
	}
}

func testKeyAgreement(t testing.TB, id uint8, pkA, prA, pkB, prB string) {
	var e error

	// KeyPairs
	alicePublic := convToPub(pkA, KeyVariant_SIDH_A, id)
	bobPublic := convToPub(pkB, KeyVariant_SIDH_B, id)
	alicePrivate := convToPrv(prA, KeyVariant_SIDH_A, id)
	bobPrivate := convToPrv(prB, KeyVariant_SIDH_B, id)

	// Do actual test
	s1, e := DeriveSecret(bobPrivate, alicePublic)
	checkErr(t, e, "derivation s1")
	s2, e := DeriveSecret(alicePrivate, bobPublic)
	checkErr(t, e, "derivation s1")

	if !bytes.Equal(s1[:], s2[:]) {
		t.Fatalf("two shared keys: %d, %d do not match", s1, s2)
	}

	// Negative case
	dec, e := hex.DecodeString(tdata[id].PkA_sike)
	if e != nil {
		t.FailNow()
	}
	dec[0] = ^dec[0]
	e = alicePublic.Import(dec)
	if e != nil {
		t.FailNow()
	}

	s1, e = DeriveSecret(bobPrivate, alicePublic)
	checkErr(t, e, "derivation of s1 failed")
	s2, e = DeriveSecret(alicePrivate, bobPublic)
	checkErr(t, e, "derivation of s2 failed")

	if bytes.Equal(s1[:], s2[:]) {
		t.Fatalf("The two shared keys: %d, %d match", s1, s2)
	}
}

func testDerivationRoundTrip(t testing.TB, id uint8) {
	var err error

	prvA := NewPrivateKey(id, KeyVariant_SIDH_A)
	prvB := NewPrivateKey(id, KeyVariant_SIDH_B)

	// Generate private keys
	err = prvA.Generate(rand.Reader)
	checkErr(t, err, "key generation failed")
	err = prvB.Generate(rand.Reader)
	checkErr(t, err, "key generation failed")

	// Generate public keys
	pubA := prvA.GeneratePublicKey()
	pubB := prvB.GeneratePublicKey()

	// Derive shared secret
	s1, err := DeriveSecret(prvB, pubA)
	checkErr(t, err, "")

	s2, err := DeriveSecret(prvA, pubB)
	checkErr(t, err, "")

	if !bytes.Equal(s1[:], s2[:]) {
		t.Fatalf("Tthe two shared keys: \n%X, \n%X do not match", s1, s2)
	}
}

// Encrypt, Decrypt, check if input/output plaintext is the same
func testPKERoundTrip(t testing.TB, id uint8) {
	// Message to be encrypted
	var params = Params(id)
	var msg = make([]byte, params.MsgLen)
	for i, _ := range msg {
		msg[i] = byte(i)
	}

	// Import keys
	pkB := NewPublicKey(params.Id, KeyVariant_SIKE)
	skB := NewPrivateKey(params.Id, KeyVariant_SIKE)
	pk_hex, err := hex.DecodeString(tdata[id].PkB_sike)
	if err != nil {
		t.Fatal(err)
	}
	sk_hex, err := hex.DecodeString(tdata[id].PrB_sike)
	if err != nil {
		t.Fatal(err)
	}
	if pkB.Import(pk_hex) != nil || skB.Import(sk_hex) != nil {
		t.Error("Import")
	}

	ct, err := Encrypt(rand.Reader, pkB, msg[:])
	if err != nil {
		t.Fatal(err)
	}
	pt, err := Decrypt(skB, ct)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(pt[:], msg[:]) {
		t.Errorf("Decryption failed \n got : %X\n exp : %X", pt, msg)
	}
}

// Generate key and check if can encrypt
func testPKEKeyGeneration(t testing.TB, id uint8) {
	// Message to be encrypted
	var params = Params(id)
	var msg = make([]byte, params.MsgLen)
	var err error
	for i, _ := range msg {
		msg[i] = byte(i)
	}

	sk := NewPrivateKey(id, KeyVariant_SIKE)
	err = sk.Generate(rand.Reader)
	checkErr(t, err, "PEK key generation")
	pk := sk.GeneratePublicKey()

	// Try to encrypt
	ct, err := Encrypt(rand.Reader, pk, msg[:])
	checkErr(t, err, "PEK encryption")
	pt, err := Decrypt(sk, ct)
	checkErr(t, err, "PEK key decryption")

	if !bytes.Equal(pt[:], msg[:]) {
		t.Fatalf("Decryption failed \n got : %X\n exp : %X", pt, msg)
	}
}

func testNegativePKE(t testing.TB, id uint8) {
	var msg [40]byte
	var err error
	var params = Params(id)

	// Generate key
	sk := NewPrivateKey(params.Id, KeyVariant_SIKE)
	err = sk.Generate(rand.Reader)
	checkErr(t, err, "key generation")

	pk := sk.GeneratePublicKey()

	// bytelen(msg) - 1
	ct, err := Encrypt(rand.Reader, pk, msg[:params.KemSize+8-1])
	if err == nil {
		t.Fatal("Error hasn't been returned")
	}
	if ct != nil {
		t.Fatal("Ciphertext must be nil")
	}

	// KemSize - 1
	pt, err := Decrypt(sk, msg[:params.KemSize+8-1])
	if err == nil {
		t.Fatal("Error hasn't been returned")
	}
	if pt != nil {
		t.Fatal("Ciphertext must be nil")
	}
}

func testKEMRoundTrip(t testing.TB, pkB, skB []byte, id uint8) {
	// Import keys
	pk := NewPublicKey(id, KeyVariant_SIKE)
	sk := NewPrivateKey(id, KeyVariant_SIKE)
	if pk.Import(pkB) != nil || sk.Import(skB) != nil {
		t.Error("Import failed")
	}

	ct, ss_e, err := Encapsulate(rand.Reader, pk)
	if err != nil {
		t.Error("Encapsulate failed")
	}

	ss_d, err := Decapsulate(sk, pk, ct)
	if err != nil {
		t.Error("Decapsulate failed")
	}
	if !bytes.Equal(ss_e, ss_d) {
		t.Error("Shared secrets from decapsulation and encapsulation differ")
	}
}

func TestKEMRoundTrip(t *testing.T) {
	for id, val := range tdata {
		fmt.Printf("\tTesting: %s\n", val.name)
		pk, err := hex.DecodeString(tdata[id].PkB_sike)
		checkErr(t, err, "public key B not a number")
		sk, err := hex.DecodeString(tdata[id].PrB_sike)
		checkErr(t, err, "private key B not a number")
		testKEMRoundTrip(t, pk, sk, id)
	}
}

func testKEMKeyGeneration(t testing.TB, id uint8) {
	// Generate key
	sk := NewPrivateKey(id, KeyVariant_SIKE)
	checkErr(t, sk.Generate(rand.Reader), "error: key generation")
	pk := sk.GeneratePublicKey()

	// calculated shared secret
	ct, ss_e, err := Encapsulate(rand.Reader, pk)
	checkErr(t, err, "encapsulation failed")
	ss_d, err := Decapsulate(sk, pk, ct)
	checkErr(t, err, "decapsulation failed")

	if !bytes.Equal(ss_e, ss_d) {
		t.Fatalf("KEM failed \n encapsulated: %X\n decapsulated: %X", ss_d, ss_e)
	}
}

func testNegativeKEM(t testing.TB, id uint8) {
	sk := NewPrivateKey(id, KeyVariant_SIKE)
	checkErr(t, sk.Generate(rand.Reader), "error: key generation")
	pk := sk.GeneratePublicKey()

	ct, ss_e, err := Encapsulate(rand.Reader, pk)
	checkErr(t, err, "pre-requisite for a test failed")

	ct[0] = ct[0] - 1
	ss_d, err := Decapsulate(sk, pk, ct)
	checkErr(t, err, "decapsulation returns error when invalid ciphertext provided")

	if bytes.Equal(ss_e, ss_d) {
		// no idea how this could ever happen, but it would be very bad
		t.Error("critical error")
	}

	// Try encapsulating with SIDH key
	pkSidh := NewPublicKey(id, KeyVariant_SIDH_B)
	prSidh := NewPrivateKey(id, KeyVariant_SIDH_B)
	_, _, err = Encapsulate(rand.Reader, pkSidh)
	if err == nil {
		t.Error("encapsulation accepts SIDH public key")
	}
	// Try decapsulating with SIDH key
	_, err = Decapsulate(prSidh, pk, ct)
	if err == nil {
		t.Error("decapsulation accepts SIDH private key key")
	}
}

// In case invalid ciphertext is provided, SIKE's decapsulation must
// return same (but unpredictable) result for a given key.
func testNegativeKEMSameWrongResult(t testing.TB, id uint8) {
	sk := NewPrivateKey(id, KeyVariant_SIKE)
	checkErr(t, sk.Generate(rand.Reader), "error: key generation")
	pk := sk.GeneratePublicKey()

	ct, encSs, err := Encapsulate(rand.Reader, pk)
	checkErr(t, err, "pre-requisite for a test failed")

	// make ciphertext wrong
	ct[0] = ct[0] - 1
	decSs1, err := Decapsulate(sk, pk, ct)
	checkErr(t, err, "pre-requisite for a test failed")

	// second decapsulation must be done with same, but imported private key
	expSk := sk.Export()

	// creat new private key
	sk = NewPrivateKey(id, KeyVariant_SIKE)
	err = sk.Import(expSk)
	checkErr(t, err, "import failed")

	// try decapsulating again. ss2 must be same as ss1 and different than
	// original plaintext
	decSs2, err := Decapsulate(sk, pk, ct)
	checkErr(t, err, "pre-requisite for a test failed")

	if !bytes.Equal(decSs1, decSs2) {
		t.Error("decapsulation is insecure")
	}

	if bytes.Equal(encSs, decSs1) || bytes.Equal(encSs, decSs2) {
		// this test requires that decapsulation returns wrong result
		t.Errorf("test implementation error")
	}
}

func readAndCheckLine(r *bufio.Reader) []byte {
	// Read next line from buffer
	line, isPrefix, err := r.ReadLine()
	if err != nil || isPrefix {
		panic("Wrong format of input file")
	}

	// Function expects that line is in format "KEY = HEX_VALUE". Get
	// value, which should be a hex string
	hexst := strings.Split(string(line), "=")[1]
	hexst = strings.TrimSpace(hexst)
	// Convert value to byte string
	ret, err := hex.DecodeString(hexst)
	if err != nil {
		panic("Wrong format of input file")
	}
	return ret
}

func testKeygenSIKE(pk, sk []byte, id uint8) bool {
	// Import provided private key
	var prvKey = NewPrivateKey(id, KeyVariant_SIKE)
	if prvKey.Import(sk) != nil {
		panic("sike test: can't load KAT")
	}

	// Generate public key
	pubKey := prvKey.GeneratePublicKey()
	return bytes.Equal(pubKey.Export(), pk)
}

func testDecapsulation(pk, sk, ct, ssExpected []byte, id uint8) bool {
	var pubKey = NewPublicKey(id, KeyVariant_SIKE)
	var prvKey = NewPrivateKey(id, KeyVariant_SIKE)
	if pubKey.Import(pk) != nil || prvKey.Import(sk) != nil {
		panic("sike test: can't load KAT")
	}

	ssGot, err := Decapsulate(prvKey, pubKey, ct)
	if err != nil {
		panic("sike test: can't perform decapsulation KAT")
	}

	if err != nil {
		return false
	}
	return bytes.Equal(ssGot, ssExpected)
}

func TestKeyAgreement(t *testing.T) {
	for id, val := range tdata {
		fmt.Printf("\tTesting: %s\n", val.name)
		testKeyAgreement(t, id, tdata[id].PkA_sike, tdata[id].PrA_sike, tdata[id].PkB_sidh, tdata[id].PrB_sidh)
	}
}

/* -------------------------------------------------------------------------
   Wrappers for 'testing' module
   -------------------------------------------------------------------------*/
func TestPKEKeyGeneration(t *testing.T)           { Do(testPKEKeyGeneration, t) }
func TestPKERoundTrip(t *testing.T)               { Do(testPKERoundTrip, t) }
func TestNegativePKE(t *testing.T)                { Do(testNegativePKE, t) }
func TestKEMKeyGeneration(t *testing.T)           { Do(testKEMKeyGeneration, t) }
func TestNegativeKEM(t *testing.T)                { Do(testNegativeKEM, t) }
func TestNegativeKEMSameWrongResult(t *testing.T) { Do(testNegativeKEMSameWrongResult, t) }
func TestKeygen(t *testing.T)                     { Do(testKeygen, t) }
func TestDerivationRoundTrip(t *testing.T)        { Do(testDerivationRoundTrip, t) }
func TestImportExport(t *testing.T)               { Do(testImportExport, t) }

/* -------------------------------------------------------------------------
   Benchmarking
   -------------------------------------------------------------------------*/

func BenchmarkSidhKeyAgreementP503(b *testing.B) {
	// KeyPairs
	alicePublic := convToPub(tdata[FP_503].PkA_sike, KeyVariant_SIDH_A, FP_503)
	alicePrivate := convToPrv(tdata[FP_503].PrA_sike, KeyVariant_SIDH_A, FP_503)
	bobPublic := convToPub(tdata[FP_503].PkB_sidh, KeyVariant_SIDH_B, FP_503)
	bobPrivate := convToPrv(tdata[FP_503].PrB_sidh, KeyVariant_SIDH_B, FP_503)

	for i := 0; i < b.N; i++ {
		// Derive shared secret
		DeriveSecret(bobPrivate, alicePublic)
		DeriveSecret(alicePrivate, bobPublic)
	}
}

func BenchmarkAliceKeyGenPrvP503(b *testing.B) {
	prv := NewPrivateKey(FP_503, KeyVariant_SIDH_A)
	for n := 0; n < b.N; n++ {
		prv.Generate(rand.Reader)
	}
}

func BenchmarkBobKeyGenPrvP503(b *testing.B) {
	prv := NewPrivateKey(FP_503, KeyVariant_SIDH_B)
	for n := 0; n < b.N; n++ {
		prv.Generate(rand.Reader)
	}
}

func BenchmarkAliceKeyGenPubP503(b *testing.B) {
	prv := NewPrivateKey(FP_503, KeyVariant_SIDH_A)
	prv.Generate(rand.Reader)
	for n := 0; n < b.N; n++ {
		prv.GeneratePublicKey()
	}
}

func BenchmarkBobKeyGenPubP503(b *testing.B) {
	prv := NewPrivateKey(FP_503, KeyVariant_SIDH_B)
	prv.Generate(rand.Reader)
	for n := 0; n < b.N; n++ {
		prv.GeneratePublicKey()
	}
}

func BenchmarkSharedSecretAliceP503(b *testing.B) {
	aPr := convToPrv(tdata[FP_503].PrA_sike, KeyVariant_SIDH_A, FP_503)
	bPk := convToPub(tdata[FP_503].PkB_sike, KeyVariant_SIDH_B, FP_503)
	for n := 0; n < b.N; n++ {
		DeriveSecret(aPr, bPk)
	}
}

func BenchmarkSharedSecretBobP503(b *testing.B) {
	// m_B = 3*randint(0,3^238)
	aPk := convToPub(tdata[FP_503].PkA_sike, KeyVariant_SIDH_A, FP_503)
	bPr := convToPrv(tdata[FP_503].PrB_sidh, KeyVariant_SIDH_B, FP_503)
	for n := 0; n < b.N; n++ {
		DeriveSecret(bPr, aPk)
	}
}
