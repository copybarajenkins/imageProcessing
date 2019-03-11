// Jiaxin Dong
// CPSC 5200, WIN 2019
// Assignment 3, image processing server

package main

import (
	"bytes"
	"encoding/json"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/gorilla/mux"
)

// ProcessImageRequest -
type ProcessImageRequest struct {
	URL        string   `json:"url"`
	Operations []string `json:"operations"`
}

// ErrorMessage -
type ErrorMessage struct {
	Error string `json:"error"`
}

func main() {
	server()
}

func server() {
	r := mux.NewRouter()
	r.HandleFunc("/process", processImageHandler).Methods("POST")
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func processImageHandler(w http.ResponseWriter, r *http.Request) {
	img, operations := decodeRequest(w, r)
	result := img
	t1 := time.Now()
	for _, op := range operations {
		result = processImage(result, op)
	}
	t2 := time.Now()
	log.Printf("Time taken for %d operations: %s",
		len(operations), t2.Sub(t1))
	writeImageAsResponse(w, &result)
}

func decodeRequest(w http.ResponseWriter, r *http.Request) (image.Image, []string) {
	// handle json
	req, e := unmarshalJSON(r)
	if e != nil {
		handleError(w, e.Error())
	}

	// handle image download
	url := req.URL // example url: https://i.imgur.com/AOXD2P4.jpg
	response, e := http.Get(url)
	if e != nil {
		handleError(w, e.Error())
	}
	defer response.Body.Close()
	imgBytes, e := ioutil.ReadAll(response.Body)
	if e != nil {
		handleError(w, e.Error())
	}
	log.Printf("Image size: %d bytes", len(imgBytes))
	img, _, _ := image.Decode(bytes.NewReader(imgBytes))

	return img, req.Operations
}

func processImage(img image.Image, op string) image.Image {
	if op == "flipVertical" {
		return flipVertical(img)
	} else if op == "flipHorizontal" {
		return flipHorizontal(img)
	} else if op == "rotateRight" {
		return rotateRight(img)
	} else if op == "rotateLeft" {
		return rotateLeft(img)
	} else if op == "grayscale" {
		return convertToGrayscale(img)
	} else if op == "thumbnail" {
		return resize(img, 100)
	} else if strings.HasPrefix(op, "resize") {
		resizeParams := strings.Split(op, ",")
		if len(resizeParams) > 2 {
			width, _ := strconv.Atoi(resizeParams[1])
			height, _ := strconv.Atoi(resizeParams[2])
			imaging.Resize(img, width, height, imaging.Lanczos)
		} else if len(resizeParams) > 1 {
			val, _ := strconv.Atoi(resizeParams[1])
			return resize(img, val)
		}
		return resize(img, 100)
	} else if strings.HasPrefix(op, "rotate") {
		rotateParams := strings.Split(op, ",")
		if len(rotateParams) > 1 {
			val, _ := strconv.ParseFloat(rotateParams[1], 64)
			return imaging.Rotate(img, val, color.Opaque)
		}
		return rotateRight(img) // default to rotate right
	}
	return img
}

func writeImageAsResponse(w http.ResponseWriter, img *image.Image) {
	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, *img, nil); err != nil {
		handleError(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		handleError(w, err.Error())
		return
	}
}

func flipHorizontal(input image.Image) image.Image {
	return imaging.FlipH(input)
}

func flipVertical(input image.Image) image.Image {
	return imaging.FlipV(input)
}

func rotateRight(input image.Image) image.Image {
	return imaging.Rotate270(input)
}

func rotateLeft(input image.Image) image.Image {
	return imaging.Rotate90(input)
}

func convertToGrayscale(input image.Image) image.Image {
	return imaging.Grayscale(input)
}

func resize(input image.Image, scale int) image.Image {
	return imaging.Resize(input, scale, 0, imaging.Lanczos)
}

func unmarshalJSON(r *http.Request) (ProcessImageRequest, error) {
	req := ProcessImageRequest{}
	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		return req, readErr
	}
	if jsonErr := json.Unmarshal(body, &req); jsonErr != nil {
		return req, jsonErr
	}
	return req, nil
}

func handleError(w http.ResponseWriter, err string) {
	response, error := json.Marshal(&ErrorMessage{Error: err})
	if error != nil {
		log.Fatal("Parsing error while reporting error")
		return
	}
	log.Fatal("Error: " + err)
	w.WriteHeader(500)
	w.Write(response)
}
