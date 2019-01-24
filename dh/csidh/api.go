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

func (c *PrivateKey) Import(key []byte) bool {
	if len(key) < len(c.e) {
		return false
	}
	for i, v := range key {
		c.e[i] = int8(v)
	}
	return true
}

func (c PrivateKey) Export(out []byte) bool {
	if len(out) < len(c.e) {
		return false
	}
	for i, v := range c.e {
		out[i] = byte(v)
	}
	return true
}

func NewPublicKey() PublicKey {
	return PublicKey{}
}

func (c *PublicKey) Import(key []byte) bool {
	if len(key) != numWords*limbByteSize {
		return false
	}
	for i := 0; i < len(key); i++ {
		j := i / limbByteSize
		k := uint64(i % 8)
		c.A[j] |= uint64(key[i]) << (8 * k)
	}
	return true
}

func (c *PublicKey) Export(out []byte) bool {
	if len(out) != numWords*limbByteSize {
		return false
	}
	for i := 0; i < len(out); i++ {
		j := i / limbByteSize
		k := uint64(i % 8)
		out[i] = byte(c.A[j] >> (8 * k))
	}
	return true
}
