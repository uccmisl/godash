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

// ./goDASH -url "[http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/10_sec/x264/bbb/DASH_Files/full/bbb_enc_x264_dash.mpd, http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/10_sec/x264/bbb/DASH_Files/full/bbb_enc_x264_dash.mpd]" -adapt default -codec AVC -debug true -initBuffer 2 -maxBuffer 10 -maxHeight 1080 -numSegments 20  -storeDASH 347985

package player

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	algo "github.com/uccmisl/godash/algorithms"
	glob "github.com/uccmisl/godash/global"
	"github.com/uccmisl/godash/hlsfunc"
	"github.com/uccmisl/godash/http"
	"github.com/uccmisl/godash/logging"
	"github.com/uccmisl/godash/utils"
)

// play position
var playPosition = 0

// current segment number
var segmentNumber = 1
var segmentDuration int
var nextSegmentNumber int

// current buffer level
var bufferLevel = 0
var maxBufferLevel int
var waitToPlayCounter = 0
var stallTime = 0

// current mpd file
var mpdListIndex = 0
var lowestMPDrepRateIndex int
var highestMPDrepRateIndex int

// save the previous mpdIndex
var oldMPDIndex = 0

// determine if an MPD is byte-range or not
var isByteRangeMPD bool
var startRange = 0
var endRange = 0

// current representation rate
var repRate = 1

//var repRatesReversed bool

// current adaptationSet
var currentMPDRepAdaptSet int

// Segment size (in bits)
var segSize int

// baseURL for this MPD file
var baseURL string

// we need to keep a tab on the different size segments - use this for now
// we will use an array in the future
var segmentDurationTotal = 0
var segmentDurationArray []int

// the list of bandwith values (rep_rates) from the current MPD file
var bandwithList []int

// list of throughtputs - noted from downloading the segments
var thrList []int

// time values
var startTime time.Time
var nextRunTime time.Time
var arrivalTime int

// additional output logs values
var repCodec string
var repHeight int
var repWidth int
var repFps int

// used to calculate targetRate - float64
var kP = 0.01
var kI = 0.001
var staticAlgParameter = 0.0

// first step is to check the first MPD for the codec (I had problem passing a
// 2-dimensional array, so I moved the check to here)
var codecList [][]string
var codecIndexList [][]int
var usedCodec bool
var codecIndex int

var urlInput []string

// For the mapSegments of segments :
// Map with the segment number and a structure of informations
// one map contains all content
var mapSegmentLogPrintout map[int]logging.SegPrintLogInformation

// a map of maps containing segment header information
var segHeadValues map[int]map[int][]int

// default value for the exponential ratio
var exponentialRatio float64

// file download location
var fileDownloadLocation string

// printHeadersData local
var printHeadersData map[string]string

// print the log to terminal
var printLog bool

// variable to determine if we are using the goDASHbed testbed
var useTestbedBool bool

// Stream :
/*
 * get the header file for the current video clip
 * check the different arguments in order to stream
 * call streamLoop to begin to stream
 */
