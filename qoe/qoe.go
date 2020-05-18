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

package qoe

import (
	glob "github.com/uccmisl/godash/global"
	"github.com/uccmisl/godash/logging"
)

// CreateQoE : get the P1203 and clae QoE values
func CreateQoE(log *map[int]logging.SegPrintLogInformation, debugLog bool, initBuffer int, maxRepRate int, printHeadersData map[string]string, saveFilesBool bool) {

	// *log does not support indexing :(
	logMap := *log

	// the P1203 standard only works for H264 (encoder) and up to resolutions of 1920x1080
	// so make sure the received segments are compliant
	P1023Print := checkInputHeader(printHeadersData, glob.P1203Header)
	ClaePrint := checkInputHeader(printHeadersData, glob.ClaeHeader)
	DuanmuPrint := checkInputHeader(printHeadersData, glob.DuanmuHeader)
	YinPrint := checkInputHeader(printHeadersData, glob.YinHeader)
	YuPrint := checkInputHeader(printHeadersData, glob.YuHeader)

	// create channels, so the output is in the right order
	var P1023Results chan float64
	P1023Results = make(chan float64)
	var stopP1203 = true
	if P1023Print {
		logging.DebugPrint(glob.DebugFile, debugLog, "\nDEBUG: ", "checking for P1203 compatibility")
		for a := 1; a <= len(*log); a++ {

			// get the encoder and resolution
			width := logMap[a].RepWidth
			height := logMap[a].RepHeight
			codec := logMap[a].RepCodec

			if codec != glob.RepRateCodecAVC || width > glob.P1203maxWidth || height > glob.P1203maxHeight {
				logging.DebugPrint(glob.DebugFile, debugLog, "\nDEBUG: ", "Downloaded segments are not P1203 compliant")
				stopP1203 = false
				// return
			}
		}
		// create the P1203 value
		if stopP1203 {
			go createP1203(log, P1023Results, saveFilesBool)
		}
	}

	var claeResults chan float64
	claeResults = make(chan float64)
	if ClaePrint {
		// create the Claye value
		go getClaye(*log, claeResults, maxRepRate, false)
	}

	var duanmuResults chan float64
	duanmuResults = make(chan float64)
	if DuanmuPrint {
		// create the Duanmu value
		go getDuanmu(*log, duanmuResults, initBuffer, false)
	}

	var yinResults chan float64
	yinResults = make(chan float64)
	if YinPrint {
		// create the Yin value
		go getYin(*log, yinResults, initBuffer, false)
	}

	var yuResults chan float64
	yuResults = make(chan float64)
	if YuPrint {
		// create the Yu value
		go getYu(*log, yuResults, false)
	}

	// // create the P1203 value
	// go createP1203(log, P1023Results)
	//
	// // create the Claye value
	// go getClaye(*log, claeResults, maxRepRate, false)
	//
	// // create the Duanmu value
	// go getDuanmu(*log, duanmuResults, initBuffer, false)
	//
	// // create the Yin value
	// go getYin(*log, yinResults, initBuffer, false)
	//
	// // create the Yu value
	// go getYu(*log, yuResults, false)

	// create a local copy of the log and allocate the QoE values
	locallogMap := *log
	locallog := locallogMap[len(locallogMap)]
	// calculate the P1203, Claye, Duanmu, Yin and Yu values and
	// save to the last log as 3 decimal floats
	if YinPrint {
		locallog.Yin = <-yinResults
	} else {
		locallog.Yin = 0.0
	}

	if YuPrint {
		locallog.Yu = <-yuResults
	} else {
		locallog.Yu = 0.0
	}

	if DuanmuPrint {
		locallog.Duanmu = <-duanmuResults
	} else {
		locallog.Duanmu = 0.0
	}

	if ClaePrint {
		locallog.Clae = <-claeResults
	} else {
		locallog.Clae = 0.0
	}

	if P1023Print && stopP1203 {
		locallog.P1203 = <-P1023Results
	} else {
		locallog.P1203 = 0.0
	}

	locallogMap[len(locallogMap)] = locallog
	*log = locallogMap
}

func checkInputHeader(printHeadersData map[string]string, key string) bool {

	/*
		fmt.Println("printHeadersData: ", printHeadersData)
		fmt.Println("key: ", key)
		fmt.Println("extendPrintString: ", *extendPrintString)
		fmt.Println("stringDuration: ", stringDuration)
		fmt.Println("val1: ", *val1)
		fmt.Println("val2: ", val2)
	*/

	// if the map has this key
	if val, ok := printHeadersData[key]; ok {
		if val == "on" || val == "On" || val == "ON" {
			return true
		}
	}
	return false
}
