package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

type WordStat struct {
	Word  string
	Count int
}

type TextStats struct {
	charCount  int
	wordCount  int
	lineCount  int
	spaceCount int
	wordFreq   map[string]int
}

func isLetterOrDigit(r rune) bool {
	if r >= '0' && r <= '9' {
		return true
	}
	if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
		return true
	}
	if (r >= 'А' && r <= 'Я') || (r >= 'а' && r <= 'я') || r == 'Ё' || r == 'ё' {
		return true
	}
	return false
}

func analyzeFile(filename string) (TextStats, error) {
	file, err := os.Open(filename)
	if err != nil {
		return TextStats{}, err
	}
	defer file.Close()

	stats := TextStats{}
	stats.wordFreq = make(map[string]int)

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		word := scanner.Text()
		stats.wordCount++

		cleanWord := strings.ToLower(word)
		cleanWord = strings.TrimFunc(cleanWord, func(r rune) bool {
			return !isLetterOrDigit(r)
		})

		if len(cleanWord) > 0 {
			stats.wordFreq[cleanWord]++
		}
	}

	if err := scanner.Err(); err != nil {
		return TextStats{}, err
	}

	file.Seek(0, 0)
	lineScanner := bufio.NewScanner(file)
	for lineScanner.Scan() {
		stats.lineCount++
		text := lineScanner.Text()
		stats.charCount += len(text)

		for _, char := range text {
			if char == ' ' || char == '\t' || char == '\n' || char == '\r' || char == '\v' || char == '\f' {
				stats.spaceCount++
			}
		}
	}
	stats.charCount += stats.lineCount

	if err := lineScanner.Err(); err != nil {
		return TextStats{}, err
	}

	return stats, nil
}

func getTopWords(freqMap map[string]int, topN int) []WordStat {
	var wordStats []WordStat
	for word, count := range freqMap {
		wordStats = append(wordStats, WordStat{Word: word, Count: count})
	}

	sort.Slice(wordStats, func(i, j int) bool {
		return wordStats[i].Count > wordStats[j].Count
	})

	if len(wordStats) < topN {
		return wordStats
	}
	return wordStats[:topN]
}

func printStats(stats TextStats, filename string) {
	fmt.Printf("Статистика по файлу %s\n", filename)
	fmt.Printf("Количество символов: %d\n", stats.charCount)
	fmt.Printf("Количество строк: %d\n", stats.lineCount)
	fmt.Printf("Количество слов: %d\n", stats.wordCount)
	fmt.Printf("Количество пробельных символов: %d\n", stats.spaceCount)
}

func main() {
	topN := flag.Int("top", 5, "Количество слов для вывода в топе")
	linesOnly := flag.Bool("l", false, "Выводить только количество строк")

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Ошибка: не указано имя файла.")
		fmt.Println("Использование: go run main.go [флаги] <Имя файла>")
		fmt.Println("Пример: go run main.go -top=10 -l test.txt")
		os.Exit(1)
	}
	filename := args[0]

	stats, err := analyzeFile(filename)
	if err != nil {
		fmt.Println("Ошибка анализа файла:", err)
		os.Exit(1)
	}

	if *linesOnly {
		fmt.Printf("Количество строк в файле %s: %d\n", filename, stats.lineCount)
	} else {
		printStats(stats, filename)

		topWords := getTopWords(stats.wordFreq, *topN)
		fmt.Printf("Топ %d слов\n", *topN)
		for i, stat := range topWords {
			fmt.Printf("%d. %s: %d\n", i+1, stat.Word, stat.Count)
		}
	}
}
