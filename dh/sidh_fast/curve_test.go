package sike

import (
	"bytes"
	"testing"
	"testing/quick"
)

func TestOne(t *testing.T) {
	var tmp Fp2Element

	Mul(&tmp, &P503_OneFp2, &affine_xP)
	if !VartimeEqFp2(&tmp, &affine_xP) {
		t.Error("Not equal 1")
	}
}

// This test is here only to ensure that ScalarMult helper works correctly
func TestScalarMultVersusSage(t *testing.T) {
	var xP ProjectivePoint

	xP = ProjectivePoint{X: affine_xP, Z: P503_OneFp2}
	xP = ScalarMult(&curve, &xP, mScalarBytes[:]) // = x([m]P)
	affine_xQ := xP.ToAffine(kCurveOps)
	if !VartimeEqFp2(&affine_xaP, affine_xQ) {
		t.Error("\nExpected\n", affine_xaP, "\nfound\n", affine_xQ)
	}
}

func Test_jInvariant(t *testing.T) {
	var curve = ProjectiveCurveParameters{A: curve_A, C: curve_C}
	var jbufRes = make([]byte, Params(FP_503).SharedSecretSize)
	var jbufExp = make([]byte, Params(FP_503).SharedSecretSize)
	// Computed using Sage
	// j = 3674553797500778604587777859668542828244523188705960771798425843588160903687122861541242595678107095655647237100722594066610650373491179241544334443939077738732728884873568393760629500307797547379838602108296735640313894560419*i + 3127495302417548295242630557836520229396092255080675419212556702820583041296798857582303163183558315662015469648040494128968509467224910895884358424271180055990446576645240058960358037224785786494172548090318531038910933793845
	var known_j = Fp2Element{
		A: FpElement{0x2c441d03b72e27c, 0xf2c6748151dbf84, 0x3a774f6191070e, 0xa7c6212c9c800ba6, 0x23921b5cf09abc27, 0x9e1baefbb3cd4265, 0x8cd6a289f12e10dc, 0x3fa364128cf87e},
		B: FpElement{0xe7497ac2bf6b0596, 0x629ee01ad23bd039, 0x95ee11587a119fa7, 0x572fb28a24772269, 0x3c00410b6c71567e, 0xe681e83a345f8a34, 0x65d21b1d96bd2d52, 0x7889a47e58901},
	}
	kCurveOps.Jinvariant(&curve, jbufRes[:])
	kCurveOps.Fp2ToBytes(jbufExp[:], &known_j)

	if !bytes.Equal(jbufRes[:], jbufExp[:]) {
		t.Error("Computed incorrect j-invariant: found\n", jbufRes, "\nexpected\n", jbufExp)
	}
}

func TestProjectivePointVartimeEq(t *testing.T) {
	var xP ProjectivePoint

	xP = ProjectivePoint{X: affine_xP, Z: P503_OneFp2}
	xQ := xP
	// Scale xQ, which results in the same projective point
	Mul(&xQ.X, &xQ.X, &curve_A)
	Mul(&xQ.Z, &xQ.Z, &curve_A)
	if !VartimeEqProjFp2(&xP, &xQ) {
		t.Error("Expected the scaled point to be equal to the original")
	}
}

func TestPointDoubleVersusSage(t *testing.T) {
	var curve = ProjectiveCurveParameters{A: curve_A, C: curve_C}
	var params = kCurveOps.CalcCurveParamsEquiv4(&curve)
	var xP ProjectivePoint

	xP = ProjectivePoint{X: affine_xP, Z: P503_OneFp2}
	kCurveOps.Pow2k(&xP, &params, 1)
	affine_xQ := xP.ToAffine(kCurveOps)
	if !VartimeEqFp2(affine_xQ, &affine_xP2) {
		t.Error("\nExpected\n", affine_xP2, "\nfound\n", affine_xQ)
	}
}

func TestPointMul4VersusSage(t *testing.T) {
	var params = kCurveOps.CalcCurveParamsEquiv4(&curve)
	var xP ProjectivePoint

	xP = ProjectivePoint{X: affine_xP, Z: P503_OneFp2}
	kCurveOps.Pow2k(&xP, &params, 2)
	affine_xQ := xP.ToAffine(kCurveOps)
	if !VartimeEqFp2(affine_xQ, &affine_xP4) {
		t.Error("\nExpected\n", affine_xP4, "\nfound\n", affine_xQ)
	}
}

