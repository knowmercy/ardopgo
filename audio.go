package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
)

const (
	Duration   = 7     // Time in seconds
	SampleRate = 44100 // 44.1 khz
	Tone1      = 1760 * 8
)

var (
	start float64 = 1.0
	end   float64 = 1.0
	tau           = math.Pi * 2
)

func doit() {
	fmt.Fprintf(os.Stderr, "generating sine wave..\n")
	file := "out.bin"
	f, _ := os.Create(file)

	go generate(Tone1, f)
	fmt.Fprintf(os.Stderr, "done")
}

func generate(tone float64, file *File) {
	samples := Duration * SampleRate
	var angle float64 = tau / float64(samples)

	decayfac := math.Pow(end/start, 1.0/float64(nsamps))
	var buf [8]byte

	for i := 0; i < samples; i++ {
		sample := math.Sin(angle * tone * float64(i))
		sample *= start
		//start *= decayfac
		binary.LittleEndian.PutUint32(buf[:], math.Float32bits(float32(sample)))
		bw, err := file.Write(buf[:])
		if err != nil {
			panic(err)
		}
		fmt.Printf("\rWrote: %v bytes to %s", bw, file)
	}
}
