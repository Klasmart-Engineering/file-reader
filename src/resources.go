package src

import (
	"embed"
)

//go:embed avros
var avrosSchemaDir embed.FS

//go:embed protos/onboarding
var protoSchemaDir embed.FS
