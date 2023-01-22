package main

import (
	"os"
	"testing"
)

func Test_devFunc(t *testing.T) {
	want := "https://vk.com/xsrghy"
	got := os.Args[1]
	if want != got {
		t.Error("Не прошло,братуха")
	}
}
