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

package http

import (
	"bufio"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"

	glob "github.com/uccmisl/godash/global"
	"github.com/uccmisl/godash/logging"
	"github.com/uccmisl/godash/utils"
)

// MPD structure
type MPD struct {
	XMLName xml.Name `xml:"MPD"`

	Xmlns                     string `xml:"xmlns,attr"`
	MinBufferTime             string `xml:"minBufferTime,attr"`
	MediaPresentationDuration string `xml:"mediaPresentationDuration,attr"`
	MaxSegmentDuration        string `xml:"maxSegmentDuration,attr"`
	Profiles                  string `xml:"profiles,attr"`

	Periods            []Period           `xml:"Period"`
	ProgramInformation ProgramInformation `xml:"ProgramInformation"`

	AvailabilityStartTime string `xml:"availabilityStartTime,attr"`
	ID                    string `xml:"id,attr"`
	MinimumUpdatePeriod   string `xml:"minimumUpdatePeriod,attr"`
	PublishTime           string `xml:"publishTime,attr"`
	TimeShiftBufferDepth  string `xml:"timeShiftBufferDepth,attr"`
	Type                  string `xml:"type,attr"`
	NS1schemaLocation     string `xml:"ns1:schemaLocation,attr"`
	BaseURL               string `xml:"BaseURL"`
}

// ProgramInformation in MPD
type ProgramInformation struct {
	XMLName            xml.Name `xml:"ProgramInformation"`
	MoreInformationURL string   `xml:"moreInformationURL,attr"`
	Title              string   `xml:"Title"`
}

// Period in MPD
type Period struct {
	XMLName       xml.Name        `xml:"Period"`
	Duration      string          `xml:"duration,attr"`
	AdaptationSet []AdaptationSet `xml:"AdaptationSet"`
	ID            string          `xml:"id,attr"`
	Start         string          `xml:"start,attr"`
}

// AdaptationSet in MPD
type AdaptationSet struct {
	XMLName            xml.Name `xml:"AdaptationSet"`
	SegmentAlignment   bool     `xml:"segmentAlignment,attr"`
	BitstreamSwitching bool     `xml:"bitstreamSwitching,attr"`
	MaxWidth           int      `xml:"maxWidth,attr"`
	MaxHeight          int      `xml:"maxHeight"`
	MaxFrameRate       int      `xml:"maxFrameRate"`

	Par string `xml:"par,attr"`

	Lang                      string                    `xml:"lang,attr"`
	BaseURL                   string                    `xml:"BaseURL"`
	Representation            []Representation          `xml:"Representation"`
	SegmentTemplate           []SegmentTemplate         `xml:"SegmentTemplate"`
	SegmentList               SegmentList               `xml:"SegmentList"`
	SubsegmentStartsWithSAP   int                       `xml:"subsegmentStartsWithSAP"`
	AudioChannelConfiguration AudioChannelConfiguration `xml:"AudioChannelConfiguration"`
	Role                      Role                      `xml:"Role"`
	ContentType               string                    `xml:"contentType,attr"`
	MimeType                  string                    `xml:"mimeType,attr"`
	StartWithSAP              int                       `xml:"startWithSAP,attr"`

	FrameRate int    `xml:"frameRate,attr"`
	Height    string `xml:"height,attr"`
	ScanType  string `xml:"scanType,attr"`
	Width     int    `xml:"width,attr"`
}

// Representation in MPD
type Representation struct {
	XMLName                   xml.Name                  `xml:"Representation"`
	ID                        string                    `xml:"id,attr"`
	MimType                   string                    `xml:"mimType,attr"`
	Codecs                    string                    `xml:"codecs,attr"`
	Width                     int                       `xml:"width,attr"`
	Height                    int                       `xml:"height,attr"`
	FrameRate                 int                       `xml:"frameRate,attr"`
	Sar                       string                    `xml:"sar,attr"`
	StartWithSap              int                       `xml:"startWithSap,attr"`
	BandWidth                 int                       `xml:"bandwidth,attr"`
	BaseURL                   string                    `xml:"BaseURL"`
	SegmentTemplate           SegmentTemplate           `xml:"SegmentTemplate"`
	SegmentList               SegmentList               `xml:"SegmentList"`
	SegmentBase               SegmentBase               `xml:"SegmentBase"`
	AudioSamplingRate         int                       `xml:"audioSamplingRate,attr"`
	AudioChannelConfiguration AudioChannelConfiguration `xml:"AudioChannelConfiguration"`
}

// SegmentTemplate in MPD
type SegmentTemplate struct {
	XMLName        xml.Name `xml:"SegmentTemplate"`
	Media          string   `xml:"media,attr"`
	Timescale      int      `xml:"timescale,attr"`
	StartNumber    int      `xml:"startNumber,attr"`
	Duration       int      `xml:"duration,attr"`
	Initialization string   `xml:"initialization,attr"`
}

// SegmentList in MPD
type SegmentList struct {
	XMLName            xml.Name       `xml:"SegmentList"`
	Timescale          int            `xml:"timescale,attr"`
	Duration           int            `xml:"duration,attr"`
	SegmentURL         []segmentURL   `xml:"SegmentURL"`
	SegmentInitization Initialization `xml:"Initialization"`
}

// AudioChannelConfiguration in MPD
type AudioChannelConfiguration struct {
	XMLName     xml.Name `xml:"AudioChannelConfiguration"`
	SchemeIDURI string   `xml:"schemeIdUri,attr"`
	Value       int      `xml:"value,attr"`
}

// SegmentBase in MPD
type SegmentBase struct {
	XMLName            xml.Name       `xml:"SegmentBase"`
	IndexRangeExact    string         `xml:"indexRangeExact,attr"`
	IndexRange         string         `xml:"indexRange,attr"`
	SegmentInitization Initialization `xml:"Initialization"`
}

