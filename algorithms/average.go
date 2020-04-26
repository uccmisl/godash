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

//MeanAverageAlgo : "normal average" -> take all the throughtputs and make the average
//call the func meanAverage with all the values of throughtput to make a "standard" average
func MeanAverageAlgo(thrList *[]int, newThr int, repRate *int, bandwithList []int, lowestMPDrepRateIndex int) {

	var average float64

	*thrList = append(*thrList, newThr)

	//if there is not enough throughtputs in the list, can't calculate the average
	if len(*thrList) < 2 {
		//if there is not enough throughtput, we call selectRepRate() with the newThr
		*repRate = SelectRepRateWithThroughtput(newThr, bandwithList, lowestMPDrepRateIndex)
		return
	}

	// average of the last throughtputs
	meanAverage(*thrList, &average)

	//fmt.Println("AVERAGE: ", int(average))

	//We select the reprate with the calculated throughtput
	*repRate = SelectRepRateWithThroughtput(int(average), bandwithList, lowestMPDrepRateIndex)
}

// meanAverage :
//calculate the average of all the elements in the list
func meanAverage(thrList []int, average *float64) {

	*average = 0.0

	for i := 0; i < len(thrList); i++ {
		*average += float64(thrList[i])
	}
	*average = *average / float64(len(thrList))

}
