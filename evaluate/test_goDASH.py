#!/usr/bin/python
# /*
#  *	goDASH, golang client emulator for DASH video streaming
#  *	Copyright (c) 2019, Jason Quinlan, Darijo Raca, University College Cork
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

# python3 ./test_goDASH.py --numClients=1 --terminalPrint="off" --debug="off" --collaborative="off"

from argparse import ArgumentParser
import os
import datetime
from subprocess import Popen
from urls.mpdURL import *
from random import randint
from time import sleep


parser = ArgumentParser(description="goDASH - a player of infinite quality")

parser.add_argument('--numClients',
                    dest="numClients",
                    help="number of clients to create and stream",
                    default=1)

parser.add_argument('--terminalPrint',
                    dest="terminalPrint",
                    help="print output of goDASH to the terminal screen",
                    default="on")

parser.add_argument('--debug',
                    dest="debug",
                    help="print output of goDASH to the log file",
                    default="on")

parser.add_argument('--collaborative',
                    dest="collaborative",
                    help="run network in collaborative mode",
                    default="off")

# Expt parameters
args = parser.parse_args()

# ouptut folder structure
# output
output_folder_name = "/output"
# - config
config_folder_name = "/config"
# - files
log_folder_name = "/files"

# get all the possible DASH MPD files from the H264 UHD dataset
if args.collaborative != "on":
    urls = full_url_list+main_url_list+live_url_list + \
        full_byte_range_url_list+main_byte_range_url_list
else:
    # lets start consul
    os.system("consul agent -dev &")
    # lets sleep until consul is set up
    sleep(5)
    urls = full_url_list_2

# create a dictionary from the default config file


def create_dict(config_file):
    # lets read in the original config file and create a dictionary we can use
    dict = {}
    # open the original config file
    with open(config_file, encoding='utf-8-sig') as fp:
        # read line by line
        line = fp.readline().strip()
        while line:
            # do not split around the brackets
            if line != "{" and line != "}":
                # split around the colon, but not the colon in http :)
                key, val = line.split(' : ')
                key = key.strip()
                val = val.strip()
                dict[key] = val
                line = fp.readline().strip()
            else:
                # otherwise, just read the line
                line = fp.readline().strip()

    # return the dictionary
    return dict

# lets get some work done


def eval_goDASH():

    # lets read in the goDASH config file
    cwd = os.getcwd()
    config_direct = cwd + "/config"
    config_file = config_direct+"/configure.json"

    # lets read in the original config file and create a dictionary we can use
    dict = create_dict(config_file)

    # print the current segment log ocation
    # print(dict['"storeDash"'])
    # print the current log file location
    # print (dict['"logFile"'])

    # lets create the log and config folder locations
    output_folder = cwd+output_folder_name
    # create a folder based on date and time
    current_folder = "/" + datetime.datetime.now().strftime('%Y-%m-%d-%H-%M-%S')
    # - config
    config_folder = output_folder+current_folder+config_folder_name

    # lets create the output folder structure
    if not os.path.exists(config_folder):
        os.makedirs(config_folder)

    # our array of processes
    processes = []

    for i in range(1, int(args.numClients)+1):

        # lets create name for this client
        client_name = "client"+str(i)+"/"
        client_config = "/configure_"+str(i)+".json"

        # - files
        log_folder = output_folder+current_folder+log_folder_name+"/"+client_name
        # lets create the file output folder structure
        if not os.path.exists(log_folder):
            os.makedirs(log_folder)

        fout = config_folder+client_config
        getHeaders = False
        with open(fout, "w") as fo:
            fo.write('{\n')
            for k, v in dict.items():

                # write the key to the config file
                fo.write('\t\t\t\t' + str(k) + ' : ')

                # set the segmemt storage location value
                if k == '"outputFolder"':
                    fo.write(str("\""+client_name+"\","))
                # set the log file location value
                elif k == '"logFile"':
                    fo.write(str(str(v)[:-2]+"_client"+str(i)+"\","))
                # set terminal print value
                elif k == '"terminalPrint"':
                    fo.write(str("\""+args.terminalPrint+"\","))
                # set debug value
                elif k == '"debug"':
                    fo.write(str("\""+args.debug+"\","))
                # set the collaborative clients
                elif k == '"serveraddr"':
                    fo.write(str("\""+args.collaborative+"\""))
                # set url value
                elif k == '"url"':
                    # generate a random number
                    value = randint(0, len(urls)-1)
                    # call the url that corresponds to the index of the random number
                    fo.write(str("\"[" + urls[value] + "]\","))
                # check the getHeaders setting
                elif k == '"getHeaders"':
                    if v != '"off",':
                        getHeaders = True
                        print(True)
                    # write the value
                    fo.write(str(v))
                else:
                    # write the value
                    fo.write(str(v))
                # write a return carriage
                fo.write('\n')
            fo.write('}')

        # lets call each client from within its output folder
        os.chdir(log_folder+"../../")

        # lets call goDASH and get some output
        cmd = cwd+"/../godash --config " + \
            output_folder+current_folder+config_folder_name+client_config
        p = Popen(cmd, shell=True)
        # add this command to our list of processes
        processes.append(p)

    # will this tell us when all processes are complete
    for p in processes:
        if p.wait() != 0:
            if not getHeaders:
                print("There was an error with test_goDASH.py")
                return

    # if we previously got the headers, now lets stream
    if getHeaders:
        print("all header files have been downloaded, now lets stream")

        # reset the processes list
        processes = []

        # for each of our clients
        for i in range(1, int(args.numClients)+1):

            # lets create name for this client
            client_name = "client"+str(i)+"/"
            client_config = "/configure_"+str(i)+".json"

            # - files
            log_folder = output_folder+current_folder+log_folder_name+"/"+client_name

            # lets read in the current config file and create a dictionary we can use
            dict = create_dict(config_folder+client_config)

            # keep the same output config file
            fout = config_folder+client_config
            getHeaders = False
            with open(fout, "w") as fo:
                fo.write('{\n')
                for k, v in dict.items():

                    # write the key to the config file
                    fo.write('\t\t\t\t' + str(k) + ' : ')

                    # check the getHeaders setting
                    if k == '"getHeaders"':
                        fo.write(str("\"off\","))
                    else:
                        # write the value
                        fo.write(str(v))
                    # write a return carriage
                    fo.write('\n')
                fo.write('}')

            # lets call each client from within its output folder
            os.chdir(log_folder+"../")

            # lets call goDASH and get some output
            cmd = cwd+"/../godash --config " + \
                output_folder+current_folder+config_folder_name+client_config
            p = Popen(cmd, shell=True)
            # add this command to our list of processes
            processes.append(p)

            #sleep(2)

        # will this tell us when all processes are complete
        for p in processes:
            if p.wait() != 0:
                if not getHeaders:
                    print("There was an error with test_goDASH.py")
                    return

        # now all clients have finished
        print("all goDASH clients have finished streaming")

    # otherwise we are done
    else:
        print("all goDASH clients have finished streaming")

    # lets stop consul
    os.system("killall -9 consul")


# let's call the main function
if __name__ == '__main__':
    eval_goDASH()
