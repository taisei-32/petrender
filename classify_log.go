package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func visualize(correct, misread string) (string, string) {
	if len(misread) == 0 || len(correct) != len(misread) {
		return misread + " (不一致)", "99:"
	}

	vis := make([]byte, len(correct))
	changes := []string{}
	firstChangePos := -1

	for i := 0; i < len(correct); i++ {
		if correct[i] != misread[i] {
			vis[i] = correct[i]
			changes = append(changes, fmt.Sprintf("%c->%c", correct[i], misread[i]))
			if firstChangePos == -1 {
				firstChangePos = i
			}
		} else {
			vis[i] = '-'
		}
	}

	changeStr := strings.Join(changes, " ")
	sortKey := fmt.Sprintf("%02d:%s", firstChangePos, changeStr)
	return string(vis) + " " + changeStr, sortKey
}

func classifyChange(correct, misread string) string {
	if len(correct) != len(misread) || len(correct) != 13 {
		return "その他"
	}

	changePositions := []int{}
	for i := 0; i < len(correct); i++ {
		if correct[i] != misread[i] {
			changePositions = append(changePositions, i)
		}
	}

	if len(changePositions) == 0 {
		return "その他"
	}

	leftChanges := 0
	rightChanges := 0
	for _, pos := range changePositions {
		if pos < 7 {
			leftChanges++
		} else {
			rightChanges++
		}
	}

	side := ""
	if leftChanges > 0 && rightChanges == 0 {
		side = "左"
	} else if leftChanges == 0 && rightChanges > 0 {
		side = "右"
	} else {
		side = "両方"
	}

	if len(changePositions) == 1 {
		return side + ":1桁"
	}

	isConsecutive := true
	for i := 1; i < len(changePositions); i++ {
		if changePositions[i] != changePositions[i-1]+1 {
			isConsecutive = false
			break
		}
	}

	if isConsecutive {
		return side + ":連続複数桁"
	}
	return side + ":非連続複数桁"
}

type Category struct {
	Vis     string
	SortKey string
	Lines   []string
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
		if len(code) == 0 {
			continue
		}
		if code == correctCode {
			continue
		}

		classify := classifyChange(correctCode, code)
		vis, sortKey := visualize(correctCode, code)
		key := classify + ":" + vis

		if !seen[key] {
			seen[key] = true
			*categoryOrder = append(*categoryOrder, key)
			categories[key] = &Category{
				Vis:     vis,
				SortKey: classify + ":" + sortKey,
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

	order := []string{"左:1桁", "右:1桁", "左:連続複数桁", "右:連続複数桁", "両方:連続複数桁", "左:非連続複数桁", "右:非連続複数桁", "両方:非連続複数桁", "その他"}

	groups := map[string][]string{}
	for _, key := range categoryOrder {
		parts := strings.SplitN(key, ":", 3)
		if len(parts) >= 2 {
			classify := parts[0] + ":" + parts[1]
			groups[classify] = append(groups[classify], key)
		}
	}

	for _, group := range order {
		keys, ok := groups[group]
		if !ok {
			continue
		}

		sort.Slice(keys, func(i, j int) bool {
			return categories[keys[i]].SortKey < categories[keys[j]].SortKey
		})

		if group == "その他" {
			fmt.Fprintf(writer, "その他\n")
			for _, key := range keys {
				cat := categories[key]
				for _, line := range cat.Lines {
					fmt.Fprintf(writer, "%s\n", line)
				}
			}
		} else {
			fmt.Fprintf(writer, "%s\n\n", group)
			for _, key := range keys {
				cat := categories[key]
				fmt.Fprintf(writer, "%s\n", cat.Vis)
				for _, line := range cat.Lines {
					fmt.Fprintf(writer, "%s\n", line)
				}
				fmt.Fprintf(writer, "\n")
			}
		}
	}
}
