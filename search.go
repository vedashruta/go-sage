package gosage

import (
	"fmt"

	"github.com/vedashruta/go-sage.git/services"
)

func Get(opts ...*FindOptions) (res []map[string]interface{}) {
	ch := make(chan FindOptions, 1)
	go applyFindOptions(opts, ch)
	options := <-ch
	total := len(model.order)
	start := total - options.Limit
	if start < 0 {
		start = 0
	}
	seen := make(map[string]struct{})
	for i := total - 1; i >= start; i-- {
		docID := model.order[i]
		if _, exists := seen[docID]; exists {
			continue
		}
		if doc, ok := model.store[docID]; ok {
			res = append(res, doc)
			seen[docID] = struct{}{}
		}
	}
	return
}

func GetStats() (res int) {
	model.mu.Lock()
	defer model.mu.Unlock()
	res = len(model.store)
	return
}

func Find(filter map[string]interface{}, opts ...*FindOptions) (res []map[string]interface{}, err error) {
	if len(model.index) == 0 {
		return nil, nil
	}
	ch := make(chan FindOptions)
	go applyFindOptions(opts, ch)
	var intersected map[string]struct{}
	for _, value := range filter {
		strVal, ok := value.(string)
		if !ok {
			err = fmt.Errorf("non-string filter value: %v", value)
			return
		}
		tokens := services.StemTokens(services.RemoveStopWords(services.Tokenize(services.Normalize(strVal))))
		if len(tokens) == 0 {
			continue
		}
		currentSet := make(map[string]struct{})
		for _, term := range tokens {
			for _, docID := range model.index[term] {
				currentSet[docID] = struct{}{}
			}
		}
		if intersected == nil {
			intersected = currentSet
		} else {
			for id := range intersected {
				if _, ok := currentSet[id]; !ok {
					delete(intersected, id)
				}
			}
		}
	}
	count := 0
	options := <-ch
	for docID := range intersected {
		if doc, ok := model.store[docID]; ok {
			res = append(res, doc)
			count++
			if options.Limit > 0 && count >= options.Limit {
				break
			}
		}
	}
	return
}

func applyFindOptions(opts []*FindOptions, resultChan chan FindOptions) {
	var finalOptions FindOptions
	if len(opts) > 0 && opts[0] != nil {
		finalOptions = *opts[0]
		if finalOptions.MatchType == "" {
			finalOptions.MatchType = DefaultOptions.MatchType
		}
		if finalOptions.Limit == 0 {
			finalOptions.Limit = DefaultOptions.Limit
		}
	} else {
		finalOptions = DefaultOptions
	}
	resultChan <- finalOptions
}
