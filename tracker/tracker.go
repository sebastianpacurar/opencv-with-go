package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"gocv.io/x/gocv/contrib"
	"image/color"
	"log"
)

func main() {
	deviceID := 0

	// open webcam
	webcam, err := gocv.OpenVideoCapture(deviceID)
	if err != nil {
		fmt.Printf("Error opening video capture device: %v\n", deviceID)
		return
	}
	defer func(webcam *gocv.VideoCapture) {
		err := webcam.Close()
		if err != nil {
			log.Fatal("error when closing webcam")
		}
	}(webcam)

	// open display window
	w := gocv.NewWindow("Tracking")
	defer func(w *gocv.Window) {
		err := w.Close()
		if err != nil {
			log.Fatalln("error when closing window")
		}
	}(w)

	// create a tracker instance
	tracker := contrib.NewTrackerKCF()
	defer func(tracker gocv.Tracker) {
		err := tracker.Close()
		if err != nil {
			log.Fatalln("error when closing tracker")
		}
	}(tracker)

	// prepare image matrix
	img := gocv.NewMat()
	defer func(img *gocv.Mat) {
		err := img.Close()
		if err != nil {
			log.Fatalln("error when closing image matrix")
		}
	}(&img)

	// read an initial image
	if ok := webcam.Read(&img); !ok {
		fmt.Printf("cannot read device %v\n", deviceID)
		return
	}

	// let the user mark a ROI to track
	rect := w.SelectROI(img)
	if rect.Max.X == 0 {
		fmt.Printf("user cancelled roi selection\n")
		return
	}

	// initialize the tracker with the image & the selected roi
	init := tracker.Init(img, rect)
	if !init {
		fmt.Printf("Could not initialize the Tracker")
		return
	}

	// color for the rect to draw
	blue := color.RGBA{B: 255, A: 255}
	fmt.Printf("Start reading device: %v\n", deviceID)
	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		// update the roi
		rect, _ := tracker.Update(img)

		// draw it.
		gocv.Rectangle(&img, rect, blue, 3)

		// show the image in the window, and wait 10 millisecond
		w.IMShow(img)
		if w.WaitKey(10) >= 0 {
			break
		}
	}
}
