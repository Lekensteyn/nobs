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

func square_multiply(x, y *Fp, exp uint64) {
	var res1 = r_squared_mod_p
	var res2 = r_squared_mod_p

	for i := exp; i != 0; i >>= 1 {
		// TODO: that's not constant time
		if (exp & 1) == 1 {
			mulRdc(&res1, &res1, x)
			mulRdc(&res2, &res2, y)
		}
		mulRdc(x, x, x)
		mulRdc(y, y, y)
	}
	for i := 0; i < len(*x); i++ {
		(*x)[i] = res1[i]
		(*y)[i] = res2[i]
	}
}

func MapPoint(img *Point, co *Coeff, kern *Point, order uint64) {
	var t0, t1, t2, S, D Fp
	var Q, prod Point
	var Aed Coeff

	// Compute twisted Edwards coefficients
	addRdc(&Aed.c, &co.c, &co.c)
	addRdc(&Aed.a, &co.a, &co.c)
	subRdc(&Aed.c, &co.a, &Aed.c)

	addRdc(&S, &img.x, &img.z)
	subRdc(&D, &img.x, &img.z)

	subRdc(&prod.x, &kern.x, &kern.z)
	addRdc(&prod.z, &kern.x, &kern.z)

	mulRdc(&t1, &prod.x, &S)
	mulRdc(&t0, &prod.z, &D)
	addRdc(&Q.x, &t0, &t1)
	subRdc(&Q.z, &t0, &t1)

	var M [3]Point
	var coTmp = Point{x: co.a, z: co.c} // OZAPTF: crap
	xDbl(&M[1], &coTmp, kern)

	// TODO: Not constant time
	for i := uint64(1); i < uint64(order/2); i++ {
		// TODO: Not constant time
		if i >= 2 {
			xAdd(&M[i%3], &M[(i-1)%3], kern, &M[(i-2)%3])
		}
		subRdc(&t1, &M[i%3].x, &M[i%3].z)
		addRdc(&t0, &M[i%3].x, &M[i%3].z)
		mulRdc(&prod.x, &prod.x, &t1)
		mulRdc(&prod.z, &prod.z, &t0)
		mulRdc(&t1, &t1, &S)
		mulRdc(&t0, &t0, &D)
		addRdc(&t2, &t0, &t1)
		mulRdc(&Q.x, &Q.x, &t2)
		subRdc(&t2, &t0, &t1)
		mulRdc(&Q.z, &Q.z, &t2)
	}

	mulRdc(&Q.x, &Q.x, &Q.x)
	mulRdc(&Q.z, &Q.z, &Q.z)
	mulRdc(&img.x, &img.x, &Q.x)
	mulRdc(&img.z, &img.z, &Q.z)

	// Aed.a^order and Aed.c^order
	square_multiply(&Aed.a, &Aed.c, order)

	// prod^8
	mulRdc(&prod.x, &prod.x, &prod.x)
	mulRdc(&prod.x, &prod.x, &prod.x)
	mulRdc(&prod.x, &prod.x, &prod.x)
	mulRdc(&prod.z, &prod.z, &prod.z)
	mulRdc(&prod.z, &prod.z, &prod.z)
	mulRdc(&prod.z, &prod.z, &prod.z)

	// compute image curve params
	mulRdc(&Aed.c, &Aed.c, &prod.x)
	mulRdc(&Aed.a, &Aed.a, &prod.z)

	// convert back to Montgomery
	addRdc(&co.a, &Aed.a, &Aed.c)
	subRdc(&co.c, &Aed.a, &Aed.c)
	addRdc(&co.a, &co.a, &co.a)
}