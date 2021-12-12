// Package h265reader implements a H265 Annex-B Reader
package h265reader

import (
	"bytes"
	"errors"
	"io"

	"github.com/pion/rtp/codecs"
)

// H265Reader reads data from stream and constructs h265 nal units
type H265Reader struct {
	stream                      io.Reader
	nalBuffer                   []byte
	countOfConsecutiveZeroBytes int
	nalPrefixParsed             bool
	readBuffer                  []byte
}

var (
	errNilReader           = errors.New("stream is nil")
	errDataIsNotH265Stream = errors.New("data is not a H265 bitstream")
)

// NewReader creates new H265Reader
func NewReader(in io.Reader) (*H265Reader, error) {
	if in == nil {
		return nil, errNilReader
	}

	reader := &H265Reader{
		stream:          in,
		nalBuffer:       make([]byte, 0),
		nalPrefixParsed: false,
		readBuffer:      make([]byte, 0),
	}

	return reader, nil
}

// NAL H.265 Network Abstraction Layer
type NAL struct {
	// NAL header
	ForbiddenZeroBit bool
	UnitType         uint8
	LayerID          uint8
	TID              uint8

	Data []byte // header byte + rbsp
}

func (reader *H265Reader) read(numToRead int) (data []byte) {
	for len(reader.readBuffer) < numToRead {
		buf := make([]byte, 4096)
		n, err := reader.stream.Read(buf)
		if n == 0 || err != nil {
			break
		}
		buf = buf[0:n]
		reader.readBuffer = append(reader.readBuffer, buf...)
	}
	var numShouldRead int
	if numToRead <= len(reader.readBuffer) {
		numShouldRead = numToRead
	} else {
		numShouldRead = len(reader.readBuffer)
	}
	data = reader.readBuffer[0:numShouldRead]
	reader.readBuffer = reader.readBuffer[numShouldRead:]
	return data
}

func (reader *H265Reader) bitStreamStartsWithH265Prefix() (prefixLength int, e error) {
	nalPrefix3Bytes := []byte{0, 0, 1}
	nalPrefix4Bytes := []byte{0, 0, 0, 1}

	prefixBuffer := reader.read(4)

	n := len(prefixBuffer)

	if n == 0 {
		return 0, io.EOF
	}

	if n < 3 {
		return 0, errDataIsNotH265Stream
	}

	nalPrefix3BytesFound := bytes.Equal(nalPrefix3Bytes, prefixBuffer[:3])
	if n == 3 {
		if nalPrefix3BytesFound {
			return 0, io.EOF
		}
		return 0, errDataIsNotH265Stream
	}

	// n == 4
	if nalPrefix3BytesFound {
		reader.nalBuffer = append(reader.nalBuffer, prefixBuffer[3])
		return 3, nil
	}

	nalPrefix4BytesFound := bytes.Equal(nalPrefix4Bytes, prefixBuffer)
	if nalPrefix4BytesFound {
		return 4, nil
	}
	return 0, errDataIsNotH265Stream
}

// NextNAL reads from stream and returns then next NAL,
// and an error if there is incomplete frame data.
// Returns all nil values when no more NALs are available.
func (reader *H265Reader) NextNAL() (*NAL, error) {
	if !reader.nalPrefixParsed {
		_, err := reader.bitStreamStartsWithH265Prefix()
		if err != nil {
			return nil, err
		}

		reader.nalPrefixParsed = true
	}

	for {
		buffer := reader.read(1)
		n := len(buffer)

		if n != 1 {
			break
		}
		readByte := buffer[0]
		nalFound := reader.processByte(readByte)
		if nalFound {
			nal := newNal(reader.nalBuffer)
			nal.parseHeader()
			if nal.UnitType == uint8(NalUnitTypeSEI) { // TODO: 仅跳过SEI，没有其他的情况了吗？
				reader.nalBuffer = nil
				continue
			} else {
				break
			}
		}

		reader.nalBuffer = append(reader.nalBuffer, readByte)
	}

	if len(reader.nalBuffer) == 0 {
		return nil, io.EOF
	}

	nal := newNal(reader.nalBuffer)
	reader.nalBuffer = nil
	nal.parseHeader()

	return nal, nil
}

func (reader *H265Reader) processByte(readByte byte) (nalFound bool) {
	nalFound = false

	switch readByte {
	case 0:
		reader.countOfConsecutiveZeroBytes++
	case 1:
		if reader.countOfConsecutiveZeroBytes >= 2 {
			countOfConsecutiveZeroBytesInPrefix := 2
			if reader.countOfConsecutiveZeroBytes > 2 {
				countOfConsecutiveZeroBytesInPrefix = 3
			}

			if nalUnitLength := len(reader.nalBuffer) - countOfConsecutiveZeroBytesInPrefix; nalUnitLength > 0 {
				reader.nalBuffer = reader.nalBuffer[0:nalUnitLength]
				nalFound = true
			}
		}

		reader.countOfConsecutiveZeroBytes = 0
	default:
		reader.countOfConsecutiveZeroBytes = 0
	}

	return nalFound
}

func newNal(data []byte) *NAL {
	return &NAL{ForbiddenZeroBit: false, UnitType: 0, LayerID: 0, TID: 0, Data: data}
}

func (h *NAL) parseHeader() {
	header := codecs.H265NALUHeader((uint16(h.Data[0]) << 8) | uint16(h.Data[1]))

	h.ForbiddenZeroBit = header.F()
	h.UnitType = header.Type()
	h.LayerID = header.LayerID()
	h.TID = header.TID()
}

// func (h *NAL) IsKeyFrame() bool {
// 	// IDRs are key frames
// 	return h.UnitType == uint8(NalUnitTypeCodedSliceIdr) || h.UnitType == uint8(NalUnitTypeCodedSliceIdrNLp) ||
// 		h.UnitType == uint8(NalUnitTypeVPS) || h.UnitType == uint8(NalUnitTypeSPS) || h.UnitType == uint8(NalUnitTypeCodedSliceIdr) || h.UnitType == uint8(NalUnitTypePPS)
// }
