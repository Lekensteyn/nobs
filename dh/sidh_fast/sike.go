// [SIKE] http://www.sike.org/files/SIDH-spec.pdf
// [REF] https://github.com/Microsoft/PQCrypto-SIDH
package sike

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"io"
)

// Constants used for cSHAKE customization
// Those values are different than in [SIKE] - they are encoded on 16bits. This is
// done in order for implementation to be compatible with [REF] and test vectors.
var G = []byte{0x00, 0x00}
var H = []byte{0x01, 0x00}
var F = []byte{0x02, 0x00}

// Generates HMAC-SHA256 sum
func HMAC(out, in, S []byte) {
	h := hmac.New(sha256.New, in)
	h.Write(S)
	copy(out, h.Sum(nil))
	//	fmt.Printf("> %X\n", out)
}

func Zeroize(fp *Fp2Element) {
	// Zeroizing in 2 seperated loops tells compiler to
	// use fast runtime.memclr()
	for i := range fp.A {
		fp.A[i] = 0
	}
	for i := range fp.B {
		fp.B[i] = 0
	}
}

// -----------------------------------------------------------------------------
// Functions for traversing isogeny trees acoording to strategy. Key type 'A' is
//

// Traverses isogeny tree in order to compute xR, xP, xQ and xQmP needed
// for public key generation.
func traverseTreePublicKeyA(curve *ProjectiveCurveParameters, xR, phiP, phiQ, phiR *ProjectivePoint, pub *PublicKey) {
	var points = make([]ProjectivePoint, 0, 8)
	var indices = make([]int, 0, 8)
	var i, sidx int
	var op = CurveOperations{Params: pub.params}

	cparam := op.CalcCurveParamsEquiv4(curve)
	phi := NewIsogeny4()
	strat := pub.params.A.IsogenyStrategy
	stratSz := len(strat)

	for j := 1; j <= stratSz; j++ {
		for i <= stratSz-j {
			points = append(points, *xR)
			indices = append(indices, i)

			k := strat[sidx]
			sidx++
			op.Pow2k(xR, &cparam, 2*k)
			i += int(k)
		}

		cparam = phi.GenerateCurve(xR)
		for k := 0; k < len(points); k++ {
			points[k] = phi.EvaluatePoint(&points[k])
		}

		*phiP = phi.EvaluatePoint(phiP)
		*phiQ = phi.EvaluatePoint(phiQ)
		*phiR = phi.EvaluatePoint(phiR)

		// pop xR from points
		*xR, points = points[len(points)-1], points[:len(points)-1]
		i, indices = int(indices[len(indices)-1]), indices[:len(indices)-1]
	}
}

// Traverses isogeny tree in order to compute xR needed
// for public key generation.
func traverseTreeSharedKeyA(curve *ProjectiveCurveParameters, xR *ProjectivePoint, pub *PublicKey) {
	var points = make([]ProjectivePoint, 0, 8)
	var indices = make([]int, 0, 8)
	var i, sidx int
	var op = CurveOperations{Params: pub.params}

	cparam := op.CalcCurveParamsEquiv4(curve)
	phi := NewIsogeny4()
	strat := pub.params.A.IsogenyStrategy
	stratSz := len(strat)

	for j := 1; j <= stratSz; j++ {
		for i <= stratSz-j {
			points = append(points, *xR)
			indices = append(indices, i)

			k := strat[sidx]
			sidx++
			op.Pow2k(xR, &cparam, 2*k)
			i += int(k)
		}

		cparam = phi.GenerateCurve(xR)
		for k := 0; k < len(points); k++ {
			points[k] = phi.EvaluatePoint(&points[k])
		}

		// pop xR from points
		*xR, points = points[len(points)-1], points[:len(points)-1]
		i, indices = int(indices[len(indices)-1]), indices[:len(indices)-1]
	}
}

