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

/* replacement for old strings to handle direction, iota is special */
type direction int
const (
	UP direction = iota
	DOWN
	LEFT
	RIGHT
	NONE
)


/* item id's */
type ident int
const (
	EMPTY ident = iota
	UNKNOWN
	ROCK
	TORCH
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

/* item struct, stats then value */
type item struct {
	icon rune
	id ident
	fill rune
}

/* methods for items */
func (it item) getDesc()(desc string) {
	switch it.id {
		case ROCK: desc = "A rock. It's pimpin."
		case TORCH:
	}
	return
}

/* get "selling" value of item */
func (it item) getValue()(val int) {
	return
}


/* inventory struct for storage and tracking */
const max_inventory = 50
type inventory struct {
	num int
	size int /* must be `<=` slot's cap */
	slot [max_inventory]item
}

/* inventory methods */
/* remove item from inventory; num is the actual array position */
func (inv *inventory) remove(num int) {
	if num < inv.size {
		for i:=num;i<inv.size;i+=1 {
				inv.slot[i].icon = inv.slot[i+1].icon
				inv.slot[i].id = inv.slot[i+1].id
			}
	}
	inv.num = inv.num - 1
}

/* set item id and icon and fill */
func (inv *inventory) add(icon rune) {
	/* perhaps add an else and form of return which states no more space */
	if num:=inv.num; inv.num+1 < inv.size {
		inv.num+=1
		inv.slot[num].icon = icon
		switch icon {
			case '*':
				inv.slot[num].id = ROCK
			default:
				inv.slot[num].id = UNKNOWN
		}
	}
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
var backpack inventory
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

/* does the standard print routine used in main() without moving/interfering with anything */
func onlyPrint(usrin string) {
	clearscrn()
	printRoom()
	printStats(usrin, message)
}

/* prints the curroom[] buf to screen */
func printRoom() {
	if height > 24 {
		for i:=0;i<(height-24);i+=1 {
			clearln(0)
		}
	}
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
		dirX, dirY := make([]direction, numsprites), make([]direction, numsprites)
		bufX[i] = sprites[i].fut.x
		bufY[i] = sprites[i].fut.y
		bufX[numsprites] = fut.x
		bufY[numsprites] = fut.y //add the player's coords

		placeRune(sprites[i].pos.x, sprites[i].pos.y, sprites[i].f.fill, i)
		if sprites[i].fut.x > fut.x {
			dirX[i] = LEFT
		} else if sprites[i].fut.x < fut.x {
			dirX[i] = RIGHT
		} else {
			dirX[i] = NONE
		}
		if sprites[i].fut.y > fut.y {
			dirY[i] = UP
		} else if sprites[i].fut.y < fut.y {
			dirY[i] = DOWN
		} else {
			dirY[i] = NONE
		}

		/* check for icons next to each other via fills */
		for h := 0; h < numsprites; h += 1 {
			if sprites[i].f.fillU == sprites[h].f.icon && dirY[i] == UP {
				dirY[i] = DOWN
			}
			if sprites[i].f.fillD == sprites[h].f.icon && dirY[i] == DOWN {
				dirY[i] = UP
			}
			if sprites[i].f.fillL == sprites[h].f.icon && dirX[i] == LEFT {
				dirX[i] = RIGHT
			}
			if sprites[i].f.fillR == sprites[h].f.icon && dirX[i] == RIGHT {
				dirX[i] = LEFT
			}
		}

		/* nullify movement if check fails, character adjacent, or movement close */
		if sprites[i].fut.x-1 > 0 {
			if check(sprites[i].fut.x-1, sprites[i].fut.y, sprites[i].f.fillL) == true {
				if dirX[i] == LEFT {
					dirX[i] = NONE
				}
			}
		} else if sprites[i].fut.x-1 <= 0 {
			if dirX[i] == LEFT {
				dirX[i] = NONE
			}
			edgeL = true
		}
		if sprites[i].fut.x+1 < 79 {
			if check(sprites[i].fut.x+1, sprites[i].fut.y, sprites[i].f.fillR) == true {
				if dirX[i] == RIGHT {
					dirX[i] = NONE
				}
			}
		} else if sprites[i].fut.x+1 >= 79 {
			if dirX[i] == RIGHT {
				dirX[i] = NONE
			}
			edgeR = true
		}
		if sprites[i].fut.y-1 > 0 {
			if check(sprites[i].fut.x, sprites[i].fut.y-1, sprites[i].f.fillU) == true {
				if dirY[i] == UP {
					dirY[i] = NONE
				}
			}
		} else if sprites[i].fut.y-1 <= 0 {
			if dirY[i] == UP {
				dirY[i] = NONE
			}
			edgeU = true
		}
		if sprites[i].fut.y+1 < 23 {
			// y+1 = 24 -> not in curroom[], thus error, account for == 23...
			if check(sprites[i].fut.x, sprites[i].fut.y+1, sprites[i].f.fillD) == true {
				if dirY[i] == DOWN {
					dirY[i] = NONE
				}
			}
		} else if sprites[i].fut.y+1 >= 23 {
			if dirY[i] == DOWN {
				dirY[i] = NONE
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
				if char:=getChar(sprites[i].fut.x+1,sprites[i].fut.y-1); char == sprites[h].f.icon && dirY[i] == UP {
					dirY[i] = DOWN
				}
				if char:=getChar(sprites[i].fut.x+1,sprites[i].fut.y-1); char == sprites[h].f.icon && dirX[i] == RIGHT {
					dirX[i] = LEFT
				}
			}
			if edgeU == false && edgeL == false {
				if char:=getChar(sprites[i].fut.x-1,sprites[i].fut.y-1); char == sprites[h].f.icon && dirY[i] == UP {
					dirY[i] = DOWN
				}
				if char:=getChar(sprites[i].fut.x-1,sprites[i].fut.y-1); char == sprites[h].f.icon && dirX[i] == LEFT {
					dirX[i] = RIGHT
				}
			}
			if edgeD == false && edgeR == false {
				if char:=getChar(sprites[i].fut.x+1,sprites[i].fut.y+1); char == sprites[h].f.icon && dirY[i] == DOWN {
					dirY[i] = UP
				}
				if char:=getChar(sprites[i].fut.x+1,sprites[i].fut.y+1); char == sprites[h].f.icon && dirX[i] == RIGHT {
					dirX[i] = LEFT
				}
			}
			if edgeD == false && edgeL == false {
				if char:=getChar(sprites[i].fut.x-1,sprites[i].fut.y+1); char == sprites[h].f.icon && dirY[i] == DOWN {
					dirY[i] = UP
				}
				if char:=getChar(sprites[i].fut.x-1,sprites[i].fut.y+1); char == sprites[h].f.icon && dirX[i] == LEFT {
					dirX[i] = RIGHT
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
				if dirX[i] == RIGHT {
					dirX[i] = LEFT
				}
				if dirY[i] == DOWN {
					dirY[i] = UP
				}
			}
			if altX == botLx && altY == botLy {
				if dirX[i] == LEFT {
					dirX[i] = RIGHT
				}
				if dirY[i] == DOWN {
					dirY[i] = UP
				}
			}
			if altX == topRx && altY == topRy {
				if dirX[i] == RIGHT {
					dirX[i] = LEFT
				}
				if dirY[i] == UP {
					dirY[i] = DOWN
				}
			}
			if altX == topLx && altY == topLy {
				if dirX[i] == LEFT {
					dirX[i] = RIGHT
				}
				if dirY[i] == UP {
					dirY[i] = DOWN
				}
			}
			/* check for `under and over it~` 5fdp yo */
			if altY == y+1 && altX == x {
				if dirY[i] == DOWN {
					dirY[i] = UP
				}
			}
			if altY == y-1 && altX == x {
				if dirY[i] == UP {
					dirY[i] = DOWN
				}
			}
			if altY == y && altX == x+1 {
				if dirX[i] == RIGHT {
					dirX[i] = LEFT
				}
			}
			if altY == y && altX == x-1 {
				if dirX[i] == LEFT {
					dirX[i] = RIGHT
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
				//dirX[i] = NONE
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
				//dirY[i] = NONE
				yCanc = true
			}
		}
		if xCanc == true && yCanc == true {
			dirX[i] = NONE
			dirY[i] = NONE
		} else if xCanc == true && yCanc == false {
			dirY[i] = dirY[i]
			dirX[i] = NONE
		} else if xCanc == false && yCanc == true {
			dirX[i] = dirX[i]
			dirY[i] = NONE
		}



		/* can only be one, thus non-specific checks are okay here */
		if sprites[i].f.fillL == char.icon || sprites[i].f.fillR == char.icon {
			dirX[i] = NONE
		}
		if sprites[i].f.fillD == char.icon || sprites[i].f.fillU == char.icon {
			dirY[i] = NONE
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
			dirY[i] = NONE
		}
		if numX < 4 && numY < 4 {
			dirX[i] = NONE
		}


		/* check for edge of map (how did this not get implemented earlier?) */
		if sprites[i].f.fillU == '⚠' && dirY[i] == UP {
			dirY[i] = NONE
		}
		if sprites[i].f.fillD == '⚠' && dirY[i] == DOWN {
			dirY[i] = NONE
		}
		if sprites[i].f.fillR == '⚠' && dirX[i] == RIGHT {
			dirX[i] = NONE
		}
		if sprites[i].f.fillL == '⚠' && dirX[i] == LEFT {
			dirX[i] = NONE
		}

		/* Pick a direction */

		if rn := svi.Random(0, 2); dirY[i] != NONE && dirX[i] != NONE {
			if rn == 0 {
				dirY[i] = NONE
			} else if rn == 1 {
				dirX[i] = NONE
			}
		}

		/* Translate dirX[i] dirY[i] */

		if dirX[i] == LEFT {
			sprites[i].fut.x -= 1
			if sprites[i].fut.x < 0 {
				sprites[i].fut.x += 1
			}
		} else if dirX[i] == RIGHT {
			sprites[i].fut.x += 1
			if sprites[i].fut.x > 79 {
				sprites[i].fut.x -= 1
			}
		} else {
			sprites[i].fut.x = sprites[i].fut.x
		}
		if dirY[i] == UP {
			sprites[i].fut.y -= 1
			if sprites[i].fut.y < 0 {
				sprites[i].fut.y += 1
			}
		} else if dirY[i] == DOWN {
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

/* checks if the target location contains an interactable item or not */
func checkItem(x, y int)(occ bool) {
	str:=curroom[y]
	constructions := []rune{'Ɵ'}
	blocked := check(x, y, getChar(x,y))
	if blocked == true {
		occ=false
	} else {
		for i:=0;len(str) > 0;i+=1 {
			schar, size := utf8.DecodeRuneInString(str)
			str = str[size:]
			if i == x {
				for _, cchar := range constructions {
					if schar == cchar {
						occ=true
					}
				}
			}
		}
	}
	return
}

/* check for background shit */
func checkBack(x, y int)(occ bool) {
	str:=curroom[y]
	constructions := []rune{' ', '░'}
	blocked := check(x, y, getChar(x,y))
	if blocked == true {
		occ=false
	} else {
		for i:=0;len(str) > 0;i+=1 {
			schar, size := utf8.DecodeRuneInString(str)
			str = str[size:]
			if i == x {
				for _, cchar := range constructions {
					if schar == cchar {
						occ=true
					}
				}
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
	var reset_message, auto_resize, playing bool = true, false, true
	flag.StringVar(&roomnum, "room", "1", "Set the initial room to begin the game in")
	flag.StringVar(&icon_string, "icon", "♔", "Set unicode character to use as player icon")
	flag.BoolVar(&auto_resize, "comp", false, "Automatically compensate for larger terminals [24x80 min]")
	flag.IntVar(&height, "height", 24, "Set height of terminal screen [24]")
	flag.IntVar(&width, "width", 80, "Set width of terminal screen [80]")
	flag.Parse()

	oldState, _ := terminal.MakeRaw(0)
	defer terminal.Restore(0, oldState)
	if 	poss_w, poss_h, _ := terminal.GetSize(0); auto_resize == true {
		height, width = poss_h, poss_w
	}

	char.icon, _ = utf8.DecodeRuneInString(icon_string)
	char.fill = ' '
	var b []byte = make([]byte, 1)
	clearscrn()

	/* set some initial values, get the number of room/map files we'll be playing with */
	dir, _ := ioutil.ReadDir("./rooms")
	numrooms = len(dir)
	setRoom(roomnum)
	plyr.hlth, plyr.atk, plyr.dfs = 10, 02, 02
	backpack.num, backpack.size, backpack.slot[0].icon, backpack.slot[0].id = 1, 12, '*', ROCK

	/* begin game loop */
	var first bool = true
	var creep_move_cnt int = 0
	for playing == true {
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
			case "q":
				message = "Really quit?: "
				onlyPrint(usrin)
				os.Stdin.Read(b)
				tmpwords := string([]byte(b)[0])
				if tmpwords == "q" {
					playing=false
				} else {
					message = "Not like you had anything better to do!"
					onlyPrint(tmpwords)
				}
				continue
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
				h:=0
				for i := 0; i < height-5; i += 1 {
					if h < backpack.num {
						// get description based on const?
						str:=backpack.slot[h].getDesc()
						size:=utf8.RuneCountInString(str)
						fmt.Print("║")
						fmt.Printf("%2d) '%c': \"%s\"", h+1, backpack.slot[h].icon, str)
						//fmt.Print(str)
						clearln(size+2+11)
						fmt.Print("║")
						h+=1
					} else {
						fmt.Print("║")
						clearln(2)
						fmt.Print("║")
					}
				}
				/* end body of inventory */
				fmt.Print("╚")
				for i := 0; i < width-2; i += 1 {
					fmt.Print("═")
				}
				fmt.Print("╝")
				clearln(0)
				fmt.Scanln()
			case "p":
				/* pickup command */
				message = "Direction to pick up from?: "
				onlyPrint(usrin)
				os.Stdin.Read(b)
				tmpwords := string([]byte(b)[0])
				x,y:=fut.x,fut.y
				nomove:=false
				switch tmpwords {
					case "w":
						y-=1
					case "a":
						x-=1
					case "s":
						y+=1
					case "d":
						x+=1
					default:
						nomove=true
				}
				if check(x, y, getChar(x,y)) == false && checkItem(x,y) == false && checkBack(x,y) == false && nomove == false {
					backpack.add(getChar(x,y))
					placeRune(x,y,' ',99)
					message="Picked up item!"
					onlyPrint(tmpwords)
				} else {
					message="Nothing to pick up!"
					onlyPrint(tmpwords)
				}
			case "P":
				/* put command (from inventory) */
					message = "Put which item?: "
					onlyPrint(usrin)
					os.Stdin.Read(b)
					tmpwords := string([]byte(b)[0])
					tmpnum, err := sc.Atoi(tmpwords)
					if err == nil {
						//tmpnum -= 1
						//fmt.Print(tmpnum)
						if tmpnum > backpack.num {
							message="You don't have that many items."
							onlyPrint(tmpwords)
						} else if tmpnum < 0 {
							message="Negative items don't exist."
							onlyPrint(tmpwords)
						} else {
							if backpack.num > 0 {
								message="Direction to place?: "
								onlyPrint(tmpwords)
								os.Stdin.Read(b)
								tmpwords = string([]byte(b)[0])
								x,y:=fut.x,fut.y
								nomove:=false
								switch tmpwords {
									case "w":
										y-=1
									case "a":
										x-=1
									case "s":
										y+=1
									case "d":
										x+=1
									default:
										nomove=true
								}
								/* should also probably save a fill of what the item's rune was placed over */
								if (check(x, y, getChar(x,y)) == false) && (checkItem(x,y) == false) && (nomove == false) {
									message="Item placed!"
									placeRune(x,y,backpack.slot[tmpnum-1].icon,99)
									(&backpack).remove(tmpnum-1)
									onlyPrint(tmpwords)
								} else {
									message="Can't place there!"
									onlyPrint(tmpwords)
								}
							} else {
								message="No items in your inventory!"
								onlyPrint(tmpwords)
							}
						}
					} else {
						message="That's not a number!"
						onlyPrint(tmpwords)
					}
					continue
			case "D":
				/* debug mode */
				if debugmode == false {
					debugmode = true
				} else {
					debugmode = false
				}
			case "C":
				/* clear screen (hard reset?) */
				onlyPrint(usrin)
				continue
			case "H":
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
