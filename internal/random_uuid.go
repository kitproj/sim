package internal

import "github.com/google/uuid"

func randomUUID() string {
	random, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	return random.String()
}
