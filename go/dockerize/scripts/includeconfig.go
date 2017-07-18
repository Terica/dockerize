package main

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
)

var conversionMap = map[string]string{".json": "JSON"}

func translateToVariableName(filename string) (string, bool) {
	for suffix, translation := range conversionMap {
		if strings.HasSuffix(filename, suffix) {
			return strings.TrimSuffix(filename, suffix) + translation, true
		}
	}
	return filename, false
}

// Reads all .json files in the current folder
// and encodes them as strings literals in jsonfiles.go
func main() {
	fs, _ := ioutil.ReadDir(".")
	out, _ := os.Create("jsonfiles.go")
	out.Write([]byte("package main \n\nconst (\n"))
	for _, f := range fs {
		if variable, replaced := translateToVariableName(f.Name()); replaced {
			out.Write([]byte(variable + " = `"))
			f, _ := os.Open(f.Name())
			io.Copy(out, f)
			out.Write([]byte("`\n"))
		}
	}
	out.Write([]byte(")\n"))
}
