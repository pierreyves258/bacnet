package objects

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/ulbios/bacnet/common"
)

func DecPriority(rawPayload APDUPayload) (uint8, error) {
	rawObject, ok := rawPayload.(*Object)
	if !ok {
		return 0, errors.Wrap(
			common.ErrWrongPayload,
			fmt.Sprintf("failed to decode Priority Object - %v", rawPayload),
		)
	}

	switch rawObject.TagClass {
	case true:
		if rawObject.Length != 1 {
			return 0, errors.Wrap(
				common.ErrWrongStructure,
				fmt.Sprintf("failed to decode Priority Object - wrong binary length - %x", rawObject.Data),)
		}
	case false:
		if rawObject.Length != 1 || !rawObject.TagClass {
			return 0, errors.Wrap(
				common.ErrWrongStructure,
				fmt.Sprintf("failed to decode Priority Object - wrong tag number - %v", rawObject.TagNumber),
			)
		}
	}

	return rawObject.Data[0], nil
}

func EncPriority(contextTag bool, tagN uint8, priority uint8) *Object {
	newObj := Object{}
	data := make([]byte, 1)
	data[0] = priority

	newObj.TagNumber = tagN
	newObj.TagClass = contextTag
	newObj.Data = data
	newObj.Length = uint8(len(data))

	return &newObj
}
