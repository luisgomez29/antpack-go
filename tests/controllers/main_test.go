package controllers

import (
	"os"
	"testing"

	"github.com/luisgomez29/antpack-go/tests"
)

func TestMain(t *testing.M) {
	tests.Init()
	code := t.Run()

	os.Exit(code)
}
