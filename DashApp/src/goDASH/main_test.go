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

package main

import (
	"goDASH/http"
	"testing"
)

//example test
func TestHello(t *testing.T) {
	expected := "coucou"
	actual := "coucou"
	if actual != expected {
		t.Error("Test failed")
	}
}

//---------------------------------------------------------------------
//TESTS FOR PARSING AND GETTING FILES

//test fileParser to see if it works with all MPDs (function called in getStructList and readURLArray)
func testFileParser(t *testing.T) {

	MPDs := http.ReadURLArray("[http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/2_sec/x265/bbb/DASH_Files/main_byte_range/bbb_enc_x265_dash.mpd, http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/4_sec/x265/bbb/DASH_Files/main_byte_range/bbb_enc_x265_dash.mpd, http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/6_sec/x265/bbb/DASH_Files/main_byte_range/bbb_enc_x265_dash.mpd, http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/8_sec/x265/bbb/DASH_Files/main_byte_range/bbb_enc_x265_dash.mpd, http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/10_sec/x265/bbb/DASH_Files/main_byte_range/bbb_enc_x265_dash.mpd]", debugLog, true, false)

	for i := 0; i < len(MPDs); i++ {
		for a := 0; a < len(MPDs[i].Periods); a++ {
			for k := 0; k < len(MPDs[i].Periods[0].AdaptationSet); k++ {
				if MPDs[i].Periods[0].AdaptationSet == nil {
					t.Error("Problem when parsing the MPD in func fileParser, the adaptationSet is nil")
				}
				for j := 0; j < len(MPDs[i].Periods[0].AdaptationSet[k].Representation); j++ {

					if MPDs[i].Periods[0].AdaptationSet[k].Representation[j].ID == 0 {
						t.Error("Problem when parsing the MPD in func fileParser, the ID of the representation is equal to 0")
					}
				}
			}
		}
	}
}

//test getFile (in functions.go) TO DO, PROBLEM
func TestGetFile(t *testing.T) {
	/*
		currentURL := "http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/2_sec/x265/bbb/DASH_Files/main_byte_range/"
		URL := " bbb_enc_x265_dash_init.mp4"
		fileDownloadLocation := "..\files\347985"
		isByteRangeMPD := false
		startRange := 0
		endRange := 0
		segmentNumber := 1
		segmentDuration := 2

		rtt := getFile(currentURL, URL, fileDownloadLocation, isByteRangeMPD, startRange, endRange, segmentNumber, segmentDuration, false)
		fmt.Println(rtt)
	*/
	//In player.go, in the streaming loop, when we are calling getFile with baseURL+segURL for here URL, baseURL is always nil

}

//---------------------------------------------------------------------------------

//tests of the main functions used in stream() in player.go
func TestStream(t *testing.T) {

	//read the string of urls passed to the application after -url
	structList := http.ReadURLArray("[http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/2_sec/x265/bbb/DASH_Files/main_byte_range/bbb_enc_x265_dash.mpd, http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/4_sec/x265/bbb/DASH_Files/main_byte_range/bbb_enc_x265_dash.mpd, http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/6_sec/x265/bbb/DASH_Files/main_byte_range/bbb_enc_x265_dash.mpd, http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/8_sec/x265/bbb/DASH_Files/main_byte_range/bbb_enc_x265_dash.mpd, http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/10_sec/x265/bbb/DASH_Files/main_byte_range/bbb_enc_x265_dash.mpd]", true, true, false)
	if structList == nil {
		t.Error("The url read should not be nil when we cass readURLArray()")
	}

	//test urlList (called in readURLArray)
	listURL := http.URLList("[http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/2_sec/x265/bbb/DASH_Files/main_byte_range/bbb_enc_x265_dash.mpd, http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/4_sec/x265/bbb/DASH_Files/main_byte_range/bbb_enc_x265_dash.mpd, http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/6_sec/x265/bbb/DASH_Files/main_byte_range/bbb_enc_x265_dash.mpd, http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/8_sec/x265/bbb/DASH_Files/main_byte_range/bbb_enc_x265_dash.mpd, http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/10_sec/x265/bbb/DASH_Files/main_byte_range/bbb_enc_x265_dash.mpd]")
	if len(listURL) != 5 {
		t.Error("URL List function in urlParser returned : ", len(listURL), " instead of 5")
	}

	//now we will try to get the codec with the structList
	codecList, codecIndexList := http.GetCodec(structList, "HEVC", debugLog)
	t.Log(codecList)
	t.Log(codecIndexList)
	if codecList[0][0] != "HEVC" {
		t.Log("Function getCodec(structList, codec)")
		t.Error("expected codecList element to be equal to HEVC but got ", codecList[0][0])
	}
	for i := 0; i < len(codecIndexList); i++ {
		if codecIndexList[i][0] == -1 {
			t.Log("Function getCodec(structList, codec)")
			t.Error("codecIndexList element should not be equal to -1")
		}
	}

}

