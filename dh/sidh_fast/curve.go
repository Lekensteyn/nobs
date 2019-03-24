package sike

// Interface for working with isogenies.
type isogeny interface {
	// Given a torsion point on a curve computes isogenous curve.
	// Returns curve coefficients (A:C), so that E_(A/C) = E_(A/C)/<P>,
	// where P is a provided projective point. Sets also isogeny constants
	// that are needed for isogeny evaluation.
	GenerateCurve(*ProjectivePoint) CurveCoefficientsEquiv
	// Evaluates isogeny at caller provided point. Requires isogeny curve constants
	// to be earlier computed by GenerateCurve.
	EvaluatePoint(*ProjectivePoint) ProjectivePoint
}

// Stores isogeny 3 curve constants
type isogeny3 struct {
	K1 Fp2Element
	K2 Fp2Element
}

// Stores isogeny 4 curve constants
type isogeny4 struct {
	isogeny3
	K3 Fp2Element
}

// Computes j-invariant for a curve y2=x3+A/Cx+x with A,C in F_(p^2). Result
// is returned in jBytes buffer, encoded in little-endian format. Caller
// provided jBytes buffer has to be big enough to j-invariant value. In case
// of SIDH, buffer size must be at least size of shared secret.
// Implementation corresponds to Algorithm 9 from SIKE.
func (c *CurveOperations) Jinvariant(cparams *ProjectiveCurveParameters, jBytes []byte) {
	var j, t0, t1 Fp2Element

	Square(&j, &cparams.A)  // j  = A^2
	Square(&t1, &cparams.C) // t1 = C^2
	Add(&t0, &t1, &t1)      // t0 = t1 + t1
	Sub(&t0, &j, &t0)       // t0 = j - t0
	Sub(&t0, &t0, &t1)      // t0 = t0 - t1
	Sub(&j, &t0, &t1)       // t0 = t0 - t1
	Square(&t1, &t1)        // t1 = t1^2
	Mul(&j, &j, &t1)        // j = j * t1
	Add(&t0, &t0, &t0)      // t0 = t0 + t0
	Add(&t0, &t0, &t0)      // t0 = t0 + t0
	Square(&t1, &t0)        // t1 = t0^2
	Mul(&t0, &t0, &t1)      // t0 = t0 * t1
	Add(&t0, &t0, &t0)      // t0 = t0 + t0
	Add(&t0, &t0, &t0)      // t0 = t0 + t0
	Inv(&j, &j)             // j  = 1/j
	Mul(&j, &t0, &j)        // j  = t0 * j

	c.Fp2ToBytes(jBytes, &j)
}

// Given affine points x(P), x(Q) and x(Q-P) in a extension field F_{p^2}, function
// recorvers projective coordinate A of a curve. This is Algorithm 10 from SIKE.
func (c *CurveOperations) RecoverCoordinateA(curve *ProjectiveCurveParameters, xp, xq, xr *Fp2Element) {
	var t0, t1 Fp2Element

	Add(&t1, xp, xq)                          // t1 = Xp + Xq
	Mul(&t0, xp, xq)                          // t0 = Xp * Xq
	Mul(&curve.A, xr, &t1)                    // A  = X(q-p) * t1
	Add(&curve.A, &curve.A, &t0)              // A  = A + t0
	Mul(&t0, &t0, xr)                         // t0 = t0 * X(q-p)
	Sub(&curve.A, &curve.A, &c.Params.OneFp2) // A  = A - 1
	Add(&t0, &t0, &t0)                        // t0 = t0 + t0
	Add(&t1, &t1, xr)                         // t1 = t1 + X(q-p)
	Add(&t0, &t0, &t0)                        // t0 = t0 + t0
	Square(&curve.A, &curve.A)                // A  = A^2
	Inv(&t0, &t0)                             // t0 = 1/t0
	Mul(&curve.A, &curve.A, &t0)              // A  = A * t0
	Sub(&curve.A, &curve.A, &t1)              // A  = A - t1
}

// Computes equivalence (A:C) ~ (A+2C : A-2C)
func (c *CurveOperations) CalcCurveParamsEquiv3(cparams *ProjectiveCurveParameters) CurveCoefficientsEquiv {
	var coef CurveCoefficientsEquiv
	var c2 Fp2Element

	Add(&c2, &cparams.C, &cparams.C)
	// A24p = A+2*C
	Add(&coef.A, &cparams.A, &c2)
	// A24m = A-2*C
	Sub(&coef.C, &cparams.A, &c2)
	return coef
}

