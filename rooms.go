package main

import (
    "svi"
    "fmt"
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
