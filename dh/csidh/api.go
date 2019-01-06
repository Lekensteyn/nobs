package csidh

// Defines operations on public key
type PublicKey struct {
	// Montgomery coefficient: represents y^2 = x^3 + Ax^2 + x
	A Fp
}

// Defines operations on private key
type PrivateKey struct {
	e [38]int8
}

// PrivateKey
func NewPrivateKey() PrivateKey {
	return PrivateKey{}
}

func (c PrivateKey) Import(key []byte) {
}

func (c PrivateKey) Export() []byte {
	return nil
}

// PublicKey
func NewPublicKey() PublicKey {
	return PublicKey{}
}

func (c PublicKey) Import(key []byte) {

}

func (c PublicKey) Export() []byte {
	return nil
}
