package utils

import (
	"log"
	"strconv"
)

func Str2Int(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal("could not parse port", err)
	}

	return i
}
