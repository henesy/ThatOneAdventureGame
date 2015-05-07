package main

import (
	"flag"
	"fmt"
	"os"
	sc "strconv"
	"strings"
	"svi"
	"unicode/utf8"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
)

/* current position in (x,y) of something in regards to the map's (x,y)
(perhaps expand to become (x,y,z)? to account for levels?) */
type position struct {
	x int
	y int
}

/* fillers are the characters stored in relation to the sprite's current (x,y) */
type fillers struct {
	icon  rune
	fill  rune
	fillU rune
	fillD rune
	fillL rune
	fillR rune
}

/* stats and characteristics of a sprite or the player */
type statistics struct {
	hlth int
	atk  int
	dfs  int
}

/* basic sprite meta-struct */
type sprite struct {
	pos position
	fut position
	f   fillers
	s   statistics
}

var height, width, numrooms int /* terminal height/width; number of room files (rooms) */
var curroom = make([]string, 23)
var roomdata = make([]string, 23)
var numsprites int
var sprites = make([]sprite, 21)
var pos, fut position
var char fillers
var plyr statistics
var debugmode bool = false
var message string

/* -- end variables -- */

/* clears a line of the screen */
func clearln(extra int) {
	for i := 0; i < (width - extra); i += 1 {
		fmt.Print(" ")
	}
}

/* clears the entire screen (if full) */
func clearscrn() {
	for h := 0; h < height; h += 1 {
		for i := 0; i < width; i += 1 {
			fmt.Print(" ")
		}
	}
}

/* clear 'num' number of spaces (on demand, limited, clearing) */
func clearnum(num int) {
	for i := 0; i < num; i += 1 {
		fmt.Print(" ")
	}
}

/* reads room from a file then places the room strings into the curroom[] buf */
func setRoom(num string) {
	var count int
	room, succ := svi.Filereader("./rooms/" + num + ".room")
	if succ == 1 {
		fmt.Print("ERROR READING ROOM FILE")
	}
	for h := 0; h < 23; h += 1 {
		curroom[h] = room[h]
		count = h
	}
	/* extract data lines*/
	count += 1 //23, but line #24
	roomdata[0] = room[count]
	/* MAYBE MAKE roomdata[0] something like POSITION COORDINATES */
	if roomdata[0] == "Data:" {
		count += 1 //24, but line #25
		/* coordinates */
		roomdata[1] = room[count]
		count +=1
		/* number of sprites */
		roomdata[2] = room[count]
		num, _ := sc.Atoi(roomdata[2]) //the number 3
		numsprites = num
		for i := 3; i < (num + 3); i += 1 { //starting at 3, until we reach 5 (2,3,4)
			count += 1 //25, but line #26
			roomdata[i] = room[count]
		}
		/* for the number of sprites, do this for reach sprite:
		   roomdata[] starts at [2] for being relevant ([0] & [1] being 'Data:'' and '3') */
		//var buf = make([]int, 5) //account for 3 plyr plus (x,y)
		for j := 0; j < num; j += 1 {
			str := roomdata[j+3] //starts at [2]
			newstr := strings.Split(str, ",")
			sprites[j].f.icon, _ = utf8.DecodeRuneInString(newstr[0])
			sprites[j].s.hlth, _ = sc.Atoi(newstr[1])
			sprites[j].s.atk, _ = sc.Atoi(newstr[2])
			sprites[j].s.dfs, _ = sc.Atoi(newstr[3])
			sprites[j].pos.x, _ = sc.Atoi(newstr[4])
			sprites[j].pos.y, _ = sc.Atoi(newstr[5])
			tmpstr := curroom[sprites[j].pos.y]
			var char rune
			for i := 0; len(tmpstr) > 0; i += 1 {
				character, size := utf8.DecodeRuneInString(tmpstr)
				tmpstr = tmpstr[size:]
				if i == sprites[j].pos.x {
					char = character
				}
			}
			sprites[j].f.fill = char

		}
		tmpstr := strings.Split(roomdata[1], ",")
		pos.x, _ = sc.Atoi(tmpstr[0])
		pos.y, _ = sc.Atoi(tmpstr[1])
	}

	str:=curroom[pos.y]
	for i:=0;i<79;i+=1 {
		character, size := utf8.DecodeRuneInString(str)
		if i == pos.x {
			char.fill = character
		}
		str=str[size:]
	}
	fut.x, fut.y = pos.x, pos.y //set initial position and reset fut.*
}

