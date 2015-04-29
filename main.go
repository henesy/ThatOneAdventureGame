package main

import (
    "fmt"
    "svi"
    "os"
    "unicode/utf8"
    "flag"
sc  "strconv"
)

type position struct {
    x int
    y int
}

type fillers struct {
    icon rune
    fill rune
    fillU rune
    fillD rune
    fillL rune
    fillR rune
}

var height, width int /* terminal height/width */
var curroom = make([]string, 23)
var pos, fut position
var char fillers
/* end variables */

/* clears a line of the screen */
func clearln(extra int) {
    for i:=0;i<(width-extra);i+=1 {
        fmt.Print(" ")
    }
}

/* clears the entire screen (if full) */
func clearscrn() {
    for h:=0;h<height;h+=1 {
        for i:=0;i<width;i+=1 {
            fmt.Print(" ")
        }
    }
}

/* reads room from a file then places the room strings into the curroom[] buf */
func setRoom(num string) {
    room, succ := svi.Filereader(num + ".room")
    if succ == 1 {
        fmt.Print("ERROR READING ROOM FILE")
    }
    for h:=0;h<23;h+=1 {
        curroom[h]=room[h]
    }
}

/* prints the curroom[] buf to screen */
func printRoom() {
    for i:=0;i<len(curroom);i+=1 {
        fmt.Printf("%s", curroom[i])
        /* clearing in case the map doesn't fill the standard 23x80 width */
        extra := utf8.RuneCountInString(curroom[i])
        clearln(extra)
    }
}

/* places a rune (pic) at coordinate (x,y) */
func placeRune(x, y int, pic rune)(filler, fU, fL, fD, fR rune) {
    str := curroom[y]
    var newstr string

    for i:=0;len(str) > 0;i+=1 {
        character, size := utf8.DecodeRuneInString(str)
		str = str[size:]
        if i == x {
            filler = character
            letter, _ := sc.Unquote(sc.QuoteRune(pic))
            newstr = newstr + letter
        } else {
            letter, _ := sc.Unquote(sc.QuoteRune(character))
            newstr = newstr + letter
        }
	}
    curroom[y] = newstr

    /* check for edge of the map */
    posU, posL, posD, posR := (fut.y -1), (fut.x -1), (fut.y +1), (fut.x +1)

    if posU < 0 {
        fU='⚠'
    } else {
        str=curroom[posU]
        for i:=0;len(str) > 0;i+=1 {
            character, size := utf8.DecodeRuneInString(str)
    		str = str[size:]
            if i == pos.x {
                fU=character
            }
        }
    }
    /* scan lower line for fillD */
    if posD > 22 {
        fD = '⚠'
    } else {
        str=curroom[posD]
        for i:=0;len(str) > 0;i+=1 {
            character, size := utf8.DecodeRuneInString(str)
            str = str[size:]
            if i == pos.x {
                fD=character
            }
        }
    }
    /* scan same line for right character */
    str=curroom[y]
    if posR > 79 {
        fR = '⚠'
    } else {
        for i:=0;len(str) > 0;i+=1 {
            character, size := utf8.DecodeRuneInString(str)
            str = str[size:]
            if i == posR {
                fR=character
            }
        }
    }
    /* scan same line for left character */
    str=curroom[y]
    if posL < 0 {
        fL = '⚠'
    } else {
        for i:=0;len(str) > 0;i+=1 {
            character, size := utf8.DecodeRuneInString(str)
            str = str[size:]
            if i == posL {
                fL=character
            }
        }
    }

    return
}

/* checks for barricades at a given coordinate */
func check(x, y int, aga rune)(occ bool) {
    str := curroom[y]
    barricades := []rune{'═', '╣', '║', '╗', '╝', '╚', '╔', '╩', '╦', '╠', '╬', '┼', '┘', '┌', '|',
         '-', '│', '┤', '┐', '└', '┴', '├', '─', '┬'}
    for i:=0;len(str) > 0;i+=1 {
        _, size := utf8.DecodeRuneInString(str)
        str = str[size:]
        if i == x {
            for _, bar := range barricades {
                if aga == bar {
                    occ = true
                    break
                } else {
                    occ = false
                }
            }
            break
        }
    }
    return
}


func main() {
    flag.IntVar(&height, "height", 24, "Set height of terminal screen [24]")
    flag.IntVar(&width, "width", 80, "Set width of terminal screen [80]")
    flag.Parse()

    char.icon, char.fill = '👱', ' '
    var b []byte = make([]byte, 1)
    clearln(0)
    setRoom("1")
    pos.x, pos.y, fut.x, fut.y = 5, 1, 5, 1
    var first bool = true

    for ;string([]byte(b)[0]) != "q"; {
        if first == false {
            os.Stdin.Read(b)
        } else {
            b[0] = 32
            first=false
        }
        usrin := string([]byte(b)[0])

        switch usrin {
            case "w":
                if char.fillU != '⚠' && (check(pos.x, pos.y-1, char.fillU) == false){
                    fut.y -=1
                    if fut.y < 0 {
                        fut.y +=1
                    }

                }
                placeRune(pos.x, pos.y, char.fill)
            case "a":
                if char.fillL != '⚠' && (check(pos.x-1, pos.y, char.fillL) == false){
                    fut.x -=1
                    if fut.x < 0 {
                        fut.x +=1
                    }

                }
                placeRune(pos.x, pos.y, char.fill)
            case "s":
                if char.fillD != '⚠' && (check(pos.x, pos.y+1, char.fillD) == false){
                    fut.y +=1
                    if fut.y > 22 {
                        fut.y -=1
                    }

                }
                placeRune(pos.x, pos.y, char.fill)

            case "d":
                if char.fillR != '⚠' && (check(pos.x+1, pos.y, char.fillR) == false){
                    fut.x +=1
                    if fut.x > 79 {
                        fut.x -=1
                    }

                }
                placeRune(pos.x, pos.y, char.fill)
            default:

        }

        char.fill, char.fillU, char.fillL, char.fillD, char.fillR = placeRune(fut.x, fut.y, char.icon)

        printRoom()
        fmt.Printf("Position: %2d,%2d; ULDR: %c,%c,%c,%c", pos.x, pos.y, char.fillU, char.fillL, char.fillD, char.fillR)
        clearln(30)
        pos.x, pos.y = fut.x, fut.y
        }

    fmt.Println("")
}