// Traverses isogeny tree in order to compute xR, xP, xQ and xQmP needed
// for public key generation.
func traverseTreePublicKeyB(curve *ProjectiveCurveParameters, xR, phiP, phiQ, phiR *ProjectivePoint, pub *PublicKey) {
	var points = make([]ProjectivePoint, 0, 8)
	var indices = make([]int, 0, 8)
	var i, sidx int
	var op = CurveOperations{Params: pub.params}

	cparam := op.CalcCurveParamsEquiv3(curve)
	phi := NewIsogeny3()
	strat := pub.params.B.IsogenyStrategy
	stratSz := len(strat)

	for j := 1; j <= stratSz; j++ {
		for i <= stratSz-j {
			points = append(points, *xR)
			indices = append(indices, i)

			k := strat[sidx]
			sidx++
			op.Pow3k(xR, &cparam, k)
			i += int(k)
		}

		cparam = phi.GenerateCurve(xR)
		for k := 0; k < len(points); k++ {
			points[k] = phi.EvaluatePoint(&points[k])
		}

		*phiP = phi.EvaluatePoint(phiP)
		*phiQ = phi.EvaluatePoint(phiQ)
		*phiR = phi.EvaluatePoint(phiR)

		// pop xR from points
		*xR, points = points[len(points)-1], points[:len(points)-1]
		i, indices = int(indices[len(indices)-1]), indices[:len(indices)-1]
	}
}

// Traverses isogeny tree in order to compute xR, xP, xQ and xQmP needed
// for public key generation.
func traverseTreeSharedKeyB(curve *ProjectiveCurveParameters, xR *ProjectivePoint, pub *PublicKey) {
	var points = make([]ProjectivePoint, 0, 8)
	var indices = make([]int, 0, 8)
	var i, sidx int
	var op = CurveOperations{Params: pub.params}

	cparam := op.CalcCurveParamsEquiv3(curve)
	phi := NewIsogeny3()
	strat := pub.params.B.IsogenyStrategy
	stratSz := len(strat)

	for j := 1; j <= stratSz; j++ {
		for i <= stratSz-j {
			points = append(points, *xR)
			indices = append(indices, i)

			k := strat[sidx]
			sidx++
			op.Pow3k(xR, &cparam, k)
			i += int(k)
		}

		cparam = phi.GenerateCurve(xR)
		for k := 0; k < len(points); k++ {
			points[k] = phi.EvaluatePoint(&points[k])
		}

		// pop xR from points
		*xR, points = points[len(points)-1], points[:len(points)-1]
		i, indices = int(indices[len(indices)-1]), indices[:len(indices)-1]
	}
}

// Generate a public key in the 2-torsion group
func publicKeyGenA(prv *PrivateKey) (pub *PublicKey) {
	var xPA, xQA, xRA ProjectivePoint
	var xPB, xQB, xRB, xR ProjectivePoint
	var invZP, invZQ, invZR Fp2Element
	var tmp ProjectiveCurveParameters

	pub = NewPublicKey(prv.params.Id, KeyVariant_SIDH_A)
	var op = CurveOperations{Params: pub.params}
	var phi = NewIsogeny4()

	// Load points for A
	xPA = ProjectivePoint{X: prv.params.A.Affine_P, Z: prv.params.OneFp2}
	xQA = ProjectivePoint{X: prv.params.A.Affine_Q, Z: prv.params.OneFp2}
	xRA = ProjectivePoint{X: prv.params.A.Affine_R, Z: prv.params.OneFp2}

	// Load points for B
	xRB = ProjectivePoint{X: prv.params.B.Affine_R, Z: prv.params.OneFp2}
	xQB = ProjectivePoint{X: prv.params.B.Affine_Q, Z: prv.params.OneFp2}
	xPB = ProjectivePoint{X: prv.params.B.Affine_P, Z: prv.params.OneFp2}

	// Find isogeny kernel
	tmp.C = pub.params.OneFp2
	xR = op.ScalarMul3Pt(&tmp, &xPA, &xQA, &xRA, prv.params.A.SecretBitLen, prv.Scalar)

	// Reset params object and travers isogeny tree
	tmp.C = pub.params.OneFp2
	Zeroize(&tmp.A)
	traverseTreePublicKeyA(&tmp, &xR, &xPB, &xQB, &xRB, pub)

	// Secret isogeny
	phi.GenerateCurve(&xR)
	xPA = phi.EvaluatePoint(&xPB)
	xQA = phi.EvaluatePoint(&xQB)
	xRA = phi.EvaluatePoint(&xRB)
	op.Fp2Batch3Inv(&xPA.Z, &xQA.Z, &xRA.Z, &invZP, &invZQ, &invZR)

	Mul(&pub.affine_xP, &xPA.X, &invZP)
	Mul(&pub.affine_xQ, &xQA.X, &invZQ)
	Mul(&pub.affine_xQmP, &xRA.X, &invZR)
	return
}

