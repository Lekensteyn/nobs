package csidh

// 128-bit number in Montgomery domain
type Fp struct {
	// TODO: 512 is maybe not best name
	v u512
}

// res = x + y mod P
func addRdc(res, x, y *Fp) {
	add512(&res.v, &x.v, &y.v)
	// TODO: check if doing it in add512 is much faster?
	crdc512(&res.v)
}

func subRdc(res, x, y *Fp) {
	borrow := sub512(&res.v, &x.v, &y.v)
	csubrdc512(&res.v, borrow)
}

func mulRdc(res, x, y *Fp) {
	mul(&res.v, &x.v, &y.v)
	crdc512(&res.v)
}