// ------------------------------------------------------------------------------------------------

// ----------------------------- Test MeanAverage and SelectRepRateWithThroughtput ----------------
/*
func TestMeanAverage(t *testing.T) {

	thrList := []int{
		2843157, 6690325, 12242549, 13067956, 15247213, 20917735, 26063698, 27587342, 26106059, 23265265,
	}

	repRate := 2
	average := 0.0

	//test MeanAverage(thrList) float64
	algo.MeanAverage(thrList, &average)
	//check the value returned
	if average != 17403129.9 {
		t.Log("Function MeanAverage(thrList) float64")
		t.Error("average should be equal to 17403129.9 and we have", average)
	}

	bandwithList := []int{
		40276548, 25312752, 15193504, 4354160, 3894826, 3046114, 2386043, 1826811, 1089489, 767717, 576208, 390172, 247230,
	}

	//test SelectRepRateWithThroughtput(thr int, repRate int, bandwithList []int, repRatesReversed bool) int
	repRate = algo.SelectRepRateWithThroughtput(int(average), bandwithList, false, 12)
	if repRate != 2 {
		t.Error("Expected repRate = 2 but got reprate choosed: ", repRate)
	}

	//test MeanAverageAlgo(
	//		thrList []int, newThr int, repRate int, bandwithList []int, repRatesReversed bool) (int, []int)
	algo.MeanAverageAlgo(&thrList, 26106059, &repRate, bandwithList, false, 12)
	if repRate != 2 {
		t.Log("test MeanAverageAlgo")
		t.Error("Expected repRate = 2 but got reprate choosed: ", repRate)
	}
	if thrList[len(thrList)-1] != 26106059 {
		t.Log("test MeanAverageAlgo")
		t.Error("the last element of the thrList should be 26106059 but it is: ", thrList[len(thrList)-1])
	}

}

// ----------------------------- Test GeomAverage -------------------------------------------------
func TestGeomAverage(t *testing.T) {

	thrList := []int{
		2843157, 6690325, 12242549, 13067956, 15247213, 20917735, 26063698, 27587342, 26106059, 23265265,
	}

	bandwithList := []int{
		40276548, 25312752, 15193504, 4354160, 3894826, 3046114, 2386043, 1826811, 1089489, 767717, 576208, 390172, 247230,
	}

	repRate := 2
	average := 0.0

	//test func GeomAverage(thrList []int) float64
	algo.GeomAverage(thrList, &average)
	//chack the value returned
	if average != 14545303.609233964 {
		t.Log("test GeomAverage(thrList []int) float64")
		t.Error("geometric average should be equal to 14545303.609233964 and we have ", average)
	}

	//test func GeomAverageAlgo(thrList []int, newThr int, repRate int, bandwithList []int, repRatesReversed bool) (int, []int)
	algo.GeomAverageAlgo(&thrList, 26106059, &repRate, bandwithList, false, 12)
	if repRate != 2 {
		t.Log("test MeanAverageAlgo")
		t.Error("Expected repRate = 2 but got reprate choosed: ", repRate)
	}
	if thrList[len(thrList)-1] != 26106059 {
		t.Log("test MeanAverageAlgo")
		t.Error("the last element of the thrList should be 26106059 but it is: ", thrList[len(thrList)-1])
	}

}

// ----------------------------- Test EMWAAverage -------------------------------------------------
func TestEMWAAverage(t *testing.T) {

	var thrList []int
	average := 0.0

	//with a list taken after a run
	thrList = append(thrList, 642)

	//test func ExpAverage(thrList []int, ratio float64) float64
	algo.ExpAverage(&thrList, 0.4, 10, &average)

	if average != 642 {
		t.Error("exponential average should be 642 and it is equal to: ", average)
	}

	//Other tests
	thrList = append(thrList, 545)
	algo.ExpAverage(&thrList, 0.4, 10, &average)
	if int(average) != 581 {
		t.Error("exponential average should be 581 and it is equal to: ", average)
	}

	thrList = append(thrList, 629)
	algo.ExpAverage(&thrList, 0.4, 10, &average)
	if int(average) != 605 {
		t.Error("exponential average should be 605 and it is equal to: ", average)
	}

	thrList = append(thrList, 721)
	algo.ExpAverage(&thrList, 0.4, 10, &average)
	if int(average) != 658 {
		t.Error("exponential average should be 658 and it is equal to: ", average)
	}

	thrList = append(thrList, 494)
	algo.ExpAverage(&thrList, 0.4, 10, &average)
	if int(average) != 587 {
		t.Error("exponential average should be 587 and it is equal to: ", average)
	}

	thrList = append(thrList, 1066)
	algo.ExpAverage(&thrList, 0.4, 10, &average)
	if int(average) != 788 {
		t.Error("exponential average should be 788 and it is equal to: ", average)
	}

	thrList = append(thrList, 761)
	algo.ExpAverage(&thrList, 0.4, 10, &average)
	if int(average) != 776 {
		t.Error("exponential average should be 776 and it is equal to: ", average)
	}

	thrList = append(thrList, 674)
	algo.ExpAverage(&thrList, 0.4, 10, &average)
	if int(average) != 735 {
		t.Error("exponential average should be 735 and it is equal to: ", average)
	}

	thrList = append(thrList, 1107)
	algo.ExpAverage(&thrList, 0.4, 10, &average)
	if int(average) != 885 {
		t.Error("exponential average should be 885 and it is equal to: ", average)
	}
}

// ----------------------------- Test HarmonicMean -------------------------------------------------
func TestHarmonicAverage(t *testing.T) {

	thrList := []int{
		2843157, 6690325, 12242549, 13067956, 15247213, 20917735, 26063698, 27587342, 26106059, 23265265,
	}

	average = 0.0

	//test func HarmonicAverage(num int, thrList []int) float64
	algo.HarmonicAverage(5, thrList, &average)
	if average != 24544673.172274478 {
		t.Log("test HarmonicAverage(num int, thrList []int) float64")
		t.Error("harmonic average should be equal to 24544673.172274478 and we have ", average)
	}
}

// ----------------------------- Test Logistic -------------------------------------------------
func TestLogistic(t *testing.T) {

	thrList := []int{
		2843157, 6690325, 12242549, 13067956, 15247213, 20917735, 26063698, 27587342, 26106059, 23265265,
	}

	bandwithList := []int{
		40276548, 25312752, 15193504, 4354160, 3894826, 3046114, 2386043, 1826811, 1089489, 767717, 576208, 390172, 247230,
	}

	repRate := 2

	//PB: we should pass the maxBufferLevel to this function
	//test func LogisticFunction(lastRateIndex int, thrList []int, bufferLevel int, highestMPDrepRateIndex int,
	//lowestMPDrepRateIndex int, maxBufferLevel int, bandwithList []int, repRatesReversed bool) int
	algo.LogisticFunction(10, &repRate, &thrList, 4000, 13, 3, 5000, bandwithList, true)
	//t.Error(retVal)
	fmt.Println("repRate TestLogistic: ", retVal)

	//test CalculateSelectedIndex(thrList []int, newThr int, bandwithList []int, bufferLevel int) int
}
*/
