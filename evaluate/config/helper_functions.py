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

from time import sleep
import os

from urls.mpdURL import *

#  mpdURL import *

def create_dict(config_file:str) -> dict:
    '''
    read in a config file and return a dictionary of the file structure
    '''
    # lets read in the original config file and create a dictionary we can use
    dict = {}
    # open the original config file
    with open(config_file, encoding='utf-8-sig') as fp:
        # read line by line
        line = fp.readline().strip()
        while line:
            # do not split around the brackets
            if len(line) > 3:
                # split around the colon, but not the colon in http :)
                key, val = line.split(': ')
                # strip out the spaces and quotes
                key = key.strip()
                # remove spaces, commas and quotes
                val = str(val.strip().strip(",").strip("\""))
                dict[key] = val
                line = fp.readline().strip()
            else:
                # otherwise, just read the next line
                line = fp.readline().strip()

    # return the dictionary
    return dict


def check_collab_and_set_url(collab_bool_val) -> list[str]:
    '''
    determine which mpd urls to use for this run
    '''
    # if we are not collaborative
    if not collab_bool_val:
        # get all the possible DASH MPD files from the H264 UHD dataset
        return full_url_list+main_url_list+live_url_list + \
            full_byte_range_url_list+main_byte_range_url_list
    else:
        # lets start consul
        os.system("consul agent -dev &")
        # let's sleep until consul is set up
        sleep(5)
        # return a single mpd url
        return full_url_list_2


def bool_to_val(bool_val:bool)->str:
    '''
    turn a boolean value to a string "on"/"off"
    '''
    if bool_val:
        return "on"
    else:
        return "off"
