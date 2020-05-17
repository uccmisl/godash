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

package logging

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/uccmisl/godash/utils"
)

/* --------------------------- Parsing and reading json config file ------------------------------------ */

// Config : Struct for reading content from the config file in json
type Config struct {
	URL            string  `json:"url"`
	Adapt          string  `json:"adapt"`
	Codec          string  `json:"codec"`
	Debug          string  `json:"debug"`
	InitBuffer     int     `json:"initBuffer"`
	MaxBuffer      int     `json:"maxBuffer"`
	MaxHeight      int     `json:"maxHeight"`
	StreamDuration int     `json:"streamDuration"`
	OutputFolder   string  `json:"outputFolder"`
	StoreDash      string  `json:"storeDash"`
	TerminalPrint  string  `json:"terminalPrint"`
	HLS            string  `json:"hls"`
	GetHeaders     string  `json:"getHeaders"`
	ExpRatio       float64 `json:"expRatio"`
	Quic           string  `json:"quic"`
	PrintHeader    string  `json:"printHeader"`
	UseTestbed     string  `json:"useTestbed"`
	QoE            string  `json:"QoE"`
	LogFile        string  `json:"logFile"`
}

// Configure : extract all parameter values from the input config file
func Configure(file string, debugFile string, debugLog bool) (urls string, adapt string, codec string, maxHeight int, streamDuration int, maxBuffer int, initBuffer int, hLS string, outputFolder string, storeDash string, getHeader string, debug string, terminalPrint string, quic string, expRatio float64, printHeader string, useTestbed string, qoe string, configLogFile string) {

	// unmarshal the json file
	config := recupStructWithConfigFile(file, debugFile, debugLog)

	// get all of the URLs from the config file
	requestedURLs := recupURLsFromConfig(config)

	// get all of the variables from the config file
	adapt, codec, maxHeight, streamDuration, maxBuffer, initBuffer, hLS, outputFolder, storeDash, getHeader, debug, terminalPrint, quic, expRatio, printHeader, useTestbed, qoe, configLogFile = recupParameters(config)

	// get list of urls
	urls = string(strings.Join(requestedURLs, ","))

	return
}

// RecupParameters : extract all of the values from the config struct (excluding url)
func recupParameters(config Config) (adapt string, codec string, maxHeight int, streamDuration int, maxBuffer int, initBuffer int, hLS string, outputFolder string, storeDash string, getHeaders string, debug string, terminalPrint string, quic string, expRatio float64, printHeader string, useTestbed string, qoe string, configLogFile string) {

	// there is no need to test conmpatibility for any of these parameters as main.go tests will check for this

	// get all passed in parameters and return them
	adapt = config.Adapt
	codec = config.Codec
	maxHeight = config.MaxHeight
	streamDuration = config.StreamDuration
	maxBuffer = config.MaxBuffer
	initBuffer = config.InitBuffer
	outputFolder = config.OutputFolder
	storeDash = config.StoreDash
	hLS = config.HLS
	getHeaders = config.GetHeaders
	debug = config.Debug
	terminalPrint = config.TerminalPrint
	quic = config.Quic
	expRatio = config.ExpRatio
	printHeader = config.PrintHeader
	useTestbed = config.UseTestbed
	qoe = config.QoE
	configLogFile = config.LogFile

	return
}

// RecupStructWithConfigFile : take the input file and generate a config struct
func recupStructWithConfigFile(file string, debugFile string, debugLog bool) (config Config) {

	// check if file exists
	e, err := os.Stat(file)
	if err != nil {
		fmt.Println("config file not found")
		utils.StopApp()
	}
	DebugPrint(debugFile, debugLog, "DEBUG: ", "File opened : "+e.Name())

	// open the file
	f, _ := os.Open(file)
	defer f.Close()

	// read the file contents
	byteValue, _ := ioutil.ReadAll(f)

	// unmarshal the file's content and save to config variable
	json.Unmarshal(byteValue, &config)

	return
}

// RecupURLsFromConfig : get url from config struct
func recupURLsFromConfig(config Config) (URLs []string) {

	//recup the different urls from the config file
	urlString := strings.TrimLeft(config.URL, "[")
	urlString = strings.TrimRight(urlString, "]")
	URLs = strings.Split(urlString, ",")

	// save to slice
	for i := 0; i < len(URLs); i++ {
		URLs[i] = strings.TrimSpace(URLs[i])
	}

	// if no MPD urls passed in, stop the app
	if URLs == nil {
		fmt.Println("No urls has been found in the config file")
		// stop the app
		utils.StopApp()
	}

	return
}
