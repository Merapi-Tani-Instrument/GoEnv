package goenv_test

import (
	"Merapi-Tani-Instrument/goenv"
	"os"
	"testing"
)

func TestLoadEnv(t *testing.T) {
	f, err := os.Open("./sample.yaml")
	if err != nil {
		panic(err)
	}
	goenv.Parse(f)
}
