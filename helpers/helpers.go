package helpers

import (
	"log"
	"strconv"
)

// PanicIfErr is a "syntax sugar" made to avoid `if err != nil {}`
func PanicIfErr(err error) {
	if err != nil {
		log.Panic("Error:", err)
	}
}

// IsFloat checks if string is convertible to Float
func IsFloat(s string) bool {
	if _, err := strconv.ParseFloat(s, 64); err == nil {
		return true
	}
	return false
}
