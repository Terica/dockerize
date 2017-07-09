package main

import (
    "io"
    "io/ioutil"
    "os"
    "strings"
)

// Reads all .json files in the current folder
// and encodes them as strings literals in jsonfiles.go
func main() {
    fs, _ := ioutil.ReadDir(".")
    out, _ := os.Create("jsonfiles.go")
    out.Write([]byte("package main \n\nconst (\n"))
    for _, f := range fs {
        if strings.HasSuffix(f.Name(), ".json") {
            out.Write([]byte(strings.Replace(f.Name(), ".", "_",-1) + " = `"))
            f, _ := os.Open(f.Name())
            io.Copy(out, f)
            out.Write([]byte("`\n"))
        }
    }
    out.Write([]byte(")\n"))
}