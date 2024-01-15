package goenv_test

import (
	"os"
	"testing"

	goenv "github.com/Merapi-Tani-Instrument/GoEnv"
)

func TestLoadEnv(t *testing.T) {
	f, err := os.Open("./sample.yaml")
	if err != nil {
		panic(err)
	}
	goenv.Parse(f)
}
