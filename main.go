package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

const NumberOfWorkers = 3

type Repo struct {
	line  string
	owner string
	name  string
}

func main() {
	client := NewGithubClient()
	lines := client.getReadmeLines(context.Background())
	jobs := make(chan Repo, NumberOfWorkers)
	results := &Results{m: make(map[int][]string, 0)}
	wgRepoFetchers := sync.WaitGroup{}

	fmt.Printf("Got %d lines from README.md\n", len(lines))

	for i := 1; i <= NumberOfWorkers; i++ {
		wgRepoFetchers.Add(1)
		go worker(&wgRepoFetchers, &client, jobs, results)
	}

	//regexp to find useful lines
	re := regexp.MustCompile(`(\* \[.*\]\(https:\/\/github\.com\/([a-zA-Z0-9-_\.]+)\/([a-zA-Z0-9-_\.]+))`)
	for _, line := range lines {
		match := re.FindAllStringSubmatch(line, -1)
		if len(match) == 0 || len(match[0]) < 4 {
			continue
		}
		jobs <- Repo{line: line, owner: match[0][2], name: match[0][3]}
	}
	close(jobs)

	wgRepoFetchers.Wait()
	fmt.Println()

	groupedRepos, sortedKeys := results.groupedRepos()

	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile("stars_numbers_"+time.Now().Format("20060102")+".md", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	for _, k := range sortedKeys {

		if _, err := f.WriteString("\n ### " + k + "\n\n"); err != nil {
			log.Fatal(err)
		}

		for _, s := range groupedRepos[k] {
			if _, err := f.WriteString(strings.TrimSpace(s) + "\n"); err != nil {
				log.Fatal(err)
			}
		}
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

}

// worker reads repos from channel and puts the number of starts in the struct
func worker(wg *sync.WaitGroup, c *Client, repos <-chan Repo, results *Results) {
	for repo := range repos {
		r, err := c.fetchRepository(context.Background(), repo.owner, repo.name)
		if err != nil {
			fmt.Printf("\n%v", err)
			continue
		}

		results.Lock()
		results.m[r.GetStargazersCount()] = append(results.m[r.GetStargazersCount()], repo.line)
		results.fetched++
		fmt.Printf("\033[2K\r%d fetched", results.fetched)
		results.Unlock()
	}
	wg.Done()
}
