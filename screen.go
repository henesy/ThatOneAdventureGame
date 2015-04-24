package main

import (
    "fmt"
    "unicode/utf8"
)

func clear(extra int) {
    for i:=0;i<(80-extra);i+=1 {
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
