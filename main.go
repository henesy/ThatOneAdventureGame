package main

import (
    "fmt"
    "svi"
    "os"
    "unicode/utf8"
    "flag"
sc  "strconv"
    "golang.org/x/crypto/ssh/terminal"
)

type position struct {
    x int
    y int
}
/*
type model interface {
    MoveUp()(int, int)
    MoveDown()(int, int)
    MoveLeft()(int, int)
    MoveRight()(int, int)
}

func (p position) moveUp()(int, int) {
    return p.x, p.y-1
}
func (p position) moveDown()(int, int) {
    return p.x, p.y+1
}
func (p position) moveLeft()(int, int) {
    return p.x-1, p.y
}
func (p position) moveRight()(int, int) {
    return p.x+1, p.y
}
*/
type fillers struct {
    icon rune
    fill rune
    fillU rune
    fillD rune
    fillL rune
    fillR rune
}

type sprite struct {
    p position
    f fillers
}

var height, width int /* terminal height/width */
var curroom = make([]string, 23)
var pos, fut position
var char fillers
var debugmode bool = false
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

func clearnum(num int) {
    for i:=0;i<num;i+=1 {
        fmt.Print(" ")
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
            if i == fut.x {
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
            if i == fut.x {
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
    var icon_string string
    flag.StringVar(&icon_string, "icon", "👱", "Set unicode character to use as player icon")
    flag.IntVar(&height, "height", 24, "Set height of terminal screen [24]")
    flag.IntVar(&width, "width", 80, "Set width of terminal screen [80]")
    flag.Parse()

    oldState, _ := terminal.MakeRaw(0)
    defer terminal.Restore(0, oldState)

    char.icon, _ = utf8.DecodeRuneInString(icon_string)
    char.fill = ' '
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
            case "a":
                if char.fillL != '⚠' && (check(pos.x-1, pos.y, char.fillL) == false){
                    fut.x -=1
                    if fut.x < 0 {
                        fut.x +=1
                    }
                }
            case "s":
                if char.fillD != '⚠' && (check(pos.x, pos.y+1, char.fillD) == false){
                    fut.y +=1
                    if fut.y > 22 {
                        fut.y -=1
                    }
                }
            case "d":
                if char.fillR != '⚠' && (check(pos.x+1, pos.y, char.fillR) == false){
                    fut.x +=1
                    if fut.x > 79 {
                        fut.x -=1
                    }
                }
            case "o":
                /* open doors */
                if char.fillU == '-' {
                    placeRune(pos.x,pos.y-1,'ˉ')
                } else if char.fillU == 'ˉ' {
                    placeRune(pos.x,pos.y-1,'-')
                }
                if char.fillL == '|' {
                    placeRune(pos.x-1,pos.y,'\\')
                } else if char.fillL == '\\' {
                    placeRune(pos.x-1,pos.y,'|')
                }
                if char.fillD == '-' {
                    placeRune(pos.x,pos.y+1,'_')
                } else if char.fillD == '_' {
                    placeRune(pos.x,pos.y+1,'-')
                }
                if char.fillR == '|' {
                    placeRune(pos.x+1,pos.y,'/')
                } else if char.fillR == '/' {
                    placeRune(pos.x+1,pos.y,'|')
                }
            case "i":
                /* read inventory */
                clearscrn()
                fmt.Print("╔")
                for i:=0;i<width-2;i+=1 {
                    fmt.Print("═")
                }
                fmt.Print("╗")
                fmt.Print("║")
                for i:=0;i<33;i+=1 {
                    fmt.Print(" ")
                }
                fmt.Print("║ Backpack ║")
                for i:=0;i<33;i+=1 {
                    fmt.Print(" ")
                }
                fmt.Print("║")
                fmt.Print("║")
                for i:=0;i<33;i+=1 {
                    fmt.Print(" ")
                }
                fmt.Print("╚")
                for i:=0;i<10;i+=1 {
                    fmt.Print("═")
                }
                fmt.Print("╝")
                for i:=0;i<33;i+=1 {
                    fmt.Print(" ")
                }
                fmt.Print("║")
                /* body of inventory */
                for i:=0;i<height-5;i+=1 {
                    fmt.Print("║")
                    clearln(2)
                    fmt.Print("║")
                }
                /* end body of inventory */
                fmt.Print("╚")
                for i:=0;i<width-2;i+=1 {
                    fmt.Print("═")
                }
                fmt.Print("╝")
                clearln(0)
                fmt.Scanln()
            case "D":
                /* debug mode */
                if debugmode == false {
                    debugmode = true
                } else {
                    debugmode = false
                }
            default:
        }
        placeRune(pos.x, pos.y, char.fill)

        char.fill, char.fillU, char.fillL, char.fillD, char.fillR = placeRune(fut.x, fut.y, char.icon)
        printRoom()
        if s:=utf8.RuneCountInString(usrin); debugmode == false {
            fmt.Printf("Position: %2d,%2d; ULDR: %c,%c,%c,%c", fut.x, fut.y, char.fillU, char.fillL, char.fillD, char.fillR)
            clearln(30)
        } else {
            fmt.Printf("Position: %2d,%2d; ULDR: %c,%c,%c,%c; Key: %s", fut.x, fut.y, char.fillU, char.fillL, char.fillD, char.fillR, usrin)
            clearln(30+7+s)
        }

        pos.x, pos.y = fut.x, fut.y
        }

    fmt.Println("")
}
