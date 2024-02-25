package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPass(plain []byte) string {
	hash, err := bcrypt.GenerateFromPassword(plain, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	return string(hash)
}

func VerifyPass(cyper, plain []byte) bool {
	byteCyper := []byte(cyper)
	err := bcrypt.CompareHashAndPassword(byteCyper, plain)
	if err != nil {
		return false
	} else {
		return true
	}

}
