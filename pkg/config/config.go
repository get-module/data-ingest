package config

import (
	"fmt"
	"os"
)

func env(key string) (string) {
	check_existence := os.LookupEnv(key)

	if (check_existence){
		return os.Get(key) 
	}
	return "Error: key not found in .env"
}
