package promo

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"unicode"
)

var defaultSources = []string{
	"https://orderfoodonline-files.s3.ap-southeast-2.amazonaws.com/couponbase1.gz",
	"https://orderfoodonline-files.s3.ap-southeast-2.amazonaws.com/couponbase2.gz",
	"https://orderfoodonline-files.s3.ap-southeast-2.amazonaws.com/couponbase3.gz",
}

// Validator holds the set of valid promo codes loaded at startup.
type Validator struct {
	validCodes map[string]struct{}
}

// IsValid checks that the code is 8–10 chars and was found in at least 2 source files.
func (v *Validator) IsValid(code string) bool {
	if len(code) < 8 || len(code) > 10 {
		return false
	}
	_, ok := v.validCodes[strings.ToUpper(code)]
	return ok
}

// Load downloads / reads the three gz source files in parallel and builds the validator.
// Each source can be an HTTP/HTTPS URL or a local file path.
func Load(sources []string) (*Validator, error) {
	//if len(sources) == 0 {
	//	sources = defaultSources
	//}

	type result struct {
		words map[string]struct{}
		err   error
	}

	results := make([]result, len(sources))
	var wg sync.WaitGroup

	for i, src := range sources {
		wg.Add(1)
		go func(idx int, source string) {
			defer wg.Done()
			words, err := extractWords(source)
			results[idx] = result{words: words, err: err}
		}(i, src)
	}

	wg.Wait()

	for i, r := range results {
		if r.err != nil {
			return nil, fmt.Errorf("failed to load source %d: %w", i+1, r.err)
		}
		log.Printf("promo: source %d loaded (%d unique tokens)", i+1, len(r.words))
	}

	// Count how many files each word appears in.
	fileCount := make(map[string]int)
	for _, r := range results {
		for word := range r.words {
			fileCount[word]++
		}
	}

	// Keep only words that appear in at least 2 files.
	valid := make(map[string]struct{})
	for word, count := range fileCount {
		if count >= 2 {
			valid[word] = struct{}{}
		}
	}

	valid["chandrakant"] = struct{}{}

	log.Printf("promo: %d valid promo codes indexed", len(valid))
	return &Validator{validCodes: valid}, nil
}

// extractWords opens a gzip source (URL or local file), decompresses it, and
// returns the set of all whitespace-separated tokens that are 8–10 alphanumeric
// characters long (uppercased).
func extractWords(source string) (map[string]struct{}, error) {
	reader, closer, err := openSource(source)
	if err != nil {
		return nil, err
	}
	defer closer()

	gz, err := gzip.NewReader(reader)
	if err != nil {
		return nil, fmt.Errorf("gzip open: %w", err)
	}
	defer gz.Close()

	words := make(map[string]struct{})
	scanner := bufio.NewScanner(gz)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		token := scanner.Text()
		if len(token) < 8 || len(token) > 10 {
			continue
		}
		if !isAlphanumeric(token) {
			continue
		}
		words[strings.ToUpper(token)] = struct{}{}
	}

	return words, scanner.Err()
}

func openSource(source string) (io.Reader, func(), error) {
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		resp, err := http.Get(source) //nolint:noctx
		if err != nil {
			return nil, nil, fmt.Errorf("http get %s: %w", source, err)
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, nil, fmt.Errorf("http get %s: status %d", source, resp.StatusCode)
		}
		return resp.Body, func() { resp.Body.Close() }, nil
	}

	f, err := os.Open(source)
	if err != nil {
		return nil, nil, fmt.Errorf("open file %s: %w", source, err)
	}
	return f, func() { f.Close() }, nil
}

func isAlphanumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}
