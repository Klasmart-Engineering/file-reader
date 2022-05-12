package src

import (
	"embed"
)

//go:embed avros
var AvrosSchemaDir embed.FS

//go:embed protos/onboarding/*.proto
var ProtoSchemaDir embed.FS
