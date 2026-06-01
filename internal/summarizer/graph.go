package summarizer

import "github.com/didasy/tldr"

// Graph implements Summarizer using the PageRank-based graph algorithm
// from github.com/didasy/tldr (didasy/tldr v0.7.0).
//
// Note: a new *tldr.Bag is created per call. The Bag type is not thread-safe;
// do not share instances across goroutines.
type Graph struct{}

func (g *Graph) Summarize(text string, n int) ([]string, error) {
	// Honor the Summarizer contract at the edges before delegating: empty input
	// and non-positive n return nil; n is capped to the sentence count so
	// "n > count ⇒ all sentences" holds, matching the other algorithms.
	sentences := TokenizeSentences(text)
	if len(sentences) == 0 || n < 1 {
		return nil, nil
	}
	if n > len(sentences) {
		n = len(sentences)
	}
	bag := tldr.New()
	return bag.Summarize(text, n)
}