// Computes equivalence (A:C) ~ (A+2C : 4C)
func (c *CurveOperations) CalcCurveParamsEquiv4(cparams *ProjectiveCurveParameters) CurveCoefficientsEquiv {
	var coefEq CurveCoefficientsEquiv

	Add(&coefEq.C, &cparams.C, &cparams.C)
	// A24p = A+2C
	Add(&coefEq.A, &cparams.A, &coefEq.C)
	// C24 = 4*C
	Add(&coefEq.C, &coefEq.C, &coefEq.C)
	return coefEq
}

// Helper function for RightToLeftLadder(). Returns A+2C / 4.
func (c *CurveOperations) CalcAplus2Over4(cparams *ProjectiveCurveParameters) (ret Fp2Element) {
	var tmp Fp2Element

	// 2C
	Add(&tmp, &cparams.C, &cparams.C)
	// A+2C
	Add(&ret, &cparams.A, &tmp)
	// 1/4C
	Add(&tmp, &tmp, &tmp)
	Inv(&tmp, &tmp)
	// A+2C/4C
	Mul(&ret, &ret, &tmp)
	return
}

// Recovers (A:C) curve parameters from projectively equivalent (A+2C:A-2C).
func (c *CurveOperations) RecoverCurveCoefficients3(cparams *ProjectiveCurveParameters, coefEq *CurveCoefficientsEquiv) {
	Add(&cparams.A, &coefEq.A, &coefEq.C)
	// cparams.A = 2*(A+2C+A-2C) = 4A
	Add(&cparams.A, &cparams.A, &cparams.A)
	// cparams.C = (A+2C-A+2C) = 4C
	Sub(&cparams.C, &coefEq.A, &coefEq.C)
	return
}

// Recovers (A:C) curve parameters from projectively equivalent (A+2C:4C).
func (c *CurveOperations) RecoverCurveCoefficients4(cparams *ProjectiveCurveParameters, coefEq *CurveCoefficientsEquiv) {
	// cparams.C = (4C)*1/2=2C
	Mul(&cparams.C, &coefEq.C, &c.Params.HalfFp2)
	// cparams.A = A+2C - 2C = A
	Sub(&cparams.A, &coefEq.A, &cparams.C)
	// cparams.C = 2C * 1/2 = C
	Mul(&cparams.C, &cparams.C, &c.Params.HalfFp2)
	return
}

// Combined coordinate doubling and differential addition. Takes projective points
// P,Q,Q-P and (A+2C)/4C curve E coefficient. Returns 2*P and P+Q calculated on E.
// Function is used only by RightToLeftLadder. Corresponds to Algorithm 5 of SIKE
func (c *CurveOperations) xDblAdd(P, Q, QmP *ProjectivePoint, a24 *Fp2Element) (dblP, PaQ ProjectivePoint) {
	var t0, t1, t2 Fp2Element
	xQmP, zQmP := &QmP.X, &QmP.Z
	xPaQ, zPaQ := &PaQ.X, &PaQ.Z
	x2P, z2P := &dblP.X, &dblP.Z
	xP, zP := &P.X, &P.Z
	xQ, zQ := &Q.X, &Q.Z

	Add(&t0, xP, zP)      // t0   = Xp+Zp
	Sub(&t1, xP, zP)      // t1   = Xp-Zp
	Square(x2P, &t0)      // 2P.X = t0^2
	Sub(&t2, xQ, zQ)      // t2   = Xq-Zq
	Add(xPaQ, xQ, zQ)     // Xp+q = Xq+Zq
	Mul(&t0, &t0, &t2)    // t0   = t0 * t2
	Mul(z2P, &t1, &t1)    // 2P.Z = t1 * t1
	Mul(&t1, &t1, xPaQ)   // t1   = t1 * Xp+q
	Sub(&t2, x2P, z2P)    // t2   = 2P.X - 2P.Z
	Mul(x2P, x2P, z2P)    // 2P.X = 2P.X * 2P.Z
	Mul(xPaQ, a24, &t2)   // Xp+q = A24 * t2
	Sub(zPaQ, &t0, &t1)   // Zp+q = t0 - t1
	Add(z2P, xPaQ, z2P)   // 2P.Z = Xp+q + 2P.Z
	Add(xPaQ, &t0, &t1)   // Xp+q = t0 + t1
	Mul(z2P, z2P, &t2)    // 2P.Z = 2P.Z * t2
	Square(zPaQ, zPaQ)    // Zp+q = Zp+q ^ 2
	Square(xPaQ, xPaQ)    // Xp+q = Xp+q ^ 2
	Mul(zPaQ, xQmP, zPaQ) // Zp+q = Xq-p * Zp+q
	Mul(xPaQ, zQmP, xPaQ) // Xp+q = Zq-p * Xp+q
	return
}

