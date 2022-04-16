#!/usr/bin/python
# /*
#  *	goDASH, golang client emulator for DASH video streaming
#  *	Copyright (c) 2022, Jason Quinlan, Darijo Raca, University College Cork
#  *											[j.quinlan,d.raca]@cs.ucc.ie)
#  *                      Maëlle Manifacier, MISL Summer of Code 2019, UCC
#  *	This program is free software; you can redistribute it and/or
#  *	modify it under the terms of the GNU General Public License
#  *	as published by the Free Software Foundation; either version 2
#  *	of the License, or (at your option) any later version.
#  *
#  *	This program is distributed in the hope that it will be useful,
#  *	but WITHOUT ANY WARRANTY; without even the implied warranty of
#  *	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
#  *	GNU General Public License for more details.
#  *
#  *	You should have received a copy of the GNU General Public License
#  *	along with this program; if not, write to the Free Software
#  *	Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA
#  *	02110-1301, USA.
#  */

from helper_functions import *

godash_run_dict = {
    # other key options to be added

    # choice of algorithm per client
    "algo_choice" : ["conventional", "arbiter", "elastic", "progressive"],


    # **** These are other config settings, which typically are not changed ****
    "codec":"\"h264\"",
    "initBuffer":2,
    "maxBuffer":60,
    "maxHeight":3000,
    "streamDuration":10,
    "printHeader":"\"{\"Algorithm\":\"on\",\"Seg_Dur\":\"on\",\"Codec\":\"on\",\"Width\":\"on\",\"Height\":\"on\",\"FPS\":\"on\",\"Play_Pos\":\"on\",\"RTT\":\"on\",\"Seg_Repl\":\"on\",\"Protocol\":\"on\",\"TTFB\":\"on\",\"TTLB\":\"on\",\"P.1203\":\"on\",\"Clae\":\"on\",\"Duanmu\":\"on\",\"Yin\":\"on\",\"Yu\":\"on\"}\"",
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
