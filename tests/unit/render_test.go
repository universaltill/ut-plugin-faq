package unit

import (
	"testing"

	faq "ut-plugin-faq/src/faq"
)

func TestSearchEntriesReturnsRelevant(t *testing.T) {
	entries := []faq.Entry{
		{ID: "1", Category: "general", Question: "How to print?", Answer: "Use printer", Keywords: []string{"print", "receipt"}},
		{ID: "2", Category: "payments", Question: "Refund", Answer: "Process refund", Keywords: []string{"refund"}},
	}
	results := faq.SearchEntries(entries, "", "print")
	if len(results) != 1 || results[0].ID != "1" {
		t.Fatalf("expected one relevant result for print")
	}
}

func TestSearchEntriesFiltersByCategory(t *testing.T) {
	entries := []faq.Entry{{ID: "1", Category: "general", Question: "A", Answer: "B"}}
	results := faq.SearchEntries(entries, "payments", "")
	if len(results) != 0 {
		t.Fatalf("expected no results for non-matching category")
	}
}

func TestSearchEntriesReturnsAllWhenNoQuery(t *testing.T) {
	entries := []faq.Entry{{ID: "1", Category: "general", Question: "A", Answer: "B"}}
	results := faq.SearchEntries(entries, "", "")
	if len(results) != 1 {
		t.Fatalf("expected all entries when no query")
	}
}
