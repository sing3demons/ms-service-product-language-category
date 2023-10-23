package customNanoid

import (
	"github.com/aidarkhanov/nanoid/v2"
)

const (
	alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

func RandomNanoID(size int) (string, error) {
	id, err := nanoid.GenerateString(alphabet, size)
	if err != nil {
		return "", err
	}
	return id, nil
}
