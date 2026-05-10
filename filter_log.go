package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: filter_log <input> <output> <correct_code>")
		os.Exit(1)
	}
	inputPath := os.Args[1]
	outputPath := os.Args[2]
	correctCode := os.Args[3]

	in, err := os.Open(inputPath)
	if err != nil {
		panic(err)
	}
	defer in.Close()

	out, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	writer := bufio.NewWriter(out)
	defer writer.Flush()

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) != 2 {
			continue
		}
		if parts[0] != correctCode {
			fmt.Fprintf(writer, "%s\n", line)
		}
	}
}
