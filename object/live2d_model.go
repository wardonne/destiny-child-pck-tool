package object

type Live2DModel struct {
	Version     *string                     `json:"version"`
	Model       *string                     `json:"model"`
	Textures    *[]string                   `json:"textures"`
	Expressions *[]*Live2DExpression        `json:"expressions"`
	Motions     *map[string][]*Live2DMotion `json:"motions"`
}
