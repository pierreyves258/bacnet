package objects

import (
	"fmt"

	"github.com/pierreyves258/bacnet/common"
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
			fmt.Sprintf("failed to unmarshal NamedTag - binary too short - %x", b),
		)
	}
	n.TagNumber = b[0] >> 4
	n.TagClass = common.IntToBool(int(b[0]) & 0x8 >> 3)
	n.Name = b[0] & 0x7

	if l := len(b); l < 1 {
		return errors.Wrap(
			common.ErrTooShortToParse,
			fmt.Sprintf("failed to unmarshal NamedTag - missing data - %v", n),
		)
	}

	return nil
}

func (n *NamedTag) MarshalBinary() ([]byte, error) {
	b := make([]byte, n.MarshalLen())
	if err := n.MarshalTo(b); err != nil {
		return nil, errors.Wrap(err, "failed to marshal binary")
	}

	return b, nil
}

func (n *NamedTag) MarshalTo(b []byte) error {
	if len(b) < n.MarshalLen() {
		return errors.Wrap(common.ErrTooShortToMarshalBinary, "failed to marshall NamedTag - marshal length too short")
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
		return false, errors.Wrap(common.ErrWrongPayload, "failed to decode OpeningTab")
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
		return false, errors.Wrap(common.ErrWrongPayload, "failed to decode ClosingTab")
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
