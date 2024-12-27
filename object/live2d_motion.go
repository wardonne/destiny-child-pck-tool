package object

import "time"

type Live2DMotion struct {
	File    string        `json:"file"`
	FadeIn  time.Duration `json:"fade_in,omitempty"`
	FadeOut time.Duration `json:"fade_out,omitempty"`
}
