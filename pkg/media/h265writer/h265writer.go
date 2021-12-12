// Package h264writer implements H264 media container writer
package h265writer

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/pion/rtp"
	"github.com/pion/rtp/codecs"
	"github.com/pion/webrtc/v3/pkg/media/h265reader"
)

type (
	// H265Writer is used to take RTP packets, parse them and
	// write the data to an io.Writer.
	H265Writer struct {
		writer       io.Writer
		hasKeyFrame  bool
		cachedPacket *codecs.H265Packet
	}
)

// New builds a new H264 writer
func New(filename string) (*H265Writer, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	return NewWith(f), nil
}

// NewWith initializes a new H264 writer with an io.Writer output
func NewWith(w io.Writer) *H265Writer {
	return &H265Writer{
		writer: w,
	}
}

// WriteRTP adds a new packet and writes the appropriate headers for it
func (h *H265Writer) WriteRTP(packet *rtp.Packet) error {
	if len(packet.Payload) == 0 {
		return nil
	}

	if !h.hasKeyFrame {
		if h.hasKeyFrame = isKeyFrame(packet.Payload); !h.hasKeyFrame {
			// key frame not defined yet. discarding packet
			return nil
		}
	}

	if h.cachedPacket == nil {
		// TODO: 需要在H265Packet外面封装一层带缓冲的packet，参考H264Packet。
		h.cachedPacket = &codecs.H265Packet{}
	}

	data, err := h.cachedPacket.Unmarshal(packet.Payload)
	if err != nil {
		return err
	}

	if len(data) < 6 {
		log.Println(fmt.Sprintf("write h265 packet head:%v size:%v", data, len(data)))
	} else {
		log.Println(fmt.Sprintf("write h265 packet head:%v size:%v", data[0:6], len(data)))
	}
	_, err = h.writer.Write(data)

	return err
}

// Close closes the underlying writer
func (h *H265Writer) Close() error {
	h.cachedPacket = nil
	if h.writer != nil {
		if closer, ok := h.writer.(io.Closer); ok {
			return closer.Close()
		}
	}

	return nil
}

func isKeyFrame(data []byte) bool {
	header := codecs.H265NALUHeader((uint16(data[0]) << 8) | uint16(data[1]))
	return header.Type() == uint8(h265reader.NalUnitTypeCodedSliceIdr) || header.Type() == uint8(h265reader.NalUnitTypeCodedSliceIdrNLp) ||
		header.Type() == uint8(h265reader.NalUnitTypeVPS) || header.Type() == uint8(h265reader.NalUnitTypeSPS) || header.Type() == uint8(h265reader.NalUnitTypePPS)
}
