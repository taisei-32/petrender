package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func classifyChange(correct, misread string) string {
	if len(correct) != len(misread) || len(correct) != 13 {
		return "その他"
	}

	leftChanges := 0
	rightChanges := 0

	for i := 0; i < 7; i++ {
		if correct[i] != misread[i] {
			leftChanges++
		}
	}
	for i := 7; i < 13; i++ {
		if correct[i] != misread[i] {
			rightChanges++
		}
	}

	if leftChanges == 1 && rightChanges == 0 {
		return "左"
	}
	if leftChanges == 0 && rightChanges == 1 {
		return "右"
	}
	return "その他"
}

func visualize(correct, misread string) string {
	if len(correct) != len(misread) {
		return misread + " (不一致)"
	}
	vis := make([]byte, len(correct))
	change := ""
	for i := 0; i < len(correct); i++ {
		if correct[i] != misread[i] {
			vis[i] = correct[i]
			change = fmt.Sprintf("%c->%c", correct[i], misread[i])
		} else {
			vis[i] = '-'
		}
	}
	return string(vis) + " " + change
}

type Category struct {
	Vis   string
	Lines []string
}

func processFile(filePath string, correctCode string, categories map[string]*Category, seen map[string]bool, categoryOrder *[]string) {
	in, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer in.Close()

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
		code := parts[0]
		if code == correctCode {
			continue
		}

		classify := classifyChange(correctCode, code)
		vis := visualize(correctCode, code)
		key := classify + ":" + vis

		if !seen[key] {
			seen[key] = true
			*categoryOrder = append(*categoryOrder, key)
			categories[key] = &Category{
				Vis: vis,
			}
		}
		categories[key].Lines = append(categories[key].Lines, line)
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: merge_classify <input_dir> <output>")
		os.Exit(1)
	}
	inputDir := os.Args[1]
	outputPath := os.Args[2]

	dir, err := os.ReadDir(inputDir)
	if err != nil {
		panic(err)
	}

	categories := map[string]*Category{}
	categoryOrder := []string{}
	seen := map[string]bool{}

	for _, entry := range dir {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		correctCode := strings.Split(name, "_")[0]
		if len(correctCode) != 13 {
			continue
		}

		processFile(inputDir+"/"+name, correctCode, categories, seen, &categoryOrder)
	}

	out, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	writer := bufio.NewWriter(out)
	defer writer.Flush()

	sorted := []string{}
	others := []string{}
	for _, key := range categoryOrder {
		if strings.HasPrefix(key, "その他") {
			others = append(others, key)
		} else {
			sorted = append(sorted, key)
		}
	}

	sort.Slice(sorted, func(i, j int) bool {
		vi := strings.ReplaceAll(categories[sorted[i]].Vis, "-", "~")
		vj := strings.ReplaceAll(categories[sorted[j]].Vis, "-", "~")
		return vi < vj
	})

	for _, key := range sorted {
		cat := categories[key]
		fmt.Fprintf(writer, "%s\n", cat.Vis)
		for _, line := range cat.Lines {
			fmt.Fprintf(writer, "%s\n", line)
		}
		fmt.Fprintf(writer, "\n")
	}

	if len(others) > 0 {
		fmt.Fprintf(writer, "その他\n")
		for _, key := range others {
			cat := categories[key]
			for _, line := range cat.Lines {
				fmt.Fprintf(writer, "%s\n", line)
			}
		}
	}
}
