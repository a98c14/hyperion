package parser

import "encoding/json"

// TODO(selim): Parse `value` based on type. Currently module
// part values are returned directly as json.RawMessage. It works
// but results are not validated. Here `value` should be converted
// to desired and return error if it can't be converted.
func ParseType(valueType EditorInputType, value string) (json.RawMessage, error) {
	switch valueType {
	case Range:
		break
	case Color:
		break
	case Animation:
		break
	case Sprite:
		break
	case Percentage:
		break
	case Vec2:
		break
	case Vec3:
		break
	case Vec4:
		break
	}
	return nil, nil
}
