package h265reader

// NalUnitType is the type of a NAL
type NalUnitType uint8

// Enums for NalUnitTypes
// ref: https://programmer.help/blogs/5c647da7dc2d1.html
const (
	NalUnitTypeCodedSliceTrailN NalUnitType = 0 //
	NalUnitTypeCodedSliceTrailR NalUnitType = 1 //
	NalUnitTypeCodedSliceTsaN   NalUnitType = 2 //
	NalUnitTypeCodedSliceTla    NalUnitType = 3 //
	NalUnitTypeCodedSliceStsaN  NalUnitType = 4 //
	NalUnitTypeCodedSliceStsaR  NalUnitType = 5 //
	NalUnitTypeCodedSliceRadlN  NalUnitType = 6 //
	NalUnitTypeCodedSliceDlp    NalUnitType = 7 //
	NalUnitTypeCodedSliceRaslN  NalUnitType = 8 //
	NalUnitTypeCodedSliceTfd    NalUnitType = 9 //
	// 10...15 // Reserved
	NalUnitTypeCodedSliceBla    NalUnitType = 16 //
	NalUnitTypeCodedSliceBlant  NalUnitType = 17 //
	NalUnitTypeCodedSliceBlaNLp NalUnitType = 18 //
	NalUnitTypeCodedSliceIdr    NalUnitType = 19 //
	NalUnitTypeCodedSliceIdrNLp NalUnitType = 20 //
	NalUnitTypeCodedSliceCra    NalUnitType = 21 //
	// 22...31 // Reserved
	NalUnitTypeVPS                 NalUnitType = 32 //
	NalUnitTypeSPS                 NalUnitType = 33 //
	NalUnitTypePPS                 NalUnitType = 34 //
	NalUnitTypeAccessUnitDelimeter NalUnitType = 35 //
	NalUnitTypeEOS                 NalUnitType = 36 //
	NalUnitTypeEOB                 NalUnitType = 37 //
	NalUnitTypeFilterData          NalUnitType = 38 //
	NalUnitTypeSEI                 NalUnitType = 39 //
	NalUnitTypeSEISuffix           NalUnitType = 40 //
	// 41...47 // Reserved
	// 48...63 // Unspecified
	// 64... // Invalid
)
