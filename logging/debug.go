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

package logging

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	glob "github.com/uccmisl/godash/global"
	"github.com/uccmisl/godash/utils"
)

// SegPrintLogInformation per segment map of print log output
type SegPrintLogInformation struct {
	ArrivalTime int
	// delivery time of file requested
	DeliveryTime    int
	StallTime       int
	Bandwidth       int
	DelRate         int
	ActRate         int
	SegSize         int
	P1203HeaderSize float64
	// buffer = difference in arr_times for adjacent segments + segment duration of this segment
	BufferLevel          int
	Adapt                string
	SegmentDuration      int
	ExtendPrintLog       bool
	RepCodec             string
	RepWidth             int
	RepHeight            int
	RepFps               int
	PlayStartPosition    int
	PlaybackTime         int
	Rtt                  float64
	FileDownloadLocation string
	RepIndex             int
	MpdIndex             int
	AdaptIndex           int
	SegmentIndex         int
	Played               bool
	SegReplace           string
	HTTPprotocol         string
	TTfb         		 float64
	TTlb         		 float64
	P1203                float64
	Clae                 float64
	Duanmu               float64
	Yin                  float64
	Yu                   float64
	P1203Kbps            float64
	SegmentFileName      string
	// QoE metrics
	SegmentRates   []float64
	SumSegRate     float64
	TotalStallDur  float64
	NumStalls      int
	NumSwitches    int
	RateDifference float64
	SumRateChange  float64
	RateChange     []float64
	MimeType       string
	Profile        string
}

// headers for the print log
const segNum = glob.SegNum
const arrTime = glob.ArrTime
const delTime = glob.DelTime
const stallDur = glob.StallDur
const repLevel = glob.RepLevel
const delRate = glob.DelRate
const actRate = glob.ActRate
const byteSize = glob.ByteSize
const buffLevel = glob.BuffLevel
const algoHeader = glob.AlgoHeader
const segDurHeader = glob.SegDurHeader
const codecHeader = glob.CodecHeader
const heightHeader = glob.HeightHeader
const widthHeader = glob.WidthHeader
const fpsHeader = glob.FpsHeader
const playHeader = glob.PlayHeader
const rttHeader = glob.RttHeader
const segReplaceHeader = glob.SegReplaceHeader
const httpProtocolHeader = glob.HTTPProtocolHeader
const ttfbHeader = glob.TTFBHeader
const ttlbHeader = glob.TTLBHeader

// QOE
const p1203Header = glob.P1203Header
const claeHeader = glob.ClaeHeader
const duanmuHeader = glob.DuanmuHeader
const yinHeader = glob.YinHeader
const yuHeader = glob.YuHeader

// DebugPrint :
// * fileLocation string - pass in fileLocation
// * printLog bool - pass in boolean to print log
// * inputPrefix string - define the prefix to use in the log file
// * inputString string - string to print to log
// * print to the debug log file
func DebugPrint(fileLocation string, printLog bool, inputPrefix string, inputString string) {

	// only print if the debug log boolean was set to true
	if printLog {
		// open the log file
		f, err := os.OpenFile(fileLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			// print an error
			log.Println(err)
			// print the flag help output
			flag.Usage()
			// exit the application
			os.Exit(3)
		}

		// create a logger, set to Debug
		logger := log.New(f, inputPrefix, log.Ldate|log.Ltime)
		// print the log string - add the Epoch time in seconds
		logger.Println("-", time.Now().Unix(), "-", inputString)

		// close the file
		f.Close()
	}
}

// DebugPrintfIntArray :
// * print an int array to the logFile
func DebugPrintfIntArray(fileLocation string, printLog bool, inputPrefix string, inputString string, arguement []int) {

	// only print if the debug log boolean was set to true
	if printLog {
		// open the log file
		f, err := os.OpenFile(fileLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			// print an error
			log.Println(err)
			// print the flag help output
			flag.Usage()
			// exit the application
			os.Exit(3)
		}

		// create a logger, set to Debug
		logger := log.New(f, inputPrefix, log.Ldate|log.Ltime)
		// print the log string
		logger.Printf(inputString, arguement)

		// close the file
		f.Close()
	}
}

// DebugPrintfStringArray :
// * print a string array to the logFile
func DebugPrintfStringArray(fileLocation string, printLog bool, inputPrefix string, inputString string, arguement []string) {

	// only print if the debug log boolean was set to true
	if printLog {
		// open the log file
		f, err := os.OpenFile(fileLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			// print an error
			log.Println(err)
			// print the flag help output
			flag.Usage()
			// exit the application
			os.Exit(3)
		}

		// create a logger, set to Debug
		logger := log.New(f, inputPrefix, log.Ldate|log.Ltime)
		// print the log string
		logger.Printf(inputString, arguement)

		// close the file
		f.Close()
	}
}

