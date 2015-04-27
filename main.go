package main

import (
    "fmt"
    "os"
    "unicode/utf8"
)

/* terminal height and width */
const height, width = 24, 80
//const width = 80
/* room info, roomb is [room num][height][width] */
var curroom = make([]string, 23)
var roomb [10][24][81]string
type position struct {
    x int
    y int
}
var pos position
var icon, fill rune = '👱', ' '
//var fill rune = ' '
//👳
var fillU, fillD, fillL, fillR rune

var extra int
var dir string
var fut int
var moved bool = false

func main() {
	var b []byte = make([]byte, 1)
    clear(0)
    setRoom("1")
    pos.x=5
    pos.y=1
    for ;string([]byte(b)[0]) != "q"; {
        dir_fill := checkObstruction()
        if safeMove == true {
            placeCharacter(pos.x, pos.y, icon)
            moved = false
        } else {
            if dir == "Up" {
                if moved == false {
                    placeCharacter(pos.x, pos.y, dir_fill)
                    pos.y +=1
                    moved = true
                }
            } else if dir == "Down" {
                if moved == false {
                    placeCharacter(pos.x, pos.y, dir_fill)
                    pos.y -=1
                    moved = true
                }
            } else if dir == "Left" {
                if moved == false {
                    placeCharacter(pos.x, pos.y, dir_fill)
                    pos.x +=1
                    moved = true
                }
            } else if dir == "Right" {
                if moved == false {
                    placeCharacter(pos.x, pos.y, dir_fill)
                    pos.x -=1
                    moved = true
                }
            } else {
                placeCharacter(pos.x, pos.y, fill)
                moved = false
                safeMove = true
            }
            placeCharacter(pos.x, pos.y, icon)
        }
        printRoom()
        //fmt.Print("User Stats:")
        //clear(11)
        fmt.Printf("Position: %2d,%2d; UDRL: %c,%c,%c,%c", pos.x, pos.y, fillU, fillD, fillR, fillL)
        //clear(15) //without "; Fills: ..."
        clear(30)
        os.Stdin.Read(b)
        char := string([]byte(b)[0])
        switch char {
            case "a":
                dir="Left"
                placeCharacter(pos.x, pos.y, fill)
                fut = pos.x - 1
                //checkBarricades(dir)
                if (fut) < 0 || safeMove == false {
                    pos.x = pos.x
                } else {
                    pos.x = fut
                }

            case "d":
                dir="Right"
                placeCharacter(pos.x, pos.y, fill)
                fut = pos.x + 1
                //checkBarricades(dir)
                if (fut) > 79 || safeMove == false {
                    pos.x = pos.x
                } else {
                    pos.x = fut
                }

            case "w":
                dir="Up"
                placeCharacter(pos.x, pos.y, fill)
                fut = pos.y - 1
                //checkBarricades(dir)
                if (fut) < 0 || safeMove == false {
                    pos.y = pos.y
                } else {
                    pos.y = fut
                }

            case "s":
                dir="Down"
                placeCharacter(pos.x, pos.y, fill)
                fut = pos.y + 1
                //checkBarricades(dir)
                if (fut) > 22 || safeMove == false {
                    pos.y = pos.y
                } else {
                    pos.y = fut
                }

            case "o":
                placeCharacter(pos.x, pos.y, fill)
                openDoors()

            default:
                placeCharacter(pos.x, pos.y, fill)
                dir=char
        }
        fmt.Print(dir)
        extra = utf8.RuneCountInString(dir) //number of runes in string
        clear(extra)
    }
    /* testing printing a room */
    printRoom()
    copyRoom("1")
    //printRoomB(1)
}
