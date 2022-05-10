package src

import (
	"embed"
)

//go:embed avros
var AvrosSchemaDir embed.FS

//go:embed protos/onboarding
var ProtoSchemaDir embed.FS