/* prints the curroom[] buf to screen */
func printRoom() {
	for i := 0; i < len(curroom); i += 1 {
		fmt.Printf("%s", curroom[i])
		/* clearing in case the map doesn't fill the standard 23x80 width */
		extra := utf8.RuneCountInString(curroom[i])
		clearln(extra)
	}
}

/* prints player stats and some debug messages if enabled q*/
func printStats(key string, usrin ...string) {
	if tots := 0; debugmode == false {
		fmt.Printf("Stats: %c%2d %c %2d %c%2d; ", '♥', plyr.hlth, '🔥', plyr.atk, '⚔', plyr.dfs)
		for _, words := range usrin {
			s := utf8.RuneCountInString(words)
			fmt.Printf("%s", words)
			tots += s
		}
		clearln(21 + tots)
	} else {
			ws:=utf8.RuneCountInString(key)
			fmt.Printf("Position: %2d,%2d; ULDR: %c,%c,%c,%c,%c; Key: %s", fut.x, fut.y, char.fill, char.fillU, char.fillL, char.fillD, char.fillR, key)
			clearln(30 + 9 + ws)
	}
}

/* returns the rune at a given position */
func getChar(x, y int)(char rune) {
	str := curroom[y]
	var size int
	for i:=0;i<len(curroom);i+=1 {
		char, size = utf8.DecodeRuneInString(str)
		if i == x {
			break
		}
		str = str[size:]
	}
	return
}

/* set original position, futures, and fills for all sprites */
func populateCreeps() {
	for i := 0; i < numsprites; i += 1 {
		sprites[i].f.fill, sprites[i].f.fillU, sprites[i].f.fillL, sprites[i].f.fillD, sprites[i].f.fillR = placeRune(sprites[i].pos.x, sprites[i].pos.y, sprites[i].f.icon, i)
		sprites[i].fut.x, sprites[i].fut.y = sprites[i].pos.x, sprites[i].pos.y
	}
}

/* decides movement for and sets sprites[i].fut.x,y for all sprites
  ABANDON ALL HOPE, YE WHO ENTER THIS FUNCTION */
