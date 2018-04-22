package model

import (
	"testing"
)

func TestTemplate(t *testing.T) {
	v, err := ApplyTemplate(`{{env "PATH"}}`)
	if err != nil {
		t.Fatal(err)
	}
	if v == "" {
		t.Fatal("PATH not found or variable not substituted.")
	}
}
