package lipsum

import (
	"math/rand"
	"regexp"
	"strings"
	"sync/atomic"

	"github.com/daichi-m/lipsum/assets"
)

/*Generator is the primary interface that exposes the different generators for word, sentence and
paragraph.*/
type Generator interface {

	// Word returns the next word from the lorem ipsum dictionary.
	Word() string

	// WordN returns the next word padded/trimmed to length n from the lorem ipsum dictionary
	WordN(n int) string

	/* Sentence returns a random length sentence (collection of words separated by a single space).
	Each word is chosen from the lorem ipsum dictionary. */
	Sentence() string

	// SentenceN returns a sentence consisting of n words from the lorem ipsum dictionary.
	SentenceN(n int) string

	/* Paragraph returns a random length paragraph (collection of sentences separated by the
	period).*/
	Paragraph() string

	// ParagraphN returns a paragraph of n sentences.
	ParagraphN(n int) string
}

const (
	space   = " "
	pad     = "x"
	period  = "."
	newline = "\n"
)

/* lipsum is a concrete implementation of Generator. It uses an atomic counter to
read the dictionary, hence it is thread safe */
type lipsum struct {
	dictionary []string
	idxCtr     int64
	dictLen    int
}

func (l *lipsum) updateIdx() int {
	idx := atomic.AddInt64(&l.idxCtr, 1)
	idx = idx % int64(l.dictLen)
	return int(idx)
}

func (l *lipsum) Word() string {
	idx := l.updateIdx()
	word := l.dictionary[idx]
	return word
}

func (l *lipsum) WordN(n int) string {
	word := l.Word()

	if len(word) < n {
		rem := n - len(word)
		word = word + strings.Repeat(pad, rem)
	} else if len(word) > n {
		word = word[0:n]
	}
	return word
}

func (l *lipsum) Sentence() string {
	n := rand.Intn(10) + 10
	return l.SentenceN(n)
}

func (l *lipsum) SentenceN(n int) string {
	words := make([]string, 0, n)
	for i := 0; i < n; i++ {
		words = append(words, l.Word())
	}
	return strings.Join(words, space)
}

func (l *lipsum) Paragraph() string {
	n := rand.Intn(10) + 5
	return l.ParagraphN(n)
}

func (l *lipsum) ParagraphN(n int) string {
	sents := make([]string, 0, n)
	for i := 0; i < n; i++ {
		sents = append(sents, l.Sentence())
	}
	return strings.Join(sents, period)
}

/*NewGenerator returns a new instance of Generator. This Generator instance's Thread-Safety
is controlled via the threadSafe parameter. */
func NewGenerator(threadSafe bool) (Generator, error) {
	wordsStr, err := assets.Asset("resources/lorem.txt")
	if err != nil {
		return nil, err
	}
	regex, err := regexp.Compile("\\s+")
	if err != nil {
		return nil, err
	}
	dict := regex.Split(string(wordsStr), -1)
	return &lipsum{
		dictionary: dict,
		idxCtr:     int64(-1),
		dictLen:    len(dict),
	}, nil
}

var _ Generator = &lipsum{}
