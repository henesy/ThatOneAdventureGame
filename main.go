package main

import (
    "fmt"
    "svi"
    "os"
    "unicode/utf8"
    "flag"
sc  "strconv"
    "golang.org/x/crypto/ssh/terminal"
    "strings"
    //"math"
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

type statistics struct {
    hlth int
    atk int
    dfs int
}

type sprite struct {
    pos position
    fut position
    f fillers
    s statistics
}

var height, width int /* terminal height/width */
var curroom = make([]string, 23)
var roomdata = make([]string, 23)
var numsprites int
var sprites = make([]sprite, 21)
var pos, fut position
var char fillers
var plyr statistics
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
    var count int
    room, succ := svi.Filereader(num + ".room")
    if succ == 1 {
        fmt.Print("ERROR READING ROOM FILE")
    }
    for h:=0;h<23;h+=1 {
        curroom[h]=room[h]
        count = h
    }
    /* extract data lines*/
    count+=1 //23, but line #24
    roomdata[0] = room[count]
    if roomdata[0] == "Data:" {
        count+=1 //24, but line #25
        roomdata[1] = room[count]
        num, _:=sc.Atoi(roomdata[1]) //the number 3
        numsprites=num
        for i:=2;i<(num+2);i+=1 { //starting at 2, until we reach 4 (2,3,4)
            count+=1//25, but line #26
            roomdata[i] = room[count]
        }
        /* for the number of sprites, do this for reach sprite:
        roomdata[] starts at [2] for being relevant (0,1 being Data: and 3) */
        //var buf = make([]int, 5) //account for 3 plyr plus (x,y)
        for j:=0;j<num;j+=1 {
            str:=roomdata[j+2] //starts at [2]
            newstr:=strings.Split(str, ",")
            sprites[j].f.icon, _ = utf8.DecodeRuneInString(newstr[0])
            sprites[j].s.hlth, _ = sc.Atoi(newstr[1])
            sprites[j].s.atk, _ = sc.Atoi(newstr[2])
            sprites[j].s.dfs, _ = sc.Atoi(newstr[3])
            sprites[j].pos.x, _ = sc.Atoi(newstr[4])
            sprites[j].pos.y, _ = sc.Atoi(newstr[5])
            tmpstr:=curroom[sprites[j].pos.y]
            var char rune
            for i:=0;len(tmpstr) > 0;i+=1 {
                character, size := utf8.DecodeRuneInString(tmpstr)
        		tmpstr = tmpstr[size:]
                if i == sprites[j].pos.x {
                    char=character
                }
            }
            sprites[j].f.fill = char

        }
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

func populateCreeps() {
    for i:=0;i<numsprites;i+=1 {
        sprites[i].f.fill, sprites[i].f.fillU, sprites[i].f.fillL, sprites[i].f.fillD, sprites[i].f.fillR = placeRune(sprites[i].pos.x, sprites[i].pos.y, sprites[i].f.icon, i)
        sprites[i].fut.x, sprites[i].fut.y = sprites[i].pos.x, sprites[i].pos.y
    }
}

func moveCreeps() {
    bufX, bufY := make([]int, numsprites+1), make([]int, numsprites+1)
    for h:=0;h<numsprites;h+=1 {
        bufX[h] = sprites[h].fut.x
        bufY[h] = sprites[h].fut.y
    }
    bufX[numsprites] = fut.x
    bufY[numsprites] = fut.y //add the player's coords

    /* set initial direction */
    for i:=0;i<numsprites;i+=1 {
        dirX, dirY := make([]string, numsprites), make([]string, numsprites)

            placeRune(sprites[i].pos.x, sprites[i].pos.y, sprites[i].f.fill, i)
            if sprites[i].fut.x > fut.x {
                dirX[i] = "Left"
            } else if sprites[i].fut.x < fut.x {
                dirX[i] = "Right"
            } else {
                dirX[i] = "None"
            }
            if sprites[i].fut.y > fut.y {
                dirY[i] = "Up"
            } else if sprites[i].fut.y < fut.y {
                dirY[i] = "Down"
            } else {
                dirY[i] = "None"
            }

        //fmt.Print(dirX[i], dirY[i])
        /* nullify movement if close to something */
        /* make math.Mod() to numX numY format */

            var res int
            var xCanc, yCanc bool
            for _, num:=range bufX {
                //fut.x and num
                //res := math.Mod(floatX, floatnum)
                if fut.x > num {
                    res = fut.x - num
                } else if num > fut.x {
                    res = num - fut.x
                }

                if num == sprites[i].fut.x {
                    dirX[i] = dirX[i]
                    xCanc = false
                } else if res < 4 {
                    //dirX[i] = "None"
                    xCanc = true
                }
            }
                //fmt.Print("X: ", floatX,floatnum,res)
            for _, num:=range bufY {
                //fut.y and num
                //res := math.Mod(floatY, floatnum)
                if fut.y > num {
                    res = fut.y - num
                } else if num > fut.y {
                    res = num - fut.y
                }

                if num == sprites[i].fut.y {
                    dirY[i] = dirY[i]
                    yCanc = false
                } else if res < 4 {
                    //dirY[i] = "None"
                    yCanc = true
                }
                //fmt.Print("Y: ", floatY,floatnum,res)
            }
            if xCanc == true && yCanc == true {
                dirX[i] = "None"
                dirY[i] = "None"
            } else if xCanc == true && yCanc == false {
                dirY[i] = dirY[i]
                dirX[i] = "None"
            } else if xCanc == false && yCanc == true {
                dirX[i] = dirX[i]
                dirY[i] = "None"
            }


        //fmt.Print(dirX[i], dirY[i])

        /* nullify movement if check fails, character adjacent, or movement close */

            if check(sprites[i].fut.x-1, sprites[i].fut.y, sprites[i].f.fillL) == true || check(sprites[i].fut.x+1, sprites[i].fut.y, sprites[i].f.fillR) == true {
                dirX[i] = "None"
            }
            if check(sprites[i].fut.x, sprites[i].fut.y-1, sprites[i].f.fillU) == true || check(sprites[i].fut.x, sprites[i].fut.y+1, sprites[i].f.fillD) == true {
                dirY[i] = "None"
            }
            if sprites[i].f.fillL == char.icon || sprites[i].f.fillR == char.icon {
                dirX[i] = "None"
            }
            if sprites[i].f.fillD == char.icon || sprites[i].f.fillU == char.icon {
                dirY[i] = "None"
            }

            var numX int
            if sprites[i].fut.x > fut.x {
                numX = sprites[i].fut.x - fut.x
            } else if sprites[i].fut.x < fut.x {
                numX = fut.x - sprites[i].fut.x
            }


            var numY int
            if sprites[i].fut.y > fut.y {
                numY = sprites[i].fut.y - fut.y
            } else if sprites[i].fut.y < fut.y {
                numY = fut.y - sprites[i].fut.y
            }
            if numY < 4  && numX < 4{
                dirY[i] = "None"
            }
            if numX < 4 && numY < 4 {
                dirX[i] = "None"
            }

        //fmt.Print(dirY[i], dirX[i])
        if rn := svi.Random(0,2); dirY[i] != "None" && dirX[i] != "None" {
            if rn == 0 {
                dirY[i] = "None"
            } else if rn == 1 {
                dirX[i] = "None"
            }
        }
        /* translate dirX[i] dirY[i] */

            if dirX[i] == "Left" {
                sprites[i].fut.x -=1
            } else if dirX[i] == "Right" {
                sprites[i].fut.x +=1
            } else {
                sprites[i].fut.x = sprites[i].fut.x
            }
            if dirY[i] == "Up" {
                sprites[i].fut.y -=1
            } else if dirY[i] == "Down" {
                sprites[i].fut.y +=1
            } else {
                sprites[i].fut.y = sprites[i].fut.y
            }

            sprites[i].f.fill, sprites[i].f.fillU, sprites[i].f.fillL, sprites[i].f.fillD, sprites[i].f.fillR = placeRune(sprites[i].fut.x, sprites[i].fut.y, sprites[i].f.icon, i)
            sprites[i].pos.x, sprites[i].pos.y = sprites[i].fut.x, sprites[i].fut.y
    }
}

/* places a rune (pic) at coordinate (x,y) */
func placeRune(x, y int, pic rune, spritenum int)(filler, fU, fL, fD, fR rune) {
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
    var posU, posL, posD, posR int
    if pic == char.icon || spritenum == 99 {
        posU, posL, posD, posR = (fut.y -1), (fut.x -1), (fut.y +1), (fut.x +1)
    } else {
        posU, posL, posD, posR = (sprites[spritenum].fut.y-1), (sprites[spritenum].fut.x-1), (sprites[spritenum].fut.y+1), (sprites[spritenum].fut.x+1)
    }

    if posU < 0 {
        fU='⚠'
    } else {
        str=curroom[posU]
        for i:=0;len(str) > 0;i+=1 {
            character, size := utf8.DecodeRuneInString(str)
    		str = str[size:]
            if pic == char.icon || spritenum == 99 {
                if i == fut.x {
                    fU=character
                }
            } else {
                if i == sprites[spritenum].fut.x {
                    fU=character
                }
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
            if pic == char.icon || spritenum == 99 {
                if i == fut.x {
                    fD=character
                }
            } else {
                if i == sprites[spritenum].fut.x {
                    fD=character
                }
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
    /* maybe re-do this to load from sprites[i].f.icon for more goodness */
    barricades := []rune{'═', '╣', '║', '╗', '╝', '╚', '╔', '╩', '╦', '╠', '╬', '┼', '┘', '┌', '|',
         '-', '│', '┤', '┐', '└', '┴', '├', '─', '┬', char.icon}
    for i:=0;len(str) > 0;i+=1 {
        _, size := utf8.DecodeRuneInString(str)
        str = str[size:]
        if i == x {
            for _, bar := range barricades {
                if aga == bar {
                    occ = true
                    return
                } else {
                    occ = false
                }
            }
            for i:=0;i<numsprites;i+=1 {
                character:=sprites[i].f.icon
                if aga == character || aga == char.icon{
                    occ = true
                    return
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
    plyr.hlth, plyr.atk, plyr.dfs = 10, 02, 02
    var first bool = true

    for ;string([]byte(b)[0]) != "q"; {
        if first == false {
            os.Stdin.Read(b)
        } else {
            b[0] = 32
            first=false
            populateCreeps()
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
                    placeRune(pos.x,pos.y-1,'ˉ', 99)
                } else if char.fillU == 'ˉ' {
                    placeRune(pos.x,pos.y-1,'-', 99)
                }
                if char.fillL == '|' {
                    placeRune(pos.x-1,pos.y,'\\', 99)
                } else if char.fillL == '\\' {
                    placeRune(pos.x-1,pos.y,'|', 99)
                }
                if char.fillD == '-' {
                    placeRune(pos.x,pos.y+1,'_', 99)
                } else if char.fillD == '_' {
                    placeRune(pos.x,pos.y+1,'-', 99)
                }
                if char.fillR == '|' {
                    placeRune(pos.x+1,pos.y,'/', 99)
                } else if char.fillR == '/' {
                    placeRune(pos.x+1,pos.y,'|', 99)
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
        placeRune(pos.x, pos.y, char.fill, 99)
        char.fill, char.fillU, char.fillL, char.fillD, char.fillR = placeRune(fut.x, fut.y, char.icon, 99)
        moveCreeps()
        printRoom()
        if s:=utf8.RuneCountInString(usrin); debugmode == false {
            fmt.Printf("Stats: %c%2d %c %2d %c%2d", '♥', plyr.hlth, '🔥', plyr.atk, '⚔', plyr.dfs)
            clearln(19)
        } else {
            fmt.Printf("Position: %2d,%2d; ULDR: %c,%c,%c,%c; Key: %s", fut.x, fut.y, char.fillU, char.fillL, char.fillD, char.fillR, usrin)
            clearln(30+7+s)
            fmt.Printf("%c,%c,%c,%c",sprites[0].f.fillU,sprites[0].f.fillL,sprites[0].f.fillD,sprites[0].f.fillR)
            clearln(7)
            fmt.Print(sprites[0])
            clearln(57)
            fmt.Print(sprites[1])
            clearln(57)
        }
        pos.x, pos.y = fut.x, fut.y
    }

    fmt.Println("NEEIIII")
}
