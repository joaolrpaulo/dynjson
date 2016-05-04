package dynjson

import (
	"fmt"
	"os"
)

func handleErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func merge(mapA map[string]string, mapB map[string]string) (mapM map[string]string) {
	length := len(mapA) + len(mapB)
	mapM = make(map[string]string, length)
	for key, value := range mapA {
		mapM[key] = value
	}

	for key, value := range mapB {
		mapM[key] = value
	}

	return
}
