package csidh

// Defines operations on public key
type PublicKey struct {
	// Montgomery coefficient: represents y^2 = x^3 + Ax^2 + x
	A Fp
}

// Defines operations on private key
type PrivateKey struct {
	e [37]int8
}

// PrivateKey
func NewPrivateKey() PrivateKey {
	return PrivateKey{}
}

func (c PrivateKey) Import(key []byte) bool {
    if len(key) < len(c.e) {
        return false
    }
    for i,v := range key {
        c.e[i] = int8(v)
    }
    return true
}

func (c PrivateKey) Export(out []byte) bool {
    if len(out) < len(c.e) {
        return false
    }
    for i,v := range c.e {
        out[i] = byte(v)
    }
    return true
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
