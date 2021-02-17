package helpers

import (
	"github.com/microcosm-cc/bluemonday"
)

func CleanStringSpecials(input string) string {
	return bluemonday.UGCPolicy().Sanitize(input)
}
