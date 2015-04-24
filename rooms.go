package main

import (
    "svi"
    "fmt"
    "unicode/utf8"
sc  "strconv"
)



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

/* need to set what was under player to var fill rune */
func placeCharacter(x, y int, pic rune) {
    str := curroom[y]
    var newstr string = ""

	for i:=0;len(str) > 0;i+=1 {
        character, size := utf8.DecodeRuneInString(str)
		str = str[size:]
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