// Generate a public key in the 3-torsion group
func publicKeyGenB(prv *PrivateKey) (pub *PublicKey) {
	var xPB, xQB, xRB, xR ProjectivePoint
	var xPA, xQA, xRA ProjectivePoint
	var invZP, invZQ, invZR Fp2Element
	var tmp ProjectiveCurveParameters

	pub = NewPublicKey(prv.params.Id, prv.keyVariant)
	var op = CurveOperations{Params: pub.params}
	var phi = NewIsogeny3()

	// Load points for B
	xRB = ProjectivePoint{X: prv.params.B.Affine_R, Z: prv.params.OneFp2}
	xQB = ProjectivePoint{X: prv.params.B.Affine_Q, Z: prv.params.OneFp2}
	xPB = ProjectivePoint{X: prv.params.B.Affine_P, Z: prv.params.OneFp2}

	// Load points for A
	xPA = ProjectivePoint{X: prv.params.A.Affine_P, Z: prv.params.OneFp2}
	xQA = ProjectivePoint{X: prv.params.A.Affine_Q, Z: prv.params.OneFp2}
	xRA = ProjectivePoint{X: prv.params.A.Affine_R, Z: prv.params.OneFp2}

	tmp.C = pub.params.OneFp2
	xR = op.ScalarMul3Pt(&tmp, &xPB, &xQB, &xRB, prv.params.B.SecretBitLen, prv.Scalar)

	tmp.C = pub.params.OneFp2
	Zeroize(&tmp.A)
	traverseTreePublicKeyB(&tmp, &xR, &xPA, &xQA, &xRA, pub)

	phi.GenerateCurve(&xR)
	xPB = phi.EvaluatePoint(&xPA)
	xQB = phi.EvaluatePoint(&xQA)
	xRB = phi.EvaluatePoint(&xRA)
	op.Fp2Batch3Inv(&xPB.Z, &xQB.Z, &xRB.Z, &invZP, &invZQ, &invZR)

	Mul(&pub.affine_xP, &xPB.X, &invZP)
	Mul(&pub.affine_xQ, &xQB.X, &invZQ)
	Mul(&pub.affine_xQmP, &xRB.X, &invZR)
	return
}

// -----------------------------------------------------------------------------
// Key agreement functions
//

// Establishing shared keys in in 2-torsion group
func deriveSecretA(prv *PrivateKey, pub *PublicKey) []byte {
	var sharedSecret = make([]byte, pub.params.SharedSecretSize)
	var cparam ProjectiveCurveParameters
	var xP, xQ, xQmP ProjectivePoint
	var xR ProjectivePoint
	var op = CurveOperations{Params: prv.params}
	var phi = NewIsogeny4()

	// Recover curve coefficients
	cparam.C = pub.params.OneFp2
	op.RecoverCoordinateA(&cparam, &pub.affine_xP, &pub.affine_xQ, &pub.affine_xQmP)

	// Find kernel of the morphism
	xP = ProjectivePoint{X: pub.affine_xP, Z: pub.params.OneFp2}
	xQ = ProjectivePoint{X: pub.affine_xQ, Z: pub.params.OneFp2}
	xQmP = ProjectivePoint{X: pub.affine_xQmP, Z: pub.params.OneFp2}
	xR = op.ScalarMul3Pt(&cparam, &xP, &xQ, &xQmP, pub.params.A.SecretBitLen, prv.Scalar)

	// Traverse isogeny tree
	traverseTreeSharedKeyA(&cparam, &xR, pub)

	// Calculate j-invariant on isogeneus curve
	c := phi.GenerateCurve(&xR)
	op.RecoverCurveCoefficients4(&cparam, &c)
	op.Jinvariant(&cparam, sharedSecret)
	return sharedSecret
}

