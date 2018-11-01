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
TEXT ·fp512Mul3(SB), NOSPLIT, $0-24
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
