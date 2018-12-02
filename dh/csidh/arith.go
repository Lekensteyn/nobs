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
	mul(z, x, x)
	crdc512(z)
}