// Role in MPD
type Role struct {
	XMLName     xml.Name `xml:"Role"`
	SchemeIDURI string   `xml:"schemeIdUri,attr"`
	Value       string   `xml:"value,attr"`
}

// Initialization in MPD
type Initialization struct {
	XMLName   xml.Name `xml:"Initialization"`
	SourceURL string   `xml:"sourceURL,attr"`
}

// segmentURL in MPD
type segmentURL struct {
	XMLName    xml.Name `xml:"SegmentURL"`
	MediaRange string   `xml:"mediaRange,attr"`
	IndexRange string   `xml:"indexRange,attr"`
}

// the current Codec
var mpdCodec string
var mpdCodecIndex int
var repRateCodec string

// SegHeadValues store the seg header maps
var SegHeadValues map[int]map[int][]int

// getStructList :
// * Take an array of string that might correspond to URLs
// * For each URL, call the method GET and parse the result with the function fileParser()
// * If the URL doesn't match, displays an error, then continue with the other strings
// * Add each structure
func getStructList(requestedURLs []string, debugFile string, debugLog bool, useTestbedBool bool, quicbool bool) (mpds []MPD) {

	// for each of the requested URLs
	for i := 0; i < len(requestedURLs); i++ {

		urls, _, _ := GetURL(requestedURLs[i], false, 0, 0, quicbool, debugFile, debugLog, useTestbedBool)

		// Call the fileParser in parser.go
		mpd := fileParser(urls)

		//Add the list of mpd structures to the list that will be returned
		mpds = append(mpds, mpd)

	}
	// return the MPD list
	return
}

/*
* Function fileParser :
*
* take an xml file
* Extract the file read
* returns an MPD struct
*
 */
func fileParser(mpdBody []byte) MPD {

	var mpd *MPD

	//extract everything from the file read in bytes to the structures
	xml.Unmarshal(mpdBody, &mpd)

	return *mpd
}

// func getSegmentSizes() {
//
// 	var mpd *MPD
//
// 	//extract everything from the file read in bytes to the structures
// 	xml.Unmarshal(mpdBody, &mpd)
//
// 	return *mpd
// }

// GetNextSegmentDuration :
// * returns an index for the MPD and the next segment we can use
// * currently randomised - to illustrate functionality
func GetNextSegmentDuration(segmentDurations []int, lastSegmentDuration int, totalSegmentDuration int, debugFile string, debugLog bool, segmentDuration int, streamDuration int) (stopApp bool, mpdIndex int, nextSegmentNumber int) {

	// variables
	var r int

	// lets switch segment number if the segment duration can actually be used
	// this gives back an array from the second segment onwards
	// we use the first segment from the first MPD as the first segment
	usableMPDindexes, usableSegmentNumbers := getPlayableSegmentMPDindex(segmentDurations, lastSegmentDuration, totalSegmentDuration, streamDuration)

	// fmt.Printf("MPD indexes - %v\n", usableMPDindexes)
	// fmt.Printf("Segment numbers - %v\n", usableSegmentNumbers)

	// this randomised the MPD index we can use - replace this with your function
	stopApp = false
	if len(usableMPDindexes) > 0 {
		r = rand.Int() % len(usableMPDindexes)
	} else {
		logging.DebugPrint(debugFile, debugLog, "DEBUG: ", "no more usable segments for this MPD - stopping App")
		// lets stop the app
		stopApp = true
		return
	}

	// this gives us the next MPD index
	mpdIndex = usableMPDindexes[r]

	// this gives us the segment number for the next chunk
	nextSegmentNumber = usableSegmentNumbers[r]

	// fmt.Println("randomised number - " + strconv.Itoa(r))
	// fmt.Println("the MPD index we will use - " + strconv.Itoa(mpdIndex))
	// fmt.Println("the next segment number - " + strconv.Itoa(nextSegmentNumber))
	// fmt.Println()

	logging.DebugPrint(debugFile, debugLog, "\nDEBUG: ", "Current segment duration: "+strconv.Itoa(segmentDuration*glob.Conversion1000))
	logging.DebugPrint(debugFile, debugLog, "DEBUG: ", "Total segment duration: "+strconv.Itoa(totalSegmentDuration))
	logging.DebugPrint(debugFile, debugLog, "DEBUG: ", "Next chunk duration: "+strconv.Itoa(segmentDurations[mpdIndex]))
	logging.DebugPrint(debugFile, debugLog, "DEBUG: ", "Next chunk index: "+strconv.Itoa(nextSegmentNumber))

	// this returns the index for the next MPD and next segment number for the MPD
	return

}

// getPlayableSegmentMPDindex :
// returns the arrays of indexes for the MPD we can use for the next segment
func getPlayableSegmentMPDindex(segmentDurations []int, lastSegmentDuration int, totalSegmentDuration int, streamDuration int) (usableSegmentDurations []int, usableSegmentNumbers []int) {
	// note I defined the arrays above, and don't need to inclue them in the return

	for i := 0; i < len(segmentDurations); i++ {
		// if the last segment or the total segment duration can be modulo, then use
		if (lastSegmentDuration%(segmentDurations[i]*glob.Conversion1000)) == 0 || (totalSegmentDuration%(segmentDurations[i]*glob.Conversion1000)) == 0 {

			if totalSegmentDuration+(segmentDurations[i]*glob.Conversion1000) <= streamDuration {

				// save an array of usable MPD indexes
				usableSegmentDurations = append(usableSegmentDurations, i)

				// we also need to know which segment number this would be
				usableSegmentNumbers = append(usableSegmentNumbers, (totalSegmentDuration/(segmentDurations[i]*glob.Conversion1000))+1)
			}
		}
	}
	// return the indexes for the MPD we can use
	return
}

