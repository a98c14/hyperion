package types

import "encoding/json"

/*

const (
	Object     EditorInputType = 0
	Range      EditorInputType = 1
	Color      EditorInputType = 2
	Animation  EditorInputType = 3
	Sprite     EditorInputType = 4
	Percentage EditorInputType = 5
	Vec2       EditorInputType = 6
	Vec3       EditorInputType = 7
	Vec4       EditorInputType = 8
	Nested     EditorInputType = 9
	Bool       EditorInputType = 10
	Number     EditorInputType = 11
	Text       EditorInputType = 12
	TextArea   EditorInputType = 13
)

*/
func ParseType(valueType EditorInputType, value json.RawMessage) {
	switch valueType {
	case Range:

	}
}
