package csidh

import "io"

func (c *PrivateKey) Generate(rand io.Reader) error {
	for i, _ := range c.e {
		c.e[i] = 0
	}

	for i := 0; i < len(primes); {
		var buf [64]byte
		_, err := io.ReadFull(rand, buf[:])
		if err != nil {
			return err
		}

		for j, _ := range buf {
			if int8(buf[j]) <= expMax && int8(buf[j]) >= -expMax {
				c.e[i>>1] |= int8((buf[j] & 0xf) << uint(i%2*4))
				i = i + 1
				if i == len(primes) {
					break
				}
			}
		}
	}
	return nil
}

// Assumes lower<upper
// TODO: non constant time
// TODO: this needs to be rewritten - function called recursivelly
/* compute [(p+1)/l] P for all l in our list of primes. */
/* divide and conquer is much faster than doing it naively,
 * but uses more memory. */
func cofactorMultiples(P []Point, A *Coeff, lower, upper uint64) {
	// OZAPTF: Needed?
	if upper-lower == 1 {
		return
	}

	// TODO: double check
	var mid = lower + ((upper - lower + 1) >> 1)
	var cl = Fp{1}
	var cu = Fp{1}

	// TODO: one loop would be OK
	for i := lower; i < mid; i++ {
		mul512(&cu, &cu, primes[i])
	}
	for i := mid; i < upper; i++ {
		mul512(&cl, &cl, primes[i])
	}

	xMul512(&P[mid], &P[lower], A, &cu)
	xMul512(&P[lower], &P[lower], A, &cl)

	cofactorMultiples(P, A, lower, mid)
	cofactorMultiples(P, A, mid, upper)
}

// evaluates x^3 + Ax^2 + x
func montEval(res, A, x *Fp) {
	var t Fp

	*res = *x
	mulRdc(res, res, res)
	mulRdc(&t, A, x)
	addRdc(res, res, &t)
	addRdc(res, res, &fp_1)
	mulRdc(res, res, x)
}

// assumes len(x) == len(y)
// return 1 if equal 0 if not
// OZAPTF: I actually need to know if x is zero
func ctEq64(x, y []uint64) uint {
	var t uint64
	var h, l uint64
	for i := 0; i < len(x); i++ {
		t |= x[i] ^ y[i]
	}

	h = ((t >> 32) - 1) >> 63
	l = ((t & 0xFFFFFFFF) - 1) >> 63
	return uint(h & l & 1)
}

// Key validation
func (c *PublicKey) Validate() bool {
	var A = Coeff{a: c.A, c: fp_1}
	var zero [8]uint64

	// TODO: how long it will loop?
	for {
		// OZAPTF: heap?
		var P [kPrimeCount]Point
		/* TODO: fp_random(P.x) to port */
		P[0].z = fp_1

		/* maximal 2-power in p+1 */
		// OZAPTF
		var t = Point{x: A.a, z: A.c}
		xDbl(&P[0], &P[0], &t)
		xDbl(&P[0], &P[0], &t)
		A.a = t.x
		A.c = t.z

		cofactorMultiples(P[:], &A, 0, kPrimeCount)
		var order = Fp{1}

		for i := kPrimeCount - 1; i >= 0; i-- {
			if ctEq64(P[i].z[:], zero[:]) == 1 {
				var t Fp
				t[0] = primes[i]
				xMul512(&P[i], &P[i], &A, &t)
				if ctEq64(P[i].z[:], zero[:]) != 1 {
					return false
				}
				mul512(&order, &order, primes[i])

				if sub512(&t, &fourSqrtP, &order) != 0 {
					return true
				}
			}
		}
	}
}

func (c *PublicKey) Action(pub *PublicKey, prv *PrivateKey) {
	var k [2]Fp
	var e [2][kPrimeCount]uint8
	var done = [2]bool{false, false}
	var A = Coeff{a: pub.A, c: fp_1}
	var zero [8]uint64

	k[0][0] = 4
	k[1][0] = 4

	for i, v := range primes {
		t := int8((prv.e[uint(i)>>1] << (uint(i) % 2) * 4) >> 4)
		if t > 0 {
			e[0][i] = uint8(t)
			e[1][i] = 0
			mul512(&k[1], &k[1], v)
		} else if t < 0 {
			e[1][i] = uint8(-t) // OZAPTF: OK?
			e[0][i] = 0
			mul512(&k[0], &k[0], v)
		} else {
			e[0][i] = 0
			e[1][i] = 0
			mul512(&k[0], &k[0], v)
			mul512(&k[1], &k[1], v)
		}
	}

	for {
		var P Point
		var rhs Fp
		// Randomize P.x
		P.z = fp_1
		montEval(&rhs, &A.a, &P.x)
		var sign = isNotSqr(&rhs)

		if done[sign] {
			continue
		}

		xMul512(&P, &P, &A, &k[sign])
		done[sign] = true

		for i, v := range primes {
			if e[sign][i] != 0 {
				var cof = Fp{1}
				var K Point

				for j := i + 1; j < len(primes); j++ {
					if e[sign][j] != 0 {
						mul512(&cof, &cof, primes[j])
					}
				}

				xMul512(&K, &P, &A, &cof)
				if ctEq64(K.z[:], zero[:]) == 1 {
					MapPoint(&P, &A, &K, v)
					e[sign][i] = e[sign][i] - 1
					if e[sign][i] == 0 {
						mul512(&k[sign], &k[sign], primes[i])
					}
				}
			}
			done[sign] = done[sign] && (e[sign][i] == 0)
		}

		modExp(&A.c, &A.c, &pMin2)
		mulRdc(&A.a, &A.a, &A.c)
		A.c = fp_1
		if done[0] && done[1] {
			break
		}
	}
	c.A = A.a
}

// todo: probably should be similar to some other interface
func (c *PublicKey) csidh(pub *PublicKey, prv *PrivateKey) bool {
	if !c.Validate() {
		// TODO: randomize out->A
		return false
	}
	c.Action(pub, prv)
	return true
}

// TODO:
func init() {
	if len(primes) != kPrimeCount {
		panic("Wrong number of primes")
	}
}
