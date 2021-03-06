package rtcp

import (
	"encoding/binary"

	"github.com/pkg/errors"
)

// RTCP packet types registered with IANA. See: https://www.iana.org/assignments/rtp-parameters/rtp-parameters.xhtml#rtp-parameters-4
const (
	TypeSenderReport       = 200 // RFC 3550, 6.4.1
	TypeReceiverReport     = 201 // RFC 3550, 6.4.2
	TypeSourceDescription  = 202 // RFC 3550, 6.5
	TypeGoodbye            = 203 // RFC 3550, 6.6
	TypeApplicationDefined = 204 // RFC 3550, 6.7
)

// A Header is the common header shared by all RTCP packets
type Header struct {
	// Identifies the version of RTP, which is the same in RTCP packets
	// as in RTP data packets.
	Version uint8
	// If the padding bit is set, this individual RTCP packet contains
	// some additional padding octets at the end which are not part of
	// the control information but are included in the length field.
	Padding bool
	// The number of reception report blocks contained in this packet.
	ReportCount uint8
	// The RTCP packet type for this packet
	Type uint8
	// The length of this RTCP packet in 32-bit words minus one,
	// including the header and any padding.
	Length uint16
}

const (
	headerLength     = 4
	versionShift     = 6
	versionMask      = 0x3
	paddingShift     = 5
	paddingMask      = 0x1
	reportCountShift = 0
	reportCountMask  = 0x1f
)

var (
	errInvalidVersion     = errors.New("invalid version")
	errInvalidReportCount = errors.New("invalid report count")
	errHeaderTooShort     = errors.New("rtcp header too short")
)

// Marshal encodes the Header in binary
func (h Header) Marshal() ([]byte, error) {
	/*
	 *  0                   1                   2                   3
	 *  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	 * |V=2|P|    RC   |   PT=SR=200   |             length            |
	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	 */
	rawPacket := make([]byte, headerLength)

	if h.Version > 3 {
		return nil, errInvalidVersion
	}
	rawPacket[0] |= h.Version << versionShift

	if h.Padding {
		rawPacket[0] |= 1 << paddingShift
	}

	if h.ReportCount > 31 {
		return nil, errInvalidReportCount
	}
	rawPacket[0] |= h.ReportCount << reportCountShift

	rawPacket[1] = h.Type

	binary.BigEndian.PutUint16(rawPacket[2:], h.Length)

	return rawPacket, nil
}

// Unmarshal decodes the Header from binary
func (h *Header) Unmarshal(rawPacket []byte) error {
	if len(rawPacket) < headerLength {
		return errHeaderTooShort
	}

	/*
	 *  0                   1                   2                   3
	 *  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	 * |V=2|P|    RC   |   PT=SR=200   |             length            |
	 * +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	 */

	h.Version = rawPacket[0] >> versionShift & versionMask
	h.Padding = (rawPacket[0] >> paddingShift & paddingMask) > 0
	h.ReportCount = rawPacket[0] >> reportCountShift & reportCountMask

	h.Type = rawPacket[1]

	h.Length = binary.BigEndian.Uint16(rawPacket[2:])

	return nil
}
