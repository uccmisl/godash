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

func getDuanmu(log map[int]logging.SegPrintLogInformation, c chan float64, initBuffer int, printOutput bool) {

	var sessionDuration int
	var totalStall float64
	var nStalls int
	var segRates []float64
	var sumSegRate float64
	var rateChange []float64
	var sumRateChange float64
	var rateDifference float64
	var nSwitches int
	var avgRateSwitchMagnitude float64

	// size of the log map
	logMapSize := len(log)

	if printOutput {
		fmt.Println("Duanmu Logs")
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
		fmt.Println("sum of seg rates: ", sumSegRate)
		fmt.Println("Sum of stall times: ", totalStall)
		fmt.Println("Number of stalls: ", nStalls)
		fmt.Println("Number of switches: ", nSwitches)
		fmt.Println("Current Rate difference: ", rateDifference)
		fmt.Println("Sum of rate changes: ", sumRateChange)
		fmt.Println("List of Rate differences: ", rateChange)
		fmt.Println("\n=================")
	}

	// get the current session duration and rebuffer percentage
	sessionDuration = log[logMapSize].PlaybackTime
	//  for the inital 2 segments this is zero, so to catch a NAN error, reset to 1
	if sessionDuration == 0 {
		sessionDuration = 1
	}
	rebufferPercentage := totalStall / float64(sessionDuration)
	if printOutput {
		fmt.Println("Current session Duration: ", sessionDuration)
		fmt.Println("Rebuffer Percentage: ", rebufferPercentage)
	}

	if nSwitches == 0 {
		avgRateSwitchMagnitude = 0
	} else {
		avgRateSwitchMagnitude = (sumRateChange / 1000) / float64(nSwitches)
	}
	if printOutput {
		fmt.Printf("Average Rate Switch Magnitude: %f = Sum of rate changes: %f divided by Number of switches: %f\n", avgRateSwitchMagnitude, sumRateChange, float64(nSwitches))
	}

	// *** THESE TO BE DETERMINED ***
	// inital delay in seconds => segment duration * segment number?
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

	// sumSegRate divided by logMapSize?
	avgBitrate := (sumSegRate / 1000) / float64(logMapSize)
	if printOutput {
		fmt.Printf("Initial Delay: %d = Segment duration in seconds: %d multiplied by inital buffered number of segments: %d\n", initDelay, log[1].SegmentDuration, initBuffer)
		fmt.Printf("Average BitRate: %f = Sum of segment rate: %f divided by Number of segments: %f\n", avgBitrate, sumSegRate, float64(logMapSize))
	}

	qoe := -2.3*float64(initDelay) - 56.5*rebufferPercentage + 0.0070*avgBitrate + 0.0007*avgRateSwitchMagnitude + 54.0

	if printOutput {
		fmt.Println("Duanmu value: ", qoe)
	}

	// returned Duanmu value
	returnedQoE := qoe

	// calculate the Duanmu value and return to the channel
	c <- returnedQoE

}
