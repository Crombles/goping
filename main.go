package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func readInput(scanner *bufio.Scanner) []string {
	var lines []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			break
		}

		lines = append(lines, line)
	}

	if len(lines) == 0 {
		fmt.Println("[!] Error: Why are you giving me an empty input? Provide valid configs.")
		return nil
	}

	return lines
}

func printResults(results []ServerResult, isSubscription bool) {
	if len(results) == 0 {
		fmt.Println("[!] No valid configs found to check.")
		fmt.Println("__________________________________________________________")
		return
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].IsAlive != results[j].IsAlive {
			return results[i].IsAlive
		}
		return results[i].Latency < results[j].Latency
	})

	var deadCount int

	for _, res := range results {
		if res.IsAlive {
			fmt.Printf("[+] Server: %s | Status: ALIVE | Latency: %v\n", res.Name, res.Latency)
		} else {
			fmt.Printf("[-] Server: %s | Status: DEAD\n", res.Name)
			deadCount++
		}
	}
	
	if isSubscription {
		fmt.Printf("\n[i] Total configs: %d | Dead: %d | Alive: %d\n", len(results), deadCount, len(results)-deadCount)
	}

	fmt.Println("__________________________________________________________")
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n& Action Menu &")
		fmt.Println("[1] - Exit")
		fmt.Println("[2] - Check configs/subscription")
		fmt.Print("Your choice: ")

		if !scanner.Scan() {
			break
		}

		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "2":
			fmt.Print("Paste your config/subscription and press Enter twice: ")

			inputLines := readInput(scanner)
			if inputLines == nil {
				continue
			}

			for _, item := range inputLines {
				var configs []string
				isSubscription := false

				if strings.HasPrefix(item, "http://") || strings.HasPrefix(item, "https://") {
					isSubscription = true
					fetched, err := fetchSubscription(item)
					if err != nil {
						fmt.Println("[!] Failed to fetch subscription data:", err)
						fmt.Println("__________________________________________________________")
						continue
					}
					configs = fetched         
				} else {
					configs = append(configs, item)
				}

				results := processParallel(configs, 10)
				printResults(results, isSubscription)
			}

		case "1":
			fmt.Println("Shutting down.")
			return

		default:
			fmt.Println("Invalid choice, please try again.")
		}
	}
}
