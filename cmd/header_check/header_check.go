//
//  header_check.go
//  header_check
//
//  Created by Karl Kraft on 12/29/2022
//  Copyright 2022-2023 Karl Kraft. All rights reserved.
//

package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/go-git/go-git/v5"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

func allowedTargetNames(filePath string) []string {
	cwd, _ := os.Getwd()
	parent := path.Dir(filePath)
	gitCheck := parent + "/.git"
	if _, err := os.Stat(gitCheck); err == nil {
		a := []string{path.Base(parent)}
		r, err := git.PlainOpen(parent)
		if err != nil {
			log.Errorf("Could not open local git respository")
			return a
		}
		c, err := r.Config()
		if c.Remotes != nil && c.Remotes["origin"] != nil {
			for _, urlString := range c.Remotes["origin"].URLs {
				split := strings.Split(urlString, "/")
				lastPath := split[len(split)-1]
				if strings.HasSuffix(lastPath, ".git") {
					lastPath = lastPath[0 : len(lastPath)-4]
				}
				return append(a, lastPath)
			}
		}
		return a
	}
	if parent == cwd || parent == "." {
		a := []string{path.Base(cwd)}
		return a
	}
	a := []string{path.Base(parent)}
	return append(a, allowedTargetNames(parent)...)
}

func analyzeFile(filePath string, license string, autodate bool) bool {

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

	// 2 - target name
	lineIndex++

	var foundTargetName = false
	for _, aTargetName := range allowedTargetNames(filePath) {
		target = "//  " + aTargetName
		if lines[lineIndex] == target {
			foundTargetName = true
			break
		}
	}
	if !foundTargetName {
		if !fileReported {
			log.Warningf("%s", filePath)
			fileReported = true
		}
		log.Warningf("at line %d", lineIndex)
		log.Warningf("- %s", lines[lineIndex])
		for _, aTargetName := range allowedTargetNames(filePath) {
			log.Warningf("+ %s", aTargetName)
		}
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
	copyrightLine := lines[lineIndex]
	validCopyright, correctedLine := isCopyrightValid(copyrightLine, license)
	if validCopyright != Correct {
		if !fileReported {
			log.Warningf("%s", filePath)
			fileReported = true
		}
		log.Warningf("at line %d", lineIndex)
		log.Warningf("- %s", lines[lineIndex])
		log.Warningf("+ %s", correctedLine)
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

	if autodate && validCopyright == Expired {
		rewriteFileWithCopyright(filePath, copyrightLine, correctedLine)
	}

	return fileReported
}

func rewriteFileWithCopyright(filePath string, expiredLine string, replaceLine string) {
	expiredLine = expiredLine + "\n"
	replaceLine = replaceLine + "\n"
	dir := path.Dir(filePath)
	output, err := os.CreateTemp(dir, "_header_check_*.tmp")
	if err != nil {
		log.Errorf("Could not create temp file for updating dates")
		return
	}
	f, err := os.Open(filePath)
	if err != nil {
		log.Errorf("Could not open " + filePath + " for reading")
		return
	}

	reader := bufio.NewReader(f)
	for {
		aLine, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Errorf("Error reading  " + filePath)
			return

		}
		if aLine == expiredLine {
			_, err := output.WriteString(replaceLine)
			if err != nil {
				log.Errorf("Error writing replacement date line")
				return
			}
		} else {
			_, err := output.WriteString(aLine)
			if err != nil {
				log.Errorf("Error writing source code while fixing date")
				return
			}
		}
	}
	err = output.Close()
	if err != nil {
		log.Errorf("Could not close temp output file")
		return
	}
	err = f.Close()
	if err != nil {
		log.Errorf("Could not close source file")
		return
	}
	err = os.Rename(output.Name(), filePath)
	if err != nil {
		log.Errorf("Could not rename tempt to source file")
		return
	}
}

type CopyrightStatus int

const (
	Correct CopyrightStatus = iota
	Expired
	Invalid
)

func isCopyrightValid(theLine string, license string) (CopyrightStatus, string) {
	year, _, _ := time.Now().Date()
	formattedYear := fmt.Sprintf("%04d", year)
	singleYearCorrectPattern := "^//  Copyright " + formattedYear + " Karl Kraft. " + license + "(\\.){0,1}$"
	singleYearExpiredPattern := "^//  Copyright (\\d{4}) Karl Kraft. " + license + "(\\.){0,1}$"
	rangeCorrectPattern := "^//  Copyright \\d{4}-" + formattedYear + " Karl Kraft. " + license + "(\\.){0,1}$"
	rangeExpiredPattern := "^//  Copyright (\\d{4})-\\d{4} Karl Kraft. " + license + "(\\.){0,1}$"

	r, _ := regexp.Compile(singleYearCorrectPattern)
	if r.MatchString(theLine) {
		return Correct, ""
	}
	r, _ = regexp.Compile(rangeCorrectPattern)
	if r.MatchString(theLine) {
		return Correct, ""
	}
	r, _ = regexp.Compile(singleYearExpiredPattern)
	if r.MatchString(theLine) {
		matches := r.FindStringSubmatch(theLine)
		return Expired, "//  Copyright " + matches[1] + "-" + formattedYear + " Karl Kraft. " + license
	}
	r, _ = regexp.Compile(rangeExpiredPattern)
	if r.MatchString(theLine) {
		matches := r.FindStringSubmatch(theLine)
		return Expired, "//  Copyright " + matches[1] + "-" + formattedYear + " Karl Kraft. " + license
	}

	return Invalid, "//  Copyright " + formattedYear + " Karl Kraft. " + license

}

func main() {
	log.Printf("Startup %s(%s) Built: %s", version, revision, builtDate)
	license := flag.String("license", "arr", "License mode (arr,apache)")
	autodate := flag.Bool("autodate", false, "Auto update copyright lines")
	infoplist := flag.Bool("infoplist", false, "Scan for Info.plist")

	flag.Parse()
	log.Infof("License is set to %s", *license)
	var failed bool
	var licenseString = "All rights reserved"

	if *license == "apache" {
		licenseString = "Licensed under Apache License, Version 2.0"
	}
	for _, s := range flag.Args() {
		log.Infof("Reading %s", s)
		failed = analyzeFile(s, licenseString, *autodate) || failed
	}
	if *infoplist {
		// scan and maybe update
	}
	if failed {
		os.Exit(-1)
	} else {
		os.Exit(0)
	}
}
