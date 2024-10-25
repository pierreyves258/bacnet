package services

import (
	"fmt"

	"github.com/pierreyves258/bacnet/common"
	"github.com/pierreyves258/bacnet/objects"
	"github.com/pierreyves258/bacnet/plumbing"
	"github.com/pkg/errors"
)

// UnconfirmedIAm is a BACnet message.
type SegmentedAck struct {
	*plumbing.BVLC
	*plumbing.NPDU
	*plumbing.APDU
}

func NewSegmentAck(bvlc *plumbing.BVLC, npdu *plumbing.NPDU) *SegmentedAck {
	e := &SegmentedAck{
		BVLC: bvlc,
		NPDU: npdu,
		APDU: plumbing.NewAPDU(plumbing.SegmentAck, ServiceConfirmedReadProperty, nil),
	}
	e.SetLength()

	return e
}

// UnmarshalBinary sets the values retrieved from byte sequence in a UnconfirmedIAm frame.
func (e *SegmentedAck) UnmarshalBinary(b []byte) error {
	if l := len(b); l < e.MarshalLen() {
		return errors.Wrap(
			common.ErrTooShortToParse,
			fmt.Sprintf("failed to unmarshal Error - marshal length %d binary length %d", e.MarshalLen(), l),
		)
	}

	var offset int = 0

	if err := e.BVLC.UnmarshalBinary(b[offset:]); err != nil {
		return errors.Wrap(
			common.ErrTooShortToParse,
			fmt.Sprintf("unmarshalling Error %v", e),
		)
	}
	offset += e.BVLC.MarshalLen()

	if err := e.NPDU.UnmarshalBinary(b[offset:]); err != nil {
		return errors.Wrap(
			common.ErrTooShortToParse,
			fmt.Sprintf("unmarshalling Error %v", e),
		)
	}
	offset += e.NPDU.MarshalLen()

	if err := e.APDU.UnmarshalBinary(b[offset:]); err != nil {
		return errors.Wrap(
			common.ErrTooShortToParse,
			fmt.Sprintf("unmarshalling Error %v", e),
		)
	}

	return nil
}

// MarshalBinary returns the byte sequence generated from a UnconfirmedIAm instance.
func (e *SegmentedAck) MarshalBinary() ([]byte, error) {
	b := make([]byte, e.MarshalLen())
	if err := e.MarshalTo(b); err != nil {
		return nil, errors.Wrap(err, "failed to marshal binary")
	}
	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (e *SegmentedAck) MarshalTo(b []byte) error {
	if len(b) < e.MarshalLen() {
		return errors.Wrap(
			common.ErrTooShortToMarshalBinary,
			fmt.Sprintf("failed to marshal Error - marshal length %d binary length %d", e.MarshalLen(), len(b)),
		)
	}
	var offset = 0
	if err := e.BVLC.MarshalTo(b[offset:]); err != nil {
		return errors.Wrap(err, "marshalling Error")
	}
	offset += e.BVLC.MarshalLen()

	if err := e.NPDU.MarshalTo(b[offset:]); err != nil {
		return errors.Wrap(err, "marshalling Error")
	}
	offset += e.NPDU.MarshalLen()

	if err := e.APDU.MarshalTo(b[offset:]); err != nil {
		return errors.Wrap(err, "marshalling Error")
	}

	return nil
}

// MarshalLen returns the serial length of UnconfirmedIAm.
func (e *SegmentedAck) MarshalLen() int {
	l := e.BVLC.MarshalLen()
	l += e.NPDU.MarshalLen()
	l += e.APDU.MarshalLen()

	return l
}

// SetLength sets the length in Length field.
func (e *SegmentedAck) SetLength() {
	e.BVLC.Length = uint16(e.MarshalLen())
}

func (e *SegmentedAck) Decode() (ErrorDec, error) {
	decErr := ErrorDec{}

	if len(e.APDU.Objects) != 4 {
		return decErr, errors.Wrap(
			common.ErrWrongObjectCount,
			fmt.Sprintf("failed to decode Error - object count: %d", len(e.APDU.Objects)),
		)
	}

	e.APDU.Objects = e.APDU.Objects[2:]

	for i, obj := range e.APDU.Objects {
		switch i {
		case 0:
			errClass, err := objects.DecEnumerated(obj)
			if err != nil {
				return decErr, errors.Wrap(err, "failed to decode Enumerated Object")
			}
			decErr.ErrorClass = uint8(errClass)
		case 1:
			errCode, err := objects.DecEnumerated(obj)
			if err != nil {
				return decErr, errors.Wrap(err, "failed to decode Enumerated Object")
			}
			decErr.ErrorCode = uint8(errCode)
		}
	}

	return decErr, nil
}
