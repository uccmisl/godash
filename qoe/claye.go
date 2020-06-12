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
	//"fmt"

	"fmt"
	"math"

	"github.com/uccmisl/godash/logging"
	"github.com/uccmisl/godash/utils"

	"gonum.org/v1/gonum/stat"
)

// getClaye : claye qoe estimation
func getClaye(log map[int]logging.SegPrintLogInformation, c chan float64, maxRepRate int, printOutput bool) {

	// I don't think we need *totVar* and *stallFree*, so I've removed these
	// I also moved *stallDurElem* to inside the *if nStalls > 0* as it is not used in the *else*

	// variables
	// bits or kilobit needed?
	// using kilobits for now
	var avgRate float64
	var maxRate = float64(maxRepRate)
	var nStalls int
	var totalStallDur float64
	var stallPenalty float64
	var segRates []float64
	var sumSegRate float64

	// size of the log map
	logMapSize := len(log)
	// fmt.Printf("The number of log %d\n", logMapSize)

	segRate := log[logMapSize].Bandwidth
	// list of segment rates
	segRates = log[logMapSize].SegmentRates
	// sum of the seg rates
	sumSegRate = log[logMapSize].SumSegRate
	// sum the total stall duration
	totalStallDur = log[logMapSize].TotalStallDur
	// number stalls
	nStalls = log[logMapSize].NumStalls

	if printOutput {
		fmt.Println()
		fmt.Println("current segment bitrate", segRate)
		fmt.Println("list of seg rates: ", segRates)
		fmt.Println("sum of seg rates: ", sumSegRate)
		fmt.Println("Sum of stall times: ", totalStallDur)
		fmt.Println("Number of stalls: ", nStalls)
		fmt.Println("\n=================")
	}

	//avgRate = sumSegRate / float64(logMapSize)
	avgRate = stat.Mean(segRates, nil)

	varianceG := stat.Variance(segRates, nil)
	sd := math.Sqrt(varianceG)

	// qoe rate
	rateQoE := 5.67 * (avgRate) / maxRate

	if printOutput {
		fmt.Printf("The average segment rates is %.4f\n", avgRate)
		fmt.Printf("The maxRate segment rates is %.4f\n", maxRate)
		fmt.Printf("The varianceG segment rates is %.4f\n", varianceG)
		fmt.Printf("The sd value is %.4f\n", sd)
		fmt.Printf("The qoeRate is %.4f\n", rateQoE)
	}

	// standard deviation of the segment duration
	// standrd deviation seems to be:
	// find the mean of the segRates by summing all the rates and dividing by the number of rates
	//mean := sumSegRate / float64(len(segRates))

	//fmt.Printf("The sum of segment rates is %.4f\n", sumSegRate)
	//fmt.Printf("The mean is %.4f\n", mean)

	//sd := 0.0
	// then loop over the rates again, adding the difference between the mean and actual rate
	// start at index 0 for append slice
	//for k := 0; k < logMapSize; k++ {
	// The use of Pow math function func Pow(x, y float64) float64
	//	sd += math.Pow(segRates[k]-mean, 2)
	//}
	//fmt.Printf("The sd pre-Sqrt is %.4f\n", sd)
	// The use of Sqrt math function func Sqrt(x float64) float64
	//sd = math.Sqrt(sd / float64(logMapSize))
	//fmt.Printf("The standard deviation is %.4f\n", sd)
	//fmt.Printf("The standard deviation (gonum) is %.4f\n", stddev_g)

	// if stalls occured
	if nStalls > 0 {
		var stallDurElem float64
		if totalStallDur/float64(nStalls) > 15.0 {
			stallDurElem = 1.0
		} else {
			stallDurElem = totalStallDur / float64(nStalls)
			stallDurElem = stallDurElem / 15.0
		}
		stallPenalty = 4.95 * (0.875*(1+math.Log(float64(nStalls)/float64(logMapSize))/6.0) + 0.125*stallDurElem)
	} else {
		stallPenalty = 0
	}

	switchingPenalty := 6.72 * float64(sd) / float64(maxRate)
	valRate := 0.17 + rateQoE - switchingPenalty - stallPenalty
	totQoE := utils.MaxFloat(0, valRate)

	if printOutput {
		fmt.Printf("The stallPenalty is %.4f\n", stallPenalty)
		fmt.Printf("The switchingPenalty is %.4f\n", switchingPenalty)
		fmt.Printf("The comparison claye QoE is %.4f\n", valRate)
		fmt.Printf("The Claye QoE is %.4f\n", totQoE)
	}

	// calculate the claye value and return to the channel
	c <- totQoE

}

/*
## CLAYES QoE estimation start here ####
	rateQoE = 5.67 * float(avgRate) / maxRate;
	totVar = 0;
	# standard deviation
	rateStd = np.std(segRate)
	stallFree = 1
	stallDurElem = 1.0
	if nStalls != 0:
		stallFree = 0
		if (float(totalStallDur) / nStalls > 15):
			stallDurElem = 1.0;
		else:
			stallDurElem = float(totalStallDur) / float(nStalls)
			stallDurElem = stallDurElem/15.0
		stallPenalty = 4.95 * (0.875*(1 + np.log( nStalls / float(numSeg) ) / 6.0) + 0.125 * stallDurElem);
	else:
		stallPenalty = 0
		stallFree = 1

	switchingPenalty = 6.72 * float(rateStd) / float(maxRate)
	totQoE = max(0,0.17 + rateQoE -  switchingPenalty - stallPenalty)
*/
