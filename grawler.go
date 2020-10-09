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

var (
	MAX_FOLLOW_THROUGH	int = -1   // -1 for all
	MAX_DEPTH			int = 3   // negative for infinite crawling (not recommended)
	REPOSITORY_PATH		string = "./repository/"
	counter				SafeCounter = SafeCounter{ count: 1 }
	mutex				sync.Mutex
	logger				*log.Logger = log.New(os.Stdout, "[GRAWL] ", log.Lshortfile)
	wg sync.WaitGroup
)

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

type SaveFunc func(page *Page)

func find_urls(page string) []string {
	reg := regexp.MustCompile("https?://\"?.*?(>|\"| |\\\\)")
	res := reg.FindAllString(page, MAX_FOLLOW_THROUGH)
	for index, str := range res {
		res[index] = str[:len(str) - 1]
	}

	return res
}

func crawl(url string, depth int, treatments []SaveFunc) {
	counter.SafeInc()
	str := fmt.Sprintf("Crawling: %s | count: %d", url, counter.SafeCount())
	logger.Println(str)
	resp, err := http.Get(url);
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	p := Page{url, body}
	if depth > 0 {
		strings := find_urls(string(body))
		for i := 0; i < len(strings); i++ {
			wg.Add(1)
			go func (s string) {
				crawl(s, depth - 1, treatments)
				wg.Done()
			} (strings[i])
		}
	}

	for _, t := range treatments {
		mutex.Lock()
		t(&p)
		mutex.Unlock()
	}
}

func StartCrawl(starts []string, treatments ...SaveFunc) {
	for _, str := range starts {
		wg.Add(1)
		go func (s string) {
			crawl(s, MAX_DEPTH, treatments)
			wg.Done()
		} (str)
	}
	wg.Wait()
}
