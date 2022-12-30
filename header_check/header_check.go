//
//  header_check.go
//  header_check
//
//  Created by Karl Kraft on 12/29/2022
//  Copyright 2022 Karl Kraft. All rights reserved.
//

package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

func analyzeFile(filePath string, license string) bool {

	basename := path.Base(filePath)

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	buffer := make([]byte, 1024)
	numberRead, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		log.Errorf("Could not open for reading: %s", filePath)
		log.Fatal(err)
	}

	if numberRead == 0 {
		log.Warningf("%s", filePath)
		log.Warningf(">> File is empty")
		return true
	}

	lines := strings.Split(string(buffer), "\n")

	lineIndex := 0
	minLines := 7

	if basename == "Package.swift" {
		minLines = minLines + 2
		lineIndex = lineIndex + 2
	}
	if len(lines) < minLines {
		log.Warningf("%s", filePath)
		log.Warningf(">> Less than seven lines in source file")
		return true
	}

	var fileReported bool
	var target string

	// 0 - blank line
	target = "//"
	if lines[lineIndex] != target {
		if !fileReported {
			log.Warningf("%s", filePath)
			fileReported = true
		}
		log.Warningf("at line %d", lineIndex)
		log.Warningf("- %s", lines[lineIndex])
		log.Warningf("+ %s", target)
	}

	// 1 - name of file
	lineIndex++
	target = "//  " + basename
	if lines[lineIndex] != target {
		if !fileReported {
			log.Warningf("%s", filePath)
			fileReported = true
		}
		log.Warningf("at line %d", lineIndex)
		log.Warningf("- %s", lines[lineIndex])
		log.Warningf("+ %s", target)
	}

	// 2 - project name
	lineIndex++
	cwd, _ := os.Getwd()
	target = "//  " + path.Base(cwd)
	if lines[lineIndex] != target {
		if !fileReported {
			log.Warningf("%s", filePath)
			fileReported = true
		}
		log.Warningf("at line %d", lineIndex)
		log.Warningf("- %s", lines[lineIndex])
		log.Warningf("+ %s", target)
	}

	// 3 - blank line
	lineIndex++
	target = "//"
	if lines[lineIndex] != target {
		if !fileReported {
			log.Warningf("%s", filePath)
			fileReported = true
		}
		log.Warningf("at line %d", lineIndex)
		log.Warningf("- %s", lines[lineIndex])
		log.Warningf("+ %s", target)
	}

	// 4 - Create by First Last on 01/01/01
	// allow optional period
	// allow 1 or 2 digit for day/month
	// allow 2 or 4 digit for year
	lineIndex++
	pattern := "^//  Created by .+ on \\d{1,2}/\\d{1,2}/(\\d{2}|\\d{4})(\\.){0,1}$"
	r, _ := regexp.Compile(pattern)
	if !r.MatchString(lines[lineIndex]) {
		if !fileReported {
			log.Warningf("%s", filePath)
			fileReported = true
		}
		log.Warningf("at line %d", lineIndex)
		log.Warningf("- %s", lines[lineIndex])
		log.Warningf("+ %s", pattern)
	}

	// 5 - Copyright and License
	// allow optional period
	// 4 digit year or range between a pair of 4 digit years
	// license indicator
	lineIndex++
	year, _, _ := time.Now().Date()
	formattedYear := fmt.Sprintf("%04d", year)
	pattern = "^//  Copyright (\\d{4}-){0,1}" + formattedYear + " Karl Kraft. " + license + "(\\.){0,1}$"
	r, _ = regexp.Compile(pattern)
	if !r.MatchString(lines[lineIndex]) {
		if !fileReported {
			log.Warningf("%s", filePath)
			fileReported = true
		}
		log.Warningf("at line %d", lineIndex)
		log.Warningf("- %s", lines[lineIndex])
		log.Warningf("+ %s", pattern)
	}

	// 6 - blank line
	lineIndex++
	target = "//"
	if lines[lineIndex] != target {
		if !fileReported {
			log.Warningf("%s", filePath)
			fileReported = true
		}
		log.Warningf("at line %d", lineIndex)
		log.Warningf("- %s", lines[lineIndex])
		log.Warningf("+ %s", target)
	}

	// 7 - empty line
	lineIndex++
	target = ""
	if lines[lineIndex] != target {
		if !fileReported {
			log.Warningf("%s", filePath)
			fileReported = true
		}
		log.Warningf("at line %d", lineIndex)
		log.Warningf("- %s", lines[lineIndex])
		log.Warningf("+ %s", target)
	}

	return fileReported
}

func main() {
	license := flag.String("license", "arr", "License mode (arr,apache)")
	flag.Parse()
	log.Infof("License is set to %s", *license)
	var failed bool
	var licenseString = "All rights reserved"

	if *license == "apache" {
		licenseString = "Licensed under Apache License, Version 2.0"
	}
	for _, s := range flag.Args() {
		log.Infof("Reading %s", s)
		failed = analyzeFile(s, licenseString) || failed
	}
	if failed {
		os.Exit(-1)
	} else {
		os.Exit(0)
	}
}