// Given the curve parameters, xP = x(P), computes xP = x([2^k]P)
// Safe to overlap xP, x2P.
func (c *CurveOperations) Pow2k(xP *ProjectivePoint, params *CurveCoefficientsEquiv, k uint32) {
	var t0, t1 Fp2Element

	x, z := &xP.X, &xP.Z
	for i := uint32(0); i < k; i++ {
		Sub(&t0, x, z)           // t0  = Xp - Zp
		Add(&t1, x, z)           // t1  = Xp + Zp
		Square(&t0, &t0)         // t0  = t0 ^ 2
		Square(&t1, &t1)         // t1  = t1 ^ 2
		Mul(z, &params.C, &t0)   // Z2p = C24 * t0
		Mul(x, z, &t1)           // X2p = Z2p * t1
		Sub(&t1, &t1, &t0)       // t1  = t1 - t0
		Mul(&t0, &params.A, &t1) // t0  = A24+ * t1
		Add(z, z, &t0)           // Z2p = Z2p + t0
		Mul(z, z, &t1)           // Zp  = Z2p * t1
	}
}

// Given the curve parameters, xP = x(P), and k >= 0, compute xP = x([3^k]P).
//
// Safe to overlap xP, xR.
func (c *CurveOperations) Pow3k(xP *ProjectivePoint, params *CurveCoefficientsEquiv, k uint32) {
	var t0, t1, t2, t3, t4, t5, t6 Fp2Element

	x, z := &xP.X, &xP.Z
	for i := uint32(0); i < k; i++ {
		Sub(&t0, x, z)           // t0  = Xp - Zp
		Square(&t2, &t0)         // t2  = t0^2
		Add(&t1, x, z)           // t1  = Xp + Zp
		Square(&t3, &t1)         // t3  = t1^2
		Add(&t4, &t1, &t0)       // t4  = t1 + t0
		Sub(&t0, &t1, &t0)       // t0  = t1 - t0
		Square(&t1, &t4)         // t1  = t4^2
		Sub(&t1, &t1, &t3)       // t1  = t1 - t3
		Sub(&t1, &t1, &t2)       // t1  = t1 - t2
		Mul(&t5, &t3, &params.A) // t5  = t3 * A24+
		Mul(&t3, &t3, &t5)       // t3  = t5 * t3
		Mul(&t6, &t2, &params.C) // t6  = t2 * A24-
		Mul(&t2, &t2, &t6)       // t2  = t2 * t6
		Sub(&t3, &t2, &t3)       // t3  = t2 - t3
		Sub(&t2, &t5, &t6)       // t2  = t5 - t6
		Mul(&t1, &t2, &t1)       // t1  = t2 * t1
		Add(&t2, &t3, &t1)       // t2  = t3 + t1
		Square(&t2, &t2)         // t2  = t2^2
		Mul(x, &t2, &t4)         // X3p = t2 * t4
		Sub(&t1, &t3, &t1)       // t1  = t3 - t1
		Square(&t1, &t1)         // t1  = t1^2
		Mul(z, &t1, &t0)         // Z3p = t1 * t0
	}
}

// Set (y1, y2, y3)  = (1/x1, 1/x2, 1/x3).
//
// All xi, yi must be distinct.
func (c *CurveOperations) Fp2Batch3Inv(x1, x2, x3, y1, y2, y3 *Fp2Element) {
	var x1x2, t Fp2Element

	Mul(&x1x2, x1, x2) // x1*x2
	Mul(&t, &x1x2, x3) // 1/(x1*x2*x3)
	Inv(&t, &t)
	Mul(y1, &t, x2) // 1/x1
	Mul(y1, y1, x3)
	Mul(y2, &t, x1) // 1/x2
	Mul(y2, y2, x3)
	Mul(y3, &t, &x1x2) // 1/x3
}

