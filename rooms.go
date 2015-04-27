package main

import (
    "svi"
    "fmt"
    "unicode/utf8"
sc  "strconv"
)

/* end list in a comma */
/* var barricades = []rune{
    '═', '╣', '║', '╗', '╝', '╚', '╔', '╩', '╦', '╠', '╬', '┼', '┘', '┌', '|',
     '-', '│', '┤', '┐', '└', '┴', '├', '─', '┬'
} */

var safeMove bool = true

/* sets current room map */
func setRoom(num string) {
    room, succ := svi.Filereader(num + ".room")
    if succ == 1 {
        fmt.Print("ERROR READING ROOM FILE")
    }
    for h:=0;h<23;h+=1 {
        curroom[h]=room[h]
    }
}

func copyRoom(num string) {
    /* room, succ := svi.Filereader(num + ".room")
    if succ == 1 {
        fmt.Print("ERROR READING ROOM FILE")
    } */
    for h:=0;h<23;h+=1 {
        for i:=0;i<80;i+=1 {
                place, _ := sc.Atoi(num)
                character, _ := utf8.DecodeRuneInString(curroom[h])
                roomb[place][h][i] = sc.QuoteRune(character)
        }
    }
}


/* need to set what was under player to var fill rune */
func placeCharacter(x, y int, pic rune) {
    str := curroom[y]
    var newstr string = ""
    //var clearL, clearR, clearU, clearD bool
    var posL, posR, posU, posD int
    posL = pos.x -1
    posR = pos.x +1
    posU = pos.y -1
    posD = pos.y +1

    /* take care of character placement */
	for i:=0;len(str) > 0;i+=1 {
        character, size := utf8.DecodeRuneInString(str)
		str = str[size:]
        //findBarricades(x, y, character)
        if i == x {
            fill=character
            letter, _ := sc.Unquote(sc.QuoteRune(pic))
            newstr = newstr + letter
        } else {
            letter, _ := sc.Unquote(sc.QuoteRune(character))
            newstr = newstr + letter
        }
	}
    curroom[y] = newstr

    if pic == icon {
        /* scan upper line for fillU */
        if posU < 0 {
            fillU='⚠'
        } else {
            str=curroom[posU]
            for i:=0;len(str) > 0;i+=1 {
                character, size := utf8.DecodeRuneInString(str)
        		str = str[size:]
                //findBarricades(x, y, character)
                if i == pos.x {
                    fillU=character
                }
            }
        }
        /* scan lower line for fillD */
        if posD > 22 {
            fillD = '⚠'
        } else {
            str=curroom[posD]
            for i:=0;len(str) > 0;i+=1 {
                character, size := utf8.DecodeRuneInString(str)
                str = str[size:]
                //findBarricades(x, y, character)
                if i == pos.x {
                    fillD=character
                }
            }
        }
        /* scan same line for right character */
        str=curroom[y]
        if posR > 79 {
            fillR = '⚠'
        } else {
            for i:=0;len(str) > 0;i+=1 {
                character, size := utf8.DecodeRuneInString(str)
                str = str[size:]
                //findBarricades(x, y, character)
                if i == posR {
                    fillR=character
                }
            }
        }
        /* scan same line for left character */
        str=curroom[y]
        if posL < 0 {
            fillL = '⚠'
        } else {
            for i:=0;len(str) > 0;i+=1 {
                character, size := utf8.DecodeRuneInString(str)
                str = str[size:]
                //findBarricades(x, y, character)
                if i == posL {
                    fillL=character
                }
            }
        }
    }
}


// use pos.x and pos.y, then direction, to find the curroom line, parse, and check
func checkBarricades(direction string) {
    var fut int
    var str string
    switch direction {
        case "Up":
            fut = pos.y - 1
            str = curroom[fut]
        case "Down":
            fut = pos.y + 1
            str = curroom[fut]
        case "Left":
            fut = pos.x - 1
            str = curroom[pos.y]
        case "Right":
            fut = pos.x + 1
            str = curroom[pos.y]
        default:
    }

    for i:=0;len(str) > 0;i+=1 {
        character, size := utf8.DecodeRuneInString(str)
        obstruction := findBarricades_old(character)
        if direction == "Up" && i == pos.x && obstruction == true {
            safeMove = false
        } else {
            safeMove = true
        }
        str = str[size:]
    }
}


// ¦ <- open door; | and - are closed doors
func findBarricades_old(char rune)(bool) {
    if char == '═' || char == '╣' || char == '║' || char == '╗' || char == '╝' || char == '╚' || char == '╔' || char == '╩' || char == '╦' || char == '╠' || char == '╬' {
        return true
    } else if char == '┼' || char == '┘' || char == '┌' || char == '|' || char == '-' || char == '│' || char == '┤' || char == '┐' || char == '└' || char == '┴' {
        return true
    } else if char == '├' || char == '─' || char == '┬'{
        return true
    } else {
        return false
    }
}

func checkObstruction()(rune) {
    /* use dir and fill? and for to scan */
    barricades := []rune{'═', '╣', '║', '╗', '╝', '╚', '╔', '╩', '╦', '╠', '╬', '┼', '┘', '┌', '|',
         '-', '│', '┤', '┐', '└', '┴', '├', '─', '┬'}
    var dir_fill rune

    if dir == "Left" {
        dir_fill = fillL
    } else if dir == "Right" {
        dir_fill = fillR
    } else if dir == "Up" {
        dir_fill = fillU
    } else if dir == "Down" {
        dir_fill = fillD
    }

    //inc := 0
    for _, char:=range barricades {
        //fmt.Printf("%c %c", dir_fill, barricades[i])
        if dir_fill == char {
            safeMove = false
            break
        } else {
            safeMove = true
        }
        //fmt.Printf("%c, %c, %v", char, dir_fill, safeMove)
    }
    return dir_fill
}

func openDoors() {
    if fillU == '-' {
        placeCharacter(pos.x, pos.y-1, 'ˉ')
        fillU = 'ˉ'
        fill = ' '
    } else if fillD == '-' {
        placeCharacter(pos.x, pos.y+1, '_')
        fillD = '_'
        fill = ' '
    } else if fillL == '|' {
        placeCharacter(pos.x-1, pos.y, '\\')
        fillL = '\\'
        fill = ' '
    } else if fillR == '|' {
        placeCharacter(pos.x+1, pos.y, '/')
        fillR = '/'
        fill = ' '
    } 

}
