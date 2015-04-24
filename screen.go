package main

import (
    "fmt"
    "unicode/utf8"
)

func clear(extra int) {
    for i:=0;i<(width-extra);i+=1 {
        fmt.Print(" ")
    }
}

/* would be printed after inserting player into the map */
func printRoom() {
    for i:=0;i<len(curroom);i+=1 {
        fmt.Printf("%s", curroom[i])
        extra = utf8.RuneCountInString(curroom[i])
        clear(extra)
    }
}

/* broken */
func printRoomB(loc int) {
    for i:=0;i<22;i+=1 {
        for h:=0;h<79;i+=1 {
            fmt.Printf("%s", roomb[loc][i][h])
            extra = 80
            clear(extra)
        }
    }
}
