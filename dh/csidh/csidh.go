package csidh

// Implements differential arithmetic in P^1
// for montgomery curves.

// Implements a mapping: x(P),x(Q),x(P-Q) -> x(P+Q)
// In: P,Q,PdQ
// Out: PaQ
func xAdd(PaQ, P, Q, PdQ *Point) {
	var t0, t1, t2, t3 Fp
	addRdc(&t0, &P.x, &P.z)
	subRdc(&t1, &P.x, &P.z)
	addRdc(&t2, &Q.x, &Q.z)
	subRdc(&t3, &Q.x, &Q.z)
	mulRdc(&t0, &t0, &t3)
	mulRdc(&t1, &t1, &t2)
	addRdc(&t2, &t0, &t1)
	subRdc(&t3, &t0, &t1)
	sqrRdc(&t2, &t2)
	sqrRdc(&t3, &t3)
	mulRdc(&PaQ.x, &PdQ.z, &t2)
	mulRdc(&PaQ.z, &PdQ.x, &t3)
}

func xDbl(Q, P, A *Point) {
	var t0, t1, t2 Fp
	addRdc(&t0, &P.x, &P.z)
	sqrRdc(&t0, &t0)
	subRdc(&t1, &P.x, &P.z)
	sqrRdc(&t1, &t1)
	subRdc(&t2, &t0, &t1)
	mulRdc(&t1, &four, &t1)
	mulRdc(&t1, &t1, &A.z)
	mulRdc(&Q.x, &t0, &t1)
	addRdc(&t0, &A.z, &A.z)
	addRdc(&t0, &t0, &A.x)
	mulRdc(&t0, &t0, &t2)
	addRdc(&t0, &t0, &t1)
	mulRdc(&Q.z, &t0, &t2)
}

// TODO: This can be improved I think (as for SIDH)
// Pap, PaQ
func xDblAdd(PaP, PaQ, P, Q, PdQ *Point, A24 *Coeff) {
	var t0, t1, t2 Fp
	addRdc(&t0, &P.x, &P.z)
	subRdc(&t1, &P.x, &P.z)
	mulRdc(&PaP.x, &t0, &t0)
	subRdc(&t2, &Q.x, &Q.z)
	addRdc(&PaQ.x, &Q.x, &Q.z)
	mulRdc(&t0, &t0, &t2)
	mulRdc(&PaP.z, &t1, &t1)
	mulRdc(&t1, &t1, &PaQ.x)
	subRdc(&t2, &PaP.x, &PaP.z)
	mulRdc(&PaP.z, &PaP.z, &A24.c)
	mulRdc(&PaP.x, &PaP.x, &PaP.z)
	mulRdc(&PaQ.x, &A24.a, &t2)
	subRdc(&PaQ.z, &t0, &t1)
	addRdc(&PaP.z, &PaP.z, &PaQ.x)
	addRdc(&PaQ.x, &t0, &t1)
	mulRdc(&PaP.z, &PaP.z, &t2)
	mulRdc(&PaQ.z, &PaQ.z, &PaQ.z)
	mulRdc(&PaQ.x, &PaQ.x, &PaQ.x)
	mulRdc(&PaQ.z, &PaQ.z, &PdQ.x)
	mulRdc(&PaQ.x, &PaQ.x, &PdQ.z)
}

func cswapPoint(P1, P2 *Point, choice uint8) {
	cswap512(&P1.x, &P2.x, choice)
	cswap512(&P1.z, &P2.z, choice)
}

// kP = [k]P
// TODO: Only one swap should be enough
func xMul512(kP, P *Point, co *Coeff, k *Fp) {
	var A24 Coeff
	kP.x = fp_1
	R := *P
	PdQ := *P

	// Precompyte A24 = (A+2C:4C) => (A24.x = A.x+2A.z; A24.z = 4*A.z)
	addRdc(&A24.a, &co.c, &co.c)
	addRdc(&A24.a, &A24.a, &co.a)
	mulRdc(&A24.c, &co.c, &four)

	for i := uint(512); i > 0; {
		i--
		bit := uint8(k[i>>6] >> (i & 63) & 1)
		cswapPoint(kP, &R, bit)
		xDblAdd(kP, &R, kP, &R, &PdQ, &A24)
		cswapPoint(kP, &R, bit)
	}
}
