package main

import (
	"os"
	"strings"
)

func getEnvironment() map[string]string {
	var envMap = make(map[string]string)

	env := os.Environ()
	for _, e := range env {
		splitted := strings.Split(e, "=")

		key := splitted[0]
		value := ""

		if len(splitted) > 1 {
			value = strings.Join(splitted[1:], "=")
		}

		envMap[key] = value
	}

	return envMap
}
