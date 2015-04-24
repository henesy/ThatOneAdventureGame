package main

import (
    "fmt"
    "os"
    "unicode/utf8"
)

/* terminal height and width */
const height = 24
const width = 80
var curroom = make([]string, 23)
var extra int
var dir string

func main() {
	var b []byte = make([]byte, 1)
    clear(0)
    setRoom("1")
    for ;string([]byte(b)[0]) != "q"; {
        os.Stdin.Read(b)
        char := string([]byte(b)[0])
        switch char {
            case "a":
                dir="Left"
            case "d":
                dir="Right"
            case "w":
                dir="Up"
            case "s":
                dir="Down"
            default:
                dir=char
        }
        fmt.Print(dir)
        extra = utf8.RuneCountInString(dir)
        clear(extra)
    }
    /* testing printing a room */
    printRoom()
}
