package services

import (
	"fmt"

	"github.com/jonalfarlinga/bacnet/common"
	"github.com/jonalfarlinga/bacnet/objects"
	"github.com/jonalfarlinga/bacnet/plumbing"
	"github.com/pkg/errors"
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
	PresentValue interface{}
}

func ComplexACKObjects(objectType uint16, instN uint32, propertyId uint8, value interface{}) []objects.APDUPayload {
	objs := make([]objects.APDUPayload, 5)
	fmt.Println("ComplexACKObjects")
	objs[0] = objects.EncObjectIdentifier(true, 0, objectType, instN)
	objs[1] = objects.EncPropertyIdentifier(true, 1, propertyId)
	objs[2] = objects.EncOpeningTag(3)

	switch v := value.(type) {
	case int:
		objs[3] = objects.EncReal(float32(v))
	case uint8:
		objs[3] = objects.EncUnsignedInteger8(v)
	case uint16:
		objs[3] = objects.EncUnsignedInteger16(v)
	case float32:
		objs[3] = objects.EncReal(v)
	case string:
		objs[3] = objects.EncString(v)
	default:
		panic(
			fmt.Sprintf("Unsupported PresentValue type %T", value),
		)
	}

	objs[4] = objects.EncClosingTag(3)
	for _, o := range objs {
		fmt.Printf("%v\n", o)
	}
	fmt.Println("ComplexACKObjects end")
	return objs
}

func NewComplexACK(bvlc *plumbing.BVLC, npdu *plumbing.NPDU) *ComplexACK {
	fmt.Println("NewComplexACK")
	c := &ComplexACK{
		BVLC: bvlc,
		NPDU: npdu,
		// TODO: Consider to implement parameter struct to an argment of New functions.
		APDU: plumbing.NewAPDU(plumbing.ComplexAck, ServiceConfirmedReadProperty, ComplexACKObjects(
			objects.ObjectTypeAnalogOutput, 1, objects.PropertyIdPresentValue, 0)),
	}
	c.SetLength()
	fmt.Printf("Type: %d, Service: %d\n", c.APDU.Type, c.APDU.Service)
	fmt.Println("NewComplexACK end")
	return c
}

func (c *ComplexACK) UnmarshalBinary(b []byte) error {
	fmt.Printf("unmarshalBinary %x\n", b)
	if l := len(b); l < c.MarshalLen() {
		return errors.Wrap(
			common.ErrTooShortToParse,
			fmt.Sprintf("unmarshal ComplexACK bin length %d, marshal length %d", l, c.MarshalLen()),
		)
	}

	var offset int = 0
	if err := c.BVLC.UnmarshalBinary(b[offset:]); err != nil {
		return errors.Wrap(
			common.ErrTooShortToParse,
			fmt.Sprintf("unmarshal BVLC %x", b[offset:]),
		)
	}
	offset += c.BVLC.MarshalLen()

	if err := c.NPDU.UnmarshalBinary(b[offset:]); err != nil {
		return errors.Wrap(
			common.ErrTooShortToParse,
			fmt.Sprintf("unmarshal NPDU %x", b[offset:]),
		)
	}
	offset += c.NPDU.MarshalLen()

	fmt.Printf("\n\nAPDU binary %x\n", b[offset+11:])

	if err := c.APDU.UnmarshalBinary(b[offset:]); err != nil {
		return errors.Wrap(
			common.ErrTooShortToParse,
			fmt.Sprintf("unmarshal APDU %x", b[offset:]),
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
			fmt.Sprintf("marshal ComplexACK bin length %d, marshal length %d", len(b), c.MarshalLen()),
		)
	}
	var offset = 0
	if err := c.BVLC.MarshalTo(b[offset:]); err != nil {
		return errors.Wrap(
			err,
			fmt.Sprintf("marshal BVLC %x", b[offset:]),
		)
	}
	offset += c.BVLC.MarshalLen()

	if err := c.NPDU.MarshalTo(b[offset:]); err != nil {
		return errors.Wrap(
			err,
			fmt.Sprintf("marshal NPDU %x", b[offset:]),
		)
	}
	offset += c.NPDU.MarshalLen()

	if err := c.APDU.MarshalTo(b[offset:]); err != nil {
		return errors.Wrap(
			err,
			fmt.Sprintf("marshal APDU %x", b[offset:]),
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
			fmt.Sprintf("complexACK object count %d", len(c.APDU.Objects)),
		)
	}

	for i, obj := range c.APDU.Objects {
		enc_obj, ok := obj.(*objects.Object)
		if !ok {
			return decCACK, errors.Wrap(
				common.ErrInvalidObjectType,
				fmt.Sprintf("ComplexACK object at index %d is not Object type", i),
			)
		}
		fmt.Printf(
			"Object i %d tagnum %d tagclass %v data %x\n",
			i, enc_obj.TagNumber, enc_obj.TagClass, enc_obj.Data,
		)
		if enc_obj.TagClass {
			switch enc_obj.TagNumber {
			case 0:
				objId, err := objects.DecObjectIdentifier(obj)
				if err != nil {
					return decCACK, errors.Wrap(err, "decode Context object case 0")
				}
				decCACK.ObjectType = objId.ObjectType
				decCACK.InstanceId = objId.InstanceNumber
			case 1:
				propId, err := objects.DecPropertyIdentifier(obj)
				if err != nil {
					return decCACK, errors.Wrap(err, "decode Context object case 1")
				}
				decCACK.PropertyId = propId
			}
		} else {
			switch enc_obj.TagNumber {
			case 4:
				value, err := objects.DecReal(obj)
				if err != nil {
					return decCACK, errors.Wrap(err, "decode Application object case 4")
				}
				decCACK.PresentValue = value
			case 7:
				value, err := objects.DecString(obj)
				if err != nil {
					return decCACK, errors.Wrap(err, "decode Application object case 7")
				}
				fmt.Printf("String value %s\n", value)
				decCACK.PresentValue = value
			}
		}
	}

	return decCACK, nil
}
