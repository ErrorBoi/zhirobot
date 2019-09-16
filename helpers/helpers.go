package helpers

import "log"

// PanicIfErr is a "syntax sugar" made to avoid `if err != nil {}`
func PanicIfErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
