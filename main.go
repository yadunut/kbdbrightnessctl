package main

import (
	"flag"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type iFlag struct {
	set  bool
	data int
}

func (f *iFlag) String() string {
	if !f.set {
		return "unset"
	}

	return strconv.Itoa(f.data)
}

func (f *iFlag) Set(s string) error {
	data, err := strconv.Atoi(s)
	if err != nil {
		return err
	}

	f.data = data
	f.set = true

	return nil
}

const (
	bFile    = "/sys/class/leds/spi::kbd_backlight/brightness"
	maxBFile = "/sys/class/leds/spi::kbd_backlight/max_brightness"
)

func main() {
	brightness := &iFlag{}
	inc := &iFlag{}
	dec := &iFlag{}

	flag.Var(brightness, "brightness", "set brightness of keyboard")
	flag.Var(inc, "inc", "increase brghtness of keyboard by amt")
	flag.Var(dec, "dec", "decrease brightness of keyboard by amt")

	flag.Parse()

	bb, err := ioutil.ReadFile(bFile)
	check(err)
	maxBB, err := ioutil.ReadFile(maxBFile)
	check(err)

	b, err := strconv.Atoi(strings.TrimSpace(string(bb)))
	check(err)
	maxB, err := strconv.Atoi(strings.TrimSpace(string(maxBB)))
	check(err)

	if brightness.set {
		b = brightness.data
		b = min(b, maxB)
		b = max(0, b)
		err := ioutil.WriteFile(bFile, []byte(strconv.Itoa(b)), 331)
		check(err)
		return
	}

	if inc.set {
		b += inc.data
		b = min(b, maxB)
		b = max(0, b)
		err := ioutil.WriteFile(bFile, []byte(strconv.Itoa(b)), 331)
		check(err)
		return
	}

	if dec.set {
		b -= dec.data
		b = min(b, maxB)
		b = max(0, b)
		err := ioutil.WriteFile(bFile, []byte(strconv.Itoa(b)), 331)
		check(err)
		return
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