func moveCreeps() {
	bufX, bufY := make([]int, numsprites+1), make([]int, numsprites+1)
	var edgeU, edgeD, edgeL, edgeR bool = false, false, false, false

	/* set initial direction */
	for i := 0; i < numsprites; i += 1 {
		dirX, dirY := make([]string, numsprites), make([]string, numsprites)
		bufX[i] = sprites[i].fut.x
		bufY[i] = sprites[i].fut.y
		bufX[numsprites] = fut.x
		bufY[numsprites] = fut.y //add the player's coords

		placeRune(sprites[i].pos.x, sprites[i].pos.y, sprites[i].f.fill, i)
		if sprites[i].fut.x > fut.x {
			dirX[i] = "Left"
		} else if sprites[i].fut.x < fut.x {
			dirX[i] = "Right"
		} else {
			dirX[i] = "None"
		}
		if sprites[i].fut.y > fut.y {
			dirY[i] = "Up"
		} else if sprites[i].fut.y < fut.y {
			dirY[i] = "Down"
		} else {
			dirY[i] = "None"
		}

		/* check for icons next to each other via fills */
		for h := 0; h < numsprites; h += 1 {
			if sprites[i].f.fillU == sprites[h].f.icon && dirY[i] == "Up" {
				dirY[i] = "Down"
			}
			if sprites[i].f.fillD == sprites[h].f.icon && dirY[i] == "Down" {
				dirY[i] = "Up"
			}
			if sprites[i].f.fillL == sprites[h].f.icon && dirX[i] == "Left" {
				dirX[i] = "Right"
			}
			if sprites[i].f.fillR == sprites[h].f.icon && dirX[i] == "Right" {
				dirX[i] = "Left"
			}
		}

		/* nullify movement if check fails, character adjacent, or movement close */
		if sprites[i].fut.x-1 > 0 {
			if check(sprites[i].fut.x-1, sprites[i].fut.y, sprites[i].f.fillL) == true {
				if dirX[i] == "Left" {
					dirX[i] = "None"
				}
			}
		} else if sprites[i].fut.x-1 <= 0 {
			if dirX[i] == "Left" {
				dirX[i] = "None"
			}
			edgeL = true
		}
		if sprites[i].fut.x+1 < 79 {
			if check(sprites[i].fut.x+1, sprites[i].fut.y, sprites[i].f.fillR) == true {
				if dirX[i] == "Right" {
					dirX[i] = "None"
				}
			}
		} else if sprites[i].fut.x+1 >= 79 {
			if dirX[i] == "Right" {
				dirX[i] = "None"
			}
			edgeR = true
		}
		if sprites[i].fut.y-1 > 0 {
			if check(sprites[i].fut.x, sprites[i].fut.y-1, sprites[i].f.fillU) == true {
				if dirY[i] == "Up" {
					dirY[i] = "None"
				}
			}
		} else if sprites[i].fut.y-1 <= 0 {
			if dirY[i] == "Up" {
				dirY[i] = "None"
			}
			edgeU = true
		}
		if sprites[i].fut.y+1 < 23 {
			// y+1 = 24 -> not in curroom[], thus error, account for == 23...
			if check(sprites[i].fut.x, sprites[i].fut.y+1, sprites[i].f.fillD) == true {
				if dirY[i] == "Down" {
					dirY[i] = "None"
				}
			}
		} else if sprites[i].fut.y+1 >= 23 {
			if dirY[i] == "Down" {
				dirY[i] = "None"
			}
			edgeD = true
		}

		/* -!- Add a check for upper right or upper left (diagonals) -!- */
		/*
		   get coords -> check x+1,y+1;x-1,y-1;x+1,y-1;x-1,y+1 ->
		   check against sprites[i].f.icon -> if x,y == [].icon ->
		   determine which sprite x,y is (sprites[i].p.x/y == ) -> note sprite num ->
		   check sprite[num].fut if num < sprites[i] (not sure if order is necessary
		   if all get checked) -> change dir to not move if num.fut.x/y matches dir=""
		*/
		/* -!- Might need to add storage for direction -!-*/
		/* maybe move this segment backwards to be adjusted moreso later as to not move extraneously? */
		for h := 0; h < numsprites; h += 1 {
			if edgeU == false && edgeR == false {
				if char:=getChar(sprites[i].fut.x+1,sprites[i].fut.y-1); char == sprites[h].f.icon && dirY[i] == "Up" {
					dirY[i] = "Down"
				}
				if char:=getChar(sprites[i].fut.x+1,sprites[i].fut.y-1); char == sprites[h].f.icon && dirX[i] == "Right" {
					dirX[i] = "Left"
				}
			}
			if edgeU == false && edgeL == false {
				if char:=getChar(sprites[i].fut.x-1,sprites[i].fut.y-1); char == sprites[h].f.icon && dirY[i] == "Up" {
					dirY[i] = "Down"
				}
				if char:=getChar(sprites[i].fut.x-1,sprites[i].fut.y-1); char == sprites[h].f.icon && dirX[i] == "Left" {
					dirX[i] = "Right"
				}
			}
			if edgeD == false && edgeR == false {
				if char:=getChar(sprites[i].fut.x+1,sprites[i].fut.y+1); char == sprites[h].f.icon && dirY[i] == "Down" {
					dirY[i] = "Up"
				}
				if char:=getChar(sprites[i].fut.x+1,sprites[i].fut.y+1); char == sprites[h].f.icon && dirX[i] == "Right" {
					dirX[i] = "Left"
				}
			}
			if edgeD == false && edgeL == false {
				if char:=getChar(sprites[i].fut.x-1,sprites[i].fut.y+1); char == sprites[h].f.icon && dirY[i] == "Down" {
					dirY[i] = "Up"
				}
				if char:=getChar(sprites[i].fut.x-1,sprites[i].fut.y+1); char == sprites[h].f.icon && dirX[i] == "Left" {
					dirX[i] = "Right"
				}
			}
		}

		for h:=0;h<numsprites;h+=1 {
			altX, altY := sprites[h].fut.x, sprites[h].fut.y
			x, y := sprites[i].fut.x, sprites[i].fut.y
			botRx, botRy := sprites[i].fut.x+1, sprites[i].fut.y+1
			botLx, botLy := sprites[i].fut.x-1, sprites[i].fut.y+1
			topRx, topRy := sprites[i].fut.x-1, sprites[i].fut.y-1
			topLx, topLy := sprites[i].fut.x+1, sprites[i].fut.y-1

			if altX == botRx && altY == botRy {
				if dirX[i] == "Right" {
					dirX[i] = "Left"
				}
				if dirY[i] == "Down" {
					dirY[i] = "Up"
				}
			}
			if altX == botLx && altY == botLy {
				if dirX[i] == "Left" {
					dirX[i] = "Right"
				}
				if dirY[i] == "Down" {
					dirY[i] = "Up"
				}
			}
			if altX == topRx && altY == topRy {
				if dirX[i] == "Right" {
					dirX[i] = "Left"
				}
				if dirY[i] == "Up" {
					dirY[i] = "Down"
				}
			}
			if altX == topLx && altY == topLy {
				if dirX[i] == "Left" {
					dirX[i] = "Right"
				}
				if dirY[i] == "Up" {
					dirY[i] = "Down"
				}
			}
			/* check for `under and over it~` 5fdp yo */
			if altY == y+1 && altX == x {
				if dirY[i] == "Down" {
					dirY[i] = "Up"
				}
			}
			if altY == y-1 && altX == x {
				if dirY[i] == "Up" {
					dirY[i] = "Down"
				}
			}
			if altY == y && altX == x+1 {
				if dirX[i] == "Right" {
					dirX[i] = "Left"
				}
			}
			if altY == y && altX == x-1 {
				if dirX[i] == "Left" {
					dirX[i] = "Right"
				}
			}
		}

		/* nullify movement if close to something */
		var res int
		var xCanc, yCanc bool
		for _, num := range bufX {
			//fut.x and num
			//res := math.Mod(floatX, floatnum)
			if fut.x > num {
				res = fut.x - num
			} else if num > fut.x {
				res = num - fut.x
			}

			if num == sprites[i].fut.x {
				dirX[i] = dirX[i]
				xCanc = false
			} else if res < 4 {
				//dirX[i] = "None"
				xCanc = true
			}
		}

		for _, num := range bufY {
			//fut.y and num
			//res := math.Mod(floatY, floatnum)
			if fut.y > num {
				res = fut.y - num
			} else if num > fut.y {
				res = num - fut.y
			}

			if num == sprites[i].fut.y {
				dirY[i] = dirY[i]
				yCanc = false
			} else if res < 4 {
				//dirY[i] = "None"
				yCanc = true
			}
		}
		if xCanc == true && yCanc == true {
			dirX[i] = "None"
			dirY[i] = "None"
		} else if xCanc == true && yCanc == false {
			dirY[i] = dirY[i]
			dirX[i] = "None"
		} else if xCanc == false && yCanc == true {
			dirX[i] = dirX[i]
			dirY[i] = "None"
		}



		/* can only be one, thus non-specific checks are okay here */
		if sprites[i].f.fillL == char.icon || sprites[i].f.fillR == char.icon {
			dirX[i] = "None"
		}
		if sprites[i].f.fillD == char.icon || sprites[i].f.fillU == char.icon {
			dirY[i] = "None"
		}

		var numX int
		if sprites[i].fut.x > fut.x {
			numX = sprites[i].fut.x - fut.x
		} else if sprites[i].fut.x < fut.x {
			numX = fut.x - sprites[i].fut.x
		}

		var numY int
		if sprites[i].fut.y > fut.y {
			numY = sprites[i].fut.y - fut.y
		} else if sprites[i].fut.y < fut.y {
			numY = fut.y - sprites[i].fut.y
		}
		if numY < 4 && numX < 4 {
			dirY[i] = "None"
		}
		if numX < 4 && numY < 4 {
			dirX[i] = "None"
		}


		/* check for edge of map (how did this not get implemented earlier?) */
		if sprites[i].f.fillU == '⚠' && dirY[i] == "Up" {
			dirY[i] = "None"
		}
		if sprites[i].f.fillD == '⚠' && dirY[i] == "Down" {
			dirY[i] = "None"
		}
		if sprites[i].f.fillR == '⚠' && dirX[i] == "Right" {
			dirX[i] = "None"
		}
		if sprites[i].f.fillL == '⚠' && dirX[i] == "Left" {
			dirX[i] = "None"
		}

		/* Pick a direction */

		if rn := svi.Random(0, 2); dirY[i] != "None" && dirX[i] != "None" {
			if rn == 0 {
				dirY[i] = "None"
			} else if rn == 1 {
				dirX[i] = "None"
			}
		}

		/* Translate dirX[i] dirY[i] */

		if dirX[i] == "Left" {
			sprites[i].fut.x -= 1
			if sprites[i].fut.x < 0 {
				sprites[i].fut.x += 1
			}
		} else if dirX[i] == "Right" {
			sprites[i].fut.x += 1
			if sprites[i].fut.x > 79 {
				sprites[i].fut.x -= 1
			}
		} else {
			sprites[i].fut.x = sprites[i].fut.x
		}
		if dirY[i] == "Up" {
			sprites[i].fut.y -= 1
			if sprites[i].fut.y < 0 {
				sprites[i].fut.y += 1
			}
		} else if dirY[i] == "Down" {
			sprites[i].fut.y += 1
			if sprites[i].fut.y > 22 {
				sprites[i].fut.y -= 1
			}
		} else {
			sprites[i].fut.y = sprites[i].fut.y
		}

		//fmt.Printf("Num: %d;%2d,%2d", i, sprites[i].fut.x, sprites[i].fut.y)
		//clearln(12)

		sprites[i].f.fill, sprites[i].f.fillU, sprites[i].f.fillL, sprites[i].f.fillD, sprites[i].f.fillR = placeRune(sprites[i].fut.x, sprites[i].fut.y, sprites[i].f.icon, i)
		sprites[i].pos.x, sprites[i].pos.y = sprites[i].fut.x, sprites[i].fut.y
	}
}

