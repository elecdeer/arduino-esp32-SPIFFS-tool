package common

import (
	"log"
)

func ApplyLogStyle() {
	log.SetFlags(0)
	log.SetPrefix(" ")
}
