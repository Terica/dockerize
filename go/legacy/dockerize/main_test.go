package main

import (
	"testing"
)

var mockFileContent = ""

func mockReadFile(path string) ([]byte,error) {
	return []byte(mockFileContent), nil
}

func TestOsdetect(t *testing.T) {
	oldReadFile := ReadFile
	defer func() { ReadFile = oldReadFile }()
	goos := "windows"
	ReadFile = mockReadFile
	theos := osdetect()
	if theos != goos {
		t.Fatalf("expected %s, got %s", goos, theos)
	}
	mockFileContent="Ubuntu"
	goos = "linux"
	theos = osdetect()
	if theos != goos {
		t.Fatalf("expected %s, got %s", goos, theos)
	}
	mockFileContent="Microsoft"
	theos = osdetect()
	if theos != "linux/windows" {
		t.Fatalf("expected %s/windows, got %s", goos, theos)
	}
}

func TestHashify(t *testing.T) {
	list := []string{"first","second","third"}
	hash := hashify(list)
	if len(list) != len(hash) {
		t.Fatalf("Expected the length of list and hash to be the same")
	}
	for _, s := range list {
		if _, ok := hash[s]; !ok {
			t.Fatalf("expected hash to contain %s, it did not", s)
		}
	}
}

func TestExclude(t *testing.T) {
	list := []string{"first","second","third"}
	exclmap := hashify(list)
	list = []string{"not=","in=","here="}
	out := exclude(list, exclmap)
	if len(out) != len(list) {
		t.Fatalf("Expected not to remove items")
	}
	for pos := range list {
		if list[pos] != out[pos] {
			t.Fatalf("Expected %s = %s", list[pos],out[pos])
		}
	}
	list = []string{"not=","first=","here="}
	out = exclude(list, exclmap)
	if len(out) != len(list)-1 {
		t.Fatalf("Expected to remove item")
	}
	exclmap = hashify(out)
	if _, ok := exclmap["first="]; ok {
		t.Fatalf("Expected \"first=\" to be the removed item.")
	}
}