// GetAllSegmentHeaders :
// get all segment headers for all MPD urls
func GetAllSegmentHeaders(mpdList []MPD, codecIndexList [][]int,
	maxHeight int,
	segmentNumber int, streamDuration int,
	isByteRangeMPD bool,
	maxBuffer int,
	headerURL string, codec string, urlInput []string, debugLog bool, printToFile bool) map[int]map[int][]int {

	// store the seg header maps
	var segHeadValues map[int]map[int][]int
	segHeadValues = make(map[int]map[int][]int)

	// loop over all of the passed in MPD files
	for mpdListIndex := 0; mpdListIndex < len(mpdList); mpdListIndex++ {

		// check if the codec is in the MPD urls passed in
		codecList, codecIndexList, _ := GetCodec(mpdList, codec, debugLog)
		// determine if the passed in codec is one of the codecs we use
		usedCodec, codecIndex := utils.FindInStringArray(codecList[mpdListIndex], codec)

		// check the codec and print error is false
		if !usedCodec {
			// print error message
			fmt.Printf("*** - " + codec + " is not in the provided MPD\n")
			// stop the app
			utils.StopApp()
		}
		// save the current MPD Rep_rate Adaptation Set
		currentMPDRepAdaptSet := codecIndexList[mpdListIndex][codecIndex]

		// get the current URL
		currentURL := strings.TrimSpace(urlInput[mpdListIndex])

		// get the segment headers for this MPD url
		segHeadValues[mpdListIndex] = getSegmentHeaders(mpdList, mpdListIndex, currentMPDRepAdaptSet, maxHeight, segmentNumber, streamDuration, isByteRangeMPD, maxBuffer, currentURL, headerURL, debugLog, printToFile)
	}
	return segHeadValues
}

// GetNSegmentHeaders :
// get N segment headers for all MPD urls (based on stream time)
func GetNSegmentHeaders(mpdList []MPD, codecIndexList [][]int,
	maxHeight int,
	segmentNumber int, streamDuration int,
	isByteRangeMPD bool,
	maxBuffer int,
	headerURL string, codec string, urlInput []string, debugLog bool, useHeaderFile bool) map[int]map[int][]int {

	// store the seg header maps
	var segHeadValues map[int]map[int][]int
	segHeadValues = make(map[int]map[int][]int)

	// loop over all of the passed in MPD files
	for mpdListIndex := 0; mpdListIndex < len(mpdList); mpdListIndex++ {

		// check if the codec is in the MPD urls passed in
		codecList, codecIndexList, _ := GetCodec(mpdList, codec, debugLog)
		// determine if the passed in codec is one of the codecs we use
		usedCodec, codecIndex := utils.FindInStringArray(codecList[mpdListIndex], codec)

		// check the codec and print error is false
		if !usedCodec {
			// print error message
			fmt.Printf("*** - " + codec + " is not in the provided MPD\n")
			// stop the app
			utils.StopApp()
		}

		// save the current MPD Rep_rate Adaptation Set
		currentMPDRepAdaptSet := codecIndexList[mpdListIndex][codecIndex]

		// get the current URL
		currentURL := strings.TrimSpace(urlInput[mpdListIndex])

		// get the segment headers for this MPD url from a file or from the webserver
		if useHeaderFile {
			segHeadValues[mpdListIndex] = getNSegmentHeadersFromFile(mpdList, mpdListIndex, currentMPDRepAdaptSet, maxHeight, segmentNumber, streamDuration, isByteRangeMPD, maxBuffer, currentURL, headerURL, debugLog)
		} else {
			segHeadValues[mpdListIndex] = getSegmentHeaders(mpdList, mpdListIndex, currentMPDRepAdaptSet, maxHeight, segmentNumber, streamDuration, isByteRangeMPD, maxBuffer, currentURL, headerURL, debugLog, useHeaderFile)
		}
	}
	SegHeadValues = segHeadValues
	return segHeadValues
}

// getNSegmentHeadersFromFile :
/*
 * get the segment headers for a given mpd, based on mpdIndex and adaptationSet index
 * from a given file
 * we also only want to get the rep_rates between certain index (due to max height)
 * we can also pass in the segment number to start at and the number of segments to get
 * we need to know if the MPD is byte range
 * we also want the url header
 * we also pass a few values, we don't need but the functions do:
 * maxBuffer, currentURL
 */
func getNSegmentHeadersFromFile(mpdList []MPD, mpdListIndex int, currentMPDRepAdaptSet int,
	maxHeight int,
	segmentNumber int, streamDuration int,
	isByteRangeMPD bool,
	maxBuffer int, currentURL string,
	headerURL string, debugLog bool) map[int][]int {

	// file name
	var fileName string

	//Map [rate]listduration
	var contentLengthDictionary map[int][]int
	contentLengthDictionary = make(map[int][]int)

	// we need some info from the MPD file, so get these:
	maxStreamDuration, _, highestMPDrepRateIndex, lowestMPDrepRateIndex, segmentDurationArray, _, _ := GetMPDValues(mpdList, mpdListIndex, maxHeight, streamDuration, maxBuffer, currentMPDRepAdaptSet, isByteRangeMPD, debugLog)

	// now get the maximum number of segments
	maxSegments := maxStreamDuration / (segmentDurationArray[mpdListIndex] * glob.Conversion1000)

	// get current segment duration
	segmentDuration := segmentDurationArray[mpdListIndex]

	// remove the tail from the mpd file name
	mpdTitle := (strings.Split(headerURL, "."))[0]

	// get the profile from the MPD file
	profiles := strings.Split(mpdList[mpdListIndex].Profiles, ":")
	numProfile := len(profiles) - 2
	profile := profiles[numProfile]

	// create the output log file name
	// we need clip name, codec, profile and segment duration
	fileName = glob.DebugFolder + strconv.Itoa(segmentDuration) + "sec_" + mpdTitle
	// if byte-range add this
	if isByteRangeMPD {
		fileName += glob.ByteRangeString
	}
	// add the tail to the file
	fileName += "_" + profile + ".csv"

	// check if the file already exists
	_, err := os.Stat(fileName)
	if err != nil {
		logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "The segment header file for this MPD does not exist")
		fmt.Println("The MPD header file: " + fileName + " does not exist, please change the -" + glob.GetHeaderName + " flag to on")
		utils.StopApp()
	} else {
		logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "The segment header file for this MPD already exists")
	}

	// create the file with the fileName
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error when opening the file for segment lengths")
		utils.StopApp()
	}
	defer f.Close()

	reader := csv.NewReader(bufio.NewReader(f))

	//for {
	lines, error := reader.ReadAll()
	if error != nil {
		log.Fatal(error)
	}

	for rowIndex, columnIndex := range lines {
		//skip the first 2 rows - due to indexes we only need to add 1
		if rowIndex < segmentNumber+1 {
			// skip header lines
			continue
		}
		// we don't want any more segments beyond the max segment number
		if rowIndex > maxSegments+1 {
			continue
		}
		// add a counter for us to use
		var counter = 1
		// if the rep_rates use incremental indexs for lower rates
		/*
			if !repRatesReversed {
		*/
		for j := highestMPDrepRateIndex; j <= lowestMPDrepRateIndex; j++ {
			// trim any white space and convert to int
			i1, err := strconv.Atoi(strings.TrimSpace(columnIndex[counter]))
			if err == nil {
				contentLengthDictionary[j] = append(contentLengthDictionary[j], i1)
				counter++
			} else {
				fmt.Println(err)
			}
		}
		/*
			} else {
				for j := highestMPDrepRateIndex; j >= lowestMPDrepRateIndex; j-- {
					// trim any white space and convert to int
					i1, err := strconv.Atoi(strings.TrimSpace(columnIndex[counter]))
					if err == nil {
						contentLengthDictionary[j] = append(contentLengthDictionary[j], i1)
						counter++
					} else {
						fmt.Println(err)
					}
				}
			}
		*/
	}

	return contentLengthDictionary
}

