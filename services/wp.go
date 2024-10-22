package services

import (
	"fmt"

	"github.com/pierreyves258/bacnet/common"
	"github.com/pierreyves258/bacnet/objects"
	"github.com/pierreyves258/bacnet/plumbing"
	"github.com/pkg/errors"
)

// UnconfirmedIAm is a BACnet message.
type ConfirmedWriteProperty struct {
	*plumbing.BVLC
	*plumbing.NPDU
	*plumbing.APDU
}

type ConfirmedWritePropertyDec struct {
	ObjectType uint16
	InstanceId uint32
	PropertyId uint8
	Value      float32
	Priority   uint8
}

func ConfirmedWritePropertyObjects(objectType uint16, instN uint32, propertyId uint8, value float32) []objects.APDUPayload {
	objs := make([]objects.APDUPayload, 7)

	objs[0] = objects.EncObjectIdentifier(true, 0, objectType, instN)
	objs[1] = objects.EncPropertyIdentifier(true, 1, propertyId)
	objs[2] = objects.EncOpeningTag(3)
	objs[3] = objects.EncReal(value)
	objs[4] = objects.EncNull()
	objs[5] = objects.EncClosingTag(3)
	objs[6] = objects.EncPriority(true, 4, 16)

	return objs
}

func NewConfirmedWriteProperty(bvlc *plumbing.BVLC, npdu *plumbing.NPDU) *ConfirmedWriteProperty {
	c := &ConfirmedWriteProperty{
		BVLC: bvlc,
		NPDU: npdu,
		APDU: plumbing.NewAPDU(plumbing.ConfirmedReq, ServiceConfirmedWriteProperty, nil),
	}
	c.SetLength()

	return c
}

func (c *ConfirmedWriteProperty) UnmarshalBinary(b []byte) error {
	if l := len(b); l < c.MarshalLen() {
		return errors.Wrap(
			common.ErrTooShortToParse,
			fmt.Sprintf("failed to unmarshal ConfirmedWP - marshal length %d binary length %d", c.MarshalLen(), l),
		)
	}

	var offset int = 0
	if err := c.BVLC.UnmarshalBinary(b[offset:]); err != nil {
		return errors.Wrap(
			common.ErrTooShortToParse,
			fmt.Sprintf("unmarshalling ConfirmedWP %v", c),
		)
	}
	offset += c.BVLC.MarshalLen()

	if err := c.NPDU.UnmarshalBinary(b[offset:]); err != nil {
		return errors.Wrap(
			common.ErrTooShortToParse,
			fmt.Sprintf("unmarshalling ConfirmedWP %v", c),
		)
	}
	offset += c.NPDU.MarshalLen()

	if err := c.APDU.UnmarshalBinary(b[offset:]); err != nil {
		return errors.Wrap(
			common.ErrTooShortToParse,
			fmt.Sprintf("unmarshalling ConfirmedWP %v", c),
		)
	}

	return nil
}

func (c *ConfirmedWriteProperty) MarshalBinary() ([]byte, error) {
	b := make([]byte, c.MarshalLen())
	if err := c.MarshalTo(b); err != nil {
		return nil, errors.Wrap(err, "failed to marshal binary")
	}
	return b, nil
}

func (c *ConfirmedWriteProperty) MarshalTo(b []byte) error {
	if len(b) < c.MarshalLen() {
		return errors.Wrap(
			common.ErrTooShortToMarshalBinary,
			fmt.Sprintf("failed to marshal ConfirmedWP - marshal length %d binary length %d", c.MarshalLen(), len(b)),
		)
	}
	var offset = 0
	if err := c.BVLC.MarshalTo(b[offset:]); err != nil {
		return errors.Wrap(err, "failed to marshal ConfirmedWP")
	}
	offset += c.BVLC.MarshalLen()

	if err := c.NPDU.MarshalTo(b[offset:]); err != nil {
		return errors.Wrap(err, "failed to marshal ConfirmedWP")
	}
	offset += c.NPDU.MarshalLen()

	if err := c.APDU.MarshalTo(b[offset:]); err != nil {
		return errors.Wrap(err, "failed to marshal ConfirmedWP")
	}

	return nil
}

func (c *ConfirmedWriteProperty) MarshalLen() int {
	l := c.BVLC.MarshalLen()
	l += c.NPDU.MarshalLen()
	l += c.APDU.MarshalLen()

	return l
}

func (c *ConfirmedWriteProperty) SetLength() {
	c.BVLC.Length = uint16(c.MarshalLen())
}

func (c *ConfirmedWriteProperty) Decode() (ConfirmedWritePropertyDec, error) {
	decCWP := ConfirmedWritePropertyDec{}

	if len(c.APDU.Objects) != 5 {
		return decCWP, errors.Wrap(
			common.ErrWrongObjectCount,
			fmt.Sprintf("failed to decode ConfirmedWP - object count %d", len(c.APDU.Objects)),
		)
	}

	for i, obj := range c.APDU.Objects {
		switch i {
		case 0:
			objId, err := objects.DecObjectIdentifier(obj)
			if err != nil {
				return decCWP, errors.Wrap(err, "decoding ConfirmedWP")
			}
			decCWP.ObjectType = objId.ObjectType
			decCWP.InstanceId = objId.InstanceNumber
		case 1:
			propId, err := objects.DecPropertyIdentifier(obj)
			if err != nil {
				return decCWP, errors.Wrap(err, "decoding ConfirmedWP")
			}
			decCWP.PropertyId = propId
		case 2:
			value, err := objects.DecReal(obj)
			if err != nil {
				return decCWP, errors.Wrap(err, "decoding ConfirmedWP")
			}
			decCWP.Value = value
		case 4:
			priority, err := objects.DecPriority(obj)
			if err != nil {
				return decCWP, errors.Wrap(err, "decoding ConfirmedWP")
			}
			decCWP.Priority = priority
		}
	}

	return decCWP, nil
}
