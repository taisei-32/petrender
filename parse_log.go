package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
)

type Result struct {
	Code string
	Name string
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: parse_log <input> <output>")
		os.Exit(1)
	}
	inputPath := os.Args[1]
	outputPath := os.Args[2]

	in, err := os.Open(inputPath)

	re := regexp.MustCompile(`Processed:\s+(\S+)\s+->\s+([0-9]+)`)

	var results []Result

	scanner := bufio.NewScanner(in)

	for scanner.Scan() {
		line := scanner.Text()

		m := re.FindStringSubmatch(line)
		if len(m) != 3 {
			continue
		}

		results = append(results, Result{
			Name: m[1],
			Code: m[2],
		})
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Code < results[j].Code
	})

	out, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	writer := bufio.NewWriter(out)
	defer writer.Flush()

	for _, r := range results {
		fmt.Fprintf(writer, "%s: %s\n", r.Code, r.Name)
	}
}
