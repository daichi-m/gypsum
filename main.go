package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	lipsum "github.com/daichi-m/lipsum/lipsum"
	flag "github.com/spf13/pflag"
)

func main() {
	input, err := scanInputs()
	if err != nil {
		log.Fatalln(err.Error())
	}
	err = input.GenerateOut()
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println("All files has been generated")
}

func scanInputs() (*lipsum.Input, error) {

	out := flag.StringP("out", "o", "-",
		"The file or directory in which the generated text would be stored")
	size := flag.StringP("size", "s", "10M", "Size of each file generated")
	count := flag.IntP("count", "c", 1, `The number of files generated. 
	If this is 1, a single text file is generated with all the content. 
	Otherwise the content is written in multiple text files in a directory`)
	flag.Parse()

	isize, err := convertSize(*size)
	if err != nil {
		return nil, err
	}

	var input lipsum.Input = lipsum.Input{
		File:  *out,
		Size:  isize,
		Count: *count,
	}
	err = input.Validate()
	if err != nil {
		return nil, err
	}
	return &input, nil

}

func convertSize(size string) (int, error) {
	reg, _ := regexp.Compile("(?P<Size>[0-9]+)(?P<Suffix>[TGMKtgmk]?[Bb]?)")
	submatch := reg.FindStringSubmatch(size)

	match := make(map[string]string)
	for i, name := range reg.SubexpNames() {
		match[name] = submatch[i]
	}

	num, ok := match["Size"]
	if !ok {
		return 0, fmt.Errorf("Could not decode size from %s", size)
	}
	suf, ok := match["Suffix"]
	if !ok {
		suf = "B"
		log.Println("Could not decode suffix for %s, defaulting to byte", size)
	}

	switch suf {
	case "b", "B":
		return strconv.Atoi(num)
	case "k", "K", "kb", "KB", "kB", "Kb":
		val, err := strconv.Atoi(num)
		if err != nil {
			return 0, err
		}
		return val * lipsum.KB, nil

	case "m", "M", "mb", "MB", "mB", "Mb":
		val, err := strconv.Atoi(num)
		if err != nil {
			return 0, err
		}
		return val * lipsum.MB, nil
	case "g", "G", "gb", "GB", "gB", "Gb":
		val, err := strconv.Atoi(num)
		if err != nil {
			return 0, err
		}
		return val * lipsum.GB, nil
	case "t", "T", "tb", "TB", "tB", "Tb":
		val, err := strconv.Atoi(num)
		if err != nil {
			return 0, err
		}
		return val * lipsum.TB, nil
	}
	return 0, fmt.Errorf("Unable to interpret size %s", size)
}