/* places a rune (pic) at coordinate (x,y) */
func placeRune(x, y int, pic rune, spritenum int) (filler, fU, fL, fD, fR rune) {
	str := curroom[y]
	var newstr string

	for i := 0; len(str) > 0; i += 1 {
		character, size := utf8.DecodeRuneInString(str)
		str = str[size:]
		if i == x {
			filler = character
			letter, _ := sc.Unquote(sc.QuoteRune(pic))
			newstr = newstr + letter
		} else {
			letter, _ := sc.Unquote(sc.QuoteRune(character))
			newstr = newstr + letter
		}
	}
	curroom[y] = newstr

	/* check for edge of the map */
	var posU, posL, posD, posR int
	if pic == char.icon || spritenum == 99 {
		posU, posL, posD, posR = (fut.y - 1), (fut.x - 1), (fut.y + 1), (fut.x + 1)
	} else {
		posU, posL, posD, posR = (sprites[spritenum].fut.y - 1), (sprites[spritenum].fut.x - 1), (sprites[spritenum].fut.y + 1), (sprites[spritenum].fut.x + 1)
	}

	if posU < 0 {
		fU = '⚠'
	} else {
		str = curroom[posU]
		for i := 0; len(str) > 0; i += 1 {
			character, size := utf8.DecodeRuneInString(str)
			str = str[size:]
			if pic == char.icon || spritenum == 99 {
				if i == fut.x {
					fU = character
				}
			} else {
				if i == sprites[spritenum].fut.x {
					fU = character
				}
			}
		}
	}
	/* scan lower line for fillD */
	if posD > 22 {
		fD = '⚠'
	} else {
		str = curroom[posD]
		for i := 0; len(str) > 0; i += 1 {
			character, size := utf8.DecodeRuneInString(str)
			str = str[size:]
			if pic == char.icon || spritenum == 99 {
				if i == fut.x {
					fD = character
				}
			} else {
				if i == sprites[spritenum].fut.x {
					fD = character
				}
			}
		}
	}
	/* scan same line for right character */
	str = curroom[y]
	if posR > 79 {
		fR = '⚠'
	} else {
		for i := 0; len(str) > 0; i += 1 {
			character, size := utf8.DecodeRuneInString(str)
			str = str[size:]
			if i == posR {
				fR = character
			}
		}
	}
	/* scan same line for left character */
	str = curroom[y]
	if posL < 0 {
		fL = '⚠'
	} else {
		for i := 0; len(str) > 0; i += 1 {
			character, size := utf8.DecodeRuneInString(str)
			str = str[size:]
			if i == posL {
				fL = character
			}
		}
	}

	return
}

