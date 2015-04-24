package main

import (
    "fmt"
    "os"
    "unicode/utf8"
)

/* terminal height and width */
const height = 24
const width = 80
/* room info, roomb is [room num][height][width] */
var curroom = make([]string, 23)
var roomb [10][24][81]string
type position struct {
    x int
    y int
}
var pos position

var extra int
var dir string

func main() {
	var b []byte = make([]byte, 1)
    clear(0)
    setRoom("1")
    pos.x=0
    pos.y=0
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
    copyRoom("1")
    //printRoomB(1)
}