// GetContentLengthHeader :
// get the header of the next segment to have the informations about it
func GetContentLengthHeader(currentMPD MPD, currentURL string, currentMPDRepAdaptSet int, repRate int, segmentNumber int, adaptationSetBaseURL string, debugLog bool) int {

	// get the base url
	baseURL := GetNextSegment(currentMPD, segmentNumber, repRate, currentMPDRepAdaptSet)

	// join the new file location to the base url
	url := JoinURL(currentURL, adaptationSetBaseURL+baseURL, debugLog)
	logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "get header file from URL: "+url+"\n")

	if url == "" {
		fmt.Println("null urlHeader")
	}

	//Get the header of the url
	resp, err := http.Head(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	contentLen, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		// fmt.Println("can't convert the content-length response to an int")
		logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "can't convert the content-length response to an int")
	}
	return contentLen
}

// getSegmentHeaders :
/*
 * get the segment headers for a given mpd, based on mpdIndex and adaptationSet index
 * we also only want to get the rep_rates between certain index (due to max height)
 * we can also pass in the segment number to start at and the number of segments to get
 * we need to know if the MPD is byte range
 * we also want the url header
 * we also pass a few values, we don't need but the functions do:
 * maxBuffer, currentURL
 */
func getSegmentHeaders(mpdList []MPD, mpdListIndex int, currentMPDRepAdaptSet int,
	maxHeight int,
	segmentNumber int, streamDuration int,
	isByteRangeMPD bool,
	maxBuffer int, currentURL string,
	headerURL string, debugLog bool, printToFile bool) map[int][]int {

	var fileName string

	// variable for the file
	var f *os.File

	//Map [rate]listduration
	var contentLengthDictionary map[int][]int
	contentLengthDictionary = make(map[int][]int)

	// we need some info from the MPD file, so get these:
	maxStreamDuration, _, highestMPDrepRateIndex, lowestMPDrepRateIndex, segmentDurationArray, _, baseURL := GetMPDValues(mpdList, mpdListIndex, maxHeight, streamDuration, maxBuffer, currentMPDRepAdaptSet, isByteRangeMPD, debugLog)

	// now get the maximum number of segments
	maxSegments := maxStreamDuration / (segmentDurationArray[mpdListIndex] * glob.Conversion1000)

	if printToFile {

		// get current segment duration
		segmentDuration := segmentDurationArray[mpdListIndex]

		// remove the tail from the mpd file name
		mpdTitle := (strings.Split(headerURL, "."))[0]

		// get the profile from the MPD file
		profiles := strings.Split(mpdList[mpdListIndex].Profiles, ":")
		numProfile := len(profiles) - 2
		profile := profiles[numProfile]

		// create the output log file name
		// we need clip name, codec, profile and segment duration
		fileName = glob.DebugFolder + strconv.Itoa(segmentDuration) + "sec_" + mpdTitle
		// if byte-range add this
		if isByteRangeMPD {
			fileName += glob.ByteRangeString
		}
		// add the tail to the file
		fileName += "_" + profile + ".csv"

		// check if the file already exists
		_, err := os.Stat(fileName)
		if err == nil && !printToFile {
			logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "The segment header file for this MPD already exists")
			fmt.Println("The MPD header file: " + fileName + " already exists, please change the -" + glob.GetHeaderName + " flag")
			utils.StopApp()
		}
		// create the file with the fileName
		f, err = os.Create(fileName)
		if err != nil {
			fmt.Println("Error when creating the file for segment lengths")
			utils.StopApp()
		}
		defer f.Close()

		/* ------------------------- PRINT SEGMENT VALUES TO FILE AND DICTIONARY -------------------------- */

		fmt.Fprintf(f, "       ,")
		//print the headers - resolutions width
		// if the rep_rates are not reversed
		/*
			if !repRatesReversed {
		*/
		for k := highestMPDrepRateIndex; k <= lowestMPDrepRateIndex; k++ {
			if k == lowestMPDrepRateIndex {
				fmt.Fprintf(f, "%8d", mpdList[mpdListIndex].Periods[0].AdaptationSet[currentMPDRepAdaptSet].Representation[k].Width)
			} else {
				fmt.Fprintf(f, "%8d,   ", mpdList[mpdListIndex].Periods[0].AdaptationSet[currentMPDRepAdaptSet].Representation[k].Width)
			}
		}
		/*
			} else {
				for k := highestMPDrepRateIndex; k >= lowestMPDrepRateIndex; k-- {
					if k == lowestMPDrepRateIndex {
						fmt.Fprintf(f, "%8d", mpdList[mpdListIndex].Periods[0].AdaptationSet[currentMPDRepAdaptSet].Representation[k].Width)
					} else {
						fmt.Fprintf(f, "%8d,   ", mpdList[mpdListIndex].Periods[0].AdaptationSet[currentMPDRepAdaptSet].Representation[k].Width)
					}
				}
			}
		*/
		fmt.Fprintln(f, "")
		//print the headers - resolutions height
		fmt.Fprintf(f, "       ,")
		// if the rep_rates are not reversed
		/*
			if !repRatesReversed {
		*/
		for k := highestMPDrepRateIndex; k <= lowestMPDrepRateIndex; k++ {
			if k == lowestMPDrepRateIndex {
				fmt.Fprintf(f, "%8d", mpdList[mpdListIndex].Periods[0].AdaptationSet[currentMPDRepAdaptSet].Representation[k].Height)
			} else {
				fmt.Fprintf(f, "%8d,   ", mpdList[mpdListIndex].Periods[0].AdaptationSet[currentMPDRepAdaptSet].Representation[k].Height)
			}
		}
		/*
			} else {
				for k := highestMPDrepRateIndex; k >= lowestMPDrepRateIndex; k-- {
					if k == lowestMPDrepRateIndex {
						fmt.Fprintf(f, "%8d", mpdList[mpdListIndex].Periods[0].AdaptationSet[currentMPDRepAdaptSet].Representation[k].Height)
					} else {
						fmt.Fprintf(f, "%8d,   ", mpdList[mpdListIndex].Periods[0].AdaptationSet[currentMPDRepAdaptSet].Representation[k].Height)
					}
				}
			}
		*/
		fmt.Fprintln(f, "")
	}

	/* -------------------- SEGMENT VALUES ---------------------- */
	//Print the values for content length
	if isByteRangeMPD {
		// run from current segment to max segment
		for i := segmentNumber; i <= maxSegments; i++ {
			if printToFile {
				fmt.Fprintf(f, "%4d,   ", i)
			}
			for j := highestMPDrepRateIndex; j <= lowestMPDrepRateIndex; j++ {
				// get the byte range of the next segment that will be downloaded
				// we get this from the MPD struct
				_, startRange, endRange := GetNextByteRangeURL(mpdList[mpdListIndex], i, j, currentMPDRepAdaptSet)
				// now we have the ranges, just substract one from the other
				contentLength := endRange - startRange
				// save this value in a dictionary
				contentLengthDictionary[j] = append(contentLengthDictionary[j], contentLength)
				if printToFile {

					if j == lowestMPDrepRateIndex {
						fmt.Fprintf(f, "%8d   ", contentLength)
					} else {
						fmt.Fprintf(f, "%8d,   ", contentLength)
					}
				}
			}
			if printToFile {
				fmt.Fprintln(f, " ")
			}
		}

	} else {

		// adjust this to run from current segment to the maximum segment
		for i := segmentNumber; i <= maxSegments; i++ {
			//for i := 1; i < numSegments; i++ {
			if printToFile {
				fmt.Fprintf(f, "%4d,   ", i)
			}
			// if the rep_rates are not reversed
			/*
				if !repRatesReversed {
			*/
			for j := highestMPDrepRateIndex; j <= lowestMPDrepRateIndex; j++ {
				// get the content length of the next segment that will be downloaded
				contentLength := GetContentLengthHeader(mpdList[mpdListIndex], currentURL, currentMPDRepAdaptSet, j, i, baseURL, debugLog)
				// save this value in a dictionary
				contentLengthDictionary[j] = append(contentLengthDictionary[j], contentLength)
				if printToFile {
					if j == lowestMPDrepRateIndex {
						fmt.Fprintf(f, "%8d   ", contentLength)
					} else {
						fmt.Fprintf(f, "%8d,   ", contentLength)
					}
				}
			}
			/*
				} else {
					for j := highestMPDrepRateIndex; j >= lowestMPDrepRateIndex; j-- {
						// get the content length of the next segment that will be downloaded
						contentLength := GetContentLengthHeader(mpdList[mpdListIndex], currentURL, currentMPDRepAdaptSet, j, i, baseURL, debugLog)
						// save this value in a dictionary
						contentLengthDictionary[j] = append(contentLengthDictionary[j], contentLength)
						if printToFile {
							if j == lowestMPDrepRateIndex {
								fmt.Fprintf(f, "%8d   ", contentLength)
							} else {
								fmt.Fprintf(f, "%8d,   ", contentLength)
							}
						}
					}
				}
			*/
			if printToFile {
				fmt.Fprintln(f, " ")
			}
		}
	}
	return contentLengthDictionary
}

