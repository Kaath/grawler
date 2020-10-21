package grawler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync"
)

const (
	FOLLOW_THROUGH int = iota
	SKIP           int = iota
)

var (
	MAX_FOLLOW_THROUGH int    = -1 // -1 for all
	MAX_DEPTH          int    = 3  // negative for infinite crawling (not recommended)
	REPOSITORY_PATH    string = "./repository/"
	mutex              sync.Mutex
	logger             *log.Logger = log.New(os.Stdout, "[GRAWL] ", log.Lshortfile)
	wg                 sync.WaitGroup
	ACCEPT_ALL         Policy = Policy{regexs: []*regexp.Regexp{regexp.MustCompile("https?://\"?.*?(>|\"| |\\\\)")}, follow_trough: FOLLOW_THROUGH}
	REJECT_ALL         Policy = Policy{regexs: []*regexp.Regexp{regexp.MustCompile("https?://\"?.*?(>|\"| |\\\\)")}, follow_trough: SKIP}
)

type Policy struct {
	regexs         []*regexp.Regexp
	follow_trough int
}

type SafeCounter struct {
	count int
	mutex sync.Mutex
}

func (s *SafeCounter) SafeCount() int {
	s.mutex.Lock()
	count := s.count
	s.mutex.Unlock()
	return count
}

func (s *SafeCounter) SafeInc() {
	s.mutex.Lock()
	s.count++
	s.mutex.Unlock()
}

func find_urls(page *Page, pol Policy) []string {
	var res []string

	for _, pol := range pol.regexs {
		res = append(res, (*pol).FindAllString(string(page.body), MAX_FOLLOW_THROUGH)...)
		for index, str := range res {
			res[index] = str[:len(str)-1]
		}
	}

	return res
}

type Crawler struct {
	Starters    []string
	Treatments []SaveFunc

	counter  SafeCounter
	LogLevel int

	// Followthrough policies
	StartPolicies []Policy
	NodePolicies  []Policy
	LeafPolicies  []Policy
	DefaultPolicy Policy
}

func New(starts []string, treatments ...SaveFunc) *Crawler {
	return &Crawler{
		Starters: starts,
		Treatments: treatments,
		counter: SafeCounter{ count: 0 },
		LogLevel: 0,
	}
}

func (c *Crawler) crawl(url string, depth int, treatments []SaveFunc) {
	c.counter.SafeInc()
	str := fmt.Sprintf("Crawling: %s | count: %d", url, c.counter.SafeCount())
	if c.LogLevel > 2 {
		logger.Println(str)
	}

	resp, err := http.Get(url)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	p := Page{url, body}
	if depth > 0 {
		strings := find_urls(&p, c.DefaultPolicy)
		for i := 0; i < len(strings); i++ {
			wg.Add(1)
			go func(s string) {
				c.crawl(s, depth-1, treatments)
				wg.Done()
			}(strings[i])
		}
	}

	for _, t := range treatments {
		mutex.Lock()
		t(&p)
		mutex.Unlock()
	}
}

func (c *Crawler) StartCrawl() {
	for _, str := range c.Starters {
		wg.Add(1)
		go func(s string) {
			c.crawl(s, MAX_DEPTH, c.Treatments)
			wg.Done()
		}(str)
	}
	wg.Wait()
}
