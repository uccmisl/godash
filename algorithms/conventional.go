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

// the throughtput is equal to -1 at the beginning by default
var thr = -1

//Conventional :
/*
* calculate of the throughtput with the ancient one and the new one
* call the func to select the repRate from the throughtput
* return the repRate and the list of throughtput
 */
func Conventional(thrList *[]int, newThr int, repRate *int, bandwithList []int, lowestMPDrepRateIndex int) {

	//if it is the first throughtput in the list, add it to the list
	if thr == -1 {
		thr = newThr
		*thrList = append(*thrList, thr)
	} else {
		//if there is already one thr, calculate the thr that will be added to the list
		//with 80% of the last thr and 20% of the new one
		thr = (8*thr)/10 + (2*newThr)/10
		*thrList = append(*thrList, thr)
	}

	*repRate = SelectRepRateWithThroughtput(thr, bandwithList, lowestMPDrepRateIndex)
}
