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

//import (
//	"math"
//)

//EMWA AVERAGE -> exponential average

// EMWAAverageAlgo :
/*
call the func expAverage with all the values of the throughtput list and a window size
to make an exponential average
*/
func EMWAAverageAlgo(thrList *[]int, repRate *int, exponentialRatio float64, window int, newThr int, bandwithList []int, lowestMPDrepRateIndex int) {

	var average float64

	*thrList = append(*thrList, newThr)

	//if there is not enough throughtputs in the list, can't calculate the average
	if len(*thrList) < 2 {
		//if there is not enough throughtput, we call selectRepRate() with the newThr
		*repRate = SelectRepRateWithThroughtput(newThr, bandwithList, lowestMPDrepRateIndex)
		return
	}

	// average of the last throughtputs
	ExpAverage(*thrList, exponentialRatio, window, &average)

	//fmt.Println("AVERAGE: ", int(average))

	//We select the reprate with the calculated throughtput
	*repRate = SelectRepRateWithThroughtput(int(average), bandwithList, lowestMPDrepRateIndex)
}

// expAverage :
//calculate the geometric average of the last *window* elements in the list