func Stream(mpdList []http.MPD, debugFile string, debugLog bool, codec string, codecName string,
	maxHeight int, streamDuration int, maxBuffer int, initBuffer int,
	adapt string, urlString string, fileDownloadLocationIn string, extendPrintLog bool, hls string, hlsBool bool, quic string, quicBool bool,
	getHeaderBool bool, getHeaderReadFromFile string, exponentialRatioIn float64, printHeadersDataIn map[string]string, printLogIn bool,
	useTestbedBoolIn bool) {

	// check if the codec is in the MPD urls passed in
	codecList, codecIndexList = http.GetCodec(mpdList, codec, debugLog)
	// determine if the passed in codec is one of the codecs we use (checking the first MPD only)
	usedCodec, codecIndex = utils.FindInStringArray(codecList[0], codec)

	// set local val
	exponentialRatio = exponentialRatioIn
	fileDownloadLocation = fileDownloadLocationIn
	printHeadersData = printHeadersDataIn
	printLog = printLogIn
	useTestbedBool = useTestbedBoolIn

	// check the codec and print error is false
	if !usedCodec {
		// print error message
		fmt.Printf("*** -" + codecName + " " + codec + " is not in the first provided MPD, please check " + urlString + " ***\n")
		// stop the app
		utils.StopApp()
	}

	// save the current MPD Rep_rate Adaptation Set
	currentMPDRepAdaptSet = codecIndexList[mpdListIndex][codecIndex]

	// determine if we are using a byte-range or standard MPD profile
	// the xml Representation>BaseURL is saved in the same location
	// for byte range full, main and onDemand
	// so check for BaseURL, if not empty, then its a byte-range
	baseURL := http.GetRepresentationBaseURL(mpdList[mpdListIndex], 0)
	if baseURL != glob.RepRateBaseURL {
		isByteRangeMPD = true
		logging.DebugPrint(debugFile, debugLog, "DEBUG: ", "Byte-range MPD: ")
	}

	// get the relevent values from this MPD
	// maxSegments - number of segments to download
	// maxBufferLevel - maximum buffer level in seconds
	// highestMPDrepRateIndex - index with the highest rep_rate
	// lowestMPDrepRateIndex - index with the lowest rep_rate
	// segmentDuration - segment duration
	// bandwithList - get all the range of representation bandwiths of the MPD

	// maxSegments was the first value
	_, maxBufferLevel, highestMPDrepRateIndex, lowestMPDrepRateIndex, segmentDurationArray, bandwithList, baseURL = http.GetMPDValues(mpdList, mpdListIndex, maxHeight, streamDuration, maxBuffer, currentMPDRepAdaptSet, isByteRangeMPD, debugLog)

	// print values to debug log
	logging.DebugPrint(debugFile, debugLog, "\nDEBUG: ", "streaming has begun")
	logging.DebugPrint(debugFile, debugLog, "DEBUG: ", "Input values to streaming algorithm: "+adapt)
	logging.DebugPrint(debugFile, debugLog, "DEBUG: ", "maxHeight: "+strconv.Itoa(maxHeight))
	logging.DebugPrint(debugFile, debugLog, "DEBUG: ", "streamDuration in seconds: "+strconv.Itoa(streamDuration))
	logging.DebugPrint(debugFile, debugLog, "DEBUG: ", "maxBuffer: "+strconv.Itoa(maxBuffer))
	logging.DebugPrint(debugFile, debugLog, "DEBUG: ", "initBuffer: "+strconv.Itoa(initBuffer))
	logging.DebugPrint(debugFile, debugLog, "DEBUG: ", "url: "+urlString)
	logging.DebugPrint(debugFile, debugLog, "DEBUG: ", "fileDownloadLocation: "+fileDownloadLocation)
	logging.DebugPrint(debugFile, debugLog, "DEBUG: ", "HLS: "+hls)
	logging.DebugPrint(debugFile, debugLog, "DEBUG: ", "extend: "+strconv.FormatBool(extendPrintLog))

	// get the stream header from the required MPD (first index in the mpdList)
	headerURL := http.GetFullStreamHeader(mpdList[mpdListIndex], isByteRangeMPD)
	logging.DebugPrint(debugFile, debugLog, "DEBUG: ", "stream initialise URL header: "+headerURL)

	// convert the url strings to a list
	urlInput = http.URLList(urlString)

	// get the current url - trim any white space
	currentURL := strings.TrimSpace(urlInput[mpdListIndex])
	logging.DebugPrint(debugFile, debugLog, "\nDEBUG: ", "current URL header: "+currentURL)

	// set the segmentDuration to the first passed in URL
	segmentDuration = segmentDurationArray[0]

	// determine the inital variables to set, based on the algorithm choice
	switch adapt {
	case glob.ConventionalAlg:
		// there is no byte range in this file, so we set byte-range bool to false
		// we don't want to add the seg duration to this file, so 'addSegDuration' is false
		http.GetFile(currentURL, baseURL+headerURL, fileDownloadLocation, false, startRange, endRange, segmentNumber,
			segmentDuration, false, quicBool, debugFile, debugLog, useTestbedBool, repRate)
		// set the inital rep_rate to the lowest value index
		repRate = lowestMPDrepRateIndex
	case glob.ElasticAlg:
		//fmt.Println("Elastic / in player.go")
		//fmt.Println("currentURL: ", currentURL)
		http.GetFile(currentURL, baseURL+headerURL, fileDownloadLocation, false, startRange, endRange, segmentNumber,
			segmentDuration, false, quicBool, debugFile, debugLog, useTestbedBool, repRate)
		repRate = lowestMPDrepRateIndex
		///fmt.Println("MPD file repRate index: ", repRate)
		//fmt.Println("MPD file bandwithList[repRate]", bandwithList[repRate])
	case glob.ProgressiveAlg:
		// get the header file
		// there is no byte range in this file, so we byte-range bool to false
		http.GetFileProgressively(currentURL, baseURL+headerURL, fileDownloadLocation, false, startRange, endRange, segmentNumber, segmentDuration, false, debugLog)
	case glob.TestAlg:
		fmt.Println("testAlg / in player.go")
		http.GetFile(currentURL, baseURL+headerURL, fileDownloadLocation, false, startRange, endRange, segmentNumber,
			segmentDuration, false, quicBool, debugFile, debugLog, useTestbedBool, repRate)

		//fmt.Println("lowestmpd: ", lowestMPDrepRateIndex)
		repRate = lowestMPDrepRateIndex

	case glob.ArbiterAlg:
		//fmt.Println("ArbiterAlg / in player.go")
		http.GetFile(currentURL, baseURL+headerURL, fileDownloadLocation, false, startRange, endRange, segmentNumber,
			segmentDuration, false, quicBool, debugFile, debugLog, useTestbedBool, repRate)

		repRate = lowestMPDrepRateIndex

	case glob.LogisticAlg:
		http.GetFile(currentURL, baseURL+headerURL, fileDownloadLocation, false, startRange, endRange, segmentNumber,
			segmentDuration, false, quicBool, debugFile, debugLog, useTestbedBool, repRate)
		repRate = lowestMPDrepRateIndex
	case glob.MeanAverageAlg:
		http.GetFile(currentURL, baseURL+headerURL, fileDownloadLocation, false, startRange, endRange, segmentNumber,
			segmentDuration, false, quicBool, debugFile, debugLog, useTestbedBool, repRate)
	case glob.GeomAverageAlg:
		http.GetFile(currentURL, baseURL+headerURL, fileDownloadLocation, false, startRange, endRange, segmentNumber,
			segmentDuration, false, quicBool, debugFile, debugLog, useTestbedBool, repRate)
	case glob.EMWAAverageAlg:
		http.GetFile(currentURL, baseURL+headerURL, fileDownloadLocation, false, startRange, endRange, segmentNumber,
			segmentDuration, false, quicBool, debugFile, debugLog, useTestbedBool, repRate)
	}
	// debug logs
	logging.DebugPrint(debugFile, debugLog, "\nDEBUG: ", "We are using repRate: "+strconv.Itoa(repRate))
	logging.DebugPrint(debugFile, debugLog, "DEBUG: ", "We are using : "+adapt+" for streaming")

	//create the map for the print log
	mapSegmentLogPrintout = make(map[int]logging.SegPrintLogInformation)

	//StartTime of downloading
	startTime = time.Now()
	nextRunTime = time.Now()

	// get the segment headers and stop this run
	if getHeaderBool {
		// get the segment headers for all MPD url passed as arguments - print to file
		http.GetAllSegmentHeaders(mpdList, codecIndexList, maxHeight, 1, streamDuration, isByteRangeMPD, maxBuffer, headerURL, codec, urlInput, debugLog, true)

		// print error message
		fmt.Printf("*** - All segment header have been downloaded to " + glob.DebugFolder + " - ***\n")
		// exit the application
		os.Exit(3)
	} else {
		if getHeaderReadFromFile == glob.GetHeaderOnline {
			// get the segment headers for all MPD url passed as arguments - not from file
			segHeadValues = http.GetAllSegmentHeaders(mpdList, codecIndexList, maxHeight, 1, streamDuration, isByteRangeMPD, maxBuffer, headerURL, codec, urlInput, debugLog, false)
		} else if getHeaderReadFromFile == glob.GetHeaderOffline {
			// get the segment headers for all MPD url passed as arguments - yes from file
			// get headers from file for a given number of seconds of stream time
			// let's assume every n seconds
			segHeadValues = http.GetNSegmentHeaders(mpdList, codecIndexList, maxHeight, 1, streamDuration, isByteRangeMPD, maxBuffer, headerURL, codec, urlInput, debugLog, true)

		}
	}

	// print the output log headers
	logging.PrintHeaders(extendPrintLog, fileDownloadLocation, glob.LogDownload, debugFile, debugLog, printLog, printHeadersData)

	// Streaming loop function - using the first MPD index - 0, and hlsUsed false
	segmentNumber, mapSegmentLogPrintout = streamLoop(segmentNumber, currentURL, initBuffer, maxBuffer, codecName, codec, urlString, urlInput,
		mpdList, adapt, maxHeight, isByteRangeMPD, startTime, nextRunTime, arrivalTime, 0, 0, hls, hlsBool, mapSegmentLogPrintout, streamDuration,
		extendPrintLog, false, bufferLevel, segmentDurationTotal, quic, quicBool, baseURL, debugLog)

	// print sections of the map to the debug log - if debug is true
	if debugLog {
		logging.PrintsegInformationLogMap(debugFile, debugLog, mapSegmentLogPrintout)
	}

	// print out the rest of the play out segments - based on playStartPosition of the last segment streamed
	// and an end time that includes for the original initial buffer size in seconds
	logging.PrintPlayOutLog(mapSegmentLogPrintout[segmentNumber-1].PlayStartPosition+mapSegmentLogPrintout[initBuffer].PlayStartPosition, initBuffer, mapSegmentLogPrintout, glob.LogDownload, printLog, printHeadersData)

}