// Establishing shared keys in in 3-torsion group
func deriveSecretB(prv *PrivateKey, pub *PublicKey) []byte {
	var sharedSecret = make([]byte, pub.params.SharedSecretSize)
	var xP, xQ, xQmP ProjectivePoint
	var xR ProjectivePoint
	var cparam ProjectiveCurveParameters
	var op = CurveOperations{Params: prv.params}
	var phi = NewIsogeny3()

	// Recover curve coefficients
	cparam.C = pub.params.OneFp2
	op.RecoverCoordinateA(&cparam, &pub.affine_xP, &pub.affine_xQ, &pub.affine_xQmP)

	// Find kernel of the morphism
	xP = ProjectivePoint{X: pub.affine_xP, Z: pub.params.OneFp2}
	xQ = ProjectivePoint{X: pub.affine_xQ, Z: pub.params.OneFp2}
	xQmP = ProjectivePoint{X: pub.affine_xQmP, Z: pub.params.OneFp2}
	xR = op.ScalarMul3Pt(&cparam, &xP, &xQ, &xQmP, pub.params.B.SecretBitLen, prv.Scalar)

	// Traverse isogeny tree
	traverseTreeSharedKeyB(&cparam, &xR, pub)

	// Calculate j-invariant on isogeneus curve
	c := phi.GenerateCurve(&xR)
	op.RecoverCurveCoefficients3(&cparam, &c)
	op.Jinvariant(&cparam, sharedSecret)
	return sharedSecret
}

func encrypt(skA *PrivateKey, pkA, pkB *PublicKey, ptext []byte) ([]byte, error) {
	var n [40]byte // n can is max 320-bit (see 1.4 of [SIKE])
	var ptextLen = len(ptext)

	if pkB.Variant() != KeyVariant_SIKE {
		return nil, errors.New("wrong key type")
	}

	j, err := DeriveSecret(skA, pkB)
	if err != nil {
		return nil, err
	}

	HMAC(n[:ptextLen], j, F)
	for i, _ := range ptext {
		n[i] ^= ptext[i]
	}

	ret := make([]byte, pkA.Size()+ptextLen)
	copy(ret, pkA.Export())
	copy(ret[pkA.Size():], n[:ptextLen])
	return ret, nil
}

// -----------------------------------------------------------------------------
// PKE interface
//

// Uses SIKE public key to encrypt plaintext. Requires cryptographically secure PRNG
// Returns ciphertext in case encryption succeeds. Returns error in case PRNG fails
// or wrongly formated input was provided.
func Encrypt(rng io.Reader, pub *PublicKey, ptext []byte) ([]byte, error) {
	var params = pub.Params()
	var ptextLen = len(ptext)
	// c1 must be security level + 64 bits (see [SIKE] 1.4 and 4.3.3)
	if ptextLen != (params.KemSize + 8) {
		return nil, errors.New("Unsupported message length")
	}

	skA := NewPrivateKey(params.Id, KeyVariant_SIDH_A)
	err := skA.Generate(rng)
	if err != nil {
		return nil, err
	}

	pkA := skA.GeneratePublicKey()
	return encrypt(skA, pkA, pub, ptext)
}

// Uses SIKE private key to decrypt ciphertext. Returns plaintext in case
// decryption succeeds or error in case unexptected input was provided.
// Constant time
func Decrypt(prv *PrivateKey, ctext []byte) ([]byte, error) {
	var params = prv.Params()
	var n [40]byte // n can is max 320-bit (see 1.4 of [SIKE])
	var c1_len int
	var pk_len = params.PublicKeySize

	if prv.Variant() != KeyVariant_SIKE {
		return nil, errors.New("wrong key type")
	}

	// ctext is a concatenation of (pubkey_A || c1=ciphertext)
	// it must be security level + 64 bits (see [SIKE] 1.4 and 4.3.3)
	c1_len = len(ctext) - pk_len
	if c1_len != (int(params.KemSize) + 8) {
		return nil, errors.New("wrong size of cipher text")
	}

	c0 := NewPublicKey(params.Id, KeyVariant_SIDH_A)
	err := c0.Import(ctext[:pk_len])
	if err != nil {
		return nil, err
	}
	j, err := DeriveSecret(prv, c0)
	if err != nil {
		return nil, err
	}

	HMAC(n[:c1_len], j, F)
	for i, _ := range n[:c1_len] {
		n[i] ^= ctext[pk_len+i]
	}
	return n[:c1_len], nil
}

