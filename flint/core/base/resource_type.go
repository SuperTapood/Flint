package base

import "github.com/google/uuid"

type ResourceType interface {
	String() string
	Synth() (uuid.UUID, map[string]any)
}
