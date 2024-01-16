package goenv_test

import (
	"fmt"
	"os"
	"testing"

	goenv "github.com/Merapi-Tani-Instrument/GoEnv"
)

type EnvTest struct {
	Coba  string `env:"spring.application.name"`
	Coba1 bool   `env:"spring.main.web-environment"`
}

func TestLoadEnv(t *testing.T) {
	f, err := os.Open("./sample.yaml")
	if err != nil {
		panic(err)
	}
	gE, err := goenv.Parse(f)
	if err != nil {
		panic(err)
	}

	var gt EnvTest
	gE.Inject(&gt)
	fmt.Println("gt", gt)
}