// GetCodec :
/*
 * for the list of passed in MPD urls
 * return an array of the codecs offered in the MPDs
 * return the index for the codec provided, -1 for all codecs
 */
func GetCodec(mpdList []MPD, codec string, debugLog bool) ([][]string, [][]int, bool) {

	var tempCodecList [][]string
	var tempIndexList [][]int
	var audioContent bool

	// fmt.Println(mpdList)

	// for a given set of representations
	for i := 0; i < len(mpdList); i++ {

		var codecList []string
		var codecIndexList []int

		// loop over adaptation sets
		for j := 0; j < len(mpdList[i].Periods[0].AdaptationSet); j++ {

			// check the current codec
			mpdCodec := mpdList[i].Periods[0].AdaptationSet[j].Representation[0].Codecs

			fmt.Println(mpdCodec)

			var repRateCodec string
			// save the codec in a name we know
			switch {
			case strings.Contains(mpdCodec, "avc"):
				repRateCodec = glob.RepRateCodecAVC
			case strings.Contains(mpdCodec, "hev"):
				repRateCodec = glob.RepRateCodecHEVC
			case strings.Contains(mpdCodec, "hvc1"):
				repRateCodec = glob.RepRateCodecHEVC
			case strings.Contains(mpdCodec, "vp"):
				repRateCodec = glob.RepRateCodecVP9
			case strings.Contains(mpdCodec, "av1"):
				repRateCodec = glob.RepRateCodecAV1
			case strings.Contains(mpdCodec, "mp4a"):
				repRateCodec = glob.RepRateCodecAudio
				audioContent = true
			case strings.Contains(mpdCodec, "ac-3"):
				repRateCodec = glob.RepRateCodecAudio
				audioContent = true
			default:
				repRateCodec = "Unknown"
			}

			// if the provided codec is the same as the current codec, save index and name
			if repRateCodec == codec || repRateCodec == "Audio/MP4" {
				codecList = append(codecList, repRateCodec)
				codecIndexList = append(codecIndexList, j)
			} else {
				codecList = append(codecList, repRateCodec)
				codecIndexList = append(codecIndexList, -1)
			}
		}

		logging.DebugPrintfIntArray(glob.DebugFile, debugLog, "DEBUG: ", "Codec Index List : %v for MPD "+strconv.Itoa(i+1), codecIndexList)
		logging.DebugPrintfStringArray(glob.DebugFile, debugLog, "DEBUG: ", "Codec List : %v for MPD "+strconv.Itoa(i+1), codecList)

		// save an array of array
		tempCodecList = append(tempCodecList, codecList)
		tempIndexList = append(tempIndexList, codecIndexList)

	}
	// return the codec and index arrays
	return tempCodecList, tempIndexList, audioContent
}

