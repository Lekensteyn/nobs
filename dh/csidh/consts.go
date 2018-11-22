package csidh

// TODO: This all could simply go to read-only segment
// 	 then any write will crash the program which
//       is actually right thing to do

// P511

var (
	pbits = 511
	p     = u512{
		0x1B81B90533C6C87B, 0xC2721BF457ACA835,
		0x516730CC1F0B4F25, 0xA7AAC6C567F35507,
		0x5AFBFCC69322C9CD, 0xB42D083AEDC88C42,
		0xFC8AB0D15E3E4C4A, 0x65B48E8F740F89BF,
	}

	// OZAPTF: this is useless
	// Zero
	fp_0 = []byte{}

	/* Montgomery R = 2^512 mod p */
	fp_1 = u512{
		0xC8FC8DF598726F0A, 0x7B1BC81750A6AF95,
		0x5D319E67C1E961B4, 0xB0AA7275301955F1,
		0x4A080672D9BA6C64, 0x97A5EF8A246EE77B,
		0x06EA9E5D4383676A, 0x3496E2E117E0EC80,
	}

	/* Montgomery R^2 mod p */
	r_squared_mod_p = u512{
		0x36905B572FFC1724, 0x67086F4525F1F27D,
		0x4FAF3FBFD22370CA, 0x192EA214BCC584B1,
		0x5DAE03EE2F5DE3D0, 0x1E9248731776B371,
		0xAD5F166E20E4F52D, 0x4ED759AEA6F3917E,
	}

	// -p^-1 mod 2^64
	pNegInv = u512{
		0x66c1301f632e294d,
	}

	// Elkies primes up to 374 + prime 587
	// p = 4 * product(primes) - 1
	primes = []uint64{
		3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59,
		61, 67, 71, 73, 79, 83, 89, 97, 101, 103, 107, 109, 113, 127, 131, 137,
		139, 149, 151, 157, 163, 167, 173, 179, 181, 191, 193, 197, 199, 211, 223, 227,
		229, 233, 239, 241, 251, 257, 263, 269, 271, 277, 281, 283, 293, 307, 311, 313,
		317, 331, 337, 347, 349, 353, 359, 367, 373, 587,
	}

//	base = []byte{
//		0,
//	}
)