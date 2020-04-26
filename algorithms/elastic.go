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
)

// ElasticAlgo call the func harmonicAverage with the last 5 values of throughtput
// to have a better estimate of the throughtput
func ElasticAlgo(thrList *[]int, newThr int, delTime int, maxBuffer int, repRate *int, bandwithList []int, staticAlgParameter *float64, bufferLevel int, kP float64, kI float64, lowestMPDrepRateIndex int) {

	//number of last averages we are going to take

	//we should take the average of the last 5 chunks
	//var harmonicAverageValue = 3
	var harmonicAverageValue = 5

	var averageRateEstimate float64
	*thrList = append(*thrList, newThr)

	/*
		//if there is not enough throughtputs in the list, can't calculate the average
		if len(*thrList) < harmonicAverageValue {
			*thrList = append(*thrList, newThr)
			averageRateEstimate = float64(newThr)

			//fmt.Println("newThr: ",newThr,"bandwithList: ", bandwithList)

			//if there is not enough throughtput, we call selectRepRate() with the newThr
			*repRate = SelectRepRateWithThroughtput(newThr, bandwithList, repRatesReversed, lowestMPDrepRateIndex)
			return
		}
	*/
	// fmt.Println("maxbuffer", maxBuffer)
	// fmt.Println("currbuffer", bufferLevel)
	// harmonic average of the last throughtputs
	harmonicAverage(harmonicAverageValue, *thrList, &averageRateEstimate)
	// fmt.Println("all", averageRateEstimate/1000)
	*staticAlgParameter += (float64(delTime) / glob.Conversion1000) * (float64(bufferLevel)/glob.Conversion1000 - float64(maxBuffer))
	targetRate := averageRateEstimate / (1 - kP*float64(bufferLevel/glob.Conversion1000) - kI*float64(*staticAlgParameter))
	// fmt.Println("target thr: ", int(targetRate), "bandwithList: ", bandwithList)

	*repRate = SelectRepRateWithThroughtput(int(targetRate), bandwithList, lowestMPDrepRateIndex)
}