/* checks for barricades at a given coordinate */
func check(x, y int, aga rune) (occ bool) {
		str := curroom[y]
	/* maybe re-do this to load from sprites[i].f.icon for more goodness */
	barricades := []rune{'═', '╣', '║', '╗', '╝', '╚', '╔', '╩', '╦', '╠', '╬', '┼', '┘', '┌', '|',
		'-', '│', '┤', '┐', '└', '┴', '├', '─', '┬', char.icon}
	for i := 0; len(str) > 0; i += 1 {
		_, size := utf8.DecodeRuneInString(str)
		str = str[size:]
		if i == x {
			for _, bar := range barricades {
				if aga == bar {
					occ = true
					return
				} else {
					occ = false
				}
			}
			for i := 0; i < numsprites; i += 1 {
				character := sprites[i].f.icon
				if aga == character || aga == char.icon {
					occ = true
					return
				} else {
					occ = false
				}
			}
			break
		}
	}
	return
}

func main() {
	var icon_string, roomnum string //👱
	var reset_message bool = true
	flag.StringVar(&roomnum, "room", "1", "Set the initial room to begin the game in")
	flag.StringVar(&icon_string, "icon", "♔", "Set unicode character to use as player icon")
	flag.IntVar(&height, "height", 24, "Set height of terminal screen [24]")
	flag.IntVar(&width, "width", 80, "Set width of terminal screen [80]")
	flag.Parse()

	oldState, _ := terminal.MakeRaw(0)
	defer terminal.Restore(0, oldState)

	char.icon, _ = utf8.DecodeRuneInString(icon_string)
	char.fill = ' '
	var b []byte = make([]byte, 1)
	clearscrn()

	/* set some initial values, get the number of room/map files we'll be playing with */
	dir, _ := ioutil.ReadDir("./rooms")
	numrooms = len(dir)
	setRoom(roomnum)
	plyr.hlth, plyr.atk, plyr.dfs = 10, 02, 02

	/* begin game loop */
	var first bool = true
	var creep_move_cnt int = 0
	for string([]byte(b)[0]) != "q" {
		if first == false {
			os.Stdin.Read(b)
		} else {
			b[0] = 32
			first = false
			populateCreeps()
		}
		usrin := string([]byte(b)[0])

		switch usrin {
		case "w":
			if char.fillU != '⚠' && (check(pos.x, pos.y-1, char.fillU) == false) {
				fut.y -= 1
				if fut.y < 0 {
					fut.y += 1
				}
			}
		case "a":
			if char.fillL != '⚠' && (check(pos.x-1, pos.y, char.fillL) == false) {
				fut.x -= 1
				if fut.x < 0 {
					fut.x += 1
				}
			}
		case "s":
			if char.fillD != '⚠' && (check(pos.x, pos.y+1, char.fillD) == false) {
				fut.y += 1
				if fut.y > 22 {
					fut.y -= 1
				}
			}
		case "d":
			if char.fillR != '⚠' && (check(pos.x+1, pos.y, char.fillR) == false) {
				fut.x += 1
				if fut.x > 79 {
					fut.x -= 1
				}
			}
		case "o":
			/* open doors */
			if char.fillU == '-' {
				placeRune(pos.x, pos.y-1, 'ˉ', 99)
			} else if char.fillU == 'ˉ' {
				placeRune(pos.x, pos.y-1, '-', 99)
			}
			if char.fillL == '|' {
				placeRune(pos.x-1, pos.y, '\\', 99)
			} else if char.fillL == '\\' {
				placeRune(pos.x-1, pos.y, '|', 99)
			}
			if char.fillD == '-' {
				placeRune(pos.x, pos.y+1, '_', 99)
			} else if char.fillD == '_' {
				placeRune(pos.x, pos.y+1, '-', 99)
			}
			if char.fillR == '|' {
				placeRune(pos.x+1, pos.y, '/', 99)
			} else if char.fillR == '/' {
				placeRune(pos.x+1, pos.y, '|', 99)
			}
		case "i":
			/* read inventory */
			clearscrn()
			fmt.Print("╔")
			for i := 0; i < width-2; i += 1 {
				fmt.Print("═")
			}
			fmt.Print("╗")
			fmt.Print("║")
			for i := 0; i < 33; i += 1 {
				fmt.Print(" ")
			}
			fmt.Print("║ Backpack ║")
			for i := 0; i < 33; i += 1 {
				fmt.Print(" ")
			}
			fmt.Print("║")
			fmt.Print("║")
			for i := 0; i < 33; i += 1 {
				fmt.Print(" ")
			}
			fmt.Print("╚")
			for i := 0; i < 10; i += 1 {
				fmt.Print("═")
			}
			fmt.Print("╝")
			for i := 0; i < 33; i += 1 {
				fmt.Print(" ")
			}
			fmt.Print("║")
			/* body of inventory */
			for i := 0; i < height-5; i += 1 {
				fmt.Print("║")
				clearln(2)
				fmt.Print("║")
			}
			/* end body of inventory */
			fmt.Print("╚")
			for i := 0; i < width-2; i += 1 {
				fmt.Print("═")
			}
			fmt.Print("╝")
			clearln(0)
			fmt.Scanln()
		case "D":
			/* debug mode */
			if debugmode == false {
				debugmode = true
			} else {
				debugmode = false
			}
		case "C":
			/* clear screen (hard reset?) */
			clearscrn()
			printRoom()
			printStats(usrin, message)
			continue
			/* set case "H": to be a help screen in spirit of Inventory, must add command to make
			   a screen box and whatnot, perhaps push ncurses-replacement library derived from this PoC? */
		case "<", ">":
			num:=1
			if char.fill == 'Ɵ' {
				if usrin == "<" {
					message = "Teleporting down!"
					num, _ =sc.Atoi(roomnum)
					if num - 1 > 0 {
						num-=1
						roomnum = sc.Itoa(num)
						setRoom(roomnum)
						first=true
						continue
					} else {
						message = "You are at the lowest level."
					}
				} else {
					message = "Teleporting up!"
					num, _=sc.Atoi(roomnum)
					if num < numrooms {
						num+=1
						roomnum = sc.Itoa(num)
						setRoom(roomnum)
						first=true
						continue
					} else {
						message = "You are at the highest level."
					}
				}
			} else {
				message = "You can't teleport here."
			}
		default:
		}

		/* perform movement of sprites and player */
		placeRune(pos.x, pos.y, char.fill, 99)
		char.fill, char.fillU, char.fillL, char.fillD, char.fillR = placeRune(fut.x, fut.y, char.icon, 99)
		/* creeps move every 4 turns */
		if creep_move_cnt == 3 {
			moveCreeps()
			creep_move_cnt = 0
		} else {
			creep_move_cnt += 1
		}

		/* print the map and other such things, perhaps make this its own function then goroutine it */
		printRoom()
		printStats(usrin, message)
		if reset_message == true {
			message = ""
		}
		pos.x, pos.y = fut.x, fut.y
	}

	fmt.Print("\n\nNEEIIII\n")
}
