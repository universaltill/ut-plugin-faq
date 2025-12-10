package faq

import "strings"

// SearchEntries filters FAQ entries by locale, optional category, and keyword query.
// Query matches question, answer, or keywords (case-insensitive).
func SearchEntries(entries []Entry, category, query string) []Entry {
	q := strings.ToLower(strings.TrimSpace(query))
	cat := strings.TrimSpace(category)

	var result []Entry
	for _, e := range entries {
		if cat != "" && !strings.EqualFold(cat, e.Category) {
			continue
		}
		if q == "" {
			result = append(result, e)
			continue
		}
		if strings.Contains(strings.ToLower(e.Question), q) ||
			strings.Contains(strings.ToLower(e.Answer), q) ||
			containsKeyword(e.Keywords, q) {
			result = append(result, e)
		}
	}
	return result
}

func containsKeyword(keywords []string, q string) bool {
	for _, k := range keywords {
		if strings.Contains(strings.ToLower(k), q) {
			return true
		}
	}
	return false
}
