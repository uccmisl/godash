# goDash Application

Current release version : 2.4.2 - updated for go 1.20+ (as of 22-02-2023) and associated security vulnerabilities in older versions of "golang.org/x"

We kindly ask that should you mention [goDASH](https://github.com/uccmisl/godash) or [goDASHbed](https://github.com/uccmisl/godashbed) or use our code in your publication, that you would reference the following papers:

D. Raca, M. Manifacier, and J.J. Quinlan.  goDASH - GO accelerated HAS framework for rapid prototyping. 12th International Conference on Quality of Multimedia Experience (QoMEX), Athlone, Ireland. 26th to 28th May, 2020 [CORA](http://hdl.handle.net/10468/9845 "CORA")

John O’Sullivan, D. Raca, and Jason J. Quinlan.  Demo Paper: godash 2.0 - The Next Evolution of HAS Evaluation. 21st IEEE International Symposium On A World Of Wireless, Mobile And Multimedia Networks (IEEE WoWMoM 2020), Cork, Ireland. August 31 to September 03, 2020 [CORA](https://cora.ucc.ie/handle/10468/10145 "CORA")

--------------------------------------------------------
## Docker Containers

With the release of version 2.4.0+, we are also releasing amd64 docker containers for both [goDASH](https://hub.docker.com/r/jjq52021/godash) or [goDASHbed](https://hub.docker.com/r/jjq52021/godashbed).

An arm64 version of [goDASH](https://hub.docker.com/r/jjq52021/godash_arm64) is also available.

In the coming weeks/months/years we will also release a network build script, so as to permit a full evaluation of DASH algorithms and associated TCP and QUIC transport protocols within a Docker test framework.

--------------------------------------------------------
## Operating System Compatibility

godash is NOT COMPATIBLE with Windows system. godash must be run on a Linux or MAC O/S.


### General Description

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
- determining "time to first byte" (TTFB) and "time to last byte" TTLB - logged in milliseconds
- update debug logs with Epoch timestamp

### Legacy
Version 2.0 of godash is a major write of the code, and versions of godash from version 2.0 onwards only work with versions of goDASHbed from version 2.0 onwards.  If you are using a  version 1 release of godash, please use a version 1 release of godash.  


### Install Steps
The easiest way to install goDASH is to use the install script available at the UCC Mobile and Internet System Lab [MISL](http://cs1dev.ucc.ie/misl/godash2.0/)

--------------------------------------------------------
# Examples to launch the app :
```
./godash -url "[http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/2_se
c/x265/bbb/DASH_Files/main_byte_range/bbb_enc_x265_dash.mpd]" -adapt conventional -codec h265 -initBuffer 2 -maxBuffer 20 -maxHeight 1080 -streamDuration 10 -storeDASH on -debug on -terminalPrint on -outputFolder "123456" -logFile "log_file_2"
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

In linux/mac :
```
export GOPATH=/home/path/to/godash
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```
or add these commands to you ~/.profile (remove "export" as this is not needed in .profile)

--------------------------------------------------------
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
If using collaborative, first set `-serveraddr` to `on` in the godash config file

Then run Consul in a separate terminal using the command :

>consul agent -dev

Then call single cooperative goDASH client using:

>./goDASH -config ../config/configure.json

Call three clients using the evaluate framework using (see below for more info on 'evaluate'):

```
python3 ./test_goDASH.py --numClients=3 --terminalPrint="off" --debug="off"  --collaborative="on"
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

As of release version : 2.4.1 - The evaluate folder has been completely rewritten - it now uses a single command and a settings file

The evaluate folder offers a means of running multiple goDASH clients during one streaming session, either natively or in the goDASHbed framework
```
python3 ./test_goDASH.py
```

Contents of the settings file:
```
godash_run_dict = {
    # other key options to be added

    # choice of algorithm per client
    "algo_choice" : ["conventional", "progressive", "elastic", "logistic"],


    # **** These are other config settings, which typically are not changed ****
    "codec":"\"h264\"",
    "initBuffer":2,
    "maxBuffer":60,
    "maxHeight":3000,
    "streamDuration":10,
    "printHeader":"\"{\"Algorithm\":\"on\",\"Seg_Dur\":\"on\",\"Codec\":\"on\",\"Width\":\"on\",\"Height\":\"on\",\"FPS\":\"on\",\"Play_Pos\":\"on\",\"RTT\":\"on\",\"Seg_Repl\":\"off\",\"Protocol\":\"on\",\"TTFB\":\"on\",\"TTLB\":\"on\",\"P.1203\":\"on\",\"Clae\":\"on\",\"Duanmu\":\"on\",\"Yin\":\"on\",\"Yu\":\"on\"}\"",
    "expRatio":0.2,
    "quic":"\"off\"",
    "useTestbed":"\"off\"",
    "QoE":"\"on\""
}

# print output of godash to the terminal screen
terminalPrint=False
terminalPrintval=bool_to_val(terminalPrint)

# print output of godash to the log file
debug=True
debugval=bool_to_val(debug)

# run network in collaborative mode
collaborative=False
collaborativeval=bool_to_val(collaborative)

# number of clients, is defined by the number of algorithm choices
numClients = len(godash_run_dict["algo_choice"])

# ouptut folder structure
# output
output_folder_name = "/output"

# - config
config_folder_name = "/config"

# - files
log_folder_name = "/files"

# - config file
config_file="/configure.json"

# use a single clip for all clients, or randomly choose a clip for each user
single_clip_choice=False
```

The evaluate folder contains a number of sub-folders:

-- "config" - contains the original configure.json file for these goDASH clients, as well as the settings.py and helper_functions.py files.  
-- "urls" - contains a list of the possible urls to choose from the five profiles of the AVC and HEVC UHD DASH datasets, provided at [DATASETS](https://www.ucc.ie/en/misl/research/datasets/ivid_uhd_dataset/)

Once "test_goDASH.py" is run, new folder content is created within the "output" folder
For each run, the "output" folder will contain a new folder defined by a time stamp

Within this folder, e.g.: "2022-04-09-06-42-20", 3 folders will be created:
"config" - contains a newly generated config file for each client (numbered per client).  The url for each client, will be randomly chosen from the list of MPDs contained within the "urls" folder.  
"files" - contains a folder per client, within which, is each downloaded segment and the requested MPD file.  Each client folder will also contain a "logDownload.txt" file, which contains the per segment download information.
"logs" - will contain the debug logs if "debug" is set to on.  This folder will also contain the header information for all segments across all of the MPD urls, if "getHeaders" is set to on.
Note: if getHeaders is set to "on", the headers will be downloaded to the log folder, then the script will auto-run, to re-call the clients with the requested algorithm.

test_goDASH.py has been tested with up to 50 goDASH clients with no loss in output log content
