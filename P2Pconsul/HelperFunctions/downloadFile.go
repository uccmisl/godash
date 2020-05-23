package HelperFunctions

// Collaborative Code - Start

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

//func main(){
//	fileUrl := "http://cs1dev.ucc.ie/misl/4K_non_copyright_dataset/2_sec/x264/bbb/DASH_Files/full/bbb_enc_x264_dash.mpd"
//
//	if err := DownloadFile("bbb_enc_x264_dash.mpd",fileUrl); err != nil{
//		panic(err)
//	}
//}

func DownloadFile(filepath string, url string) error {

	fmt.Println(filepath)
	fmt.Println(url)

	//download data
	response, err := http.Get(url)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	//create file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}

	defer out.Close()

	//Write to file
	_, err = io.Copy(out, response.Body)
	return err
}

// Collaborative Code - End
