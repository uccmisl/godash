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

import "math"

//GEOM AVERAGE -> geometric average (square root of the thr)

// GeomAverageAlgo :
// call the func geomAverage with all the values of the throughtput list
// to make a geometric average
func GeomAverageAlgo(thrList *[]int, newThr int, repRate *int, bandwithList []int, lowestMPDrepRateIndex int) {

	var average float64

	*thrList = append(*thrList, newThr)

	//if there is not enough throughtputs in the list, can't calculate the average
	if len(*thrList) < 2 {
		//if there is not enough throughtput, we call selectRepRate() with the newThr
		*repRate = SelectRepRateWithThroughtput(newThr, bandwithList, lowestMPDrepRateIndex)
		return
	}

	// average of the last throughtputs
	geomAverage(*thrList, &average)

	//We select the reprate with the calculated throughtput
	*repRate = SelectRepRateWithThroughtput(int(average), bandwithList, lowestMPDrepRateIndex)

}

// geomAverage :
//calculate the geometric average of all the elements in the list
func geomAverage(thrList []int, average *float64) {

	*average = 1.0

	//calculate the geometric average : nb *= nextNb / nb^(1/ len(thrList))
	for i := 0; i < len(thrList); i++ {
		*average *= float64(thrList[i])
	}
	*average = math.Pow(*average, 1/float64(len(thrList)))

}
