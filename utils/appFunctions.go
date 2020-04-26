/*
 *	goDASH, golang client emulator for DASH video streaming
 *	Copyright (c) 2019, Jason Quinlan, Darijo Raca, University College Cork
 *											[j.quinlan,d.raca]@cs.ucc.ie)
 *                      Maëlle Manifacier, MISL Summer of Code 2019, UCC
 *	This program is free software; you can redistribute it and/or
 *	modify it under the terms of the GNU General Public License
 *	as published by the Free Software Foundation; either version 2
 *	of the License, or (at your option) any later version.
 *
 *	This program is distributed in the hope that it will be useful,
 *	but WITHOUT ANY WARRANTY; without even the implied warranty of
 *	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *	GNU General Public License for more details.
 *
 *	You should have received a copy of the GNU General Public License
 *	along with this program; if not, write to the Free Software
 *	Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA
 *	02110-1301, USA.
 */

package utils

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// StopApp : Print out flag usage and Stop the application
func StopApp() {
	// print the flag help output
	flag.Usage()
	// exit the application
	os.Exit(3)
}

// IsFlagSet :
// * Determine if a Flag has been passed to the applicaiton
// * return true if the flag has been called as a parameter to the app
func IsFlagSet(name string) (passed bool) {

	passed = false

	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			passed = true
		}
	})
	return
}

// RecoverPanic :
// * Called after go panic, when a parameter is not found in the main function
func RecoverPanic() {
	if err := recover(); err != nil {
		fmt.Println("Panic error recovered : ", err)
		fmt.Println("Parameters needed :\n Path to the file needed after '-config' : for example : ./goDASH -config ../config/config\n or URL missing after -url.\n To have more informations, execute ./goDASH -help")
	}
}

// WriteFile :
// * fileLocation string - pass in fileLocation
// * Setup the debug log file
func WriteFile(fileLocation string) {

	// create the debug log file
	f, err := os.Create(fileLocation)
	if err != nil {
		// print an error
		log.Println(err)
		// print the flag help output
		flag.Usage()
		// exit the application
		os.Exit(3)
	}
	f.Close()
}

// FindInStringArray :
// * return true if an array contains a string value
// * return index for this string value (-1 if value not found)
func FindInStringArray(array []string, item string) (bool, int) {

	// the input must be a value
	for index, value := range array {
		if item == value {
			return true, index
		}
	}
	return false, -1
}

// FindInIntArray :
// * return true if an array contains a string value
// * return index for this string value (-1 if value not found)
func FindInIntArray(array []int, item int) (bool, int) {

	// the input must be a defined value
	for index, value := range array {
		if item == value {
			return true, index
		}
	}
	return false, -1
}

// CheckStringVal : assign a value to a string if the assigning string is not empty
func CheckStringVal(p *string, name *string) {

	//fmt.Println(*p)
	//fmt.Println(*name)
	//fmt.Println()

	// if the value is not empty, then save to the original variable pointer
	if *p != "" {
		*name = *p
	}
}

// CheckIntVal : assign a value to a int if the assigning int is non zero
func CheckIntVal(p *int, name *int) {

	//fmt.Println(*p)
	//fmt.Println(*name)
	//fmt.Println()

	// if the value is not empty, then save to the original variable pointer
	if *p != 0 {
		*name = *p
	}
}

// CheckFloatVal : assign a value to a float if the assigning float is non zero
func CheckFloatVal(p *float64, name *float64) {

	//fmt.Println(*p)
	//fmt.Println(*name)
	//fmt.Println()

	// if the value is not empty, then save to the original variable pointer
	if *p != 0.0 {
		*name = *p
	}
}