// GetMPDValues :
// get important values from the provided MPD
func GetMPDValues(mpd []MPD, mpdListIndex int, maxHeight int, streamDuration int, maxBuffer int, currentMPDRepAdaptSet int, isByteRangeMPD bool, debugLog bool) (int, int, int, int, []int, []int, string) {

	var maxStreamDuration int
	var segmentDurationArray []int
	var maxBufferLevel int
	var minMPDlistIndex int
	var maxMPDlistIndex int
	var bandwithList []int
	var baseURL = mpd[mpdListIndex].Periods[0].AdaptationSet[0].BaseURL
	var maxSegments int

	if isByteRangeMPD {
		// if this is a byte-range MPD, get byte range metrics
		maxSegments, segmentDurationArray = GetByteRangeSegmentDetails(mpd, mpdListIndex)
	} else {
		// if not, get standard profile metrics
		maxSegments, segmentDurationArray = GetSegmentDetails(mpd, mpdListIndex)
	}
	// if the numSegments is greater than zero, use it as maxSegments
	if streamDuration != 0 {
		maxStreamDuration = streamDuration
	} else {
		// get the segment duration of the last segment (typically larger than normal)
		lastSegmentDuration := SplitMPDSegmentDuration(mpd[mpdListIndex].MaxSegmentDuration)
		// current segment duration for the first MPS in the url list
		segmentDuration := segmentDurationArray[mpdListIndex]
		// get MPD stream duration in segments
		maxStreamDuration = segmentDuration*(maxSegments-1) + lastSegmentDuration
	}
	maxBufferLevel = maxBuffer
	minMPDlistIndex = GetMPDheightIndex(mpd[mpdListIndex], maxHeight, currentMPDRepAdaptSet, debugLog)
	maxMPDlistIndex = GetMaxListIndex(mpd[mpdListIndex], currentMPDRepAdaptSet)
	bandwithList = GetRepresentationBandwidth(mpd[mpdListIndex], currentMPDRepAdaptSet)

	// determine if the MPD rep_rates are highest to lowest or lowest to highest
	if bandwithList[minMPDlistIndex] > bandwithList[maxMPDlistIndex] {
		logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "Rep_Rate "+strconv.Itoa(minMPDlistIndex)+" @ "+strconv.Itoa(bandwithList[minMPDlistIndex])+" is bigger than Rep_Rate "+strconv.Itoa(maxMPDlistIndex)+" @ "+strconv.Itoa(bandwithList[maxMPDlistIndex]))
	} else {
		// rep_rates are lowest to highest
		// max index is reversed
		logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "Rep_Rate "+strconv.Itoa(maxMPDlistIndex)+" @ "+strconv.Itoa(bandwithList[maxMPDlistIndex])+" is smaller than Rep_Rate "+strconv.Itoa(minMPDlistIndex)+" @ "+strconv.Itoa(bandwithList[minMPDlistIndex]))
		fmt.Println("There is a problem with the indexes set for the representation rates in the MPD file, so stop")
		os.Exit(3)
	}

	logging.DebugPrint(glob.DebugFile, debugLog, "\nDEBUG: ", "Collect metrics from selected MPD file")
	logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "Maximum stream duration: "+strconv.Itoa(maxStreamDuration))
	logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "Maximum number of segments: "+strconv.Itoa(maxSegments))
	logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "Maximum buffer: "+strconv.Itoa(maxBufferLevel))
	logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "Minimum MPD index: "+strconv.Itoa(minMPDlistIndex))
	logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "Maximum MPD index: "+strconv.Itoa(maxMPDlistIndex))
	logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "Segment Duration: "+strconv.Itoa(segmentDurationArray[mpdListIndex]))
	logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "Base URL: "+baseURL)
	logging.DebugPrintfIntArray(glob.DebugFile, debugLog, "\nDEBUG: ", "\nBandwidth List : %v\n\n", bandwithList)

	// return the values
	return maxStreamDuration, maxBufferLevel, minMPDlistIndex, maxMPDlistIndex, segmentDurationArray, bandwithList, baseURL
}

