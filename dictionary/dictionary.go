package dictionary

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"gonum.org/v1/gonum/stat/combin"
)

const (
	DictionaryFile = "/usr/share/dict/words"
)

type Dictionary struct {
	Entries map[string][]string
	mutex   sync.Mutex
}

type Anagrams struct {
	results []string
	mutex   sync.Mutex
}

func (res *Anagrams) JSON() []byte {
	resMap := make(map[int][]string)
	for _, word := range res.results {
		resMap[len(word)] = append(resMap[len(word)], word)
	}
	tmp, err := json.Marshal(resMap)
	if err != nil {
		log.Fatalln("Anagrams::JSON error marshalling results", err)
	}
	return tmp
}

func (res *Anagrams) Insert(words ...string) {
	res.mutex.Lock()
	defer res.mutex.Unlock()
	for _, word := range words {
		i := sort.Search(len(res.results), func(i int) bool { return res.results[i] >= word })
		if i >= len(res.results) {
			res.results = append(res.results, word)
		} else if res.results[i] != word {
			res.results = append(res.results[:i+1], res.results[i:]...)
			res.results[i] = word
		}
	}
}

func (dict *Dictionary) GetPartialAnagrams(word string) *Anagrams {
	results := &Anagrams{}
	split := strings.Split(word, "")
	sort.Strings(split)
	SortedWord := strings.Join(split, "")
	for i := 1; i <= len(SortedWord); i++ {
		permutations := combin.Permutations(len(SortedWord), i)
		for _, indexes := range permutations {
			var tmp []byte
			for _, index := range indexes {
				tmp = append(tmp, SortedWord[index])
			}
			results.Insert(dict.GetStraightAnagrams(string(tmp), word).results...)
		}
	}
	return results
}

func (dict *Dictionary) GetStraightAnagrams(word string, exclude ...string) *Anagrams {
	results := &Anagrams{}
	excludeStrings := func() string {
		if len(exclude) == 0 {
			return ""
		}
		return strings.Join(exclude, " ")
	}()
	split := strings.Split(word, "")
	sort.Strings(split)
	SortedWord := strings.Join(split, "")
	if words, found := dict.Entries[SortedWord]; found {
		for _, w := range words {
			if w == word || strings.Contains(excludeStrings, w) {
				continue
			}
			results.Insert(w)
		}
	}
	return results
}

func New() *Dictionary {
	start := time.Now()
	dict := &Dictionary{Entries: make(map[string][]string, 0)}
	f, err := os.Open(DictionaryFile)
	if err != nil {
		log.Fatalln("NewDictionary could not read dictionary file", DictionaryFile, err)
	}
	defer f.Close()
	var wg sync.WaitGroup
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if len(text) <= 1 && text != "a" && text != "i" {
			continue
		}
		wg.Add(1)
		go func(word string) {
			defer wg.Done()
			dict.Insert(word)
		}(text)
	}
	wg.Wait()
	fmt.Println("finished building new Dictionary in", time.Since(start), "found", len(dict.Entries), "SortedWords")
	return dict
}

func (dict *Dictionary) Insert(word string) {
	split := strings.Split(word, "")
	sort.Strings(split)
	SortedWord := strings.Join(split, "")
	dict.mutex.Lock()
	defer dict.mutex.Unlock()
	if IsValidWord(word) {
		words := dict.Entries[SortedWord]
		i := sort.Search(len(words), func(i int) bool { return words[i] >= word })
		if i >= len(words) {
			words = append(words, word)
		} else if words[i] != word {
			words = append(words[:i+1], words[i:]...)
			words[i] = word
		}
		dict.Entries[SortedWord] = words
	}
}

func IsValidChar(c byte) bool {
	return ((c >= 97 && c <= 122) || (c >= 65 && c <= 90))
}

func IsValidWord(word string) bool {
	for _, c := range []byte(word) {
		if !IsValidChar(c) {
			return false
		}
	}
	return true
}
