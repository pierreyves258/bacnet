package bacnet

import (
	"fmt"
	"reflect"

	"log"

	"github.com/pierreyves258/bacnet/common"
	"github.com/pierreyves258/bacnet/plumbing"
	"github.com/pierreyves258/bacnet/services"
	"github.com/pkg/errors"
)

const bacnetLenMin = 8

func combine(t, s uint8) uint16 {
	// 0001, 0010
	// return:
	// 0001 0000 OR
	// 0000 0010
	// = 0001 0010
	return uint16(t)<<8 | uint16(s)
}

// Parse decodes the given bytes.
func Parse(b []byte) (plumbing.BACnet, error) {
	if len(b) < bacnetLenMin {
		return nil, errors.Wrap(
			common.ErrTooShortToParse,
			fmt.Sprintf("Parsing length %d", len(b)),
		)
	}

	var bvlc plumbing.BVLC
	var npdu plumbing.NPDU
	var bacnet plumbing.BACnet

	offset := 0
	log.Println("parsing")
	if err := bvlc.UnmarshalBinary(b); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Parsing BVLC %x", b))
	}
	log.Println("bvlc done")
	offset += bvlc.MarshalLen()

	if err := npdu.UnmarshalBinary(b[offset:]); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Parsing NPDU %x", b[offset:]))
	}
	log.Println("npdu done")
	offset += npdu.MarshalLen()

	var c uint16
	// We can use b[offset] >> 4 & 0xF
	// PDU Types are [0x0, ..., 0x7]
	// We can copmplete the list of PDUs with 0x6 and 0x4
	PDUType := b[offset] >> 4 & 0xFF
	switch PDUType {
	case plumbing.UnConfirmedReq:
		c = combine(b[offset], b[offset+1])
	case plumbing.ConfirmedReq:
		c = combine(b[offset], b[offset+3]) // We need to skip the PDU flags and the InvokeID
	case plumbing.ComplexAck, plumbing.SimpleAck, plumbing.Error:
		c = combine(b[offset], 0) // We need to skip the PDU flags and the InvokeID
	}

	// why don't we create c using PDUType instead of b[offset]?
	// then below cases don't have to be left-shifted
	switch c {
	case combine(plumbing.UnConfirmedReq<<4, services.ServiceUnconfirmedWhoIs):
		bacnet = services.NewUnconfirmedWhoIs(&bvlc, &npdu)
	case combine(plumbing.UnConfirmedReq<<4, services.ServiceUnconfirmedIAm):
		bacnet = services.NewUnconfirmedIAm(&bvlc, &npdu)
	case combine(plumbing.ConfirmedReq<<4, services.ServiceConfirmedReadProperty):
		bacnet = services.NewConfirmedReadProperty(&bvlc, &npdu)
	case combine(plumbing.ConfirmedReq<<4, services.ServiceConfirmedWriteProperty):
		bacnet = services.NewConfirmedWriteProperty(&bvlc, &npdu)
	case combine(plumbing.ComplexAck<<4, 0):
		bacnet = services.NewComplexACK(&bvlc, &npdu)
	case combine(plumbing.SimpleAck<<4, 0):
		bacnet = services.NewSimpleACK(&bvlc, &npdu)
	case combine(plumbing.Error<<4, 0):
		bacnet = services.NewError(&bvlc, &npdu)
	default:
		return nil, errors.Wrap(
			common.ErrNotImplemented,
			fmt.Sprintf("Parsing service: %x", c),
		)
	}

	log.Printf("type %s\n", reflect.TypeOf(bacnet))
	if err := bacnet.UnmarshalBinary(b); err != nil {
		return nil, errors.Wrap(
			err,
			fmt.Sprintf("Parsing BACnet %x", b[offset:]),
		)
	}

	log.Println("processed")
	return bacnet, nil
}