// PrintsegInformationLogMap :
// * print the elements of mapSegments to the logFile
func PrintsegInformationLogMap(debugFile string, debugLog bool, mapSegments map[int]SegPrintLogInformation) {
	// print to debug
	DebugPrint(debugFile, debugLog, "\n", "segments map :")

	// for _, mapSegments := range newMapSegment {

	// print map header
	mainPrintString := "%7s  %10s  %8s  %12s  %8s  %12s  %8s  %8s  %10s"
	extendPrintString := "  %12s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s%s\n"
	PrintToFile("seg_Num", "size", "downTime", "thr", "duration", "playbackTime", "repIndex", "MPDIndex", "adaptIndex", "bandwith", "", 
	true, "", "", "", "", "", "", mainPrintString, extendPrintString, debugFile, "", "", "", "", "", "", "", "", "")

	for k := 1; k <= len(mapSegments); k++ {
		// print out each segment map
		PrintToFile(strconv.Itoa(k), strconv.Itoa(mapSegments[k].SegSize), strconv.Itoa(mapSegments[k].DeliveryTime), 
		strconv.Itoa(mapSegments[k].DelRate), strconv.Itoa(mapSegments[k].SegmentDuration*glob.Conversion1000), 
		strconv.Itoa(mapSegments[k].PlaybackTime), strconv.Itoa(mapSegments[k].RepIndex), strconv.Itoa(mapSegments[k].MpdIndex), 
		strconv.Itoa(mapSegments[k].AdaptIndex), strconv.Itoa(mapSegments[k].Bandwidth), "", true, "", "", "", "", "", "", mainPrintString, 
		extendPrintString, debugFile, "", "", "", "", "", "", "", "", "")
	}
	// }
}

// PrintToFile :
// * print a line to the file logDownload
func PrintToFile(segNum string, arrTime string, delTime string, stallDur string,
	repLevel string, delRate string, actRate string, byteSize string,
	buffLevel string, algo string, segDuration string, extendPrintLog bool, codec string, 
	width string, height string, fps string, playHeader string, rttHeader string, 
	mainPrintString string, extendPrintString string, fileLocation string, segReplace string, 
	httpProtocol string, TTFB string, TTLB string, p1203 string, clae string, duanmu string, 
	yin string, yu string) {

	// open the logfile and print to it
	f, err := os.OpenFile(fileLocation, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("error here?")
		fmt.Println(err)
		return
	}

	// print to file
	fmt.Fprintf(f, mainPrintString, segNum, arrTime,
		delTime, stallDur, repLevel, delRate, actRate, byteSize, buffLevel)

	if extendPrintLog {
		//fmt.Fprint(f, algo+"\t"+segDuration+"\t"+codec+"\t"+height+"\t"+width+"\t"+fps+"\t"+playHeader+"\t"+rttHeader+"\t\n")
		fmt.Fprintf(f, extendPrintString, algo, segDuration, codec, width, height, fps, playHeader, rttHeader, 
			segReplace, httpProtocol, TTFB, TTLB, p1203, clae, duanmu, yin, yu)
	} else {
		fmt.Fprint(f, "\n")
	}

	defer f.Close()
}

// PrintHeaders :
// * print headers to the output log
// * create a logFile of the print output
func PrintHeaders(extendPrintLog bool, fileLocation string, logDownload string, debugFile string, debugLog bool, printLog bool, printHeadersData map[string]string) {

	// create the log file
	f, err := os.Create(fileLocation + "/" + logDownload)
	if err != nil {
		DebugPrint(debugFile, debugLog, "DEBUG: ", "can't create the file logDownload.txt in files")
	}
	defer f.Close()

	// print a line of the log file to terminal
	PrintLog(segNum, arrTime, delTime, stallDur, repLevel, delRate, actRate,
		byteSize, buffLevel, algoHeader, segDurHeader, extendPrintLog, codecHeader, heightHeader, 
		widthHeader, fpsHeader, playHeader, rttHeader, fileLocation, logDownload, printLog, 
		printHeadersData, segReplaceHeader, httpProtocolHeader, ttfbHeader, ttlbHeader, p1203Header, 
		claeHeader, duanmuHeader, yinHeader, yuHeader)
}

