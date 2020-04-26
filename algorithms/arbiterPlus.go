package algorithms

import (
	//"fmt"

	glob "../global"
	"../http"
	"../utils"
	//"math"
)

var DEFAULT_EXPONENT float64 = 0.4
var DEFAULT_MIN_BUFFER_FACTOR float64 = 0.75
var DEFAULT_MAX_BUFFER_FACTOR float64 = 1.15
var DEFAULT_HISTORIC_ESTIMATION_WINDOW int = 10
var DEFAULT_MAXIMUM_SWITCH = 2
var DEFAULT_PREDICTIVE_ESTIMATION_WINDOW = 5

var detectSuddenDrop = false
var exponent float64

var bufferScaling = true
var bufferingFactor float64 = 1
var minBufferFactor float64
var maxBufferFactor float64

var selectedIndex int

var historicEstimationWindow int

var netScaling = false
var netScalingFactor float64 = 1

var switchingControl = true

var actualRateQuality = true

var predictiveEstimationWindow int

var maximumSwitch int

// CalculateSelectedIndexArbiter :
/*
 * return the index of the segment which should be selected
 *
 */
func CalculateSelectedIndexArbiter(newThr int, lastDuration int, lastIndex int, maxBufferLevel int,
	lastRate int, thrList *[]int, mpdDuration int, currentMPD http.MPD, currentURL string,
	currentMPDRepAdaptSet int, segmentNumber int, baseURL string, debugLog bool, downloadTime int, bufferLevel int,
	highestMPDrepRateIndex int, lowestMPDrepRateIndex int, bandwithList []int,
	segmentSize int) int {

	//Does not work if repRatesReversed
	//the typical default buffer should be 60 seconds, however this is set in the config json files

	//add the last throughput to the list
	*thrList = append(*thrList, newThr)

	//samples are the throughput
	//chunks are the segments

	//need to check whether maxBufferLevel can be used for typicalBufferMS
	//check if in ms or s
	//maxbufferlevel is in seconds
	//bufferLevel is in ms
	//fmt.Println("buff", bufferLevel, maxBufferLevel*1000)
	//fmt.Println(float64(bufferLevel)/float64(maxBufferLevel*1000))
	bufferFullness := FloatMin(1.0, (float64(bufferLevel) / float64(maxBufferLevel*1000)))

	exponent = DEFAULT_EXPONENT

	historicEstimationWindow = DEFAULT_HISTORIC_ESTIMATION_WINDOW

	//if playeractivity.predictedValue < 0
	var exponentialAverageRate float64

	//exponentialAverageRate := SampleExponentialAverage(historicEstimationWindow, exponent, *thrList)

	ExpAverage(*thrList, exponent, historicEstimationWindow, &exponentialAverageRate)

	//fmt.Println(exponentialAverageRate)
	//fmt.Println("exponentialavgrate: ", exponentialAverageRate)

	//---------------------------------------------------------------------------------------
	//minBufferFactor and maxBufferFactor
	//normally the values for these variables would be passed in via the class constructor
	//For now the assumption is that the default values are used

	minBufferFactor = DEFAULT_MIN_BUFFER_FACTOR
	maxBufferFactor = DEFAULT_MAX_BUFFER_FACTOR

	if bufferScaling {
		bufferingFactor = minBufferFactor + (maxBufferFactor-minBufferFactor)*bufferFullness
	}

	targetRate := exponentialAverageRate * bufferingFactor

	targetIndex := SelectRepRateWithThroughtput(int(targetRate), bandwithList, lowestMPDrepRateIndex)

	if switchingControl {

		maximumSwitch = DEFAULT_MAXIMUM_SWITCH

		//NOTE lastIndex is lastrate here

		if lastRate-targetIndex == 1 && targetRate < (1.065-(float64((len(bandwithList))-targetIndex))*0.015)*float64(bandwithList[targetIndex]/glob.Conversion1000) {
			targetIndex++
			//fmt.Println("here")
		} else if lastRate-targetIndex > maximumSwitch {

			//fmt.Println("also here")
			targetIndex = utils.Max(lastRate-maximumSwitch, 0)
		}

	}
	//fmt.Println("targetIndex 2: ", targetIndex)

	predictiveEstimationWindow = DEFAULT_PREDICTIVE_ESTIMATION_WINDOW
	//segHeadValues := http.GetNSegmentHeaders(mpdList, codecIndexList, maxHeight, 1, streamDuration, isByteRangeMPD, maxBuffer, headerURL, codec, urlInput, debugLog, true)
	//fmt.Println("test", http.SegHeadValues)

	if actualRateQuality {
		videoChunks := mpdDuration / lastDuration
		//fmt.Println("vidchunks", videoChunks)

		videoWindow := utils.Min(videoChunks-lastIndex, predictiveEstimationWindow)

		//fmt.Println("videoWindow", videoWindow)
		if http.SegHeadValues == nil {
			for targetIndex < lowestMPDrepRateIndex && !SmartConvHelper(targetIndex, videoWindow, targetRate, currentMPD, currentURL, currentMPDRepAdaptSet, lastRate, segmentNumber, baseURL, debugLog, lastDuration) {
				targetIndex++
			}
		} else {

			for targetIndex < lowestMPDrepRateIndex && !SmartConvHelperFromFile(videoWindow, targetRate, targetIndex, segmentNumber-1, lastDuration) {
				targetIndex++
			}

		}

	}

	//fmt.Println("targetIndex 3: ", targetIndex)
	return targetIndex
}

//doesnt actually matter, never will be called
func suddenDropHappening(rateCV float64, chunkSizeRatio float64, bufferFullness float64, lastDuration int, throughputDecreasing bool) bool {

	return rateCV > 0.5 && throughputDecreasing && float64(lastDuration) > 2*chunkSizeRatio*float64(lastDuration) && bufferFullness < 0.25

	//if indeed lastChunkDuration is equal to lastDuration
}

//helper functions

//ended up not being used
/*
func sampleCV(window int, thrList []int) float64 {

  thrSamples := ThroughputSamples(window, thrList)
  //now calculate the coefficient of variation
  return coefficientOfVariation(thrSamples)
  //return throughputSamples
  //thrList *[]int
}

func coefficientOfVariation(throughputSamples []int) float64 {
  //get the arithemetic average
  average := arithemeticAverage(throughputSamples)
  //then get the arithmetic variance
  variance := arithmeticVariance(throughputSamples, average)

  return math.Sqrt(variance)/average
}

func arithemeticAverage(throughputSamples []int) float64 {
  var average float64
  for i:= 0; i < len(throughputSamples); i++ {
    average += float64(throughputSamples[i])
  }

  return average/float64(len(throughputSamples))
}


func arithmeticVariance(throughputSamples []int, inputAverage float64) float64 {
  var totalDeviation float64
  var result float64

  for i:=0; i < len(throughputSamples); i++ {
    val := throughputSamples[i]
    totalDeviation += math.Pow(inputAverage-float64(val), 2)
  }

  if len(throughputSamples) > 1 {
    result = totalDeviation/float64((len(throughputSamples)-1))
  }

  return result
}

func isThroughputDecreasing(thrList []int) bool {
  if len(thrList) < 2 {
    return false
  }
  last := thrList[len(thrList) -1]
  secondlast := thrList[len(thrList) -2]

  return last < secondlast
}

//Could be used to reverse the rep rates
func reverseList(List []int) []int {

  for i, j := 0, len(List) -1; i<j ; i, j = i+1, j-1 {
    List[i], List[j] = List[j], List[i]
  }

  return List
}
*/
