package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

var (
	// Column 1: only lowercase a–z
	reCol1 = regexp.MustCompile(`^[a-z']+$`)
	// Columns 2..n: lowercase a–z, underscore, caret, digits 1–6
	reRest = regexp.MustCompile(`^[a-z_^1-6*']+$`)
	// For concat rule: keep ONLY lowercase a–z from columns 2..n
	reLettersOnly = regexp.MustCompile(`[^a-z']+`)
	// Valid phonogram check 
	// only one phonogram per line including *, '
)

func main() {
	file, err := os.Open("data/phonograms_examples.csv")
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // allow variable number of fields per row

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("read error: %v", err)
		}

		// ---- Rule 1: valid characters check ----
		if !validChars(record) {
			fmt.Printf("valid characters check: %s\n", strings.Join(record, ","))
			continue // stop evaluating further checks for this row
		}

		// ---- Rule 2: correct character check ----
		if !concatMatches(record) {
			fmt.Printf("match character check: %s\n", strings.Join(record, ","))
			continue
		}

		// Passes all rules -> no output
	}
}

func validChars(rec []string) bool {
	// first column strict a–z
	if !reCol1.MatchString(rec[0]) {
		return false
	}
	// remaining columns allowed charset (no normalization)
	for _, f := range rec[1:] {
		if !reRest.MatchString(f) {
			return false
		}
	}
	return true
}

func concatMatches(rec []string) bool {
	want := rec[0]

	var b strings.Builder
	for _, f := range rec[1:] {
		b.WriteString(reLettersOnly.ReplaceAllString(f, ""))
	}
	return b.String() == want
}