// PrintLog :
// * print a line to the output log
func PrintLog(segNum string, arrTime string, delTime string, stallDur string,
	repLevel string, delRate string, actRate string, byteSize string,
	buffLevel string, algoIn string, segDurationIn string, extendPrintLog bool, 
	codecIn string, widthIn string, heightIn string, fpsIn string, playIn string, 
	rttIn string, fileLocation string, logDownload string, printLog bool, 
	printHeadersData map[string]string, segReplaceIn string, httpProtocolIn string, ttfbIn string, ttlbIn string, 
	p1203In string, claeIn string, duanmuIn string, yinIn string, yuIn string) {

	const mainPrintString = "%10s   %10s   %10s   %10s   %10s   %10s   %10s   %10s   %10s"
	// update this string for new logs columns
	const fileExtendPrintString = "   %12s   %7s   %10s   %5s   %6s   %5s   %8s   %8s   %8s   %8s   %8s   %8s   %8s   %8s   %8s   %12s   %12s\n"
	var extendPrintString = ""
	const fiveString = "   %5s"
	const eightString = "   %8s"
	const twelveString = "   %12s"
	var algo = ""
	var segDuration = ""
	var codec = ""
	var width = ""
	var height = ""
	var fps = ""
	var play = ""
	var rtt = ""
	var segReplace = ""
	var httpProtocol = ""
	var ttfb = ""
	var ttlb = ""
	var p1203 = ""
	var clae = ""
	var duanmu = ""
	var yin = ""
	var yu = ""

	//"   %12s   %7s   %5s   %5s   %6s   %5s   %8s   %8s\n"
	//"Algorithm\":\"off\",\"Seg_Dur\":\"on\",\"Codec\":\"on\",\"Width\":\"on\",\"Height\":\"on\",\"FPS\":\"on\",\"Play_Pos\":\"on\",\"RTT\"

	if printLog {
		fmt.Printf(mainPrintString, segNum, arrTime,
			delTime, stallDur, repLevel, delRate, actRate, byteSize, buffLevel)

		if extendPrintLog {
			// these must be in the same order as print to log
			checkInputHeader(printHeadersData, algoHeader, &extendPrintString, "   %12s", &algo, algoIn)
			checkInputHeader(printHeadersData, segDurHeader, &extendPrintString, "   %7s", &segDuration, segDurationIn)
			checkInputHeader(printHeadersData, codecHeader, &extendPrintString, twelveString, &codec, codecIn)
			checkInputHeader(printHeadersData, widthHeader, &extendPrintString, fiveString, &width, widthIn)
			checkInputHeader(printHeadersData, heightHeader, &extendPrintString, "   %6s", &height, heightIn)
			checkInputHeader(printHeadersData, fpsHeader, &extendPrintString, fiveString, &fps, fpsIn)
			checkInputHeader(printHeadersData, playHeader, &extendPrintString, eightString, &play, playIn)
			checkInputHeader(printHeadersData, rttHeader, &extendPrintString, eightString, &rtt, rttIn)
			checkInputHeader(printHeadersData, segReplaceHeader, &extendPrintString, eightString, &segReplace, segReplaceIn)
			checkInputHeader(printHeadersData, httpProtocolHeader, &extendPrintString, eightString, &httpProtocol, httpProtocolIn)
			checkInputHeader(printHeadersData, ttfbHeader, &extendPrintString, eightString, &ttfb, ttfbIn)
			checkInputHeader(printHeadersData, ttlbHeader, &extendPrintString, eightString, &ttlb, ttlbIn)
			checkInputHeader(printHeadersData, p1203Header, &extendPrintString, eightString, &p1203, p1203In)
			checkInputHeader(printHeadersData, claeHeader, &extendPrintString, eightString, &clae, claeIn)
			checkInputHeader(printHeadersData, duanmuHeader, &extendPrintString, twelveString, &duanmu, duanmuIn)
			checkInputHeader(printHeadersData, yinHeader, &extendPrintString, twelveString, &yin, yinIn)
			checkInputHeader(printHeadersData, yuHeader, &extendPrintString, twelveString, &yu, yuIn)

			// one of these has to be true, so print a new line at the end
			extendPrintString += "\n"
			fmt.Printf(extendPrintString, algo, segDuration, codec, width, height, fps, play, rtt, segReplace, httpProtocol, ttfb, ttlb, p1203, clae, duanmu, yin, yu)
		} else {
			fmt.Printf("\n")
		}
	}

	printLocal := fileLocation + "/" + logDownload

	PrintToFile(segNum, arrTime, delTime, stallDur, repLevel, delRate, actRate, byteSize, buffLevel, algoIn, 
		segDurationIn, extendPrintLog, codecIn, widthIn, heightIn, fpsIn, playIn, rttIn, mainPrintString, 
		fileExtendPrintString, printLocal, segReplaceIn, httpProtocolIn, ttfbIn, ttlbIn, p1203In, claeIn, 
		duanmuIn, yinIn, yuIn)
}

