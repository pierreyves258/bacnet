package objects

import (
	"fmt"

	"github.com/jonalfarlinga/bacnet/common"
	"github.com/pkg/errors"
)

type NamedTag struct {
	TagNumber uint8
	TagClass  bool
	Name      uint8
}

func NewNamedTag(number uint8, class bool, name uint8) *NamedTag {
	return &NamedTag{
		TagNumber: number,
		TagClass:  class,
		Name:      name,
	}
}

func (n *NamedTag) UnmarshalBinary(b []byte) error {
	if l := len(b); l < objLenMin {
		return errors.Wrap(
			common.ErrTooShortToParse,
			fmt.Sprintf("Unmarshal NamedTag bin length %d, min length %d", l, objLenMin),
		)
	}
	n.TagNumber = b[0] >> 4
	n.TagClass = common.IntToBool(int(b[0]) & 0x8 >> 3)
	n.Name = b[0] & 0x7

	if l := len(b); l < 1 {
		return errors.Wrap(
			common.ErrTooShortToParse,
			fmt.Sprintf("Unmarshal NamedTag bin length %d", l),
		)
	}

	return nil
}

func (n *NamedTag) MarshalBinary() ([]byte, error) {
	b := make([]byte, n.MarshalLen())
	if err := n.MarshalTo(b); err != nil {
		return nil, errors.Wrap(err, "unable to marshal NamedTag")
	}

	return b, nil
}

func (n *NamedTag) MarshalTo(b []byte) error {
	if len(b) < n.MarshalLen() {
		return errors.Wrap(
			common.ErrTooShortToMarshalBinary,
			fmt.Sprintf("marshall NamedTag length %d is less than %d", len(b), n.MarshalLen()),
		)
	}
	b[0] = n.TagNumber<<4 | uint8(common.BoolToInt(n.TagClass))<<3 | n.Name

	return nil
}

func (n *NamedTag) MarshalLen() int {
	return 1
}

func DecOpeningTab(rawPayload APDUPayload) (bool, error) {
	rawTag, ok := rawPayload.(*NamedTag)
	if !ok {
		return false, errors.Wrap(
			common.ErrWrongPayload,
			fmt.Sprintf("DecOpeningTab not ok %T", rawPayload),
		)
	}
	return rawTag.Name == 0x6 && rawTag.TagClass, nil
}

func EncOpeningTag(tagN uint8) *NamedTag {
	oTag := NamedTag{}

	oTag.TagClass = true
	oTag.TagNumber = tagN
	oTag.Name = 0x6

	return &oTag
}

func DecClosingTab(rawPayload APDUPayload) (bool, error) {
	rawTag, ok := rawPayload.(*NamedTag)
	if !ok {
		return false, errors.Wrap(
			common.ErrWrongPayload,
			fmt.Sprintf("DecClosingTab not ok %T", rawPayload),
		)
	}
	return rawTag.Name == 0x7 && rawTag.TagClass, nil
}

func EncClosingTag(tagN uint8) *NamedTag {
	cTag := NamedTag{}

	cTag.TagClass = true
	cTag.TagNumber = tagN
	cTag.Name = 0x7

	return &cTag
}
