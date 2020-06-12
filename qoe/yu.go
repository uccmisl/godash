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
	"fmt"

	"github.com/uccmisl/godash/logging"
)

func getYu(log map[int]logging.SegPrintLogInformation, c chan float64, printOutput bool) {

	var totalStall float64
	var segRates []float64
	var sumSegRate float64
	var rateChange []float64
	var sumRateChange float64
	var rateDifference float64
	var nSwitches int
	var avgRateSwitchMagnitude float64

	// I think this goes to segRate -1, to not have an index out of range error
	// stopValue = len(segRate) - 1
	// I think we can use log size for the rates length
	// size of the log map
	logMapSize := len(log)

	//model paramters
	w1 := float64(1 / 3)
	w2 := 20.0
	if printOutput {
		fmt.Println("Yu Logs")
	}

	// get the values direct from the logs for this segment
	segRate := log[logMapSize].Bandwidth
	// list of segment rates
	segRates = log[logMapSize].SegmentRates
	// sum of the seg rates
	sumSegRate = log[logMapSize].SumSegRate
	// sum the total stall duration
	totalStall = log[logMapSize].TotalStallDur
	// number stalls
	// nStalls = log[logMapSize].NumStalls
	// number switches
	// nSwitches = log[logMapSize].NumSwitches
	// rate changes
	rateDifference = log[logMapSize].RateDifference
	// sum of rate changes
	sumRateChange = log[logMapSize].SumRateChange
	// list of rate differences
	rateChange = log[logMapSize].RateChange

	if printOutput {
		fmt.Println()
		fmt.Println("current segment bitrate", segRate)
		fmt.Println("list of seg rates: ", segRates)
		fmt.Println("sum of seg rates: ", sumSegRate)
		fmt.Println("Sum of stall times: ", totalStall)
		// fmt.Println("Number of stalls: ", nStalls)
		// fmt.Println("Number of switches: ", nSwitches)
		fmt.Println("Current Rate difference: ", rateDifference)
		fmt.Println("Sum of rate changes: ", sumRateChange)
		fmt.Println("List of Rate differences: ", rateChange)
		fmt.Println("\n=================")
	}

	// 	avgMediaQuality = float(np.mean(segRate))
	avgBitrate := (sumSegRate / 1000000) / float64(logMapSize)
	if printOutput {
		fmt.Printf("Average BitRate: %f = Sum of segment rate: %f divided by Number of segments: %d\n", avgBitrate, sumSegRate, logMapSize)
	}

	// qualitySwitchingFrequency = float(np.mean(qualitySwitching))
	avgRateSwitchMagnitude = (sumRateChange / 1000000) / float64(logMapSize)

	if printOutput {
		fmt.Printf("Average Rate Switch Magnitude: %f = Sum of rate changes: %f divided by Number of switches: %d\n", avgRateSwitchMagnitude, sumRateChange, nSwitches)
	}

	// totalDisplayTime = len(segRate)*segDuration + totalStall
	// totalDisplayTime, considering changes in segment duration
	totalDisplayTime := 0.0
	i := 1
	for i <= logMapSize {
		totalDisplayTime += float64(log[i].SegmentDuration)
		i++
	}
	totalDisplayTime += totalStall
	if printOutput {
		fmt.Printf("Total Display Time: %f = number of segment: %d * segment duration %d plus total stall duration: %f\n", totalDisplayTime, logMapSize, log[1].SegmentDuration, totalStall)
	}

	starvation := totalStall / totalDisplayTime
	if printOutput {
		fmt.Printf("Total Starvation Time: %f = Sum of stall times: %f divided by Total Display Time: %f\n", starvation, totalStall, totalDisplayTime)
	}

	// switchingQoE = w1 * qualitySwitchingFrequency
	switchingQoE := w1 * float64(avgRateSwitchMagnitude)
	// starvationQoE = w2 * starvation
	starvationQoE := float64(w2 * starvation)
	if printOutput {
		fmt.Printf("Switching QoE: %f, Starvation QoE: %f\n", switchingQoE, starvationQoE)
	}

	// yuOoE = avgMediaQuality - w1*qualitySwitchingFrequency - w2*starvation
	// returned QoE value
	returnedQoE := float64(avgBitrate) - switchingQoE - starvationQoE

	if printOutput {
		fmt.Println("Yu value: ", returnedQoE)
	}

	// calculate the Yu value and return to the channel
	c <- returnedQoE

}
