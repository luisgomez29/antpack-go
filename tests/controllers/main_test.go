package controllers

import (
	"os"
	"testing"

	"github.com/luisgomez29/antpack-go/tests"
)

// TestMain runs all the tests within the package.
func TestMain(t *testing.M) {
	tests.Init()
	os.Exit(t.Run())
}
