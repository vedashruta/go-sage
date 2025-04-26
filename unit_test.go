package gosage

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/vedashruta/go-sage.git/services"
)

func TestNormalize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello, World!", "hello world"},
		{"It's 2025!", "its 2025"},
		{"ALL CAPS!!", "all caps"},
		{"Go-lang.", "golang"},
	}

	for _, test := range tests {
		result := services.Normalize(test.input)
		if result != test.expected {
			t.Fatal()
		}
	}
}

func TestTokenize(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"hello world", []string{"hello", "world"}},
		{"one two  three", []string{"one", "two", "three"}},
		{"   trim   space   ", []string{"trim", "space"}},
		{"", []string{}},
	}

	for _, test := range tests {
		result := services.Tokenize(test.input)
		if len(result) != len(test.expected) {
			t.Fatal()
		}
		for i := range result {
			if result[i] != test.expected[i] {
				t.Fatal()
			}
		}
	}
}

func TestRemoveStopWords(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "No stop words",
			input:    []string{"hello", "world"},
			expected: []string{"hello", "world"},
		},
		{
			name:     "Only stop words",
			input:    []string{"the", "is", "and"},
			expected: []string{},
		},
		{
			name:     "Mixed tokens",
			input:    []string{"the", "quick", "brown", "fox"},
			expected: []string{"quick", "brown", "fox"},
		},
		{
			name:     "Empty input",
			input:    []string{},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := services.RemoveStopWords(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Fatal()
			}
		})
	}
}

func TestStemTokens(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "Basic stemming",
			input:    []string{"running", "jumps", "lazy", "easily"},
			expected: []string{"run", "jump", "lazi", "easili"},
		},
		{
			name:     "Empty input",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "No change after stemming",
			input:    []string{"dog", "cat"},
			expected: []string{"dog", "cat"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := services.StemTokens(tt.input)
			if !reflect.DeepEqual(output, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, output)
			}
		})
	}
}

func TestIndex(t *testing.T) {
	doc := map[string]interface{}{
		"title": "The quick brown fox jumps over the lazy dog",
		"meta":  "In the forest",
	}
	docID := uuid.NewString()
	err := Index(doc, docID)
	if err != nil {
		t.Fatal()
	}
	expectedIndex := map[string][]string{
		"quick":  {docID},
		"brown":  {docID},
		"fox":    {docID},
		"jump":   {docID},
		"over":   {docID},
		"lazi":   {docID},
		"dog":    {docID},
		"forest": {docID},
	}
	if !reflect.DeepEqual(model.index, expectedIndex) {
		t.Errorf("Expected index %v, got %v", expectedIndex, model.index)
	}
	expectedDoc := doc
	storedDoc, exists := model.store[docID]
	if !exists {
		t.Fatalf("Document with ID %s not found in docStore", docID)
	}
	if !reflect.DeepEqual(storedDoc, expectedDoc) {
		t.Errorf("Expected stored document %v, got %v", expectedDoc, storedDoc)
	}
}

func TestSearch(t *testing.T) {
	docs := []map[string]interface{}{
		{
			"title": "The quick brown fox jumps over the lazy dog",
			"meta":  "In the forest",
		},
		{
			"title": "The quick white bird flies over the lazy fox",
			"meta":  "In the sky",
		},
		{
			"title": "An agile fox sprints in the daylight",
			"meta":  "In the wild",
		},
	}
	for _, doc := range docs {
		docID := uuid.NewString()
		err := Index(doc, docID)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Test case 1: Filter for "lazy fox" — should return only doc[1]
	filter := map[string]interface{}{
		"title": "bird",
	}
	result, err := Find(filter)
	if err != nil {
		t.Fatal(err)
	}
	expected := []map[string]interface{}{docs[1]}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Search failed.\nExpected: %v\nGot: %v", expected, result)
	}

	// Test case 2: Filter for "fox" — should return all 3
	filter2 := map[string]interface{}{
		"title": "fox",
	}
	result2, err := Find(filter2)
	if err != nil {
		t.Fatal(err)
	}
	expected2 := []map[string]interface{}{docs[0], docs[1], docs[2]}
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("Search failed.\nExpected: %v\nGot: %v", expected2, result2)
	}

	// Test case 3: Filter for "daylight" — should return only doc[2]
	filter3 := map[string]interface{}{
		"title": "daylight",
	}
	result3, err := Find(filter3)
	if err != nil {
		t.Fatal(err)
	}
	expected3 := []map[string]interface{}{docs[2]}
	if !reflect.DeepEqual(result3, expected3) {
		t.Errorf("Search failed.\nExpected: %v\nGot: %v", expected3, result3)
	}

	// Test case 4: Filter for "fox" — should return all 3, limit the results to 1
	filter4 := map[string]interface{}{
		"title": "fox",
	}
	opts := FindOptions{
		Limit: 1,
	}
	result4, err := Find(filter4, &opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(result4) != 1 {
		t.Fatal()
	}
}

func TestGet(t *testing.T) {
	for i := 1; i <= 5; i++ {
		docID := "doc" + string(rune('A'+i-1))
		doc := map[string]interface{}{
			"id":   docID,
			"name": "Doc " + docID,
		}
		Index(doc, docID)
	}

	// Limit = 3, should return latest 3 in reverse order
	limit := 3
	opts := &FindOptions{Limit: limit}
	got := Get(opts)
	want := []map[string]interface{}{
		model.store["docE"],
		model.store["docD"],
		model.store["docC"],
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Expected latest %d documents: %+v\nGot: %+v", limit, want, got)
	}
}