// ScalarMul3Pt is a right-to-left point multiplication that given the
// x-coordinate of P, Q and P-Q calculates the x-coordinate of R=Q+[scalar]P.
// nbits must be smaller or equal to len(scalar).
func (c *CurveOperations) ScalarMul3Pt(cparams *ProjectiveCurveParameters, P, Q, PmQ *ProjectivePoint, nbits uint, scalar []uint8) ProjectivePoint {
	var R0, R2, R1 ProjectivePoint
	aPlus2Over4 := c.CalcAplus2Over4(cparams)
	R1 = *P
	R2 = *PmQ
	R0 = *Q

	// Iterate over the bits of the scalar, bottom to top
	prevBit := uint8(0)
	for i := uint(0); i < nbits; i++ {
		bit := (scalar[i>>3] >> (i & 7) & 1)
		swap := prevBit ^ bit
		prevBit = bit
		CondSwap(&R1.X, &R1.Z, &R2.X, &R2.Z, swap)
		R0, R2 = c.xDblAdd(&R0, &R2, &R1, &aPlus2Over4)
	}
	CondSwap(&R1.X, &R1.Z, &R2.X, &R2.Z, prevBit)
	return R1
}

// Convert the input to wire format.
//
// The output byte slice must be at least 2*bytelen(p) bytes long.
func (c *CurveOperations) Fp2ToBytes(output []byte, fp2 *Fp2Element) {
	if len(output) < 2*c.Params.Bytelen {
		panic("output byte slice too short")
	}
	var a Fp2Element
	FromMontgomery(fp2, &a)

	// convert to bytes in little endian form
	for i := 0; i < c.Params.Bytelen; i++ {
		// set i = j*8 + k
		tmp := i / 8
		k := uint64(i % 8)
		output[i] = byte(a.A[tmp] >> (8 * k))
		output[i+c.Params.Bytelen] = byte(a.B[tmp] >> (8 * k))
	}
}

// Read 2*bytelen(p) bytes into the given ExtensionFieldElement.
//
// It is an error to call this function if the input byte slice is less than 2*bytelen(p) bytes long.
func (c *CurveOperations) Fp2FromBytes(fp2 *Fp2Element, input []byte) {
	if len(input) < 2*c.Params.Bytelen {
		panic("input byte slice too short")
	}

	for i := 0; i < c.Params.Bytelen; i++ {
		j := i / 8
		k := uint64(i % 8)
		fp2.A[j] |= uint64(input[i]) << (8 * k)
		fp2.B[j] |= uint64(input[i+c.Params.Bytelen]) << (8 * k)
	}
	ToMontgomery(fp2)
}

/* -------------------------------------------------------------------------
   Mechnisms used for isogeny calculations
   -------------------------------------------------------------------------*/

// Constructs isogeny3 objects
func NewIsogeny3() isogeny {
	return &isogeny3{}
}

// Constructs isogeny4 objects
func NewIsogeny4() isogeny {
	return &isogeny4{}
}

// Given a three-torsion point p = x(PB) on the curve E_(A:C), construct the
// three-isogeny phi : E_(A:C) -> E_(A:C)/<P_3> = E_(A':C').
//
// Input: (XP_3: ZP_3), where P_3 has exact order 3 on E_A/C
// Output: * Curve coordinates (A' + 2C', A' - 2C') corresponding to E_A'/C' = A_E/C/<P3>
//         * isogeny phi with constants in F_p^2
func (phi *isogeny3) GenerateCurve(p *ProjectivePoint) CurveCoefficientsEquiv {
	var t0, t1, t2, t3, t4 Fp2Element
	var coefEq CurveCoefficientsEquiv
	var K1, K2 = &phi.K1, &phi.K2

	Sub(K1, &p.X, &p.Z)            // K1 = XP3 - ZP3
	Square(&t0, K1)                // t0 = K1^2
	Add(K2, &p.X, &p.Z)            // K2 = XP3 + ZP3
	Square(&t1, K2)                // t1 = K2^2
	Add(&t2, &t0, &t1)             // t2 = t0 + t1
	Add(&t3, K1, K2)               // t3 = K1 + K2
	Square(&t3, &t3)               // t3 = t3^2
	Sub(&t3, &t3, &t2)             // t3 = t3 - t2
	Add(&t2, &t1, &t3)             // t2 = t1 + t3
	Add(&t3, &t3, &t0)             // t3 = t3 + t0
	Add(&t4, &t3, &t0)             // t4 = t3 + t0
	Add(&t4, &t4, &t4)             // t4 = t4 + t4
	Add(&t4, &t1, &t4)             // t4 = t1 + t4
	Mul(&coefEq.C, &t2, &t4)       // A24m = t2 * t4
	Add(&t4, &t1, &t2)             // t4 = t1 + t2
	Add(&t4, &t4, &t4)             // t4 = t4 + t4
	Add(&t4, &t0, &t4)             // t4 = t0 + t4
	Mul(&t4, &t3, &t4)             // t4 = t3 * t4
	Sub(&t0, &t4, &coefEq.C)       // t0 = t4 - A24m
	Add(&coefEq.A, &coefEq.C, &t0) // A24p = A24m + t0
	return coefEq
}

