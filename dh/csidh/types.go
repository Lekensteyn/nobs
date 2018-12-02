package csidh

// 128-bit number representing prime field element GF(p)
type Fp [8]uint64

// Represents projective point on elliptic curve E over Fp
type Point struct {
	x Fp
	z Fp
}

// Curve coefficients
type Coeff struct {
	a Fp
	c Fp
}
