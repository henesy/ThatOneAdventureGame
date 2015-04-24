#!/bin/bash
#./src/cpsvi
go build
stty raw
./ThatOneAdventureGame
stty cooked

