package random

import (
	"math/rand"
	"time"
)

const (
	letters       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lettersLength = int32(len(letters))
)

func String(randomString *string, length int) {
	for i := 0; i < length; i++ {
		*randomString += getRandomLetter()
	}
}

func getRandomLetter() string {
	source := rand.NewSource(time.Now().UTC().UnixNano())
	random := rand.New(source)
	return string(letters[random.Int31n(lettersLength)])
}
