// Package summary its exported api for making a summaries of texts
package summary

// Jak to ma działać:
// Trzymamy pliki/teksty źródłowe.
//

type summator struct {
	prompts []string
}

// constructor

func Summator(prompts ...string) summator {
	return summator{
		prompts: prompts,
	}
}

var promptsDefault = []string{
	"",
	"",
}

func DefaultPrompts() *[]string {
	return &promptsDefault
}