//
func checkInputHeader(printHeadersData map[string]string, key string, extendPrintString *string, stringDuration string, val1 *string, val2 string) {

	/*
		fmt.Println("printHeadersData: ", printHeadersData)
		fmt.Println("key: ", key)
		fmt.Println("extendPrintString: ", *extendPrintString)
		fmt.Println("stringDuration: ", stringDuration)
		fmt.Println("val1: ", *val1)
		fmt.Println("val2: ", val2)
	*/

	// if the map has this key
	if val, ok := printHeadersData[key]; ok {
		if val == "on" || val == "On" || val == "ON" {
			*extendPrintString += stringDuration
			*val1 = val2
		} else {
			*extendPrintString += "%s"
		}
		// include this incase someone removes the flags from printHeaders
	} else {
		*extendPrintString += "%s"
	}
}

// PrintPlayOutLog :
// * print the play_out logs only when the current time is >= play_out time
func PrintPlayOutLog(currentTime int, initBuffer int, mapSegments []map[int]SegPrintLogInformation, logDownload string, printLog bool, printHeadersData map[string]string) {

	for playoutSegmentNumber := 1; playoutSegmentNumber <= len(mapSegments[0]); playoutSegmentNumber++ {

		for logIndex := range mapSegments {

			if currentTime >= (mapSegments[logIndex][playoutSegmentNumber-1].PlayStartPosition+mapSegments[logIndex][initBuffer].PlayStartPosition) && !mapSegments[logIndex][playoutSegmentNumber].Played {

				// print out the content of the segment that is currently passed to the player
				PrintLog(strconv.Itoa(playoutSegmentNumber),
					strconv.Itoa(mapSegments[logIndex][playoutSegmentNumber].ArrivalTime),
					strconv.Itoa(mapSegments[logIndex][playoutSegmentNumber].DeliveryTime),
					strconv.Itoa(utils.Abs(mapSegments[logIndex][playoutSegmentNumber].StallTime)),
					strconv.Itoa(mapSegments[logIndex][playoutSegmentNumber].Bandwidth/glob.Conversion1000),
					strconv.Itoa(mapSegments[logIndex][playoutSegmentNumber].DelRate/glob.Conversion1000),
					strconv.Itoa(mapSegments[logIndex][playoutSegmentNumber].ActRate),
					strconv.Itoa(mapSegments[logIndex][playoutSegmentNumber].SegSize),
					strconv.Itoa(mapSegments[logIndex][playoutSegmentNumber].BufferLevel),
					mapSegments[logIndex][playoutSegmentNumber].Adapt,
					strconv.Itoa(mapSegments[logIndex][playoutSegmentNumber].SegmentDuration*glob.Conversion1000),
					mapSegments[logIndex][playoutSegmentNumber].ExtendPrintLog,
					mapSegments[logIndex][playoutSegmentNumber].RepCodec,
					strconv.Itoa(mapSegments[logIndex][playoutSegmentNumber].RepWidth),
					strconv.Itoa(mapSegments[logIndex][playoutSegmentNumber].RepHeight),
					strconv.Itoa(mapSegments[logIndex][playoutSegmentNumber].RepFps),
					// print out the value of the comulative segment size less the segment size of the first segment
					strconv.Itoa(mapSegments[logIndex][playoutSegmentNumber-1].PlayStartPosition),
					fmt.Sprintf("%.3f", mapSegments[logIndex][playoutSegmentNumber].Rtt),
					mapSegments[logIndex][playoutSegmentNumber].FileDownloadLocation,
					logDownload,
					printLog,
					printHeadersData,
					mapSegments[logIndex][playoutSegmentNumber].SegReplace,
					mapSegments[logIndex][playoutSegmentNumber].HTTPprotocol,
					fmt.Sprintf("%.3f", mapSegments[logIndex][playoutSegmentNumber].TTfb),
					fmt.Sprintf("%.3f", mapSegments[logIndex][playoutSegmentNumber].TTlb),
					// add the QoE model outputs
					fmt.Sprintf("%.3f", mapSegments[logIndex][playoutSegmentNumber].P1203),
					fmt.Sprintf("%.3f", mapSegments[logIndex][playoutSegmentNumber].Clae),
					fmt.Sprintf("%.3f", mapSegments[logIndex][playoutSegmentNumber].Duanmu),
					fmt.Sprintf("%.3f", mapSegments[logIndex][playoutSegmentNumber].Yin),
					fmt.Sprintf("%.3f", mapSegments[logIndex][playoutSegmentNumber].Yu))

				// update the played boolean to true
				localMap := mapSegments[logIndex][playoutSegmentNumber]
				localMap.Played = true
				mapSegments[logIndex][playoutSegmentNumber] = localMap
			}
		}
	}
}
