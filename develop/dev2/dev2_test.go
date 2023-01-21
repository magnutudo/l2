package main

import "testing"

func TestPackedString_Unpacking_First(t *testing.T) {
	var ps PackedString
	var got string
	var want string
	ps = "abcd"
	want = "abcd"
	got = ps.Unpack()
	if got != want {
		t.Error("Неа, не верно")
	}
}
func TestPackedString_Unpacking_Second(t *testing.T) {
	var ps PackedString
	var got string
	var want string
	ps = "a4bc2d5e"
	want = "aaaabccddddde"
	got = ps.Unpack()
	if got != want {
		t.Error("Неа, не верно")
	}
}
