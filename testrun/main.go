package main

import (
	"fmt"
	"time"

	gst "github.com/FistOfTheNorthStar/gstreamer-go"
	//gst "github.com/notedit/gstreamer-go"
)

func main() {
	finishedSnd := make(chan bool)
	finishedPoll := make(chan bool)

	testGst()
	go testGstSend(finishedSnd)
	go testGstPoll(finishedPoll)

	<-finishedSnd
	<-finishedPoll

}

func testGst() {
	fmt.Printf("Testing gst start\n")

	pipeline, err := gst.New("videotestsrc  ! capsfilter name=filter ! autovideosink")
	if err != nil {
		// t.Error("pipeline create error", err)
		// t.FailNow()
	}

	filter := pipeline.FindElement("filter")

	if filter == nil {
		// t.Error("pipeline find element error ")
	}

	filter.SetCap("video/x-raw,width=1280,height=720")

	pipeline.Start()

	fmt.Printf("Testing gst done\n")
}

func testGstSend(finishedSnd chan bool) {
	pipeline, err := gst.New("appsrc name=mysource format=time is-live=true do-timestamp=true ! videoconvert ! autovideosink")

	if err != nil {
		// t.Error("pipeline create error", err)
		// t.FailNow()
	}

	appsrc := pipeline.FindElement("mysource")

	appsrc.SetCap("video/x-raw,format=RGB,width=320,height=240,bpp=24,depth=24")

	pipeline.Start()

	for {
		fmt.Printf("Testing gst send, pushing data to source\n")
		time.Sleep(1 * time.Second)
		appsrc.Push(make([]byte, 320*240*3))
	}
	finishedSnd <- true

}

func testGstPoll(finishedPoll chan bool) {
	pipeline, err := gst.New("videotestsrc ! video/x-raw,format=I420,framerate=15/1 ! x264enc bframes=0 speed-preset=veryfast key-int-max=60  ! udpsink host=127.0.0.1 port=1235 name=sink")

	if err != nil {
		fmt.Println("pipeline create error", err)
		// t.FailNow()
	}

	appsink := pipeline.FindElement("sink")

	pipeline.Start()

	fmt.Printf("Testing poll\n")

	out := appsink.Poll()

	fmt.Printf("Testing gst poll done\n")

	for {
		buffer := <-out
		fmt.Println("push ", len(buffer))
	}

	finishedPoll <- true

}
