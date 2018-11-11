package csidh

// 128-bit number in Montgomery domain
type u512mont struct {
}

// u512 specific functions

// go:noescape
func mul512(a, b *u512, c uint64)

// go:noescape
func add512(x, y, z *u512) uint64

// go:noescape
func sub512(x, y, z *u512) uint64

// go:noescape
func cswap512(x, y *u512, choice uint8)

// go:noescape
//func fp_set(x u512mont, y uint64)

//
// //go:noescape
// func fp_cswap(u512mont *x, fp *y, uint8_t c)
//
// //go:noescape
// func fp_enc(u512mont *x, u512 const *y) /* encode to Montgomery representation */
//
// //go:noescape
// func fp_dec(u512 *x, fp const *y) /* decode from Montgomery representation */
//
// //go:noescape
// func fp_add2(u512mont *x, fp const *y)
//
// //go:noescape
// func fp_sub2(u512mont *x, fp const *y);
//
// //go:noescape
// func fp_mul2(u512mont *x, fp const *y);
//
// //go:noescape
// func fp_add3(u512mont *x, fp const *y, fp const *z);
//
// //go:noescape
// func fp_sub3(u512mont *x, fp const *y, fp const *z);
//
// //go:noescape
// func fp_mul3(u512mont *x, fp const *y, fp const *z);
//
// //go:noescape
// func fp_sq1(u512mont *x);
//
// //go:noescape
// func fp_sq2(u512mont *x, fp const *y);
//
// //go:noescape
// func fp_inv(u512mont *x);
//
// //go:noescape
// bool fp_issquare(fp const *x);
//
// void fp_random(u512mont *x);
//
