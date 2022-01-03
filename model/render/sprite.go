package render

import "encoding/json"

type Sprite struct {
	Name       string
	SpriteId   string
	InternalId string
	Pivot      json.RawMessage
	Border     json.RawMessage
	Rect       json.RawMessage
	Alignment  int
}