// GetNextSegment :
/*
 * select the right segment in the MPD given
 * Return the URL of this segment
 */
func GetNextSegment(mpd MPD, SegNumber int, SegQUALITY int, currentMPDRepAdaptSet int) string {

	// the base url for this segment/rep_rate
	var repRateBaseURL string

	// get the base media url for a given representation rate
	// remember index's are one less than rep_rate value
	repRateBaseURL = mpd.Periods[0].AdaptationSet[currentMPDRepAdaptSet].Representation[(SegQUALITY)].SegmentTemplate.Media

	// convert the segment int to string
	nb := strconv.Itoa(SegNumber)

	//Replace $Number$ by the segment number in the url and return it
	return strings.Replace(repRateBaseURL, "$Number$", nb, -1)
}

// GetMPDheightIndex :
//get the maximum index for a given resolution height in a provided MPD file
func GetMPDheightIndex(mpd MPD, maxHeight int, currentMPDRepAdaptSet int, debugLog bool) int {

	// define the maximum height index
	var maxHeightIndex = 0
	// determine the maximum height index based on the bandwidth for a given index
	var bitrate = 0

	// for a given set of representations
	for j := 0; j < len(mpd.Periods[0].AdaptationSet[currentMPDRepAdaptSet].Representation); j++ {

		// if the representation height is the same as the passed in height
		if mpd.Periods[0].AdaptationSet[currentMPDRepAdaptSet].Representation[j].Height <= maxHeight {
			logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "getMPDheightIndex - maxHeight: "+strconv.Itoa(maxHeight))
			logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "getMPDheightIndex - bitrate: "+strconv.Itoa(mpd.Periods[0].AdaptationSet[currentMPDRepAdaptSet].Representation[j].BandWidth))

			// determine which rep index has the highest bitrate
			if mpd.Periods[0].AdaptationSet[currentMPDRepAdaptSet].Representation[j].BandWidth > bitrate {
				// save the bitrate (some MPD files have multiple rep_rates for a given height)
				bitrate = mpd.Periods[0].AdaptationSet[currentMPDRepAdaptSet].Representation[j].BandWidth
				// save maxHeightIndex
				//maxHeightIndex = mpd.Periods[0].AdaptationSet[currentMPDRepAdaptSet].Representation[j].ID
				maxHeightIndex = j + 1

				//logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "rep_rate for maxHeight: "+strconv.Itoa(mpd.Periods[0].AdaptationSet[currentMPDRepAdaptSet].Representation[j].ID-1))
				logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "rep_rate for maxHeight:"+strconv.Itoa(j))
			}
		}
	}

	// return the maximum index for a given resolution height
	return maxHeightIndex - 1
}

// GetMaxListIndex :
// get the maximum rep_rate from the MPD
func GetMaxListIndex(mpd MPD, currentMPDRepAdaptSet int) int {

	// return the lenght of the Representations
	return len(mpd.Periods[0].AdaptationSet[currentMPDRepAdaptSet].Representation) - 1
}

// GetRepresentationBandwidth :
// get Representation Bandwidth - divided by 1000
func GetRepresentationBandwidth(mpd MPD, currentMPDRepAdaptSet int) (bandwithList []int) {
	for i := 0; i < len(mpd.Periods[0].AdaptationSet[currentMPDRepAdaptSet].Representation); i++ {
		bandwithList = append(bandwithList, mpd.Periods[0].AdaptationSet[currentMPDRepAdaptSet].Representation[i].BandWidth)
	}
	return bandwithList
}

// GetByteRangeSegmentDetails :
// get segment duration and number of segments from provided MPD file
func GetByteRangeSegmentDetails(mpd []MPD, mpdListIndex int) (int, []int) {

	//var maxSegmentNumber = 0
	var streamDuration int
	var segmentDurations []int

	// split the MPD segment durations
	// lets now get the MPD files
	streamDuration = SplitMPDSegmentDuration(mpd[mpdListIndex].MediaPresentationDuration)

	// get an array of all the segment durations
	for i := 0; i < len(mpd); i++ {

		//  mpd.MaxSegmentDuration may not be the actual segment size (just the size of the last segment)
		//segmentDuration = splitMPDSegmentDuration(mpd.MaxSegmentDuration)
		duration := (mpd[i].Periods[0].AdaptationSet[0].Representation[0].SegmentList.Duration)
		timeScale := (mpd[i].Periods[0].AdaptationSet[0].Representation[0].SegmentList.Timescale)

		// get segment duration
		segmentDurations = append(segmentDurations, duration/timeScale)
	}

	// return the number of segments and segment duration
	return streamDuration / segmentDurations[mpdListIndex], segmentDurations
}

// GetSegmentDetails :
// get segment duration and number of segments from provided MPD file
func GetSegmentDetails(mpd []MPD, mpdListIndex int, adaptationSetIndex ...int) (int, []int) {

	// default value for adaptationSetIndex
	var adaptIndex = 0
	if len(adaptationSetIndex) > 0 {
		adaptIndex = adaptationSetIndex[0]
	}

	//var maxSegmentNumber = 0
	var streamDuration int
	var segmentDurations []int

	// split the MPD segment durations
	streamDuration = SplitMPDSegmentDuration(mpd[mpdListIndex].MediaPresentationDuration)

	// get an array of all the segment durations
	for i := 0; i < len(mpd); i++ {

		//  mpd.MaxSegmentDuration may not be the actual segment size (just the size of the last segment)
		//segmentDuration = splitMPDSegmentDuration(mpd.MaxSegmentDuration)
		duration := (mpd[i].Periods[0].AdaptationSet[adaptIndex].Representation[0].SegmentTemplate.Duration)
		timeScale := (mpd[i].Periods[0].AdaptationSet[adaptIndex].Representation[0].SegmentTemplate.Timescale)

		// this might be a different type of MPD
		if duration == 0 {
			duration = (mpd[i].Periods[0].AdaptationSet[adaptIndex].SegmentTemplate[0].Duration)
		}
		if timeScale == 0 {
			timeScale = (mpd[i].Periods[0].AdaptationSet[adaptIndex].SegmentTemplate[0].Timescale)
		}

		// this might be a byte-range, so return empty if timeScale is empty
		if timeScale == 0 {
			timeScale = 1
		}

		// get segment duration
		segmentDurations = append(segmentDurations, duration/timeScale)
	}

	// return the number of segments and segment duration
	return streamDuration / segmentDurations[mpdListIndex], segmentDurations
}

