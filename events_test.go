package main

import (
	"os"
	"testing"
	"time"
)

func TestEventJSON(t *testing.T) {
	blocks, _ := getBlocks(time.Now(), time.Now().AddDate(0, 0, 10))
	/* Test makeEvents from TimeBlock slice */
	events := makeEvents(blocks)
	/* If this is only event no blocks in range (assuming blocks_test passes) */
	events = append(events, &Event{Start: time.Now(), End: time.Now().Add(150), Title: "TEST", Room: "yellow"})
	/* Serve the JSON to stdout for analysis */
	serveEventJSON(os.Stdout, events)
	t.Log("TestEventJSON Complete")
}
