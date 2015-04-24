// Package SVI for PNDTC
package svi

import (
"strings"
"fmt"
"math/rand"
"time"
//"bufio"
//"io"
"io/ioutil"
"os"
)


/* make sure to (at minimum, prompt for) create the the file if it is not found
-- this function should double as a hack-in Filecreator() function */
func Filewriter(filename string, texttowrite []string, directory string) {
	fmt.Print("Useless function at the moment, <WIP>.")
}


/* Should add an option for path to/directory of the file as per Filewriter() */
func Filereader(filename string)(lines []string, success int) {
	content, err := ioutil.ReadFile(filename)
	if err == nil {

		lines = strings.Split(string(content), "\n") // split at each line
		//newpstats = Pstats{lines[0], lines[1], lines[2], lines[3], lines[4], lines[5], lines[6]}

		success = 0 //success
	} else {
		fmt.Println("Error opening file (See: svi.go:Filereader() and/or golang's ioutil.ReadFile()")
		success = 1 //failure
		//var lines []string
		lines[0] = "!Error"
	}
	return
}


/* Do not call manually, callead by metaread only */
func filecheck(e error) {
	if e != nil {
		panic(e)
	}
}

// forgot I made this one, whoops
func metaread() {

	dat, err := ioutil.ReadFile("Instigator.meta")
	filecheck(err)
	fmt.Print(string(dat))

	buf := make([]byte, 1024)

	instigator, err := os.Open("Instigator.meta")
	filecheck(err)

	for {
		n, err := instigator.Read(buf)
		filecheck(err)
		if n == 0 {
			break
		}
	}

}


/* processes [y/n] dialogues with the yes/no user input
the values returned are the "y" or "n" and then a true/false for looping
the false value is stated false if the input is not a yes or no */

func YorN(prompt string) (newinputz string, tf bool) {
	inputz := "Mountain Goat"
	tf = false

	for tf != true {
		fmt.Printf("\n%v [y/n]: ", prompt)
		fmt.Scan(&inputz)

		if inputz = strings.ToLower(inputz); inputz == "yes" {
			newinputz = "y"
			tf = true
		} else if inputz == "no" {
			newinputz = "n"
			tf = true
		} else if inputz == "y" {
			tf = true
			newinputz = inputz
		} else if inputz == "n" {
			tf = true
			newinputz = inputz
		} else {
			newinputz = "Mountain Goat"
			tf = false
		}

	}
	return
}


func Random(min, max int)(newnum int) {
rand.Seed(time.Now().Unix())
newnum = rand.Intn(max - min) + min //non-inclusive on max
return
}