func TestPointMul9VersusSage(t *testing.T) {
	var params = kCurveOps.CalcCurveParamsEquiv3(&curve)
	var xP ProjectivePoint

	xP = ProjectivePoint{X: affine_xP, Z: P503_OneFp2}
	kCurveOps.Pow3k(&xP, &params, 2)
	affine_xQ := xP.ToAffine(kCurveOps)
	if !VartimeEqFp2(affine_xQ, &affine_xP9) {
		t.Error("\nExpected\n", affine_xP9, "\nfound\n", affine_xQ)
	}
}

func TestPointPow2kVersusScalarMult(t *testing.T) {
	var xP, xQ, xR ProjectivePoint
	var params = kCurveOps.CalcCurveParamsEquiv4(&curve)

	xP = ProjectivePoint{X: affine_xP, Z: P503_OneFp2}
	xQ = xP
	kCurveOps.Pow2k(&xQ, &params, 5)
	xR = ScalarMult(&curve, &xP, []byte{32})
	affine_xQ := xQ.ToAffine(kCurveOps) // = x([32]P)
	affine_xR := xR.ToAffine(kCurveOps) // = x([32]P)

	if !VartimeEqFp2(affine_xQ, affine_xR) {
		t.Error("\nExpected\n", affine_xQ, "\nfound\n", affine_xR)
	}
}

func TestPointTripleVersusAddDouble(t *testing.T) {
	tripleEqualsAddDouble := func(params GeneratedTestParams) bool {
		var P2, P3, P2plusP ProjectivePoint

		eqivParams4 := kCurveOps.CalcCurveParamsEquiv4(&params.Cparam)
		eqivParams3 := kCurveOps.CalcCurveParamsEquiv3(&params.Cparam)
		P2 = params.Point
		P3 = params.Point
		kCurveOps.Pow2k(&P2, &eqivParams4, 1)                   // = x([2]P)
		kCurveOps.Pow3k(&P3, &eqivParams3, 1)                   // = x([3]P)
		P2plusP = AddProjFp2(&P2, &params.Point, &params.Point) // = x([2]P + P)
		return VartimeEqProjFp2(&P3, &P2plusP)
	}

	if err := quick.Check(tripleEqualsAddDouble, quickCheckConfig); err != nil {
		t.Error(err)
	}
}

