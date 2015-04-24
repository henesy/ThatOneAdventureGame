package main

import (
    "svi"
    "fmt"
    "unicode/utf8"
sc  "strconv"
)

var barricadesX = make([]int, 1920)
var freeX int = 0
var barricadesY = make([]int, 1920)
var freeY int = 0
// default to -1

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
    //wipe out barricades
    for i:=0;i<len(barricadesX);i+=1 {
        barricadesX[i] = (-1)
    }
    for i:=0;i<len(barricadesY);i+=1 {
        barricadesY[i] = (-1)
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

func checkBarricades_old(locX, locY int) {
    for i:=0;i<len(barricadesX);i+=1 {
        if barricadesX[i] == locX && barricadesY[i] == locY {
            safeMove = false
        } else {
            safeMove = true
        }
    }
}

func placeCharacter_old() {
    var buf = make([]string, 81)
    var inserted string = ""
    up := pos.y
    side := pos.x
    line := curroom[up]
    for i:=0;len(curroom[up]) > 0;i+=1 {
        letter, size := utf8.DecodeRuneInString(line)
        buf[i] = sc.QuoteRune(letter)
        line = line[size:]
    }
    buf[side] = sc.QuoteRune(icon)
    for h:=0;h<len(buf);h+=1 {
        inserted = inserted + buf[h]
    }
    curroom[up] = inserted
}
