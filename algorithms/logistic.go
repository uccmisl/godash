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

package algorithms

import (
	"strconv"

	"github.com/uccmisl/godash/logging"
	"github.com/uccmisl/godash/utils"

	"math"
)

var debugFile string
var debugLog bool

// Logistic :
//add the last throughtput to the list and call CalculateSelectedIndex,
//return the rate and throughtput list
func Logistic(thrList *[]int, newThr int, repRate *int, bandwithList []int, bufferLevel int,
	highestMPDrepRateIndex int, lowestMPDrepRateIndex int, debugFiles string, debugLogs bool,
	maxBufferLevel int) {

	debugFile = debugFiles
	debugLog = debugLogs

	*thrList = append(*thrList, newThr)

	*repRate = calculateSelectedIndex(*thrList, newThr, bandwithList, bufferLevel, *repRate, highestMPDrepRateIndex,
		lowestMPDrepRateIndex, maxBufferLevel)

}

//----------------------------------------------------------------------------------------------------------
//----------------------------------------------------------------------------------------------------------

// calculateSelectedIndex :
//call the func LogisticFunction(lastRateIndex, thrList, bufferLevel) to calculate the rate
func calculateSelectedIndex(thrList []int, newThr int, bandwithList []int, bufferLevel int, repRate int,
	highestMPDrepRateIndex int, lowestMPDrepRateIndex int, maxBufferLevel int) int {

	//take the last rate
	//current rep rate : repRate

	//find the index of the last rate ?
	lastRateIndex := repRate

	retVal := LogisticFunction(lastRateIndex, thrList, bufferLevel, highestMPDrepRateIndex, lowestMPDrepRateIndex,
		maxBufferLevel, bandwithList)
	//fmt.Println(retVal)
	return retVal
}

//----------------------------------------------------------------------------------------------------------
//----------------------------------------------------------------------------------------------------------

// LogisticFunction :
//calculate and return the rate index
func LogisticFunction(lastRateIndex int, thrList []int, bufferLevel int, highestMPDrepRateIndex int,
	lowestMPDrepRateIndex int, maxBufferLevel int, bandwithList []int) int {

	//len(tracks) = number of rates -> of representations in the MPD

	var optRateIndex int
	var rateUindex int
	var rateLindex int

	/*
		if repRatesReversed {
			//min
			rateUindex = utils.Min(lastRateIndex+1, highestMPDrepRateIndex)
			//max
			rateLindex = utils.Max(lastRateIndex-1, lowestMPDrepRateIndex)
		} else {
	*/
	//min
	rateUindex = utils.Max(lastRateIndex-1, highestMPDrepRateIndex)
	//max
	rateLindex = utils.Min(lastRateIndex+1, lowestMPDrepRateIndex)
	//}

	//log to file
	logging.DebugPrint(debugFile, debugLog, "DEBUG: ", "highestMPDrateindex: "+strconv.Itoa(highestMPDrepRateIndex))
	logging.DebugPrint(debugFile, debugLog, "DEBUG: ", "lowest: "+strconv.Itoa(lowestMPDrepRateIndex))
	logging.DebugPrint(debugFile, debugLog, "DEBUG: ", "rateUindex: "+strconv.Itoa(rateUindex))
	logging.DebugPrint(debugFile, debugLog, "\nDEBUG: ", "rateLindex: "+strconv.Itoa(rateLindex))
	logging.DebugPrint(debugFile, debugLog, "\nDEBUG: ", "bufferLevel: "+strconv.Itoa(bufferLevel/100))
	logging.DebugPrint(debugFile, debugLog, "\nDEBUG: ", "max: "+strconv.Itoa(maxBufferLevel))

	if float64(bufferLevel/1000) >= float64(0.97)*float64(maxBufferLevel) {
		optRateIndex = 0
	} else {
		var low float64
		low = LowestBitrate(bandwithList) / 1024
		var high float64
		high = HighestBitrate(bandwithList) / 1024

		//target rate calcul
		targetRate := (math.Ceil(high / (1 + ((high/(low) - 1) * math.Exp(-0.05*float64(bufferLevel)/1000)))))

		logging.DebugPrint(debugFile, debugLog, "\nDEBUG: ", "(high/(low) - 1): "+utils.FloatToString((high/(low)-1)))
		logging.DebugPrint(debugFile, debugLog, "\nDEBUG: ", "bufferLevel: "+strconv.Itoa(int(bufferLevel)))
		logging.DebugPrint(debugFile, debugLog, "\nDEBUG: ", "Lower: "+strconv.Itoa(int(low)))
		logging.DebugPrint(debugFile, debugLog, "\nDEBUG: ", "Higher: "+strconv.Itoa(int(high)))
		logging.DebugPrint(debugFile, debugLog, "\nDEBUG: ", "targetRate: "+strconv.Itoa(int(targetRate)))

		//optRateIndex = findBestRateIndex(targetRate * 1024)
		optRateIndex = SelectRepRateWithThroughtput(int((targetRate+1)*1000), bandwithList, lowestMPDrepRateIndex)
	}

	logging.DebugPrint(debugFile, debugLog, "\nDEBUG: ", "optRateIndex: "+strconv.Itoa(optRateIndex))
	logging.DebugPrint(debugFile, debugLog, "\nDEBUG: ", "lastRateIndex: "+strconv.Itoa(lastRateIndex))

	/*
		if !repRatesReversed {
	*/
	//first rate is the lowest
	if optRateIndex == 0 {
		if lastRateIndex <= lowestMPDrepRateIndex {
			return lowestMPDrepRateIndex
		}
		return lastRateIndex - 1
	} else if optRateIndex <= rateLindex {
		return optRateIndex
	} else if optRateIndex >= rateUindex {
		return optRateIndex - 1
	}
	return optRateIndex

	//}
	/*
		//else if repRateReversed
		if optRateIndex == 0 {
			fmt.Println("repRateReversed first rate is the highest")
			if lastRateIndex <= lowestMPDrepRateIndex {
				return lastRateIndex - 1
			}
			return lowestMPDrepRateIndex
		} else if optRateIndex >= rateUindex {
			return optRateIndex
		} else if optRateIndex <= rateLindex {
			return optRateIndex + 1
		}
		return optRateIndex
	*/

}

//----------------------------------------------------------------------------------------------------------
//----------------------------------------------------------------------------------------------------------

// LowestBitrate :
//return the lowest bitrate of the list
func LowestBitrate(thrList []int) float64 {
	var lowest int
	lowest = thrList[len(thrList)-1]
	for i := 0; i < len(thrList); i++ {
		if thrList[i] < lowest {
			lowest = thrList[i]
		}
	}

	return float64(lowest)
}

// HighestBitrate :
//return the highest bitrate of the list
func HighestBitrate(thrList []int) float64 {
	var highest int
	highest = thrList[0]
	for i := 0; i < len(thrList); i++ {
		if thrList[i] > highest {
			highest = thrList[i]
		}
	}
	return float64(highest)
}
