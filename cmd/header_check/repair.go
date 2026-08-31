//
//  repair.go
//  header_check
//
//  Created by Karl Kraft on 08/19/2026
//  Copyright 2026 Karl Kraft. All rights reserved
//

package main

import (
	"bufio"
	"io"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// headerTemplate describes one expected line of the header block: how to
// recognize an existing line that already satisfies it, and how to
// synthesize a replacement when it is missing.
type headerTemplate struct {
	match    func(line string) bool
	generate func() string
}

var createdByPattern = regexp.MustCompile(`^//  Created by .+ on \d{1,2}/\d{1,2}/(\d{2}|\d{4})(\.){0,1}$`)

func blankCommentTemplate() headerTemplate {
	return headerTemplate{
		match:    func(line string) bool { return line == "//" },
		generate: func() string { return "//" },
	}
}

func headerTemplates(filePath, basename, license string, owners []string) []headerTemplate {
	targetNames, defaultTarget := allowedTargetNames(filePath)

	targetDir := path.Dir(filePath)
	if targetDir == "." || targetDir == "" {
		if cwd, err := os.Getwd(); err == nil {
			targetDir = cwd
		}
	}
	if defaultTarget == nil {
		fallback := path.Base(targetDir)
		defaultTarget = &fallback
	}

	return []headerTemplate{
		blankCommentTemplate(),
		{
			match:    func(line string) bool { return line == "//  "+basename },
			generate: func() string { return "//  " + basename },
		},
		{
			match: func(line string) bool {
				for _, aTargetName := range targetNames.ToSlice() {
					if line == "//  "+aTargetName {
						return true
					}
				}
				return false
			},
			generate: func() string { return "//  " + *defaultTarget },
		},
		blankCommentTemplate(),
		{
			match: createdByPattern.MatchString,
			generate: func() string {
				return "//  Created by " + owners[0] + " on " + time.Now().Format("1/2/2006")
			},
		},
		{
			match: func(line string) bool {
				status, _ := isCopyrightValid(line, license, owners)
				return status != Invalid
			},
			generate: func() string {
				_, generated := isCopyrightValid("", license, owners)
				return generated
			},
		},
		blankCommentTemplate(),
		{
			match:    func(line string) bool { return line == "" },
			generate: func() string { return "" },
		},
	}
}

// repairFile reads filePath, inserts any missing header lines (leaving
// present-and-valid lines untouched), and rewrites the file in place. It
// returns true if the file was modified.
func repairFile(filePath, license string, owners []string) bool {
	basename := path.Base(filePath)

	f, err := os.Open(filePath)
	if err != nil {
		log.Errorf("Could not open %s for reading", filePath)
		return false
	}

	var originalLines []string
	reader := bufio.NewReader(f)
	for {
		aLine, readErr := reader.ReadString('\n')
		if len(aLine) > 0 {
			originalLines = append(originalLines, aLine)
		}
		if readErr == io.EOF {
			break
		} else if readErr != nil {
			log.Errorf("Error reading %s", filePath)
			if closeErr := f.Close(); closeErr != nil {
				log.Warningf("Could not close %s: %s", filePath, closeErr)
			}
			return false
		}
	}
	if closeErr := f.Close(); closeErr != nil {
		log.Warningf("Could not close %s: %s", filePath, closeErr)
	}

	if len(originalLines) == 0 {
		return false
	}

	// Package.swift requires a two line tools-version preamble ahead of the
	// copyright block; header_check does not validate those lines, so
	// repair leaves them untouched rather than guessing their content.
	preambleLines := 0
	if basename == "Package.swift" {
		preambleLines = 2
		if preambleLines > len(originalLines) {
			preambleLines = len(originalLines)
		}
	}

	templates := headerTemplates(filePath, basename, license, owners)

	result := append([]string{}, originalLines[:preambleLines]...)

	actualIndex := preambleLines
	changed := false
	for _, tmpl := range templates {
		if actualIndex < len(originalLines) {
			actualLine := strings.TrimRight(originalLines[actualIndex], "\n")
			if tmpl.match(actualLine) {
				result = append(result, originalLines[actualIndex])
				actualIndex++
				continue
			}
		}
		result = append(result, tmpl.generate()+"\n")
		changed = true
	}
	result = append(result, originalLines[actualIndex:]...)

	if !changed {
		return false
	}

	dir := path.Dir(filePath)
	output, err := os.CreateTemp(dir, "_header_check_*.tmp")
	if err != nil {
		log.Errorf("Could not create temp file for repairing header")
		return false
	}
	for _, aLine := range result {
		if _, writeErr := output.WriteString(aLine); writeErr != nil {
			log.Errorf("Error writing repaired header for %s", filePath)
			if closeErr := output.Close(); closeErr != nil {
				log.Warningf("Could not close %s: %s", output.Name(), closeErr)
			}
			if removeErr := os.Remove(output.Name()); removeErr != nil {
				log.Warningf("Could not remove %s: %s", output.Name(), removeErr)
			}
			return false
		}
	}
	if err := output.Close(); err != nil {
		log.Errorf("Could not close temp output file")
		return false
	}
	if err := os.Rename(output.Name(), filePath); err != nil {
		log.Errorf("Could not rename temp to source file")
		return false
	}
	return true
}
