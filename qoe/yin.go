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

func getYin(log map[int]logging.SegPrintLogInformation, c chan float64, initBuffer int, printOutput bool) {

	var totalStall float64
	var nStalls int
	var segRates []float64
	var sumSegRate float64
	var rateChange []float64
	var sumRateChange float64
	var rateDifference float64
	var nSwitches int
	var avgRateSwitchMagnitude float64
	var returnedQoE float64

	// size of the log map
	logMapSize := len(log)

	if printOutput {
		fmt.Println("Yin Logs")
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
	nStalls = log[logMapSize].NumStalls
	// number switches
	nSwitches = log[logMapSize].NumSwitches
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
		fmt.Println("sum of seg rates: ", sumSegRate/1000)
		fmt.Println("Sum of stall times: ", totalStall)
		fmt.Println("Number of stalls: ", nStalls)
		fmt.Println("Number of switches: ", nSwitches)
		fmt.Println("Current Rate difference: ", rateDifference)
		fmt.Println("Sum of rate changes: ", sumRateChange)
		fmt.Println("List of Rate differences: ", rateChange)
		fmt.Println("PlayBackTime: ", log[logMapSize].PlaybackTime)

		fmt.Println("\n=================")
	}

	if nSwitches == 0 {
		avgRateSwitchMagnitude = 0
	} else {
		avgRateSwitchMagnitude = (sumRateChange / 1000)
	}

	if printOutput {
		fmt.Printf("Average Rate Switch Magnitude: %.2f = Sum of rate changes: %.2f divided by Number of switches: %d\n", avgRateSwitchMagnitude, sumRateChange, nSwitches)
	}

	// initial delay, considering changes in segment duration
	initDelay := 0
	i := 1
	for i <= initBuffer {
		initDelay += log[i].SegmentDuration
		i++
	}
	if printOutput {
		fmt.Printf("Initial Delay: %d\n", initDelay)
	}

	// avgRate - 1*avgQualityVariation - 3000*totalStallDur - 3000*initDelay
	// returned QoE value
	if log[logMapSize].PlaybackTime > 0 {
		returnedQoE = (sumSegRate / 1000) - 1*avgRateSwitchMagnitude - 3*totalStall
	} else {
		returnedQoE = (sumSegRate / 1000) - 1*avgRateSwitchMagnitude - 3*totalStall - 3000*float64(initDelay)
	}

	if printOutput {
		fmt.Println("Yin value: ", returnedQoE)
	}

	// calculate the Yin value and return to the channel
	c <- float64(returnedQoE)

}
