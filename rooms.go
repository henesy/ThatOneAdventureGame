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