// -----------------------------------------------------------------------------
// KEM interface
//

// Encapsulation receives the public key and generates SIKE ciphertext and shared secret.
// The generated ciphertext is used for authentication.
// The rng must be cryptographically secure PRNG.
// Error is returned in case PRNG fails or wrongly formated input was provided.
func Encapsulate(rng io.Reader, pub *PublicKey) (ctext []byte, secret []byte, err error) {
	var params = pub.Params()
	// Buffer for random, secret message
	var ptext = make([]byte, params.MsgLen)
	// r = G(ptext||pub)
	var r = make([]byte, params.A.SecretByteLen)
	// Resulting shared secret
	secret = make([]byte, params.KemSize)

	// Generate ephemeral value
	_, err = io.ReadFull(rng, ptext)
	if err != nil {
		return nil, nil, err
	}

	var hmac_key [378 + 24 + 24]byte //make([]byte, len(ptext)+pub.Size())
	copy(hmac_key[:], ptext)
	copy(hmac_key[len(ptext):], pub.Export())
	HMAC(r, hmac_key[:len(ptext)+pub.Size()], G)
	// Ensure bitlength is not bigger then to 2^e2-1
	r[len(r)-1] &= (1 << (params.A.SecretBitLen % 8)) - 1

	// (c0 || c1) = Enc(pkA, ptext; r)
	skA := NewPrivateKey(params.Id, KeyVariant_SIDH_A)
	err = skA.Import(r)
	if err != nil {
		return nil, nil, err
	}

	pkA := skA.GeneratePublicKey()
	ctext, err = encrypt(skA, pkA, pub, ptext)
	if err != nil {
		return nil, nil, err
	}

	// K = H(ptext||(c0||c1))
	copy(hmac_key[:], ptext)
	copy(hmac_key[len(ptext):], ctext)
	HMAC(secret, hmac_key[:len(ptext)+len(ctext)], H)
	return ctext, secret, nil
}

// Decapsulate given the keypair and ciphertext as inputs, Decapsulate outputs a shared
// secret if plaintext verifies correctly, otherwise function outputs random value.
// Decapsulation may fail in case input is wrongly formated.
// Constant time for properly initialized input.
func Decapsulate(prv *PrivateKey, pub *PublicKey, ctext []byte) ([]byte, error) {
	var params = pub.Params()
	var r = make([]byte, params.A.SecretByteLen)
	// Resulting shared secret
	var secret = make([]byte, params.KemSize)
	var skA = NewPrivateKey(params.Id, KeyVariant_SIDH_A)

	m, err := Decrypt(prv, ctext)
	if err != nil {
		return nil, err
	}

	// r' = G(m'||pub)
	var hmac_key [378 + 24 + 24]byte //make([]byte, len(m)+pub.Size())
	copy(hmac_key[:], m)
	copy(hmac_key[len(m):], pub.Export())
	HMAC(r, hmac_key[:len(m)+pub.Size()], G)
	// Ensure bitlength is not bigger than 2^e2-1
	r[len(r)-1] &= (1 << (params.A.SecretBitLen % 8)) - 1

	// Never fails
	skA.Import(r)

	// Never fails
	pkA := skA.GeneratePublicKey()
	c0 := pkA.Export()

	if subtle.ConstantTimeCompare(c0, ctext[:len(c0)]) == 1 {
		copy(hmac_key[:], m)
	} else {
		// S is chosen at random when generating a key and unknown to other party. It
		// may seem weird, but it's correct. It is important that S is unpredictable
		// to other party. Without this check, it is possible to recover a secret, by
		// providing series of invalid ciphertexts. It is also important that in case
		//
		// See more details in "On the security of supersingular isogeny cryptosystems"
		// (S. Galbraith, et al., 2016, ePrint #859).
		copy(hmac_key[:], prv.S)
	}
	copy(hmac_key[len(m):], ctext)
	HMAC(secret, hmac_key[:len(m)+len(ctext)], H)
	return secret, nil
}
