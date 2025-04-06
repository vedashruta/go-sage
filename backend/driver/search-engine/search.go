package searchengine

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// Get retrieves a list of documents based on the provided FindOptions.
// It applies pagination and sorting, then returns the documents along with meta-information about the query.
func Get(opts ...*FindOptions) (res []map[string]interface{}) {
	// Records the start time for the operation.
	startTime := time.Now()
	// Create a channel to apply FindOptions asynchronously.
	ch := make(chan FindOptions, 1)
	go ApplyFindOptions(opts, ch)
	options := <-ch

	// Retrieve the total number of documents in the model.
	total := len(model.order)
	// If the Limit is less than or equal to 0, return an empty result.
	if options.Limit <= 0 {
		return
	}

	// Calculate the start and end indices based on pagination options.
	end := total - options.Start
	start := end - options.Limit
	if end > total {
		end = total
	}
	if start < 0 {
		start = 0
	}

	// Map to track which documents have already been seen.
	seen := make(map[string]struct{})
	// Total number of documents in the store.
	totalRecords := len(model.store)

	// Iterate over the documents in the specified range and collect the results.
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

	// Meta information about the query.
	meta := map[string]interface{}{
		"matchedRecords":  totalRecords,
		"returnedRecords": options.Limit,
		"totalRecords":    totalRecords,
		"totalTime":       time.Since(startTime).String(),
	}

	// Prepend the meta information to the results.
	res = append([]map[string]interface{}{{"meta": meta}}, res...)
	return
}

// GetDoc retrieves documents that contain the given keyword in any of their fields.
// It applies pagination and returns documents that match the keyword.
func GetDoc(keyword string, opts ...*FindOptions) (res []map[string]interface{}) {
	// Records the start time for the operation.
	startTime := time.Now()
	// Create a channel to apply FindOptions asynchronously.
	ch := make(chan FindOptions, 1)
	go ApplyFindOptions(opts, ch)
	options := <-ch

	// Retrieve the total number of documents in the model.
	total := len(model.order)
	// If the Limit is less than or equal to 0, return an empty result.
	if options.Limit <= 0 {
		return
	}

	// Normalize and convert the keyword to lowercase.
	keyword = strings.ToLower(keyword)

	// Variable to track matched documents.
	matched := 0
	// Map to track which documents have already been seen.
	seen := make(map[string]struct{})

	// Iterate over the documents and count how many match the keyword.
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

	// Calculate the start and end indices for pagination.
	end := total - options.Start
	start := end - options.Limit
	if end > total {
		end = total
	}
	if start < 0 {
		start = 0
	}

	// Reset the seen map and iterate over the documents to return the results.
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

	// Meta information about the query.
	meta := map[string]interface{}{
		"totalRecords":    total,
		"matchedRecords":  matched,
		"returnedRecords": len(res),
		"totalTime":       time.Since(startTime).Nanoseconds(),
	}

	// Prepend the meta information to the results.
	res = append([]map[string]interface{}{{"meta": meta}}, res...)
	return
}

// GetStats returns the total number of documents stored in the model.
func GetStats() (res int) {
	model.mu.Lock()
	defer model.mu.Unlock()
	// Return the number of documents in the store.
	res = len(model.store)
	return
}

// Find performs a search based on the provided filter, applying tokenization, stop-word removal, and stemming.
// It uses the provided FindOptions for pagination, sorting, and matching logic.
func Find(filter map[string]interface{}, opts ...*FindOptions) (res []map[string]interface{}, err error) {
	// Records the start time for the operation.
	startTime := time.Now()

	// If the index is empty, return no results.
	if len(model.index) == 0 {
		return nil, nil
	}

	// Create a channel to apply FindOptions asynchronously.
	ch := make(chan FindOptions)
	go ApplyFindOptions(opts, ch)
	options := <-ch

	// Variable to store the set of intersected document IDs.
	var intersected map[string]struct{}
	// Apply filters to the model's index.
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

	// If no documents match the filter, return no results.
	if intersected == nil {
		return []map[string]interface{}{
			{
				"meta": map[string]interface{}{
					"matchedRecords":  0,
					"totalRecords":    0,
					"returnedRecords": 0,
					"totalTime":       time.Since(startTime).String(),
				},
			},
		}, nil
	}

	// Collect matched docIDs into a slice.
	docIDs := make([]string, 0, len(intersected))
	for id := range intersected {
		docIDs = append(docIDs, id)
	}

	// Optionally sort by Timestamp or another field.
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
		return docIDs[i] < docIDs[j] // Fallback to default sorting by docID
	})

	// Apply pagination to the results.
	start := options.Start
	if start >= len(docIDs) {
		return []map[string]interface{}{
			{
				"meta": map[string]interface{}{
					"matchedRecords":  0,
					"totalRecords":    len(model.store),
					"returnedRecords": 0,
					"totalTime":       time.Since(startTime).String(),
				},
			},
		}, nil
	}
	end := start + options.Limit
	if end > len(docIDs) {
		end = len(docIDs)
	}

	// Retrieve and return the results within the specified range.
	for _, docID := range docIDs[start:end] {
		if doc, ok := model.store[docID]; ok {
			res = append(res, doc)
		}
	}

	// Meta information about the query.
	meta := map[string]interface{}{
		"matchedRecords":  len(docIDs),                    // Total number of documents that matched the filter
		"totalRecords":    len(model.store),               // Total number of records in the model
		"returnedRecords": len(res),                       // Number of documents returned after pagination
		"totalTime":       time.Since(startTime).String(), // Total time taken for the query
	}
	// Ensure meta doesn't affect the paginated result count
	res = append([]map[string]interface{}{{"meta": meta}}, res...)
	return
}

// ApplyFindOptions applies default or provided FindOptions to the result channel.
func ApplyFindOptions(opts []*FindOptions, resultChan chan FindOptions) {
	var finalOptions FindOptions

	// If options are provided, apply them. Otherwise, use default options.
	if len(opts) > 0 && opts[0] != nil {
		finalOptions = *opts[0]

		// Apply defaults if specific options are missing.
		if finalOptions.MatchType == "" {
			finalOptions.MatchType = DefaultOptions.MatchType
		}
		if finalOptions.Limit == 0 {
			finalOptions.Limit = DefaultOptions.Limit
		}
		if finalOptions.Start < 0 {
			finalOptions.Start = DefaultOptions.Start
		}
		if finalOptions.Sort == "" {
			finalOptions.Sort = DefaultOptions.Sort
		}
	} else {
		finalOptions = DefaultOptions
	}

	// Send the final options to the channel.
	resultChan <- finalOptions
}
