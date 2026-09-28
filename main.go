package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func readInput(scanner *bufio.Scanner) string {
	var lines []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			break
		}

		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

func printResults(results []ServerResult) {
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

	fmt.Printf("\n[i] Total configs: %d | Dead: %d | Alive: %d", len(results), deadCount, len(results)-deadCount)
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

			input := readInput(scanner)
			inputLines := strings.Fields(input)

			for _, item := range inputLines {
				var configs []string

				if strings.HasPrefix(item, "http://") || strings.HasPrefix(item, "https://") {
					fetched, err := fetchSubscription(item)
					if err != nil {
						fmt.Println("[!] Failed to fetch subscription data:", err)
						continue
					}

					for _, cfg := range fetched {
						if cfgClean := strings.TrimSpace(cfg); cfgClean != "" {
							configs = append(configs, cfgClean)
						}
					}
				} else {
					configs = append(configs, item)
				}

				if len(configs) > 0 {
					results := processParallel(configs, 10)
					printResults(results)
				}
			}

		case "1":
			fmt.Println("Shutting down.")
			return

		default:
			fmt.Println("Invalid choice, please try again.")
		}
	}
}
