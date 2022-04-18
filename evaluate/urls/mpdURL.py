#!/usr/bin/python
# /*
#  *	goDASH, golang client emulator for DASH video streaming
#  *	Copyright (c) 2022, Jason Quinlan, University College Cork
#  *					        j.quinlan@cs.ucc.ie, 
#                           Darijo Raca, University of Sarajev, BiH
#                               draca@etf.unsa.ba, 
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

full_url_list = ["http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/4_sec/x264/bbb/DASH_Files/full/bbb_enc_x264_dash.mpd",
                 "http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/4_sec/x264/sintel/DASH_Files/full/sintel_enc_x264_dash.mpd", 
                 "http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/4_sec/x264/tearsofsteel/DASH_Files/full/tearsofsteel_enc_x264_dash.mpd", 
                 ]

main_url_list = ["http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/4_sec/x264/bbb/DASH_Files/main/bbb_enc_x264_dash.mpd",
                 "http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/4_sec/x264/sintel/DASH_Files/main/sintel_enc_x264_dash.mpd", 
                 "http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/4_sec/x264/tearsofsteel/DASH_Files/main/tearsofsteel_enc_x264_dash.mpd", 
                 ]

live_url_list = ["http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/4_sec/x264/bbb/DASH_Files/live/bbb_enc_x264_dash.mpd",
                 "http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/4_sec/x264/sintel/DASH_Files/live/sintel_enc_x264_dash.mpd", 
                 "http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/4_sec/x264/tearsofsteel/DASH_Files/live/tearsofsteel_enc_x264_dash.mpd",
                 ]

full_byte_range_url_list = []

main_byte_range_url_list = []

full_url_list_2 = [
    "http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/4_sec/x264/bbb/DASH_Files/full/bbb_enc_x264_dash.mpd"]

# [
# "http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/4_sec/x264/bbb/DASH_Files/full_byte_range/bbb_enc_x264_dash.mpd",
#  "http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/4_sec/x264/sintel/DASH_Files/full_byte_range/sintel_enc_x264_dash.mpd", 
#  "http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/4_sec/x264/tearsofsteel/DASH_Files/full_byte_range/tearsofsteel_enc_x264_dash.mpd",
# ]

# ["http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/4_sec/x264/bbb/DASH_Files/main_byte_range/bbb_enc_x264_dash.mpd",
# "http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/4_sec/x264/sintel/DASH_Files/main_byte_range/sintel_enc_x264_dash.mpd", 
# "http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/4_sec/x264/tearsofsteel/DASH_Files/main_byte_range/tearsofsteel_enc_x264_dash.mpd",
# # ] 