package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/hybridgroup/mjpeg"
	"gocv.io/x/gocv"
)

var (
	device int
	err    error
	webcam *gocv.VideoCapture
	stream *mjpeg.Stream
)

func main() {
	deviceID := 0
	host := "localhost:3000"

	// open webcam
	webcam, err = gocv.OpenVideoCapture(0)
	if err != nil {
		fmt.Printf("Error opening capture device: %v\n", deviceID)
		return
	}
	defer func(webcam *gocv.VideoCapture) {
		err := webcam.Close()
		if err != nil {
			log.Fatalln("error when closing video capture device")
		}
	}(webcam)

	// create the mjpeg stream
	stream = mjpeg.NewStream()

	// start video capture routine
	go mjpegCapture()

	fmt.Println("Capturing. Point your browser to " + host)

	// start http server
	http.Handle("/", stream)
	log.Fatal(http.ListenAndServe(host, nil))
}

func mjpegCapture() {
	img := gocv.NewMat()
	defer func(img *gocv.Mat) {
		err := img.Close()
		if err != nil {
			log.Fatalln("error when closing image matrix")
		}
	}(&img)

	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", device)
			return
		}
		if img.Empty() {
			continue
		}

		buf, _ := gocv.IMEncode(".jpg", img)
		stream.UpdateJPEG(buf.GetBytes())
		buf.Close()
	}
}
