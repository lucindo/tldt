package summarizer

// Ensemble combines LexRank and TextRank by averaging per-sentence scores,
// then returning top-N sentences in document order.
type Ensemble struct{}

// Summarize runs LexRank and TextRank independently, averages their per-sentence
// scores, and returns the top n sentences in document order.
// Returns nil, nil for empty input. Caps n to sentence count.
func (e *Ensemble) Summarize(text string, n int) ([]string, error) {
	sentences := TokenizeSentences(text)
	if len(sentences) == 0 {
		return nil, nil
	}
	if n > len(sentences) {
		n = len(sentences)
	}
	lr := lexrankScores(sentences)
	tr := textrankScores(sentences)
	combined := make([]float64, len(sentences))
	for i := range combined {
		combined[i] = (lr[i] + tr[i]) / 2.0
	}
	return selectTopN(combined, n, sentences), nil
}

// lexrankScores returns per-sentence eigenvector centrality scores (LexRank).
func lexrankScores(sentences []string) []float64 {
	return lexrankCompute(sentences, false).scores
}

// textrankScores returns per-sentence PageRank scores (TextRank).
func textrankScores(sentences []string) []float64 {
	return textrankCompute(sentences).scores
}
