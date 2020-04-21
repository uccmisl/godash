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

package global

// Conversion1000 : divider for conversion from bit to kilobit, to megabit, etc
const Conversion1000 = 1000

// Conversion1024 : divider for conversion from bit to kilobit, to megabit, etc
const Conversion1024 = 1024

// DebugFileName : debug log name
var DebugFileName = "logFile"

// DebugFolder : debug log folder location
const DebugFolder = "../logs/"

// DebugTextFile : debug log file location
const DebugTextFile = "log_file"

// FileFormat : debug file format
const FileFormat = ".txt"

// DebugFile : debug log folder + file + FileFormat
var DebugFile = DebugFolder + DebugTextFile + FileFormat

// LogDownload : where to save the log download text
const LogDownload = "logDownload.txt"

// RepRateCodecAVC : AVC constants for our encoder
const RepRateCodecAVC = "h264"

// RepRateCodecHEVC : HEVC constants for our encoder
const RepRateCodecHEVC = "h265"

// RepRateCodecVP9 : VP9 constants for our encoder
const RepRateCodecVP9 = "VP9"

// RepRateCodecAV1 : AV1 constants for our encoder
const RepRateCodecAV1 = "AV1"

// ConventionalAlg : constants for our algorithms
const ConventionalAlg = "conventional"

// ProgressiveAlg : constants for our algorithms
const ProgressiveAlg = "progressive"

// ElasticAlg : constants for our algorithms
const ElasticAlg = "elastic"

// LogisticAlg : constants for our algorithms
const LogisticAlg = "logistic"

// MeanAverageAlg : constants for our algorithms
const MeanAverageAlg = "average"

// GeomAverageAlg : constants for our algorithms
const GeomAverageAlg = "geometric"

// EMWAAverageAlg : constants for our algorithms
const EMWAAverageAlg = "exponential"

// TestAlg : test constants for our algorithms
const TestAlg = "test"

//ArbiterAlg : constants for our algorithms
const ArbiterAlg = "arbiter"

// HlsOff : constants for HLS
const HlsOff = "off"

// HlsOn : constants for HLS
const HlsOn = "on"

// TrueBool : true string for booleans
const TrueBool = "true"

// FalseBool : false string for booleans
const FalseBool = "false"

// GetHeaderOff : constants for getHeader
const GetHeaderOff = "off"

// GetHeaderOn : constants for getHeader
const GetHeaderOn = "on"

// GetHeaderOnline : constants for getHeader
const GetHeaderOnline = "online"

// GetHeaderOffline : constants for getHeader
const GetHeaderOffline = "offline"

// URLName : parameter variables
const URLName = "url"

// ConfigName : parameter variables
const ConfigName = "config"

// DebugName : parameter variables
const DebugName = "debug"

// CodecName : parameter variables
const CodecName = "codec"

// MaxHeightName : parameter variables
const MaxHeightName = "maxHeight"

// NumSegmentsName : parameter variables
const NumSegmentsName = "numSegments"

// StreamDurationName : parameter variables
const StreamDurationName = "streamDuration"

// PrintHeaderName : parameter variables
const PrintHeaderName = "printHeader"

// MaxBufferName : parameter variables
const MaxBufferName = "maxBuffer"

// InitBufferName : parameter variables
const InitBufferName = "initBuffer"

// AdaptName : parameter variables
const AdaptName = "adapt"

// FileStoreName : parameter variables
const FileStoreName = "storeDASH"

// TerminalPrintName : parameter variables
const TerminalPrintName = "terminalPrint"

// HlsName : parameter variables
const HlsName = "hls"

// QuicName : parameter variables
const QuicName = "quic"

// AppName : parameter variables
const AppName = "goDASH"

// ExpRatioName : parameter variables
const ExpRatioName = "expRatio"

// GetHeaderName : print header variables
const GetHeaderName = "getHeaders"

// RepRateBaseURL : used for determining if byte range MPD
const RepRateBaseURL = ""

// ByteRangeString : string for byte_range
const ByteRangeString = "_byte_range"

// DebugOff : constants for debug
const DebugOff = "off"

// DebugOn : constants for debug
const DebugOn = "on"

// TerminalPrintOff : constants for print
const TerminalPrintOff = "off"

// TerminalPrintOn : constants for print
const TerminalPrintOn = "on"

// QuicOff : constants for quic
const QuicOff = "off"

// QuicOn : constants for Extend
const QuicOn = "on"

// UseTestBedName : parameter variables
const UseTestBedName = "useTestbed"

// UseTestBedOff : constants for useTest
const UseTestBedOff = "off"

// UseTestBedOn : constants for useTest
const UseTestBedOn = "on"

// HTTPcertLocation : location of the http cert
const HTTPcertLocation = "http/certs/cert.pem"

// HTTPkeyLocation : location of the http key
const HTTPkeyLocation = "http/certs/key.pem"

// InsecureSSL :  "Accept/Ignore all server SSL certificates"
const InsecureSSL = false

// Serv : port for the server
const Serv = "serverPort"

// Client : port for the "client"
const Client = "clientPort"
