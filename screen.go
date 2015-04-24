package main

import (
    "fmt"
)

func clear(extra int) {
    for i:=0;i<(80-extra);i+=1 {
        fmt.Print(" ")
    }
}
