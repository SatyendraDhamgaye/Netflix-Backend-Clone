package helpers

import "log"

func ErrorCatcher(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
