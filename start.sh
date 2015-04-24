#!/bin/bash
go build
stty raw
./ThatOneAdventureGame
stty cooked

