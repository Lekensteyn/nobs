package csidh

import "io"
import "crypto/rand"

// TODO: this is weird. How do I know loop will end?
func randFp(fp *Fp) {
	mask := uint64(1<<(pbits%limbBitSize)) - 1
	for {
		*fp = Fp{}
		var buf [len(fp) * limbByteSize]byte
		if _, err := io.ReadFull(rand.Reader, buf[:]); err != nil {
			// OZAPTF: to be re-done (AES_CTR)
			panic("Can't read random number")
		}

		for i := 0; i < len(buf); i++ {
			j := i / limbByteSize
			k := uint(i % 8)
			fp[j] |= uint64(buf[i]) << (8 * k)
		}

		fp[len(fp)-1] &= mask
		if checkBigger(&p, fp) {
			return
		}
	}
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

// Assumes lower<upper
// TODO: non constant time
// TODO: this needs to be rewritten - function called recursivelly
/* compute [(p+1)/l] P for all l in our list of primes. */
/* divide and conquer is much faster than doing it naively,
 * but uses more memory. */
func cofactorMultiples(P []Point, A *Coeff, lower, upper uint64) {
	if upper-lower == 1 {
		return
	}

	// TODO: double check
	var mid = lower + ((upper - lower + 1) >> 1)
	var cl = Fp{1}
	var cu = Fp{1}

	// TODO: one loop would be OK?
	for i := lower; i < mid; i++ {
		mul512(&cu, &cu, primes[i])
	}
	for i := mid; i < upper; i++ {
		mul512(&cl, &cl, primes[i])
	}

	xMul512(&P[mid], &P[lower], A, &cu)
	xMul512(&P[lower], &P[lower], A, &cl)

	// TODO: this function call is not needed
	cofactorMultiples(P, A, lower, mid)
	cofactorMultiples(P, A, mid, upper)
}

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
				c.e[i>>1] |= int8((buf[j] & 0xf) << uint((i%2)*4))
				i = i + 1
				if i == len(primes) {
					break
				}
			}
		}
	}
	return nil
}

// Key validation
// OZAPTF: To be checked
func (c *PublicKey) validate() bool {
	var A = Coeff{a: c.A, c: fp_1}
	var zero [8]uint64

	// TODO: make sure curve is nonsingular
	// -- this needs to be tested before implementing

	// P must have big enough order to prove supersingularity. The
	// probability that this loop will be repeated is negligible.
	// TODO: do max 2 loops
	for {
		// OZAPTF: heap?
		var P [kPrimeCount]Point
		randFp(&P[0].x)
		P[0].z = fp_1

		/* maximal 2-power in p+1 */
		// OZAPTF
		var t = Point{x: A.a, z: A.c}
		xDbl(&P[0], &P[0], &t)
		xDbl(&P[0], &P[0], &t)

		// OZAPTF:that's wrong?
		A.a = t.x
		A.c = t.z

		// TODO: this can be mixed with loop below
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

				// Checks wether t>4*sqrt(p)
				if sub512(&t, &fourSqrtP, &order) != 0 {
					return true
				}
			}
		}
	}
}

func (c *PublicKey) groupAction(pub *PublicKey, prv *PrivateKey) {
	var k [2]Fp
	var e [2][kPrimeCount]uint8
	var done = [2]bool{false, false}
	var A = Coeff{a: pub.A, c: fp_1}
	var zero [8]uint64

	k[0][0] = 4
	k[1][0] = 4

	for i, v := range primes {
		t := int8((prv.e[uint(i)>>1] << ((uint(i) % 2) * 4)) >> 4)
		if t > 0 {
			e[0][i] = uint8(t)
			e[1][i] = 0
			mul512(&k[1], &k[1], v)
		} else if t < 0 {
			e[1][i] = uint8(-t)
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
		randFp(&P.x)
		P.z = fp_1
		montEval(&rhs, &A.a, &P.x)
		sign := isNonQuadRes(&rhs)

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
				if ctEq64(K.z[:], zero[:]) == 0 {
					MapPoint(&P, &A, &K, v)
					e[sign][i] = e[sign][i] - 1
					if e[sign][i] == 0 {
						mul512(&k[sign], &k[sign], primes[i])
					}
				}
			}
			done[sign] = done[sign] && (e[sign][i] == 0)
		}

		modExpRdc(&A.c, &A.c, &pMin2)
		mulRdc(&A.a, &A.a, &A.c)
		A.c = fp_1

		if done[0] && done[1] {
			break
		}
	}
	c.A = A.a
}

func (c *PublicKey) Generate(prv *PrivateKey) {
	var emptyKey PublicKey
	c.groupAction(&emptyKey, prv)
}

// todo: probably should be similar to some other interface
// OZAPTF: should be attribute of private key
func (c *PublicKey) DeriveSecret(out []byte, pub *PublicKey, prv *PrivateKey) bool {
	var ss PublicKey
	// TODO: validation doesn't work yet correctly
	//	if !pub.validate() {
	//		randFp(&pub.A)
	//		return false
	//	}
	ss.groupAction(pub, prv)
	ss.Export(out)
	return true
}

// TODO:
func init() {
	if len(primes) != kPrimeCount {
		panic("Wrong number of primes")
	}
}
