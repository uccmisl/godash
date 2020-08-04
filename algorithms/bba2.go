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
	glob "github.com/uccmisl/godash/global"
	"github.com/uccmisl/godash/http"
	"github.com/uccmisl/godash/utils"
)

// CalculateSelectedIndexBba :
/*
 * return the index of the segment which should be selected
 *
 */
func CalculateSelectedIndexBba(newThr int, lastDuration int, lastIndex int, maxBufferLevel int,
	lastRate int, thrList *[]int, mpdDuration int, currentMPD http.MPD, currentURL string,
	currentMPDRepAdaptSet int, segmentNumber int, baseURL string, debugLog bool, downloadTime int, bufferLevel int,
	highestMPDrepRateIndex int, lowestMPDrepRateIndex int, bandwithList []int, quicBool bool, useTestbedBool bool) int {

	*thrList = append(*thrList, newThr)

	//average
	var average float64

	//average of throughtputs
	meanAverage(*thrList, &average)

	//fmt.Println("average of throughtput: ", average)

	//set the rate to the lowest one
	//we need an index of the lowest value
	//the bandwithlist contains all of the available video bitrates, ordering depends on repRatesReversed
	//if repRatesReversed then from smallest to largest, otherwise largest to smallest
	var lowestBitrateIndex int

	/*
		if repRatesReversed {
			lowestBitrateIndex = 0

		} else {
	*/
	lowestBitrateIndex = len(bandwithList) - 1

	//}
	//fmt.Println("bandwtihlist: ", bandwithList)

	//lowest := LowestBitrate(*thrList)
	//fmt.Println("lowest bitrate...: ", lowest)

	reservoir := bba1UpdateReservoir(lastRate, lastIndex, mpdDuration, lastDuration, maxBufferLevel,
		currentMPD, currentURL, currentMPDRepAdaptSet, segmentNumber, baseURL, debugLog, bandwithList, quicBool, useTestbedBool)

	//fmt.Println("reservoir: ", reservoir)
	//fmt.Println("ret1:", retVal)
	//Not sure that implementation of mStaticAlgPar is correct, I think the scope could be wrong,
	//as the last 2 statements have ineff assign...
	var mStaticAlgPar int
	var retVal int
	//fmt.Println("m1:", mStaticAlgPar)
	//fmt.Println("download time", downloadTime)
	//if it took more time to download then the length of the last segment
	if downloadTime > lastDuration {
		mStaticAlgPar = 1
	}
	//fmt.Println("bufferLevel: ", bufferLevel)

	//fmt.Println("lastindex", lastIndex)
	//fmt.Println("lastrate", lastRate)
	//fmt.Println("res", reservoir)
	//fmt.Println(downloadTime, lastDuration)
	if bufferLevel < reservoir {
		//fmt.Println("buff less than")
		if mStaticAlgPar != 0 {
			retVal = lowestBitrateIndex
			//fmt.Println("retval less than", retVal)
			//retVal = LowestBitrate(*thrList)
		} else {
			if float64(downloadTime) < 0.125*float64(lastDuration) {
				//increase next segment quality
				/*
					if repRatesReversed {
						retVal = utils.Min(lastRate+1, (len(bandwithList) - 1))
					} else {
				*/
				//!repRatesReversed case
				retVal = utils.Max(lastRate-1, 0)
				//}

			} else {
				//fmt.Println("retval 3", retVal)
				retVal = lastRate

			}

		}

	} else {
		if mStaticAlgPar != 0 {
			retVal = bba1VRAA(lastRate, *thrList, bufferLevel, highestMPDrepRateIndex,
				lowestMPDrepRateIndex, maxBufferLevel, bandwithList, reservoir)
			//fmt.Println("retval 5", retVal)
		} else {
			bba1RateIndex := bba1VRAA(lastRate, *thrList, bufferLevel, highestMPDrepRateIndex,
				lowestMPDrepRateIndex, maxBufferLevel, bandwithList, reservoir)

			///fmt.Println("bba1rateindex", bba1RateIndex)
			if float64(downloadTime) <= 0.5*float64(lastDuration) {
				//increase segment quality
				/*
					if repRatesReversed {
						lowestBitrateIndex = utils.Min(lastRate+1, len(bandwithList)-1)
					} else {
				*/
				lowestBitrateIndex = utils.Max(lastRate-1, 0)

				//}
			}
			/*
				if repRatesReversed {
					retVal = utils.Max(bba1RateIndex, lowestBitrateIndex)
				} else {
			*/
			//originally here
			retVal = utils.Min(bba1RateIndex, lowestBitrateIndex)

			//}

			/*
				if repRatesReversed {
					if bba1RateIndex > lowestBitrateIndex {
						mStaticAlgPar = 1
					}
				} else {
			*/
			if bba1RateIndex < int(lowestBitrateIndex) {
				//fmt.Println("bbahigher!!!")
				mStaticAlgPar = 1
			}
			//}

		}
	}
	//fmt.Println("lastindex: ", lastIndex)
	//fmt.Println("retval:", retVal)
	//fmt.Println("m2", mStaticAlgPar)
	return retVal
}

