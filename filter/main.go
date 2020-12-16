package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

var (
	dir  = flag.String("inputDir", "input", "")
	expr = flag.String("expr", "", "")
	re   *regexp.Regexp
)

func main() {
	flag.Parse()
	if *expr != "" {
		re = regexp.MustCompile(*expr)
	}
	filepath.Walk(*dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		return processFile(path)
	})
}

func processFile(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		oldS := strings.TrimRightFunc(scanner.Text(), unicode.IsSpace)
		newS := strings.TrimSpace(oldS)
		if newS == "" {
			continue
		}
		if oldS == newS {
			processBlock(newS, scanner)
		}
	}
	return scanner.Err()
}

func processBlock(name string, scanner *bufio.Scanner) error {
	headingPrinted := false
	defer func() {
		if headingPrinted {
			fmt.Println()
		}
	}()
	for scanner.Scan() {
		t := scanner.Text()
		if !strings.HasPrefix(t, "\t") {
			return nil
		}
		if !re.MatchString(t) {
			continue
		}
		if !headingPrinted {
			fmt.Println(name)
			headingPrinted = true
		}
		fmt.Println(t)
	}
	return scanner.Err()
}
