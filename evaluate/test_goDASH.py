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

# options for this script are defined in /config/settings.py
# python3 ./test_goDASH.py

from argparse import ArgumentParser
import os
import datetime
from subprocess import Popen
from random import randint
from time import sleep

from config.settings import *

# lets get some work done
def eval_goDASH():

    # read in the goDASH config file
    cwd = os.getcwd()
    # add this, just to make sure we are in the evaluate folder
    if not "evaluate" in cwd:
        cwd += "/evaluate"
    # set up the config folder details
    config_direct = cwd + config_folder_name
    config_file_loc = config_direct+config_file

    # lets read in the original config file and create a dictionary we can use
    dict = create_dict(config_file_loc)

    # lets create the log and config folder locations
    output_folder = cwd+output_folder_name
    # create a folder based on date and time
    current_folder = "/" + datetime.datetime.now().strftime('%Y-%m-%d-%H-%M-%S')
    # - config
    config_folder = output_folder+current_folder+config_folder_name

    # lets create the output folder structure
    if not os.path.exists(config_folder):
        os.makedirs(config_folder)

    # get the possible DASH MPD files from the H264 UHD dataset
    urls = check_collab_and_set_url(collaborative, single_clip_choice)

    # our array of processes
    processes = []

    for i in range(1, int(numClients)+1):

        # lets create a name for this client
        client_name = "client"+str(i)+"/"
        client_config = "/configure_"+str(i)+".json"

        # - files
        log_folder = output_folder+current_folder+log_folder_name+"/"+client_name

        # lets create the file output folder structure, if it does not exist
        if not os.path.exists(log_folder):
            os.makedirs(log_folder)

        # create the config file for this client
        fout = config_folder+client_config
        getHeaders = False
        with open(fout, "w") as fo:
            fo.write('{\n')
            for k, v in dict.items():

                # write the key to the config file
                fo.write('\t\t' + str(k) + ' : ')

                # set the segmemt storage location value
                if k == '"outputFolder"':
                    fo.write(str("\""+client_name+"\","))
                # set the algorithm for each clients
                elif k == '"adapt"':
                    fo.write(str("\""+godash_run_dict["algo_choice"][i-1]+"\","))
                # set the log file location value
                elif k == '"logFile"':
                    fo.write(str(str(v)[:-3]+"_client"+str(i)+"\","))
                # set terminal print value
                elif k == '"terminalPrint"':
                    fo.write(str("\""+terminalPrintval+"\","))
                # set debug value
                elif k == '"debug"':
                    fo.write(str("\""+debugval+"\","))
                # set the collaborative clients - no comma on this one
                elif k == '"serveraddr"':
                    fo.write(str("\""+collaborativeval+"\""))
                # store the files, if collab is on
                elif k == '"storeDash"':
                    fo.write(str("\""+collaborativeval+"\","))
                # set url value
                elif k == '"url"':
                    # generate a random number
                    value = randint(0, len(urls)-1)
                    # call the url that corresponds to the index of the random number
                    fo.write(str("\"[" + urls[value] + "]\","))
                # check the getHeaders setting
                elif k == '"getHeaders"':
                    if v != '"off"':
                        getHeaders = True
                    # write the value
                    fo.write(str(godash_run_dict[k])+",")

                # set the kind of default values - these are changed as per the settings file
                else:
                    # if "printHeader" in k:
                    #     fo.write(str(v)+ ",")
                    # else:
                        # write the value
                    fo.write(str(godash_run_dict[k])+",")
                # write a return carriage
                fo.write('\n')
            fo.write('}')

        # lets call each client from within its output folder
        os.chdir(log_folder+"../../")

        # lets call goDASH and get some output
        cmd = cwd+"/../godash --config " + \
            output_folder+current_folder+config_folder_name+client_config
        sleep(2)
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
        for i in range(1, int(numClients)+1):

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
                    fo.write('\t\t' + str(k) + ' : ')

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
            sleep(2)
            p = Popen(cmd, shell=True)
            # add this command to our list of processes
            processes.append(p)

            #sleep(2)

        # this tell us when all processes are complete
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
    if collaborative:
        os.system("killall -9 consul")


# let's call the main function
if __name__ == '__main__':
    eval_goDASH()
