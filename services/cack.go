package services

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/ulbios/bacnet/common"
	"github.com/ulbios/bacnet/objects"
	"github.com/ulbios/bacnet/plumbing"
)

// UnconfirmedIAm is a BACnet message.
type ComplexACK struct {
	*plumbing.BVLC
	*plumbing.NPDU
	*plumbing.APDU
}

type ComplexACKDec struct {
	ObjectType   uint16
	InstanceId   uint32
	PropertyId   uint8
	PresentValue float32
}

func ComplexACKObjects(objectType uint16, instN uint32, propertyId uint8, value float32) []objects.APDUPayload {
	objs := make([]objects.APDUPayload, 5)

	objs[0] = objects.EncObjectIdentifier(true, 0, objectType, instN)
	objs[1] = objects.EncPropertyIdentifier(true, 1, propertyId)
	objs[2] = objects.EncOpeningTag(3)
	objs[3] = objects.EncReal(value)
	objs[4] = objects.EncClosingTag(3)

	return objs
}

func NewComplexACK(bvlc *plumbing.BVLC, npdu *plumbing.NPDU) *ComplexACK {
	c := &ComplexACK{
		BVLC: bvlc,
		NPDU: npdu,
		// TODO: Consider to implement parameter struct to an argment of New functions.
		APDU: plumbing.NewAPDU(plumbing.ComplexAck, ServiceConfirmedReadProperty, ComplexACKObjects(
			objects.ObjectTypeAnalogOutput, 1, objects.PropertyIdPresentValue, 0)),
	}
	c.SetLength()

	return c
}

func (c *ComplexACK) UnmarshalBinary(b []byte) error {
	if l := len(b); l < c.MarshalLen() {
		return errors.Wrap(
			common.ErrTooShortToParse,
			fmt.Sprintf("Unmarshal ComplexACK bin length %d, marshal length %d", l, c.MarshalLen()),
		)
	}

	var offset int = 0
	if err := c.BVLC.UnmarshalBinary(b[offset:]); err != nil {
		return errors.Wrap(
			common.ErrTooShortToParse,
			fmt.Sprintf("Unmarshal BVLC %x", b[offset:]),
		)
	}
	offset += c.BVLC.MarshalLen()

	if err := c.NPDU.UnmarshalBinary(b[offset:]); err != nil {
		return errors.Wrap(
			common.ErrTooShortToParse,
			fmt.Sprintf("Unmarshal NPDU %x", b[offset:]),
		)
	}
	offset += c.NPDU.MarshalLen()

	if err := c.APDU.UnmarshalBinary(b[offset:]); err != nil {
		return errors.Wrap(
			common.ErrTooShortToParse,
			fmt.Sprintf("Unmarshal APDU %x", b[offset:]),
		)
	}

	return nil
}

func (c *ComplexACK) MarshalBinary() ([]byte, error) {
	b := make([]byte, c.MarshalLen())
	if err := c.MarshalTo(b); err != nil {
		return nil, errors.Wrap(err, "unable to marshal ComplexACK")
	}
	return b, nil
}

func (c *ComplexACK) MarshalTo(b []byte) error {
	if len(b) < c.MarshalLen() {
		return errors.Wrap(
			common.ErrTooShortToMarshalBinary,
			fmt.Sprintf("Marshal ComplexACK bin length %d, marshal length %d", len(b), c.MarshalLen()),
		)
	}
	var offset = 0
	if err := c.BVLC.MarshalTo(b[offset:]); err != nil {
		return errors.Wrap(
			err,
			fmt.Sprintf("Marshal BVLC %x", b[offset:]),
		)
	}
	offset += c.BVLC.MarshalLen()

	if err := c.NPDU.MarshalTo(b[offset:]); err != nil {
		return errors.Wrap(
			err,
			fmt.Sprintf("Marshal NPDU %x", b[offset:]),
		)
	}
	offset += c.NPDU.MarshalLen()

	if err := c.APDU.MarshalTo(b[offset:]); err != nil {
		return errors.Wrap(
			err,
			fmt.Sprintf("Marshal APDU %x", b[offset:]),
		)
	}

	return nil
}

func (c *ComplexACK) MarshalLen() int {
	l := c.BVLC.MarshalLen()
	l += c.NPDU.MarshalLen()
	l += c.APDU.MarshalLen()

	return l
}

func (u *ComplexACK) SetLength() {
	u.BVLC.Length = uint16(u.MarshalLen())
}

func (c *ComplexACK) Decode() (ComplexACKDec, error) {
	decCACK := ComplexACKDec{}

	if len(c.APDU.Objects) != 3 {
		return decCACK, errors.Wrap(
			common.ErrWrongObjectCount,
			fmt.Sprintf("ComplexACK object count %d", len(c.APDU.Objects)),
		)
	}

	for i, obj := range c.APDU.Objects {
		switch i {
		case 0:
			objId, err := objects.DecObjectIdentifier(obj)
			if err != nil {
				return decCACK, errors.Wrap(err, "decode ComplexACK object case 0")
			}
			decCACK.ObjectType = objId.ObjectType
			decCACK.InstanceId = objId.InstanceNumber
		case 1:
			propId, err := objects.DecPropertyIdentifier(obj)
			if err != nil {
				return decCACK, errors.Wrap(err, "decode ComplexACK object case 1")
			}
			decCACK.PropertyId = propId
		case 2:
			value, err := objects.DecReal(obj)
			if err != nil {
				return decCACK, errors.Wrap(err, "decode ComplexACK object case 2")
			}
			decCACK.PresentValue = value
		}
	}

	return decCACK, nil
}
