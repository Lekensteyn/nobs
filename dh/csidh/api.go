package csidh

// Defines operations on public key
type PublicKey struct {
	key
	affine_xP   Fp2Element
	affine_xQ   Fp2Element
	affine_xQmP Fp2Element
}

// Defines operations on private key
type PrivateKey struct {
	key
	// Secret key
	Scalar []byte
	// Used only by KEM
	S []byte
}

// PrivateKey
func NewPrivateKey() PrivateKey {

}

func (c PrivateKey) Import(key []byte) {

}

func (c PrivateKey) Export() []byte {

}

// PublicKey
func NewPublicKey() PublicKey {

}

func (c PublicKey) Import(key []byte) {

}

func (c PublicKey) Export() []byte {

}