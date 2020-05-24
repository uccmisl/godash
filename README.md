# goDash Application - Pre-Release dev branch

# DO NOT USE!!!

Current release version : 2.0

We kindly ask that should you mention [goDASH](https://github.com/uccmisl/godash) or [goDASHbed](https://github.com/uccmisl/godashbed), or use our code, in your publication, that you would reference the following paper:

D. Raca, M. Manifacier, and J.J. Quinlan.  goDASH - GO accelerated HAS framework for rapid prototyping. 12th International Conference on Quality of Multimedia Experience (QoMEX), Athlone, Ireland. 26th to 28th May, 2020 [CORA](http://hdl.handle.net/10468/9845 "CORA") (To Appear)

## General Description

goDASH is an infrastructure for headless streaming of DASH video content, implemented in the language golang, an open-source programming language supported by Google.

goDASH is a highly dynamic application which provides options for:
- adaptation algorithms, such as conventional, elastic, progressive, logistic, average, bba, geometric, arbiter and exponential
- video codec, such as h264, h265, VP9 and AV1
- DASH profiles, such as full, main, live, full_byte_range and main_byte_range
- stream options for audio and video DASH content
- config file input
- ability to store the downloaded segments
- debug option for printing information for this video stream
- getting the header information for all segments of the MPD url
- defining the initial number of segments to download before stream starts
- defining the maximum stream buffer in seconds
- defining a maximum height resolution to stream
- printing log output to file/terminal columns based on selected print headers
- downloading the stream using the QUIC transport protocol
- defining a folder location within ../files/ to store the streamed DASH files
- utilising the goDASHbed testbed and internally setting up https certs
- log output from five QoE models: [P.1203](github.com/itu-p1203/itu-p1203.git), Yu, Yin, Claye and Duanmu
- collaborative framework for sharing DASH content between multiple clients using [consul](https://www.consul.io) and [gRPC](https://godoc.org/google.golang.org/grpc)

## Legacy
Version 2.0 of godash is a major write of the code, and versions of godash from version 2.0 onwards only work with versions of goDASHbed from version 2.0 onwards.  If you are using a  version 1 release of godash, please use a version 1 release of godash.  

--------------------------------------------------------

## Install Steps
The easiest way to install goDASH is to use the install script available at the UCC Mobile and Internet System Lab [MISL](http://cs1dev.ucc.ie/misl/goDASH/)

# Examples to launch the app :
```
./godash -url "[http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/2_se
c/x265/bbb/DASH_Files/main_byte_range/bbb_enc_x265_dash.mpd]" -adapt conventional -codec h265 -debug true -initBuffer 2 -maxBuffer 20 -maxHeight 1080 -streamDuration 10 -storeDASH 347985 -debug on -terminalPrint on
```
or use the pre-defined configure file (advised option):
```
./godash -config ./config/configure.json
```
By setting "getHeaders" to "on", you can download all of the per segment transmission costs for the provided MPD url.  This information is needed by some algorithms to maximum video quality.  This file is stored in "logs", and can be used at any time by the requesting algorithms.

--------------------------------------------------------

## Requirements - if install script not used
Install Google [GO](https://golang.org/dl/):

Clone or download this repository.  Depending on where you save goDASH, you may have to change your GOPATH.

(you can check your goPath by using `go env $GOPATH` )

Install [consul](https://www.consul.io) and follow their install instructions.

In Windows :
Open the control panel, go to "System and Security", then "System", "advanced settings", "environment var" and add a variable called GOPATH with a value of "path/to/goDash/DashApp" and a GOBIN with a value "path/to/godash/../bin".

In linux :
```
export GOPATH=/home/path/to/godash
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```
or add these commands to you ~/.profile (remove "export" as this is not needed in .profile)

## Build Instructions
In a terminal :
```
cd godash
```
Update all repositories and dependencies, and build the player using:
```
go build
```

The best option to run goDASH is to use the configure.json file
```
./godash -config ./config/configure.json
```

--------------------------------------------------------
To output the P.1203 QoE values, you will need to install the P.1203 GitHub repository
```
git clone github.com/itu-p1203/itu-p1203.git
```

Then follow the install instruction for P.1203.

Make sure that once P.1203 has been installed that you run P.1203 before using goDASH, as you will need to accept their code.
```
python3 -m itu_p1203 examples/mode0.json
```

--------------------------------------------------------
## Example DASH content
Video only MPD example:
```
http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/4_sec/x264/bbb/DASH_Files/full/bbb_enc_x264_dash.mpd
```

Audio only MPD example:
```
http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/4_sec/x264/bbb/DASH_Files/full/dash_audio.mpd
```

Audio and Video MPD example:
```
http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/4_sec/x264/bbb/DASH_Files/full/dash_video_audio.mpd
```

--------------------------------------------------------

## Print help about parameters:
```
./godash -help
```
Flags for goDASH:
```
  -adapt string :  
    	DASH algorithms - "conventional|elastic|progressive|logistic|average|geometric|exponential|arbiter|bba"
        (default "conventional")

  -codec string :  
    	video codec to use - used when accessing multi-codec MPD files
        "[h264|h265|VP9|AV1]" (default "h264")

  -config string :  
    	config file for this video stream - "[path/to/config/file]"
        values in the config file have precedence over all parameters passed via command line

  -debug string :  
    	set debug information for this video stream - "[on|off]" (default "off")

  -expRatio float :  
    	download the stream with exponential parameter:
        ratio - this only works with only a select few algorithms

  -getHeaders string :  
    	get the header information for all segments across all of the MPD urls - based on:  
        "[off|on|online|offline]"
        off: do not get headers,
        on: get all headers defined by MPD,
        online: get headers from webserver based on algorithm input
        offline: get headers from header file based on algorithm input (file created by "on")
        (default "off").
        If getHeaders is set to "on", the client will download the headers and then stop the client.  

  -initBuffer int :  
    	initial number of segments to download before stream starts
        (default 2)

  -logFile string
        Location to store the debug logs (default "./logs/log_file.txt")

  -maxBuffer int :  
    	maximum stream buffer in seconds (default 30)

  -maxHeight int :  
    	maximum height resolution to stream - defaults to maximum resolution height in MPD file (default 2160)

  -outputFolder string :  
	    folder location within ./files/ to store the streamed DASH files
        if no folder is passed, output defaults to "./files" folder

  -printHeader string :  
    	print columns based on selected print headers:

  -quic string :  
    	download the stream using the QUIC transport protocol
        "[on|off]" (default "off")

  -serveraddr string
        implement Collaborative framework for streaming clients - "[on|off]" (default "off")

  -storeDASH string :  
    	store the streamed DASH, and associated files
        "[on|off]" (default "off")

  -streamDuration int :  
    	number of seconds to stream
        defaults to maximum stream duration in MPD file

  -terminalPrint string :  
    	extend the output logs to provide additional information
        "[on|off]" (default "off")

  -url string :  
    	a list of urls specifying the location of the video clip MPD(s) files
        "[url,url]"

  -useTestbed string :  
    	setup https certs and use goDASHbed testbed
        "[on|off]" (default "off")

  -QoE string :  
    	print per segment QoE values (P1203 mode 0, Claye, Duanmu, Yin, Yu) - "[on|off]" (default "off")

  -help or -h :  
	    Print help screen
```
--------------------------------------------------------

# Evaluate Folder:

The evaluate folder offers a means of running multiple goDASH clients during one streaming session, either natively or in the goDASHbed framework
```
python3 ./test_goDASH.py --numClients=1 --terminalPrint="off" --debug="off"  --collaborative="off"
```
```
--numClients - defines the number of goDASH clients to stream
--terminalPrint - determines if the clients should output their logs to the terminal screen
--debug - defines if the debug logs should be created - note: even if "debug" is set to "off", a log file, "logDownload.txt", containing the output features of each downloaded segment will be created per client.
 --collaborative - determines if we should implement sharing of content through a collaborative framework.  These uses consul and gRPC to share dash content between clients.  Setting this is "on", mandates 'storeDash' will be set to 'on'
```
The evaluate folder contains a number of sub-folders:
"config" - contains the original configure.json file for these goDASH clients.  The "terminalPrint" and "debug" setting passed into the script will overwrite the respective "terminalPrint" and "debug" settings in this config file.
"urls" - contains a list of the possible urls to choose from the five profiles of the AVC and HEVC UHD DASH datasets, provided at [DATASETS](https://www.ucc.ie/en/misl/research/datasets/ivid_uhd_dataset/)

Once "test_goDASH.py" is run, new folder content is created within the "output" folder
For each run, the "output" folder will contain a new folder defined by a time stamp

Within this folder, e.g.: "2020-04-09-06-42-20", 3 folders will be created:
"config" - contains a newly generated config file for each client (numbered per client).  The url for each client, will be randomly chosen from the list of MPDs contained within the "urls" folder.  
"files" - contains a folder per client, within which, is each downloaded segment and the requested MPD file.  Each client folder will also contain a "logDownload.txt" file, which contains the per segment download information.
"logs" - will contain the debug logs if "debug" is set to on.  This folder will also contain the header information for all segments across all of the MPD urls, if "getHeaders" is set to on.
Note: if getHeaders is set to "on", the headers will be downloaded to the log folder, then the script will auto-run, to re-call the clients with the requested algorithm.

test_goDASH.py has been tested with up to 50 goDASH clients with no loss in output log content