// Given a 3-isogeny phi and a point pB = x(PB), compute x(QB), the x-coordinate
// of the image QB = phi(PB) of PB under phi : E_(A:C) -> E_(A':C').
//
// The output xQ = x(Q) is then a point on the curve E_(A':C'); the curve
// parameters are returned by the GenerateCurve function used to construct phi.
func (phi *isogeny3) EvaluatePoint(p *ProjectivePoint) ProjectivePoint {
	var t0, t1, t2 Fp2Element
	var q ProjectivePoint
	var K1, K2 = &phi.K1, &phi.K2
	var px, pz = &p.X, &p.Z

	Add(&t0, px, pz)   // t0 = XQ + ZQ
	Sub(&t1, px, pz)   // t1 = XQ - ZQ
	Mul(&t0, K1, &t0)  // t2 = K1 * t0
	Mul(&t1, K2, &t1)  // t1 = K2 * t1
	Add(&t2, &t0, &t1) // t2 = t0 + t1
	Sub(&t0, &t1, &t0) // t0 = t1 - t0
	Square(&t2, &t2)   // t2 = t2 ^ 2
	Square(&t0, &t0)   // t0 = t0 ^ 2
	Mul(&q.X, px, &t2) // XQ'= XQ * t2
	Mul(&q.Z, pz, &t0) // ZQ'= ZQ * t0
	return q
}

// Given a four-torsion point p = x(PB) on the curve E_(A:C), construct the
// four-isogeny phi : E_(A:C) -> E_(A:C)/<P_4> = E_(A':C').
//
// Input: (XP_4: ZP_4), where P_4 has exact order 4 on E_A/C
// Output: * Curve coordinates (A' + 2C', 4C') corresponding to E_A'/C' = A_E/C/<P4>
//         * isogeny phi with constants in F_p^2
func (phi *isogeny4) GenerateCurve(p *ProjectivePoint) CurveCoefficientsEquiv {
	var coefEq CurveCoefficientsEquiv
	var xp4, zp4 = &p.X, &p.Z
	var K1, K2, K3 = &phi.K1, &phi.K2, &phi.K3

	Sub(K2, xp4, zp4)
	Add(K3, xp4, zp4)
	Square(K1, zp4)
	Add(K1, K1, K1)
	Square(&coefEq.C, K1)
	Add(K1, K1, K1)
	Square(&coefEq.A, xp4)
	Add(&coefEq.A, &coefEq.A, &coefEq.A)
	Square(&coefEq.A, &coefEq.A)
	return coefEq
}

// Given a 4-isogeny phi and a point xP = x(P), compute x(Q), the x-coordinate
// of the image Q = phi(P) of P under phi : E_(A:C) -> E_(A':C').
//
// Input: isogeny returned by GenerateCurve and point q=(Qx,Qz) from E0_A/C
// Output: Corresponding point q from E1_A'/C', where E1 is 4-isogenous to E0
func (phi *isogeny4) EvaluatePoint(p *ProjectivePoint) ProjectivePoint {
	var t0, t1 Fp2Element
	var q = *p
	var xq, zq = &q.X, &q.Z
	var K1, K2, K3 = &phi.K1, &phi.K2, &phi.K3

	Add(&t0, xq, zq)
	Sub(&t1, xq, zq)
	Mul(xq, &t0, K2)
	Mul(zq, &t1, K3)
	Mul(&t0, &t0, &t1)
	Mul(&t0, &t0, K1)
	Add(&t1, xq, zq)
	Sub(zq, xq, zq)
	Square(&t1, &t1)
	Square(zq, zq)
	Add(xq, &t0, &t1)
	Sub(&t0, zq, &t0)
	Mul(xq, xq, &t1)
	Mul(zq, zq, &t0)
	return q
}

/* -------------------------------------------------------------------------
   Utils
   -------------------------------------------------------------------------*/
func (point *ProjectivePoint) ToAffine(c *CurveOperations) *Fp2Element {
	var affine_x Fp2Element
	Inv(&affine_x, &point.Z)
	Mul(&affine_x, &affine_x, &point.X)
	return &affine_x
}
