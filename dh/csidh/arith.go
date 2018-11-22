package csidh

// res = x + y mod P
func addRdc(res, x, y *Fp) {
	add512(res, x, y)
	// TODO: check if doing it in add512 is much faster?
	crdc512(res)
}

func subRdc(res, x, y *Fp) {
	borrow := sub512(res, x, y)
	csubrdc512(res, borrow)
}

func mulRdc(res, x, y *Fp) {
	mul(res, x, y)
	crdc512(res)
}

func sqrRdc(res, x *Fp) {
	mul(res, x, x)
	crdc512(res)
}