func bba1UpdateReservoir(lastRate int, lastRateIndex int, mpdDuration int,
	lastSegmentDuration int, maxBufferLevel int, currentMPD http.MPD, currentURL string, currentMPDRepAdaptSet int,
	segmentNumber int, baseURL string, debugLog bool, bandwithList []int, quicBool bool, useTestbedBool bool) int {
	//we need to convert the maxBufferLevel to milliseconds
	//otherwise the comparison is between seconds and milliseconds

	resvWin := utils.Min(2*maxBufferLevel*1000/lastSegmentDuration, (mpdDuration/lastSegmentDuration)-lastRateIndex)
	//fmt.Println("mpddur:", mpdDuration, "lastrateindex:", lastRateIndex)
	//fmt.Println("resvWin: ", resvWin)

	avgSegSize := (int(bandwithList[lastRate]/glob.Conversion1000) * lastSegmentDuration) / 8000
	//fmt.Println("lastrate", lastRate, "lastsegduration" ,lastSegmentDuration)
	//fmt.Println("averageSegSize: ", avgSegSize)

	largeSeg := 0
	smallSeg := 0

	// fmt.Println("contlenght", http.GetContentLengthHeader(currentMPD,
	// 	currentURL, currentMPDRepAdaptSet, lastRate, segmentNumber+1, baseURL, debugLog))

	_, client, _ := http.GetHTTPClient(quicBool, glob.DebugFile, debugLog, useTestbedBool)

	for i := 0; i < resvWin; i++ {
		//do a func getSegBySize(lastSegNumber+i, lastRateIndex) and return the size of the segment
		if http.GetContentLengthHeader(currentMPD,
			currentURL, currentMPDRepAdaptSet, lastRate, segmentNumber+i, baseURL, debugLog, client) > avgSegSize {

			largeSeg += http.GetContentLengthHeader(currentMPD,
				currentURL, currentMPDRepAdaptSet, lastRate, segmentNumber+i, baseURL, debugLog, client)
		} else {
			smallSeg += http.GetContentLengthHeader(currentMPD,
				currentURL, currentMPDRepAdaptSet, lastRate, segmentNumber+i, baseURL, debugLog, client)
		}
	}

	//fmt.Println("large", largeSeg, "small", smallSeg, "lastrate", lastRate)

	//fmt.Println(lastSegmentDuration, maxBufferLevel)
	//fmt.Println("actual representation rate", bandwithList[lastRate]/glob.Conversion1000)

	reservoir := 8 * float64(largeSeg-smallSeg) / float64(bandwithList[lastRate]/glob.Conversion1000)

	//fmt.Println("res1", reservoir)

	if reservoir < float64(2*lastSegmentDuration) {
		reservoir = 2 * float64(lastSegmentDuration)
	} else {
		if reservoir > 0.6*float64(maxBufferLevel*1000) {
			reservoir = 0.6 * float64(maxBufferLevel*1000)
		}
	}
	//fmt.Println("reservoir after calc::", reservoir)
	return int(reservoir)
}

func bba1VRAA(lastRateIndex int, thrList []int, bufferLevel int, highestMPDrepRateIndex int,
	lowestMPDrepRateIndex int, maxBufferLevel int, bandwithList []int, reservoir int) int {

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

	if float64(bufferLevel/1000) >= float64(0.9)*float64(maxBufferLevel) {
		optRateIndex = 0
	} else {
		low := LowestBitrate(bandwithList) / 1024
		high := HighestBitrate(bandwithList) / 1024

		//target rate calcul
		targetRate := (high - low) / (0.9*float64(maxBufferLevel) - float64(reservoir))

		//optRateIndex = findBestRateIndex(targetRate * 1024)
		optRateIndex = SelectRepRateWithThroughtput(int((targetRate+1)*1000), bandwithList, lowestMPDrepRateIndex)
	}

	/*
		if !repRatesReversed {
	*/
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

	/*
		}
		//else if repRateReversed
		if optRateIndex == 0 {
			//fmt.Println("repRateReversed first rate is the highest")
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
