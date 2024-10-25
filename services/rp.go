package services

import (
	"fmt"

	"github.com/pierreyves258/bacnet/common"
	"github.com/pierreyves258/bacnet/objects"
	"github.com/pierreyves258/bacnet/plumbing"
	"github.com/pkg/errors"
)

// UnconfirmedIAm is a BACnet message.
type ConfirmedReadProperty struct {
	*plumbing.BVLC
	*plumbing.NPDU
	*plumbing.APDU
}

type ConfirmedReadPropertyDec struct {
	ObjectType uint16
	InstanceId uint32
	PropertyId uint8
}

func ConfirmedReadPropertyObjects(objectType uint16, instN uint32, propertyId uint8) []objects.APDUPayload {
	objs := make([]objects.APDUPayload, 2)

	objs[0] = objects.EncObjectIdentifier(true, 0, objectType, instN)
	objs[1] = objects.EncPropertyIdentifier(true, 1, propertyId)

	return objs
}

func ConfirmedReadMultiplePropertyObjects(objectType uint16, instN uint32, propertyId []uint8) []objects.APDUPayload {
	objs := make([]objects.APDUPayload, 3+len(propertyId))

	objs[0] = objects.EncObjectIdentifier(true, 0, objectType, instN)
	objs[1] = objects.EncOpeningTag(1)

	for i := range propertyId {
		objs[i+2] = objects.EncPropertyIdentifier(true, 1, propertyId[i])
	}

	objs[len(propertyId)+2] = objects.EncClosingTag(1)

	return objs
}

func NewConfirmedReadProperty(bvlc *plumbing.BVLC, npdu *plumbing.NPDU) *ConfirmedReadProperty {
	c := &ConfirmedReadProperty{
		BVLC: bvlc,
		NPDU: npdu,
		// TODO: Consider to implement parameter struct to an argment of New functions.
		APDU: plumbing.NewAPDU(plumbing.ConfirmedReq, ServiceConfirmedReadProperty, nil),
	}
	c.SetLength()

	return c
}

func NewConfirmedReadPropertyMultiple(bvlc *plumbing.BVLC, npdu *plumbing.NPDU) *ConfirmedReadProperty {
	c := &ConfirmedReadProperty{
		BVLC: bvlc,
		NPDU: npdu,
		APDU: plumbing.NewAPDU(plumbing.ConfirmedReq, ServiceConfirmedReadPropMultiple, nil),
	}
	c.SetLength()

	return c
}

func (c *ConfirmedReadProperty) UnmarshalBinary(b []byte) error {
	if l := len(b); l < c.MarshalLen() {
		return errors.Wrap(
			common.ErrTooShortToParse,
			fmt.Sprintf("failed to unmarshal ConfirmedRP - marshal length %d binary length %d", c.MarshalLen(), l),
		)
	}

	var offset int = 0
	if err := c.BVLC.UnmarshalBinary(b[offset:]); err != nil {
		return errors.Wrap(
			common.ErrTooShortToParse,
			fmt.Sprintf("unmarshalling ConfirmedRP %v", c),
		)
	}
	offset += c.BVLC.MarshalLen()

	if err := c.NPDU.UnmarshalBinary(b[offset:]); err != nil {
		return errors.Wrap(
			common.ErrTooShortToParse,
			fmt.Sprintf("unmarshalling ConfirmedRP %v", c),
		)
	}
	offset += c.NPDU.MarshalLen()

	if err := c.APDU.UnmarshalBinary(b[offset:]); err != nil {
		return errors.Wrap(
			common.ErrTooShortToParse,
			fmt.Sprintf("unmarshalling ConfirmedRP %v", c),
		)
	}

	return nil
}

func (c *ConfirmedReadProperty) MarshalBinary() ([]byte, error) {
	b := make([]byte, c.MarshalLen())
	if err := c.MarshalTo(b); err != nil {
		return nil, errors.Wrap(err, "failed to marshal binary")
	}
	return b, nil
}

func (c *ConfirmedReadProperty) MarshalTo(b []byte) error {
	if len(b) < c.MarshalLen() {
		return errors.Wrap(
			common.ErrTooShortToMarshalBinary,
			fmt.Sprintf("failed to marshal ConfirmedRP - marshal length %d binary length %d", c.MarshalLen(), len(b)),
		)
	}
	var offset = 0
	if err := c.BVLC.MarshalTo(b[offset:]); err != nil {
		return errors.Wrap(err, "failed to marshal ConfirmedRP")
	}
	offset += c.BVLC.MarshalLen()

	if err := c.NPDU.MarshalTo(b[offset:]); err != nil {
		return errors.Wrap(err, "failed to marshal ConfirmedRP")
	}
	offset += c.NPDU.MarshalLen()

	if err := c.APDU.MarshalTo(b[offset:]); err != nil {
		return errors.Wrap(err, "failed to marshal ConfirmedRP")
	}

	return nil
}

func (c *ConfirmedReadProperty) MarshalLen() int {
	l := c.BVLC.MarshalLen()
	l += c.NPDU.MarshalLen()
	l += c.APDU.MarshalLen()

	return l
}

func (c *ConfirmedReadProperty) SetLength() {
	c.BVLC.Length = uint16(c.MarshalLen())
}

func (c *ConfirmedReadProperty) Decode() (ConfirmedReadPropertyDec, error) {
	decCRP := ConfirmedReadPropertyDec{}

	if len(c.APDU.Objects) != 2 {
		return decCRP, errors.Wrap(
			common.ErrWrongObjectCount,
			fmt.Sprintf("failed to decode ConfirmedRP - object count %d", len(c.APDU.Objects)),
		)
	}

	for i, obj := range c.APDU.Objects {
		switch i {
		case 0:
			objId, err := objects.DecObjectIdentifier(obj)
			if err != nil {
				return decCRP, errors.Wrap(err, "decoding ConfirmedRP")
			}
			decCRP.ObjectType = objId.ObjectType
			decCRP.InstanceId = objId.InstanceNumber
		case 1:
			propId, err := objects.DecPropertyIdentifier(obj)
			if err != nil {
				return decCRP, errors.Wrap(err, "decoding ConfirmedRP")
			}
			decCRP.PropertyId = propId
		}
	}

	return decCRP, nil
}
