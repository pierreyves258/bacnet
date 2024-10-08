package objects

import (
	"fmt"

	"github.com/jonalfarlinga/bacnet/common"
	"github.com/pkg/errors"
)

func DecPropertyIdentifier(rawPayload APDUPayload) (uint8, error) {
	rawObject, ok := rawPayload.(*Object)
	if !ok {
		return 0, errors.Wrap(
			common.ErrWrongPayload,
			fmt.Sprintf("DecPropertyIdentifier not ok %v", rawPayload),
		)
	}

	switch rawObject.TagClass {
	case true:
		if rawObject.Length != 1 {
			return 0, errors.Wrap(
				common.ErrWrongStructure,
				fmt.Sprintf("DecPropertyIdentifier length %d tag class  %v", rawObject.Length, rawObject.TagClass),
			)
		}
	case false:
		if rawObject.Length != 1 || !rawObject.TagClass {
			return 0, errors.Wrap(
				common.ErrWrongStructure,
				fmt.Sprintf("DecPropertyIdentifier length %d tag class  %v", rawObject.Length, rawObject.TagClass),
			)
		}
	}

	return rawObject.Data[0], nil
}

func EncPropertyIdentifier(contextTag bool, tagN uint8, propId uint8) *Object {
	newObj := Object{}
	data := make([]byte, 1)
	data[0] = propId

	newObj.TagNumber = tagN
	newObj.TagClass = contextTag
	newObj.Data = data
	newObj.Length = uint8(len(data))

	return &newObj
}
