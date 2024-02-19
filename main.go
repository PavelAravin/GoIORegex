package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func calculate(expression string) (string, error) {
	re := regexp.MustCompile(`(\d+)\s*([+-])\s*(\d+)=`)
	parts := re.FindStringSubmatch(expression)
	if parts == nil {
		return "", fmt.Errorf("invalid expression: %s", expression)
	}

	left, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid expression: %s", expression)
	}

	sign := parts[2]
	right, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid expression: %s", expression)
	}

	var result int64
	switch sign {
	case "+":
		result = left + right
	case "-":
		result = left - right
	default:
		return "", fmt.Errorf("invalid expression: %s", expression)
	}

	return fmt.Sprintf("%d%s%d=%d", int(left), sign, int(right), int(result)), nil
}

func processFile(input string, output string) error {
	data, err := ioutil.ReadFile(input)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	f, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer f.Close()
	f.Truncate(0)

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.Contains(line, "?") && strings.Count(line, "=") == 1 {
			result, err := calculate(line)
			if err == nil {
				f.WriteString(result + "\n")
			} else {
				f.WriteString(line + " = ?\n")
			}
			f.Sync()
		}
	}

	return scanner.Err()
}

func main() {

	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <input_file> <output_file>")
		os.Exit(1)
	}

	input := os.Args[1]
	output := os.Args[2]

	err := processFile(input, output)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
