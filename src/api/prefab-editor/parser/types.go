package parser

/*
	Editor input types
	{
		Object,
		Range,
		Color,
		Animation,
		Sprite,
		Percentage,
		Vec2,
		Vec3,
		Vec4,
		Nested,
		Bool,
		Number,
		Text,
		TextArea,
	}
*/
type EditorInputType int32

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

type RangeTypeInt struct {
	Value int32 `json:"value"`
}

type RangeTypeFloat struct {
	Value float32 `json:"value"`
}

type ColorType struct {
	R float32 `json:"r"`
	G float32 `json:"g"`
	B float32 `json:"b"`
	A float32 `json:"a"`
}

type AnimationType struct {
	Id int32 `json:"id"`
}

type SpriteType struct {
	Id int32 `json:"id"`
}

type PercentageType struct {
	Value float32 `json:"value"`
}

type Vector2Type struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type Vector3Type struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
}

type Vector4Type struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
	W float32 `json:"w"`
}

type BoolType struct {
	Value bool `json:"value"`
}

type NumberType struct {
	Value int32 `json:"value"`
}

type TextType struct {
	Value string `json:"value"`
}