func TestFourIsogenyVersusSage(t *testing.T) {
	var xR, xP4, resPhiXr, expPhiXr ProjectivePoint
	var phi = NewIsogeny4()

	// Given 4-torsion point on E, constructs a four-isogeny phi. Then evaluates
	// point xR on the isogeny and compares with expected result.
	// Values generated with reference implementation.

	xP4 = ProjectivePoint{
		X: Fp2Element{
			A: FpElement{0xD00E20EC32B0EC29, 0xB931E12F6B486A34, 0x96EFFCAEC541E83F, 0x415729120E95D243, 0xB15DADFFBC7986EB, 0x27E7697979D482AC, 0xB269E255C3C11421, 0x35DFF53EF1BDE4},
			B: FpElement{0x691F8D69E98BBF40, 0xAB3894C2F436C73, 0x6CE884D45D785C50, 0xBCDE642D761476C0, 0x8023EF2FCF4C9506, 0x6E3914CFCA94C8A0, 0x8AFE4F1C54EB8744, 0x163227D8890C30},
		},
		Z: Fp2Element{
			A: FpElement{0x38B48A08355253FD, 0x457A001F6F522A58, 0x1959F48231C94070, 0xDF3B4C55A3FF1202, 0x3835E8FB47E9B93F, 0x84320E41E65889B5, 0x6D4AA6D38651BE7E, 0xF50448746FF64},
			B: FpElement{0xEBBCCCBB347E448C, 0xFBC721B5DB2103C9, 0x54FD31DF0C538F18, 0xDE7B3C6CBB60C5BD, 0x86B664DCF5F4B272, 0x705CFC301B13DCD6, 0xFD250579C9257778, 0x366F73666C6C92},
		},
	}
	xR = ProjectivePoint{
		X: Fp2Element{
			A: FpElement{0x6F50E690932A1473, 0x3EC8EE10B576C790, 0x5CABB067D0648B46, 0x77EA840A4219753C, 0xBFEE6EAB2073A69A, 0x845637223AB3687B, 0x20294B44CBDC9F34, 0x59C58391A33D5C},
			B: FpElement{0x68832275EA18BDDC, 0x90039FCD522B6CCF, 0x43A97285E71B711A, 0xBCBFC2C3BCCF6135, 0xDE13C2E410DCF1FE, 0xB9B1243C7E4FC515, 0x3CE1C024813A61D, 0x2BED536959B2D},
		},
		Z: Fp2Element{
			A: FpElement{0x99C27A12675FD4CD, 0x856E300D657ADDE3, 0x156C170BB8983CD3, 0x6A366F8BA2FD7805, 0xE922609C4B80E4A4, 0xAC5A1D2EBE7F2A9A, 0x2E732DAF59AE4A03, 0x6AC91B99882D54},
			B: FpElement{0x909A822C8536612D, 0xBF579BF499C34C2D, 0xE2FAD61D94E1E60F, 0x37CB4E1F0A819D5F, 0xDBD36EA4FC053430, 0x28F3805ECA4730D8, 0x33F47EAF9ED8CEA2, 0x24FC2437192954},
		},
	}
	expPhiXr = ProjectivePoint{
		X: Fp2Element{
			A: FpElement{0x2E2D7C96BB057AE9, 0x58FF5432A90EA157, 0x6EED2543FED809C7, 0xF721E3657B17C6D3, 0xC9F8EBED3E1430AF, 0xA94DAFEC2ED7275A, 0xFC8A869CF993A64D, 0x45C8B4291BC602},
			B: FpElement{0xD5730CA5DA535196, 0x958D80511DCD695F, 0xCFDCAA016F0D6AF, 0x176FAA4414FC230B, 0x61A5CDD045B67365, 0x13AC43A5E7F0E446, 0x7BCABE9E555C2729, 0x2CA6A01B26BFEB},
		},
		Z: Fp2Element{
			A: FpElement{0x684A5999FCD11607, 0x3D0057EA6B62FC92, 0x692895B2D37F8EAA, 0xF0BB08106CCF7FDF, 0x3A521D25A431C5CF, 0x8F8DCB43E0BD2475, 0x9CF6266E32D712D3, 0x3B98B6D5C0B377},
			B: FpElement{0x8F4E4EA61ACA375, 0xE8DF168DA6349D03, 0x8DFD68ABA4AB08CC, 0x5352A227C5C6D59C, 0x45750EB03218D4D6, 0x71E2AD1F130DB05E, 0x64F35BBA642804EC, 0x26542493BF5F1C},
		},
	}

	phi.GenerateCurve(&xP4)
	resPhiXr = phi.EvaluatePoint(&xR)
	if !VartimeEqProjFp2(&expPhiXr, &resPhiXr) {
		t.Error("\nExpected\n", expPhiXr.ToAffine(kCurveOps), "\nfound\n", resPhiXr.ToAffine(kCurveOps))
	}
}

