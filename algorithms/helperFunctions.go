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
	"math"

	glob "github.com/uccmisl/godash/global"
	"github.com/uccmisl/godash/http"
	"github.com/uccmisl/godash/utils"
)

// CalculateThroughtput :
// input segment size and time, returns the throughtput
func CalculateThroughtput(segmentSize, time int) int {
	return int(float64(segmentSize) / (float64(time) / glob.Conversion1000))
}

// SelectRepRateWithThroughtput :
/*
 * Select the rate the nearest just below the throughtput
 * return the selected rate
 */
func SelectRepRateWithThroughtput(thr int, bandwithList []int, lowestMPDrepRateIndex int) int {
	var selectedBandwith = 0

	//We'll use the index of the list to have the repRate
	i := 0
	repRate := lowestMPDrepRateIndex

	// select the throughtput with the list of bandwith from the MPD file
	for i = 0; i < len(bandwithList)-1; i++ {
		//fmt.Println(selectedBandwith)
		//fmt.Println(thr)
		//fmt.Println(bandwithList[i])
		//fmt.Println("thr opti ", thr)

		/*
			if bandwithList[i] < thr {
				repRate = i
				//fmt.Println("sel index: ", repRate)
				break
			}
		*/
		/*
			if repRatesReversed {
				if bandwithList[i] > thr {
					//&& selectedBandwith < bandwithList[i] {
					//selectedBandwith = bandwithList[i]

					// select the repRate with the selectedBandwith
					repRate = i
					// no need to loop through the rest of the for loop
					break
				}
				//If the throughtput is too bad, then we select the last bandwith
				if i == len(bandwithList)-1 && selectedBandwith == 0 {
					selectedBandwith = bandwithList[len(bandwithList)-1]

					// select the repRate with the selectedBandwith
					repRate = i
				}
			} else {
		*/
		if bandwithList[i] < thr {
			//&& selectedBandwith < bandwithList[i] {
			//selectedBandwith = bandwithList[i]

			// select the repRate with the selectedBandwith
			repRate = i
			// no need to loop through the rest of the for loop
			break
		}
		//If the throughtput is too bad, then we select the last bandwith
		if i == len(bandwithList)-1 && selectedBandwith == 0 {
			selectedBandwith = bandwithList[len(bandwithList)-1]

			// select the repRate with the selectedBandwith
			repRate = len(bandwithList) - 1
		}
		//}

	}

	/*
		fmt.Println(bandwithList)
		fmt.Println("THR: ", thr, "Selected bandwith: ", selectedBandwith, "repRate: ", repRate)
		os.Exit(1)
	*/

	return repRate

}

// ThroughputSamples :
/*
 * return the last "window" number of throughtputs
 */
func ThroughputSamples(window int, thrList []int) []int {
	workingWindow := utils.Min(window, len(thrList))
	//get last "window" number of throughputSamples
	last := len(thrList)
	first := last - workingWindow

	//var throughputSamples [window]int
	throughputSamples := make([]int, window)

	for i := first; i < last; i++ {
		j := i - last + workingWindow
		throughputSamples[j] = thrList[i]
	}

	return throughputSamples
}

// SmartConvHelper :
/*
 * Checks next "videoWindow" of segments and makes sure the average rate is less than the estimated rate
 */
func SmartConvHelper(qIndex int, videoWindow int, estRate float64, currentMPD http.MPD, currentURL string, currentMPDRepAdaptSet int, lastRate int, segmentNumber int, baseURL string, debugLog bool, lastDuration int) bool {
	var totSegSize int

	for i := 0; i < videoWindow; i++ {

		totSegSize += 8 * http.GetContentLengthHeader(currentMPD,
			currentURL, currentMPDRepAdaptSet, qIndex, segmentNumber+i, baseURL, debugLog)

	}
	actualAvgRate := float64(float64(totSegSize) / (float64(lastDuration) / 1000 * float64(videoWindow)))

	return actualAvgRate <= estRate

}

// SmartConvHelperFromFile :
/*
 * Checks next "videoWindow" of segments and makes sure the average rate is less than the estimated rate
 */
func SmartConvHelperFromFile(videoWindow int, estRate float64, qRate int, segmentNumber int, lastDuration int) bool {
	var totSegSize int

	for i := 0; i < videoWindow; i++ {

		totSegSize += http.SegHeadValues[0][qRate][segmentNumber+i] * 8

	}
	actualAvgRate := float64(float64(totSegSize) / (float64(lastDuration) / 1000.0 * float64(videoWindow)))

	return actualAvgRate <= estRate

}

// FloatMin :
/*
 * return the minimum of two float numbers
 */
func FloatMin(a, b float64) float64 {
	if a < b {
		return a
	}
	return b

}

// Harmonic Mean :
/*
 * return the exponential average of the last "window" number of throughtputs
 */
func harmonicAverage(num int, thrList []int, average *float64) {

	// slice the last num of values from throughput list if possible
	// if not take the full list (cases when length of the list is smaller than window)

	if len(thrList) < num {
		thrList = thrList[len(thrList)-len(thrList):]

	} else {
		thrList = thrList[len(thrList)-num:]

	}

	//harmonic average
	*average = 0.0

	// average calcul
	for k := 0; k < len(thrList); k++ {
		*average += float64(1.0 / float64(thrList[k]))
	}
	*average = float64(len(thrList)) / float64(*average)

	//fmt.Println("avg", *average)

}

// ExpAverage :
/*
 * return the exponential average of the last "window" number of throughtputs
 */
func ExpAverage(thrList []int, ratio float64, window int, average *float64) {

	// list for the number of throughtput we are asking from thrList
	var lastThrList []int

	var j = 0
	// to have the LAST number of the throughtput list
	// we begin by the end of the tab
	for i := len(thrList) - 1; i >= 0; i-- {
		if j == window {
			break
		} else {
			lastThrList = append(lastThrList, thrList[i])
		}
		j++
	}

	// average calcul

	var thisWeight float64
	weightSum := (1 - math.Pow(float64(1-ratio), float64(len(lastThrList))))
	*average = 0.0

	for i := 0; i < len(lastThrList); i++ {
		thisWeight = ratio * math.Pow(float64(1-ratio), float64(len(lastThrList)-1-i)) / weightSum
		*average += thisWeight * float64(lastThrList[i])
	}

}

// ExponentialAverage :
/*
 * return the exponential average of the throughtputs
 */
/*
	func ExponentialAverage(thrList []int, ratio float64) float64 {
	  weightsum := (1 - math.Pow(1-ratio, float64(len(thrList))))
	  var subTotal float64

	  for i := 0; i < len(thrList);  i++ {

	    thisWeight := ratio * math.Pow(1-ratio, float64(len(thrList)-1-i))/weightsum
	    subTotal += thisWeight * float64(thrList[i])
	  }

	  return subTotal
	}

	// SampleExponentialAverage :
	/*
	 * return the exponential average of the last "window" number of throughtputs
*/
/*
	func SampleExponentialAverage(window int, exponentialAverageRatio float64, thrList []int) float64 {
	  thrSamples := ThroughputSamples(window, thrList)
	  //fmt.Println(thrSamples)

	  return ExponentialAverage(thrSamples, exponentialAverageRatio)
	}
*/
