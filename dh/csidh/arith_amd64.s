// +build amd64,!noasm

#include "textflag.h"

/*
TEXT ·fp512Set(SB), NOSPLIT, $0-8
	CLD
	MOVQ	SI, AX
	STOSQ
	XORQ	AX, AX
	MOVQ	$7, CX
	REP	STOSQ

	RET
*/

// Multipies 512-bit value by 64-bit value. Uses MULX instruction
// x = y * z
//
// func mul512(a, b *u512, c uint64)
TEXT ·mul512(SB), NOSPLIT, $0-24
	MOVQ	x+ 0(FP), DI	// result
	MOVQ	y+ 8(FP), SI	// multiplicand
	MOVQ	z+16(FP), DX	// 64 byte multiplier

	MULXQ	0(SI), AX, R10
	MOVQ	AX, 0(DI)	// x[0]

	MULXQ	8(SI), AX, R11
	ADDQ	R10, AX
	MOVQ	AX, 8(DI)	// x[1]

	MULXQ	16(SI), AX, R10
	ADCQ	R11, AX
	MOVQ	AX, 16(DI)	// x[2]

	MULXQ	24(SI), AX, R11
	ADCQ	R10, AX
	MOVQ	AX, 24(DI)	// x[3]

	MULXQ	32(SI), AX, R10
	ADCQ	R11, AX
	MOVQ	AX, 32(DI)	// x[4]

	MULXQ	40(SI), AX, R11
	ADCQ	R10, AX
	MOVQ	AX, 40(DI)	// x[5]

	MULXQ	48(SI), AX, R10
	ADCQ	R11, AX
	MOVQ	AX, 48(DI)	// x[6]

	MULXQ	56(SI), AX, R11
	ADCQ	R10, AX
	MOVQ	AX, 56(DI)	// x[7]

	RET

// x = y + z
// func add512(x, y, z *u512) uint64
TEXT ·add512(SB), NOSPLIT, $0-32
	MOVQ	x+ 0(FP), DI	// result
	MOVQ	y+ 8(FP), SI	// first summand
	MOVQ	z+16(FP), DX	// second summand

	XORQ	AX, AX

	MOVQ	 0(SI), R8
	ADDQ	 0(DX), R8
	MOVQ	R8,  0(DI)	// x[0]

	MOVQ	 8(SI), R8
	ADCQ	 8(DX), R8
	MOVQ	R8,  8(DI)	// x[1]

	MOVQ	16(SI), R8
	ADCQ	16(DX), R8
	MOVQ	R8, 16(DI)	// x[2]

	MOVQ	24(SI), R8
	ADCQ	24(DX), R8
	MOVQ	R8, 24(DI)	// x[3]

	MOVQ	32(SI), R8
	ADCQ	32(DX), R8
	MOVQ	R8, 32(DI)	// x[4]

	MOVQ	40(SI), R8
	ADCQ	40(DX), R8
	MOVQ	R8, 40(DI)	// x[5]

	MOVQ	48(SI), R8
	ADCQ	48(DX), R8
	MOVQ	R8, 48(DI)	// x[6]

	MOVQ	56(SI), R8
	ADCQ	56(DX), R8
	MOVQ	R8, 56(DI)	// x[7]

	// return carry
	ADCQ	AX, AX
	MOVQ	AX, ret+24(FP)

	RET


// x = y - z
// func sub512(x, y, z *u512) uint64
TEXT ·sub512(SB), NOSPLIT, $0-32
	MOVQ	x+ 0(FP), DI	// result
	MOVQ	y+ 8(FP), SI	// minuend
	MOVQ	z+16(FP), DX	// subtrahend

	XORQ	AX, AX

	MOVQ	 0(SI), R8
	SUBQ	 0(DX), R8
	MOVQ	R8,  0(DI)	// x[0]

	MOVQ	 8(SI), R8
	SBBQ	 8(DX), R8
	MOVQ	R8,  8(DI)	// x[1]

	MOVQ	16(SI), R8
	SBBQ	16(DX), R8
	MOVQ	R8, 16(DI)	// x[2]

	MOVQ	24(SI), R8
	SBBQ	24(DX), R8
	MOVQ	R8, 24(DI)	// x[3]

	MOVQ	32(SI), R8
	SBBQ	32(DX), R8
	MOVQ	R8, 32(DI)	// x[4]

	MOVQ	40(SI), R8
	SBBQ	40(DX), R8
	MOVQ	R8, 40(DI)	// x[5]

	MOVQ	48(SI), R8
	SBBQ	48(DX), R8
	MOVQ	R8, 48(DI)	// x[6]

	MOVQ	56(SI), R8
	SBBQ	56(DX), R8
	MOVQ	R8, 56(DI)	// x[7]

	// return borrow
	ADCQ	AX, AX
	MOVQ	AX, ret+24(FP)

	RET

TEXT ·cswap512(SB),NOSPLIT,$0-17
	MOVQ    x+0(FP), DI
	MOVQ    y+8(FP), SI
    MOVBLZX choice+16(FP), AX       // AL = 0 or 1

	// Make AX, so that either all bits are set or non
	// AX = 0 or 1
	NEGQ    AX

	// Fill xmm15. After this step first half of XMM15 is
	// just zeros and second half is whatever in AX
	MOVQ    AX, X15

	// Copy lower double word everywhere else. So that
	// XMM15=AL|AL|AL|AL. As AX has either all bits set
	// or non result will be that XMM15 has also either
	// all bits set or non of them.
	PSHUFD $0, X15, X15

#ifndef CSWAP_BLOCK
#define CSWAP_BLOCK(idx)       \
	MOVOU   (idx*16)(DI), X0 \
	MOVOU   (idx*16)(SI), X1 \
	\ // X2 = mask & (X0 ^ X1)
	MOVO     X1, X2 \
	PXOR     X0, X2 \
	PAND    X15, X2 \
	\
	PXOR     X2, X0 \
	PXOR     X2, X1 \
	\
	MOVOU    X0, (idx*16)(DI) \
	MOVOU    X1, (idx*16)(SI)
#endif

	CSWAP_BLOCK(0)
	CSWAP_BLOCK(1)
	CSWAP_BLOCK(2)
	CSWAP_BLOCK(3)

	RET
