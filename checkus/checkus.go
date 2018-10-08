package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func isSuspicious(s string) bool {
	if strings.Contains(s, " = _") {
		return true
	}
	return strings.Contains(s, " = ") && strings.Contains(s, "(")
}

func main() {
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		scanner := bufio.NewScanner(f)
		inPub, level := false, 0
		n := 0
		for scanner.Scan() {
			n++
			t := strings.TrimSpace(scanner.Text())
			if strings.Contains(t, "IPublishContext") {
				inPub = true
				level = 0
				continue
			}
			if !inPub {
				continue
			}
			switch t {
			case "{":
				level++
				continue
			case "}":
				level--
				inPub = level > 0
				continue
			}
			if isSuspicious(t) {
				fmt.Printf("suspicious: %s:%d: %s\n", info.Name(), n, t)
			}
		}
		return scanner.Err()
	})
	if err != nil {
		log.Fatal(err)
	}
}
