package searchengine

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

func Get(opts ...*FindOptions) (res []map[string]interface{}) {
	startTime := time.Now()
	ch := make(chan FindOptions, 1)
	go ApplyFindOptions(opts, ch)
	options := <-ch
	total := len(model.order)
	if options.Limit <= 0 {
		return
	}
	end := total - options.Start
	start := end - options.Limit
	if end > total {
		end = total
	}
	if start < 0 {
		start = 0
	}
	seen := make(map[string]struct{})
	totalRecords := len(model.store)
	for i := end - 1; i >= start; i-- {
		docID := model.order[i]
		if _, exists := seen[docID]; exists {
			continue
		}
		if doc, ok := model.store[docID]; ok {
			res = append(res, doc)
			seen[docID] = struct{}{}
		}
	}
	meta := map[string]interface{}{
		"matchedRecords":  totalRecords,
		"returnedRecords": options.Limit,
		"totalRecords":    totalRecords,
		"totalTime":       time.Since(startTime).String(),
	}
	res = append([]map[string]interface{}{{"meta": meta}}, res...)
	return
}

func GetDoc(keyword string, opts ...*FindOptions) (res []map[string]interface{}) {
	startTime := time.Now()
	ch := make(chan FindOptions, 1)
	go ApplyFindOptions(opts, ch)
	options := <-ch

	total := len(model.order)
	if options.Limit <= 0 {
		return
	}
	keyword = strings.ToLower(keyword)
	matched := 0
	seen := make(map[string]struct{})
	for i := len(model.order) - 1; i >= 0; i-- {
		docID := model.order[i]
		if _, exists := seen[docID]; exists {
			continue
		}
		if doc, ok := model.store[docID]; ok {
			for _, value := range doc {
				strValue := fmt.Sprintf("%v", value)
				if strings.Contains(strings.ToLower(strValue), keyword) {
					matched++
					seen[docID] = struct{}{}
					break
				}
			}
		}
	}
	end := total - options.Start
	start := end - options.Limit
	if end > total {
		end = total
	}
	if start < 0 {
		start = 0
	}
	seen = make(map[string]struct{}) // Reset for slice
	for i := end - 1; i >= start; i-- {
		docID := model.order[i]
		if _, exists := seen[docID]; exists {
			continue
		}
		if doc, ok := model.store[docID]; ok {
			for _, value := range doc {
				strValue := fmt.Sprintf("%v", value)
				if strings.Contains(strings.ToLower(strValue), keyword) {
					res = append(res, doc)
					seen[docID] = struct{}{}
					break
				}
			}
		}
	}
	meta := map[string]interface{}{
		"totalRecords":    total,
		"matchedRecords":  matched,
		"returnedRecords": len(res),
		"totalTime":       time.Since(startTime).Nanoseconds(),
	}
	res = append([]map[string]interface{}{{"meta": meta}}, res...)
	return
}

func GetStats() (res int) {
	model.mu.Lock()
	defer model.mu.Unlock()
	res = len(model.store)
	return
}

func Find(filter map[string]interface{}, opts ...*FindOptions) (res []map[string]interface{}, err error) {
	startTime := time.Now()

	if len(model.index) == 0 {
		return nil, nil
	}

	ch := make(chan FindOptions)
	go ApplyFindOptions(opts, ch)
	options := <-ch

	var intersected map[string]struct{}
	for _, value := range filter {
		strVal, ok := value.(string)
		if !ok {
			err = fmt.Errorf("non-string filter value: %v", value)
			return
		}
		tokens := stemTokens(removeStopWords(tokenize(normalize(strVal))))
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

	if intersected == nil {
		return []map[string]interface{}{
			{
				"meta": map[string]interface{}{
					"matchedRecords": 0,
					"totalTime":      time.Since(startTime).String(),
				},
			},
		}, nil
	}

	// Collect matched docIDs into slice
	docIDs := make([]string, 0, len(intersected))
	for id := range intersected {
		docIDs = append(docIDs, id)
	}

	// Optional: Sort by Timestamp or another field if needed
	sort.Slice(docIDs, func(i, j int) bool {
		di, dj := model.store[docIDs[i]], model.store[docIDs[j]]
		ti, tok := di["Timestamp"].(float64)
		tj, tjok := dj["Timestamp"].(float64)
		if tok && tjok {
			if options.Sort == ASCENDING {
				return ti < tj
			}
			return ti > tj
		}
		return docIDs[i] < docIDs[j] // fallback
	})

	// Apply pagination
	start := options.Start
	if start >= len(docIDs) {
		return []map[string]interface{}{
			{
				"meta": map[string]interface{}{
					"matchedRecords": 0,
					"totalTime":      time.Since(startTime).String(),
				},
			},
		}, nil
	}
	end := start + options.Limit
	if end > len(docIDs) {
		end = len(docIDs)
	}

	for _, docID := range docIDs[start:end] {
		if doc, ok := model.store[docID]; ok {
			res = append(res, doc)
		}
	}

	meta := map[string]interface{}{
		"matchedRecords": len(docIDs),
		"totalTime":      time.Since(startTime).String(),
	}
	res = append([]map[string]interface{}{{"meta": meta}}, res...)
	return
}

func ApplyFindOptions(opts []*FindOptions, resultChan chan FindOptions) {
	var finalOptions FindOptions

	if len(opts) > 0 && opts[0] != nil {
		finalOptions = *opts[0]

		if finalOptions.MatchType == "" {
			finalOptions.MatchType = DefaultOptions.MatchType
		}
		if finalOptions.Limit == 0 {
			finalOptions.Limit = DefaultOptions.Limit
		}
		// if Start is negative, reset to 0
		if finalOptions.Start < 0 {
			finalOptions.Start = DefaultOptions.Start
		}
		if finalOptions.Sort == "" {
			finalOptions.Sort = DefaultOptions.Sort
		}
	} else {
		finalOptions = DefaultOptions
	}

	resultChan <- finalOptions
}
