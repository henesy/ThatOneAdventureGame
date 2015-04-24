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
var icon rune = '👱'
var fill rune = ' '
//👳

var extra int
var dir string

func main() {
	var b []byte = make([]byte, 1)
    clear(0)
    setRoom("1")
    pos.x=5
    pos.y=1
    for ;string([]byte(b)[0]) != "q"; {
        placeCharacter(pos.x, pos.y, icon)
        printRoom()
        //fmt.Print("User Stats:")
        //clear(11)
        fmt.Printf("Position: %2d,%2d", pos.x, pos.y)
        clear(15)
        os.Stdin.Read(b)
        char := string([]byte(b)[0])
        switch char {
            case "a":
                dir="Left"
                placeCharacter(pos.x, pos.y, fill)
                if (pos.x - 1) < 0 {
                    pos.x = pos.x
                } else {
                    pos.x = pos.x - 1
                }
            case "d":
                dir="Right"
                placeCharacter(pos.x, pos.y, fill)
                if (pos.x + 1) > 79 {
                    pos.x = pos.x
                } else {
                    pos.x = pos.x + 1
                }
            case "w":
                dir="Up"
                placeCharacter(pos.x, pos.y, fill)
                if (pos.y - 1) < 0 {
                    pos.y = pos.y
                } else {
                    pos.y = pos.y - 1
                }
            case "s":
                dir="Down"
                placeCharacter(pos.x, pos.y, fill)
                if (pos.y + 1) > 22 {
                    pos.y = pos.y
                } else {
                    pos.y = pos.y + 1
                }
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
