package internal

import (
	"log"
)

var console = map[string]any{
	"log": func(args ...any) {
		log.Println(append([]any{"console:"}, args...)...)
	},
}