// SplitMPDSegmentDuration :
// get the per second details from the MPD segments
func SplitMPDSegmentDuration(mpdSegDuration string) int {

	var totalTimeinSeconds int
	var streamDuration string

	// lets first determine the length of the file
	// remove the "PT"
	streamDurationHMS := strings.Replace(mpdSegDuration, "PT", "", -1)

	// if streamDurationHMS contains hours
	if strings.Contains(streamDurationHMS, "H") {
		// get the hours
		H := strings.Split(streamDurationHMS, "H")
		streamDurationH := H[0]
		// if there are hours, convert to seconds
		i3, err := strconv.Atoi(streamDurationH)
		if err != nil {
			fmt.Println("*** Problem with converting segment hours to int ***")
			// stop the app
			utils.StopApp()
		}
		if i3 > 0 {
			totalTimeinSeconds = i3 * 60 * 60
		}
		streamDuration = H[1]
	} else {

		// remove the "PT0H"
		streamDuration = strings.Replace(mpdSegDuration, "PT0H", "", -1)
	}

	// split around the Minutes
	s := strings.Split(streamDuration, "M")

	// if there are minutes, convert to seconds
	i1, err := strconv.Atoi(s[0])
	if err != nil {
		fmt.Println("*** Problem with converting segment minutes to int ***")
		// stop the app
		utils.StopApp()
	}
	if i1 > 0 {
		totalTimeinSeconds = i1 * 60
	}

	// get the seconds and convert to int
	s = strings.Split(s[1], ".")
	i2, err := strconv.Atoi(s[0])
	if err != nil {
		fmt.Println("*** Problem with converting segment seconds to int ***")
	}

	// return the minutes and seconds (in seconds)
	return totalTimeinSeconds + i2
}

// URLList :
// turn the url string into a urlList
func URLList(urlString string) []string {

	// begin by removing the "[" and "]" at the end and begining of the url(s)
	urlString = strings.TrimRight(urlString, "]")
	urlString = strings.TrimLeft(urlString, "[")

	// lets now split the url(s) around the ","
	return strings.Split(urlString, ",")

}

// ReadURLArray :
/*
* Read the string of url parameters passed to the app
* split the urls to have a list
* call getStructList with the list to have the MPDs
* return a struct of MPDs
 */
func ReadURLArray(args string, debugLog bool, useTestbedBool bool, quicbool bool) (structList []MPD) {

	var requestedURLs []string

	//fmt.Println("URL array : ", args)
	// print to debug log
	logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "URL array passed : "+args)

	// lets now split the url(s) around the ","
	urlInput := URLList(args)

	// if more than one url is passed in, then stop and print error
	if len(urlInput) > 1 {
		fmt.Println("*** only one url can be passed to goDASH, please remove any additional URLs. Use -h for more info ***")
		fmt.Println(args)
		os.Exit(3)
	}

	for i := 0; i < len(urlInput); i++ {
		// remove any other white space
		urlInput[i] = strings.TrimSpace(urlInput[i])

		// add the url to the list
		requestedURLs = append(requestedURLs, urlInput[i])

		//fmt.Printf("Parameters: %s\n", urlInput[i])
		// print to debug log
		logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "Parameters: "+urlInput[i])
	}
	if len(requestedURLs) > 0 {
		// get the []struct of MPDs
		structList = getStructList(requestedURLs, glob.DebugFile, debugLog, useTestbedBool, quicbool)
	}

	return structList
}

// GetFullStreamHeader :
/*
 * get the header file for the current video clip
 * I've called this full in case the other profile have a different structure
 */
func GetFullStreamHeader(mpd MPD, isByteRangeMPD bool) string {

	// get the url base location for the header file
	if isByteRangeMPD {
		return mpd.Periods[0].AdaptationSet[0].SegmentList.SegmentInitization.SourceURL
	}
	return mpd.Periods[0].AdaptationSet[0].SegmentTemplate[0].Initialization
}

// GetNextByteRangeURL :
// Return the base URL, start and end range for the byte range MPD
func GetNextByteRangeURL(mpd MPD, SegNumber int, SegQUALITY int, currentMPDRepAdaptSet int) (string, int, int) {

	// get the base media url for a given representation rate
	// remember index's are one less than rep_rate value
	baseURL := mpd.Periods[0].AdaptationSet[currentMPDRepAdaptSet].Representation[(SegQUALITY)].BaseURL
	mediaRange := mpd.Periods[0].AdaptationSet[currentMPDRepAdaptSet].Representation[(SegQUALITY)].SegmentList.SegmentURL[SegNumber-1].MediaRange

	// get the start and end ranges
	startRange, endRange := splitByteRange(mediaRange)

	return baseURL, startRange, endRange
}

// splitByteRange :
// split and return the string range into start and end int values
func splitByteRange(byteRange string) (int, int) {

	// split the input string around the "-"
	s := strings.Split(byteRange, "-")

	// get the start range
	startRange, err := strconv.Atoi(s[0])
	if err != nil {
		fmt.Println("*** Problem with converting Byte Range to int ***")
		// stop the app
		utils.StopApp()
	}

	// get the endRange
	endRange, err := strconv.Atoi(s[1])
	if err != nil {
		fmt.Println("*** Problem with converting Byte Range to int ***")
		// stop the app
		utils.StopApp()
	}

	// return the byte ranges
	return startRange, endRange
}
