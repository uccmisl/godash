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

/*
 * Function  :
 * let's think about HLS - chunk replacement
 * before we decide what chunks to change, lets create a file for HLS
 * then add functions to switch out an old chunk
 * then add some additional logic to decide on how agressive (depends on incoming call)
 * we also need to decide where this case statments should go.
 *	Once per loop (passive), or multiple times (aggressive)
 */
// if HLS is to be used - how dynamic must it be...
// only use HLS if we have at least one segment to replacement
/*
   switch hls {
   // passive - least amount of replacement
   case hlsPassive:
     // competitive - maybe used for gRPC- for when clients compete
   case hlsCompetitive:
     // aggressive - highest amount of replacement
   case hlsAggressive:
     // dynamic - permit the player to switch between the other cases as needed
   case hlsDynamic:

   }
*/

package hlsFunc

import (
	"strings"
	"time"

	glob "github.com/uccmisl/godash/global"
	"github.com/uccmisl/godash/http"
	"github.com/uccmisl/godash/logging"
)

// GetHlsSegment :
/*
 * pass in relevent information to redownload a previous chunk at a higher quality level
 * update the segment map information for the replaced chunk
 */
func GetHlsSegment(f func(segmentNumber int, currentURL string,
	initBuffer int, maxBuffer int, codecName string, codec string, urlString string, urlInput []string,
	mpdList []http.MPD, adapt string, maxHeight int, isByteRangeMPD bool, startTime time.Time,
	nextRunTime time.Time, arrivalTime int, oldMPDIndex int, nextSegmentNumber int, hls string,
	hlsBool bool, mapSegmentLogPrintout map[int]logging.SegPrintLogInformation, numSeg int, extendPrintLog bool,
	hlsUsed bool, bufferLevel int, segmentDurationTotal int, quic string, quicBool bool, baseURL string, debugLog bool) (int, map[int]logging.SegPrintLogInformation), hlsChunkNumber int,
	mapSegmentLogPrintout map[int]logging.SegPrintLogInformation, maxHeight int, urlInput []string, initBuffer int, maxBuffer int, codecName string, codec string, urlString string, mpdList []http.MPD, nextSegmentNumber int, extendPrintLog bool, startTime time.Time, nextRunTime time.Time, arrivalTime int, hlsUsed bool, quic string, quicBool bool, baseURL string, debugFile string, debugLog bool, repRateBaseURL string) (int, map[int]logging.SegPrintLogInformation, int, int, time.Time) {

	// store the segment map details
	previousChunk := mapSegmentLogPrintout[hlsChunkNumber]

	// reset the buffer to a previous level for this hls chunk
	oldBuffer := mapSegmentLogPrintout[hlsChunkNumber-1].BufferLevel

	// get the buffer level for this previous chunk
	currentBuffer := mapSegmentLogPrintout[hlsChunkNumber].BufferLevel

	// reset the total segment duration to a previous level for this hls chunk
	oldSegmentDuration := mapSegmentLogPrintout[hlsChunkNumber-1].PlayStartPosition

	// get the current url - trim any white space
	currentURL := strings.TrimSpace(urlInput[previousChunk.MpdIndex])
	// get the algorithm from the segment map
	adapt := previousChunk.Adapt

	// determine if the MPD is byte range
	isByteRangeMPD := false
	repBaseURL := http.GetRepresentationBaseURL(mpdList[previousChunk.MpdIndex], 0)
	if repBaseURL != repRateBaseURL {
		isByteRangeMPD = true
		logging.DebugPrint(debugFile, debugLog, "DEBUG: ", "Byte-range MPD: ")
	}
	// get the old MPD Index from the segment map
	oldMPDIndex := previousChunk.MpdIndex
	// turn off hls for this call - otherwise we could have recursive calls to hls
	hls := "off"
	hlsBool := false

	// add an additional segment duration to the current Segmetn duration so it only runs once
	newStreamduration := oldSegmentDuration + (previousChunk.SegmentDuration * glob.Conversion1000)

	// reduce the initBuffer to zero, so we are constantly counting segment number
	_, newChunkMap := f(hlsChunkNumber, currentURL, 0, maxBuffer, codecName, codec, urlString, urlInput, mpdList, adapt, maxHeight, isByteRangeMPD,
		startTime, nextRunTime, arrivalTime, oldMPDIndex, nextSegmentNumber, hls, hlsBool, mapSegmentLogPrintout, newStreamduration,
		extendPrintLog, hlsUsed, oldBuffer, oldSegmentDuration, quic, quicBool, baseURL, debugLog)

	// reset the buffer to a previous level for this hls chunk
	newBuffer := mapSegmentLogPrintout[hlsChunkNumber].BufferLevel

	// get the difference between the previous and current buffer levels
	bufferDifference := currentBuffer - newBuffer

	thisRunTimeVal := int(time.Since(nextRunTime).Nanoseconds() / (glob.Conversion1000 * glob.Conversion1000))

	nextRunTime = time.Now()

	// return the next segment and the updated segment map
	return nextSegmentNumber, newChunkMap, bufferDifference, thisRunTimeVal, nextRunTime
}

// ChangeBufferLevels :
// change the buffer levels in a range of segment maps
func ChangeBufferLevels(mapSegmentLogPrintout map[int]logging.SegPrintLogInformation, segmentNumber int, chunkReplace int, bufferDifference int) map[int]logging.SegPrintLogInformation {

	// lets only find the maps between the chunk we replaced and the current chunk
	for a := chunkReplace + 1; a < segmentNumber; a++ {
		// for this chunk, get the old buffer level subtract the new buffer difference
		// and re-allocate to the old map location
		localMap := mapSegmentLogPrintout[a]
		localMap.BufferLevel = mapSegmentLogPrintout[a].BufferLevel - bufferDifference
		mapSegmentLogPrintout[a] = localMap
	}

	return mapSegmentLogPrintout
}
