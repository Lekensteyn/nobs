package csidh

// z = x + y mod P
func addRdc(z, x, y *Fp) {
	add512(z, x, y)
	// TODO: check if doing it in add512 is much faster?
	crdc512(z)
}

func subRdc(z, x, y *Fp) {
	borrow := sub512(z, x, y)
	csubrdc512(z, borrow)
}

func mulRdc(z, x, y *Fp) {
	mul(z, x, y)
	crdc512(z)
}

func sqrRdc(z, x *Fp) {
	// TODO: to be implemented faster
	mul(z, x, x)
	crdc512(z)
}

// 4 bit fixed-window exponentiation
// TODO: currently mul is interleaved with reduction step. It would be better to
// do reduction only once at the end (like in Barrett's mul)
func modExp(res, base, exp *Fp) {
	var precomp [16]Fp

	// Precompute step, computes an array of small powers of 'base'. As this
	// algorithm implements 4-bit window, we need 2^4=16 of such values.
	// base^0 = 1, which is equal to R from REDC.
	precomp[0] = fp_1
	precomp[1] = *base
	for i := 2; i < 16; i = i + 2 {
		// Interleave fast squaring with multiplication. It's currently not a case
		// but squaring can be implemented faster than multiplication.
		sqrRdc(&precomp[i], &precomp[i/2])
		mulRdc(&precomp[i+1], &precomp[i], base)
	}

	*res = fp_1
	for i := int(127); i >= 0; i-- {
		for j := 0; j < 4; j++ {
			mulRdc(res, res, res)
		}
		// TODO: not constant time
		idx := (exp[i/16] >> uint((i%16)*4)) & 15
		mulRdc(res, res, &precomp[idx])
	}
	crdc512(res)
}
