package docs

import (
	"github.com/mvrilo/go-redoc"
)

func Initialize() redoc.Redoc {
	doc := redoc.Redoc{
		Title:       "Documentation of BlogSystemAPI",
		Description: "Documentation describes working procedures of BlogSystemAPI like structs, handlers, etc.",
		SpecFile:    "./docs/proto.swagger.json",
		SpecPath:    "/docs/proto.swagger.json",
		DocsPath:    "/docs",
	}

	return doc
}