func TestThreeIsogenyVersusSage(t *testing.T) {
	var xR, xP3, resPhiXr, expPhiXr ProjectivePoint
	var phi = NewIsogeny3()

	// Given 3-torsion point on E, constructs a three-isogeny phi. Then evaluates
	// point xR on the isogeny and compares with expected result.
	// Values generated with reference implementation.

	xP3 = ProjectivePoint{
		X: Fp2Element{
			A: FpElement{0x43C64B1158DE7ED, 0xC522F8AB7DCC9247, 0xC5BFCC8EA95E9F4D, 0xA6DFCE67C53F63BC, 0x9C6329D65EBBBE44, 0x91949F2E9864BD5A, 0xC9AE7B8B0435B0AF, 0x1607E735F9E10},
			B: FpElement{0x3EEFA1A8B5D59CD9, 0x1ED7264A82282E14, 0x253309D0531054E1, 0x7557CC9966B63AB1, 0xAAB3B77A0CF3D9C, 0xF9BE0DC1977358F4, 0xC5B7AE198CF22874, 0x3250464B34AAD1},
		},
		Z: Fp2Element{
			A: FpElement{0xC06D37BCBBD98418, 0x1C7C9E348A880023, 0xB1F61CA46EA815FD, 0x7E0E5F01EAB9D7B6, 0xE8737A5EF457E188, 0xBD228FDA0BAF18D8, 0xAB7823AF7BAFD785, 0x2BCA7CCFFC1DDA},
			B: FpElement{0xBC34D39B7CBF3EDC, 0x882C3AFC4011C8E8, 0x68A2D74B0FBA196E, 0x810E59E7DD937844, 0xE796B5D4BFC3982F, 0xC7D23388B8E91883, 0x552B783D3986109F, 0x1337962318DFC0},
		},
	}

	xR = ProjectivePoint{
		X: Fp2Element{
			A: FpElement{0xBA300F2F1C010480, 0xE196ACEE08FEA8BA, 0xC1B8AB47C5D6D9A, 0x2CDFF1E375E5FAFC, 0x2D55CBA6472829AF, 0xE03ECA628015CA0E, 0xD1055B779C2DCC6C, 0x7F915F81CAF33},
			B: FpElement{0x5179F86B4F63CA45, 0x8CF33AD2D0D7E409, 0xE9065B70EB5F8680, 0xFBA717809FF35FE, 0x8E31E6EF3CAD154C, 0x65907A2B38A0B673, 0x9E5A4FFCF1F7E74, 0x3170F0C18D5F96},
		},
		Z: Fp2Element{
			A: FpElement{0x1F48F3A2DFB1C73B, 0x3E35C8CD0752F9A4, 0x88601205D0B6B8C0, 0xCFF48E40A9C200CD, 0x10E6964543C6195C, 0x6B8F141796914E13, 0xA7B5F96629DF495E, 0x6600DB36C90874},
			B: FpElement{0xAB54D5B8247FE6CA, 0xD5EE5EAE7C19E9B4, 0x16CB352BA75CB7EF, 0x6D651A77FEB51C5E, 0x2D72F65AC9D39E8A, 0xE10F942CEAD9C7EA, 0x36A5A27BE681CE7A, 0x1C500AA0D9A62F},
		},
	}

	expPhiXr = ProjectivePoint{
		X: Fp2Element{
			A: FpElement{0x61B04752330F7FFF, 0x67F7FADAE5093E06, 0xB665F1E8F70118C6, 0x4F529F9BB30AE6A0, 0xD38E0FC09717C6D1, 0xB7886970ADE8584B, 0x73D66E118BAA193C, 0x4604C634755CFB},
			B: FpElement{0x65CCBE0938AB5A99, 0x1F23B14E1548E3BC, 0x2A565624008051D8, 0xC45D118553BEA2E5, 0x7E2C027737E386EA, 0xF8EC1668C4C09CFB, 0x24CBE8F9D424021D, 0x62E99144A24A6},
		},
		Z: Fp2Element{
			A: FpElement{0x71D9A198BB845CCA, 0xB2D0A8D2168F4399, 0x9C85368AF08AC7E1, 0x76D71A16B7F4B966, 0x60821CCED03DE7DB, 0x80D404686B651216, 0x8489AF1E2E14BF8E, 0x370781CDE810FE},
			B: FpElement{0xB12EE10B6B80F65B, 0xC4C1CD99C671118D, 0xB84A2C8B2C153F37, 0x9170BAE0CE11B7A8, 0xF38DE8F9AF1BF991, 0x88612A07E7F7015A, 0x9611B2C68B94BC68, 0x5BCFB00EC5DE0},
		},
	}

	phi.GenerateCurve(&xR)
	resPhiXr = phi.EvaluatePoint(&xP3)

	if !VartimeEqProjFp2(&expPhiXr, &resPhiXr) {
		t.Error("\nExpected\n", expPhiXr.ToAffine(kCurveOps), "\nfound\n", resPhiXr.ToAffine(kCurveOps))
	}
}

func BenchmarkThreePointLadder255BitScalar(b *testing.B) {
	var mScalarBytes = [...]uint8{203, 155, 185, 191, 131, 228, 50, 178, 207, 191, 61, 141, 174, 173, 207, 243, 159, 243, 46, 163, 19, 102, 69, 92, 36, 225, 0, 37, 114, 19, 191, 0}
	for n := 0; n < b.N; n++ {
		kCurveOps.ScalarMul3Pt(&curve, &threePointLadderInputs[0], &threePointLadderInputs[1], &threePointLadderInputs[2], 255, mScalarBytes[:])
	}
}