// streamLoop :
/*
 * take the first segment number, download it with a low quality
 * call itself with the next segment number
 */
func streamLoop(segmentNumber int, currentURL string,
	initBuffer int, maxBuffer int, codecName string, codec string, urlString string, urlInput []string,
	mpdList []http.MPD, adapt string, maxHeight int, isByteRangeMPD bool, startTime time.Time,
	nextRunTime time.Time, arrivalTime int, oldMPDIndex int, nextSegmentNumber int, hls string,
	hlsBool bool, mapSegmentLogPrintout map[int]logging.SegPrintLogInformation, streamDuration int, extendPrintLog bool,
	hlsUsed bool, bufferLevel int, segmentDurationTotal int, quic string, quicBool bool, baseURL string,
	debugLog bool) (int, map[int]logging.SegPrintLogInformation) {

	// variable for rtt for this segment
	var rtt time.Duration
	// has this chunk been replaced by hls
	var hlsReplaced = "no"
	// if we undertake HLS, we need to revise the buffer values
	var bufferDifference int
	// if we set this chunk to HLS used
	if hlsUsed {
		hlsReplaced = "yes"
	}
	var segURL string

	// save point for the HTTP protocol used
	var protocol string

	//
	var segmentFileName string
	/*
	 * Function  :
	 * let's think about HLS - chunk replacement
	 * before we decide what chunks to change, lets create a file for HLS
	 * then add functions to switch out an old chunk
	 */
	// only use HLS if we have at least one segment to replacement
	if hlsBool && segmentNumber > 1 {
		switch hls {
		// passive - least amount of replacement
		case glob.HlsOn:
			if segmentNumber == 6 {
				// hlsUsed is set to true
				chunkReplace := 5
				var thisRunTimeVal int
				// replace a previously downloaded segment with this call
				nextSegmentNumber, mapSegmentLogPrintout, bufferDifference, thisRunTimeVal, nextRunTime =
					hlsfunc.GetHlsSegment(streamLoop, chunkReplace, mapSegmentLogPrintout, maxHeight, urlInput,
						initBuffer, maxBuffer, codecName, codec, urlString, mpdList, nextSegmentNumber, extendPrintLog,
						startTime, nextRunTime, arrivalTime, true, quic, quicBool, baseURL, glob.DebugFile, debugLog,
						glob.RepRateBaseURL)

				// change the current buffer to reflect the time taken to get this HLS segment
				bufferLevel -= (thisRunTimeVal + bufferDifference)

				// change the buffer levels of the previous chunks, so the printout reflects this value
				mapSegmentLogPrintout = hlsfunc.ChangeBufferLevels(mapSegmentLogPrintout, segmentNumber, chunkReplace, bufferDifference)
			}
		}
	}

	// if we have changed the MPD, we need to update some variables
	if oldMPDIndex != mpdListIndex {

		// set the new mpdListIndex
		mpdListIndex = oldMPDIndex

		// get the current url - trim any white space
		currentURL = strings.TrimSpace(urlInput[mpdListIndex])
		logging.DebugPrint(glob.DebugFile, debugLog, "\nDEBUG: ", "current URL header: "+currentURL)

		// get the relavent values from this MPD
		streamDuration, maxBufferLevel, highestMPDrepRateIndex, lowestMPDrepRateIndex, segmentDurationArray, bandwithList, baseURL = http.GetMPDValues(mpdList, mpdListIndex, maxHeight, streamDuration, maxBuffer, currentMPDRepAdaptSet, isByteRangeMPD, debugLog)

		// current segment duration
		segmentDuration = segmentDurationArray[mpdListIndex]

		// ONLY CHANGE THE NUMBER OF SEGMENTS HERE
		//	numSegments := streamDuration / segmentDuration

		//	fmt.Println(segmentNumber)
		//	fmt.Println(segmentDuration)
		//	fmt.Println(numSegments)

		// determine if the passed in codec is one of the codecs we use (checking the current MPD)
		usedCodec, codecIndex = utils.FindInStringArray(codecList[mpdListIndex], codec)
		// check the codec and print error is false
		if !usedCodec {
			// print error message
			fmt.Printf("*** -" + codecName + " " + codec + " is not in the provided MPD, please check " + urlString + " ***\n")
			// stop the app
			utils.StopApp()
		}
		// save the current MPD Rep_rate Adaptation Set
		currentMPDRepAdaptSet = codecIndexList[mpdListIndex][codecIndex]
	}

	// break out if we have downloaded all of our segments
	// which is current segment duration total plus the next segment to be downloaded
	if segmentDurationTotal+(segmentDuration*glob.Conversion1000) > streamDuration {
		return segmentNumber, mapSegmentLogPrintout
	}

	// keep rep_rate within the index boundaries
	// MISL - might cause problems
	if repRate < highestMPDrepRateIndex {
		logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "Changing rep_rate index: from "+strconv.Itoa(repRate)+" to "+strconv.Itoa(highestMPDrepRateIndex))
		repRate = highestMPDrepRateIndex
	}

	// get the segment
	if isByteRangeMPD {
		segURL, startRange, endRange = http.GetNextByteRangeURL(mpdList[mpdListIndex], segmentNumber, repRate, currentMPDRepAdaptSet)
		logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "byte start range: "+strconv.Itoa(startRange))
		logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "byte end range: "+strconv.Itoa(endRange))
	} else {
		segURL = http.GetNextSegment(mpdList[mpdListIndex], segmentNumber, repRate, currentMPDRepAdaptSet)
	}
	logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "current segment URL: "+segURL)

	// Start Time of this segment
	currentTime := time.Now()

	// Download the segment - add the segment duration to the file name
	switch adapt {
	case glob.ConventionalAlg:
		rtt, segSize, protocol, segmentFileName = http.GetFile(currentURL, baseURL+segURL, fileDownloadLocation, isByteRangeMPD, startRange, endRange, segmentNumber, segmentDuration, true, quicBool, glob.DebugFile, debugLog, useTestbedBool, repRate)
	case glob.ElasticAlg:
		rtt, segSize, protocol, segmentFileName = http.GetFile(currentURL, baseURL+segURL, fileDownloadLocation, isByteRangeMPD, startRange, endRange, segmentNumber, segmentDuration, true, quicBool, glob.DebugFile, debugLog, useTestbedBool, repRate)
	case glob.ProgressiveAlg:
		rtt, segSize = http.GetFileProgressively(currentURL, baseURL+segURL, fileDownloadLocation, isByteRangeMPD, startRange, endRange, segmentNumber, segmentDuration, true, debugLog)
	case glob.LogisticAlg:
		rtt, segSize, protocol, segmentFileName = http.GetFile(currentURL, baseURL+segURL, fileDownloadLocation, isByteRangeMPD, startRange, endRange, segmentNumber, segmentDuration, true, quicBool, glob.DebugFile, debugLog, useTestbedBool, repRate)
	case glob.MeanAverageAlg:
		rtt, segSize, protocol, segmentFileName = http.GetFile(currentURL, baseURL+segURL, fileDownloadLocation, isByteRangeMPD, startRange, endRange, segmentNumber, segmentDuration, true, quicBool, glob.DebugFile, debugLog, useTestbedBool, repRate)
	case glob.GeomAverageAlg:
		rtt, segSize, protocol, segmentFileName = http.GetFile(currentURL, baseURL+segURL, fileDownloadLocation, isByteRangeMPD, startRange, endRange, segmentNumber, segmentDuration, true, quicBool, glob.DebugFile, debugLog, useTestbedBool, repRate)
	case glob.EMWAAverageAlg:
		rtt, segSize, protocol, segmentFileName = http.GetFile(currentURL, baseURL+segURL, fileDownloadLocation, isByteRangeMPD, startRange, endRange, segmentNumber, segmentDuration, true, quicBool, glob.DebugFile, debugLog, useTestbedBool, repRate)
	case glob.TestAlg:
		rtt, segSize, protocol, segmentFileName = http.GetFile(currentURL, baseURL+segURL, fileDownloadLocation, isByteRangeMPD, startRange, endRange, segmentNumber, segmentDuration, true, quicBool, glob.DebugFile, debugLog, useTestbedBool, repRate)
	case glob.ArbiterAlg:
		rtt, segSize, protocol, segmentFileName = http.GetFile(currentURL, baseURL+segURL, fileDownloadLocation, isByteRangeMPD, startRange, endRange, segmentNumber, segmentDuration, true, quicBool, glob.DebugFile, debugLog, useTestbedBool, repRate)
	}

	// arrival and delivery times for this segment
	arrivalTime = int(time.Since(startTime).Nanoseconds() / (glob.Conversion1000 * glob.Conversion1000))
	deliveryTime := int(time.Since(currentTime).Nanoseconds() / (glob.Conversion1000 * glob.Conversion1000)) //Time in milliseconds
	thisRunTimeVal := int(time.Since(nextRunTime).Nanoseconds() / (glob.Conversion1000 * glob.Conversion1000))

	nextRunTime = time.Now()

	// some times we want to wait for an initial number of segments before stream begins
	// no need to do asny printouts when we are replacing this chunk
	// && !hlsReplaced
	if initBuffer <= waitToPlayCounter {

		// get the segment less the initial buffer
		// this needs to be based on running time and not based on number segments
		// I'll need a function for this
		//playoutSegmentNumber := segmentNumber - initBuffer

		// only print this out if we are not hls replaced
		if !hlsUsed {
			// print out the content of the segment that is currently passed to the player
			logging.PrintPlayOutLog(arrivalTime, initBuffer, mapSegmentLogPrintout, glob.LogDownload, printLog, printHeadersData)
		}

		// get the current buffer (excluding the current segment)
		currentBuffer := (bufferLevel - thisRunTimeVal)

		// if we have a buffer level then we have no stalls
		if currentBuffer >= 0 {
			stallTime = 0

			// if the buffer is empty, then we need to calculate
		} else {
			stallTime = currentBuffer
		}

		// To have the bufferLevel we take the max between the remaining buffer and 0, we add the duration of the segment we downloaded
		bufferLevel = utils.Max(bufferLevel-thisRunTimeVal, 0) + (segmentDuration * glob.Conversion1000)

		// increment the waitToPlayCounter
		waitToPlayCounter++

	} else {
		// add to the current buffer before we start to play
		bufferLevel += (segmentDuration * glob.Conversion1000)
		// increment the waitToPlayCounter
		waitToPlayCounter++
	}

	// check if the buffer level is higher than the max buffer
	if bufferLevel > maxBuffer*glob.Conversion1000 {
		// retrieve the time it is going to sleep from the buffer level
		// sleep until the max buffer level is reached
		sleepTime := bufferLevel - (maxBuffer * glob.Conversion1000)
		// sleep
		time.Sleep(time.Duration(sleepTime) * time.Millisecond)

		// reset the buffer to the new value less sleep time - should equal maxBuffer
		bufferLevel -= sleepTime
	}

	// some times we want to wait for an initial number of segments before stream begins
	// if we are going to print out some additonal log headers, then get these values
	if extendPrintLog && initBuffer < waitToPlayCounter {
		// base the play out position on the buffer level
		playPosition = segmentDurationTotal + (segmentDuration * glob.Conversion1000) - bufferLevel
		// we need to keep a tab on the different size segments - use this for now
		segmentDurationTotal += (segmentDuration * glob.Conversion1000)
	} else {
		segmentDurationTotal += (segmentDuration * glob.Conversion1000)
	}

	// if we are going to print out some additonal log headers, then get these values
	if extendPrintLog {

		// get the current codec
		repCodec = mpdList[mpdListIndex].Periods[0].AdaptationSet[currentMPDRepAdaptSet].Representation[repRate].Codecs

		// change the codec into something we can understand
		switch {
		case strings.Contains(repCodec, "avc"):
			// set the inital rep_rate to the lowest value
			repCodec = glob.RepRateCodecAVC
		case strings.Contains(repCodec, "hev"):
			repCodec = glob.RepRateCodecHEVC
		case strings.Contains(repCodec, "vp"):
			repCodec = glob.RepRateCodecVP9
		case strings.Contains(repCodec, "av1"):
			repCodec = glob.RepRateCodecAV1
		}

		// get rep_rate height, width and frames per second
		repHeight = mpdList[mpdListIndex].Periods[0].AdaptationSet[currentMPDRepAdaptSet].Representation[repRate].Height
		repWidth = mpdList[mpdListIndex].Periods[0].AdaptationSet[currentMPDRepAdaptSet].Representation[repRate].Width
		repFps = mpdList[mpdListIndex].Periods[0].AdaptationSet[currentMPDRepAdaptSet].Representation[repRate].FrameRate
	}

	// calculate the throughtput (we get the segSize while downloading the file)
	// multiple segSize by 8 to get bits and not bytes
	thr := algo.CalculateThroughtput(segSize*8, deliveryTime)

	// Print to output log
	//printLog(strconv.Itoa(segmentNumber), strconv.Itoa(arrivalTime), strconv.Itoa(deliveryTime), strconv.Itoa(Abs(stallTime)), strconv.Itoa(bandwithList[repRate]/1000), strconv.Itoa((segSize*8)/deliveryTime), strconv.Itoa((segSize*8)/(segmentDuration*1000)), strconv.Itoa(segSize), strconv.Itoa(bufferLevel), adapt, strconv.Itoa(segmentDuration*1000), extendPrintLog, repCodec, strconv.Itoa(repWidth), strconv.Itoa(repHeight), strconv.Itoa(repFps), strconv.Itoa(playPosition), strconv.FormatFloat(float64(rtt.Nanoseconds())/1000000, 'f', 3, 64), fileDownloadLocation)

	// store the current segment log output information in a map
	printInformation := logging.SegPrintLogInformation{
		ArrivalTime:          arrivalTime,
		DeliveryTime:         deliveryTime,
		StallTime:            stallTime,
		Bandwidth:            bandwithList[repRate],
		DelRate:              thr,
		ActRate:              (segSize * 8) / (segmentDuration * glob.Conversion1000),
		SegSize:              segSize,
		BufferLevel:          bufferLevel,
		Adapt:                adapt,
		SegmentDuration:      segmentDuration,
		ExtendPrintLog:       extendPrintLog,
		RepCodec:             repCodec,
		RepWidth:             repWidth,
		RepHeight:            repHeight,
		RepFps:               repFps,
		PlayStartPosition:    segmentDurationTotal,
		PlaybackTime:         playPosition,
		Rtt:                  float64(rtt.Nanoseconds()) / (glob.Conversion1000 * glob.Conversion1000),
		FileDownloadLocation: fileDownloadLocation,
		RepIndex:             repRate,
		MpdIndex:             mpdListIndex,
		AdaptIndex:           currentMPDRepAdaptSet,
		SegmentIndex:         nextSegmentNumber,
		SegReplace:           hlsReplaced,
		Played:               false,
		HTTPprotocol:         protocol,
		SegmentFileName:      segmentFileName,
	}

	// this saves per segment number so from 1 on, and not 0 on
	// remember this :)
	mapSegmentLogPrintout[segmentNumber] = printInformation

	// to calculate throughtput and select the repRate from it (in algorithm.go)
	switch adapt {
	//Conventional Algo
	case glob.ConventionalAlg:
		//fmt.Println("old: ", repRate)
		algo.Conventional(&thrList, thr, &repRate, bandwithList, lowestMPDrepRateIndex)
		//fmt.Println("new: ", repRate)
		//Harmonic Mean Algo
	case glob.ElasticAlg:
		//fmt.Println("old repRate index: ", repRate)
		//fmt.Println("old bandwithList[repRate]", bandwithList[repRate])
		algo.ElasticAlgo(&thrList, thr, deliveryTime, maxBuffer, &repRate, bandwithList, &staticAlgParameter, bufferLevel, kP, kI, lowestMPDrepRateIndex)
		//fmt.Println("new repRate index: ", repRate)
		//fmt.Println("new bandwithList[repRate]", bandwithList[repRate])
		//fmt.Println("elastic segmentNumber: ", segmentNumber)
		//fmt.Println("segURL: ", segURL)
	//Progressive Algo
	case glob.ProgressiveAlg:
		fmt.Println("old: ", repRate)
		algo.Conventional(&thrList, thr, &repRate, bandwithList, lowestMPDrepRateIndex)
		fmt.Println("new: ", repRate)
	//Logistic Algo
	case glob.LogisticAlg:
		fmt.Println("old: ", repRate)
		algo.Logistic(&thrList, thr, &repRate, bandwithList, bufferLevel,
			highestMPDrepRateIndex, lowestMPDrepRateIndex, glob.DebugFile, debugLog,
			maxBufferLevel)
		fmt.Println("new: ", repRate)
		logging.DebugPrint(glob.DebugFile, debugLog, "\nDEBUG: ", "reprate returned: "+strconv.Itoa(repRate))
	//Mean Average Algo
	case glob.MeanAverageAlg:
		//fmt.Println("old: ", repRate)
		algo.MeanAverageAlgo(&thrList, thr, &repRate, bandwithList, lowestMPDrepRateIndex)
		//fmt.Println("new: ", repRate)
	//Geometric Average Algo
	case glob.GeomAverageAlg:
		//fmt.Println("old: ", repRate)
		algo.GeomAverageAlgo(&thrList, thr, &repRate, bandwithList, lowestMPDrepRateIndex)
		//fmt.Println("new: ", repRate)
	//Exponential Average Algo
	case glob.EMWAAverageAlg:
		//fmt.Println("old: ", repRate)
		algo.EMWAAverageAlgo(&thrList, &repRate, exponentialRatio, 3, thr, bandwithList, lowestMPDrepRateIndex)

	case glob.ArbiterAlg:

		repRate = algo.CalculateSelectedIndexArbiter(thr, segmentDuration*1000, segmentNumber, maxBufferLevel,
			repRate, &thrList, streamDuration, mpdList[mpdListIndex], currentURL,
			currentMPDRepAdaptSet, segmentNumber, baseURL, debugLog, deliveryTime, bufferLevel,
			highestMPDrepRateIndex, lowestMPDrepRateIndex, bandwithList,
			segSize)
		//fmt.Println("new: ", repRate)
	case glob.TestAlg:
		//fmt.Println("segDur: ", segmentDuration*1000)
		//fmt.Println("index rate: ", repRate)
		//fmt.Println("baseURL: ", baseURL)
		//fmt.Println("downloadDurationLastSegment: ", deliveryTime)
		//fmt.Println("maxStreamDuration: ", streamDuration)
		//fmt.Println("bufferLevel: ", bufferLevel)
		//fmt.Println("")

		repRate = algo.CalculateSelectedIndexBba(thr, segmentDuration*1000, segmentNumber, maxBufferLevel,
			repRate, &thrList, streamDuration, mpdList[mpdListIndex], currentURL,
			currentMPDRepAdaptSet, segmentNumber, baseURL, debugLog, deliveryTime, bufferLevel,
			highestMPDrepRateIndex, lowestMPDrepRateIndex, bandwithList)

		//fmt.Println("")
	}
	logging.DebugPrint(glob.DebugFile, debugLog, "\nDEBUG: ", adapt+" has choosen rep_Rate "+strconv.Itoa(repRate)+" @ a rate of "+strconv.Itoa(bandwithList[repRate]/glob.Conversion1000))

	//Increase the segment number
	segmentNumber++

	// break out if we have downloaded all of our segments
	if segmentDurationTotal+(segmentDuration*glob.Conversion1000) > streamDuration {
		return segmentNumber, mapSegmentLogPrintout
	}

	// this gets the index for the next MPD and the segment number for the next chunk
	stopPlayer := false
	stopPlayer, oldMPDIndex, nextSegmentNumber = http.GetNextSegmentDuration(segmentDurationArray, segmentDuration*glob.Conversion1000, segmentDurationTotal, glob.DebugFile, debugLog, segmentDurationArray[mpdListIndex], streamDuration)

	//fmt.Println("streamLoop oldMPDIndex: ", stopPlayer)

	// stream the next chunk
	if !stopPlayer {
		segmentNumber, mapSegmentLogPrintout = streamLoop(segmentNumber, currentURL, initBuffer, maxBuffer, codecName, codec, urlString, urlInput,
			mpdList, adapt, maxHeight, isByteRangeMPD, startTime, nextRunTime, arrivalTime, oldMPDIndex, nextSegmentNumber, hls, hlsBool,
			mapSegmentLogPrintout, streamDuration, extendPrintLog, hlsUsed, bufferLevel, segmentDurationTotal, quic, quicBool, baseURL, debugLog)
	}

	return segmentNumber, mapSegmentLogPrintout

}

/*
 * Function  :
 *
 *
 *
 */

/*
   fmt.Println("test ID of each representation :")
   for a := 0; a < len(mpd.Periods); a++ {
     for i := 0; i < len(mpd.Periods[a].AdaptationSet); i++ {
       for j := 0; j < len(mpd.Periods[a].AdaptationSet[i].SegmentTemplate); j++ {
         fmt.Println(mpd.Periods[a].AdaptationSet[i].SegmentTemplate[j].Initialization)
       }
     }
   }

   if mpdList != nil {
     for l := 0; l < len(mpdList); l++ {
       //printTest(&structList[l])
     }
   }
*/
