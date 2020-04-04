package main

import (
	"log"
	"github.com/nfnt/resize"
	"net/http"
	"fmt"
	"unicode"
	"image/jpeg"
	"strconv"
	"os"
	"image"
)

var img image.Image

func main() {

	//Load the image of Väinö
	file, err := os.Open("vape.jpg")
	if err != nil { log.Fatal(err) }
	img, err = jpeg.Decode(file)
	if err != nil { log.Fatal(err) }
	file.Close()

	http.HandleFunc("/", sendImg)
	log.Fatal(http.ListenAndServe(":3003", nil))
}

func sendImg(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path[1:]
	widthStr := ""
	heightStr := ""
	onFirst := true
	//Parse image dimensions from the url
	for _, ch := range urlPath {
		if !unicode.IsDigit(ch) {
			//Skip everything until the first number appears
			if(len(widthStr) == 0) {
				continue
			}
			//if a non-number appears after we have got the
			//height, break.  
			if !onFirst && heightStr != "" {
				break
			}
			//Change to reading to the height string
			onFirst = false
			continue
		}
		if onFirst {
			widthStr += string(ch)
		} else {
			heightStr += string(ch)
		}
	}
	//If the user didn't input parameters that 
	//result in actual image dimensions
	//display the "documentation" page
	if(heightStr == "" || widthStr == "") {
		fmt.Fprintf(w, "<h1>Placeholding Väinö Dog</h1><h2>USAGE: m.ktj.st/vp/[width]x[height]</h2><h3>Example: m.ktj.st/vp/200x200</h3>")
	} else {
		width, err1 := strconv.Atoi(widthStr)
		height, err2 := strconv.Atoi(heightStr)
		if err1 != nil || err2 != nil {
			fmt.Fprintf(w, "%v %v", err1, err2)
			return
		}
		if width > 3000 || height > 3000 {
			fmt.Fprintf(w, "The max value for either dimension is 3000px")
			return
		}
		newImg := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
		jpeg.Encode(w, newImg, nil)
	}

}