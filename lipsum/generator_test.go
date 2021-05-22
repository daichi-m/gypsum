package lipsum

import (
	"strings"
	"testing"
)

var sample string = "Lorem ipsum dolor sit amet consectetuer adipiscing elit Maecenas porttitor congue massa"

var dictionary []string = strings.Split(sample, " ")

func Test_lipsum_Word(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "single", want: "Lorem"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &lipsum{
				dictionary: dictionary,
				idxCtr:     int64(-1),
				dictLen:    len(dictionary),
			}
			if got := l.Word(); got != tt.want {
				t.Errorf("lipsum.Word() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lipsum_WordN(t *testing.T) {
	tests := []struct {
		name string
		N    int
		want string
	}{
		{
			name: "same size",
			N:    5,
			want: "Lorem",
		},
		{
			name: "pad",
			N:    8,
			want: "Loremxxx",
		},
		{
			name: "trim",
			N:    3,
			want: "Lor",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &lipsum{
				dictionary: dictionary,
				idxCtr:     int64(-1),
				dictLen:    len(dictionary),
			}
			if got := l.WordN(tt.N); got != tt.want {
				t.Errorf("lipsum.Word() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mkset(list []string) map[string]bool {
	mp := make(map[string]bool, len(list))
	for _, x := range list {
		mp[x] = true
	}
	return mp
}

func Test_lipsumSentence(t *testing.T) {

	tests := []struct {
		name        string
		expectWords map[string]bool
	}{
		{name: "sentence", expectWords: mkset(dictionary)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &lipsum{
				dictionary: dictionary,
				idxCtr:     int64(-1),
				dictLen:    len(dictionary),
			}
			got := l.Sentence()
			parts := strings.Split(got, " ")
			if len(parts) < 10 {
				t.Errorf("Length of sentence less than 10 words")
			}
			for _, x := range parts {
				if !tt.expectWords[x] {
					t.Errorf("Did not expect word %s", x)
				}
			}
		})

	}
}

func Test_lipsumSentenceN(t *testing.T) {

	tests := []struct {
		name        string
		N           int
		expectCt    int
		expectWords map[string]bool
	}{
		{name: "single", N: 1, expectCt: 1, expectWords: mkset(dictionary)},
		{name: "ten", N: 10, expectCt: 10, expectWords: mkset(dictionary)},
		{name: "hundred", N: 100, expectCt: 100, expectWords: mkset(dictionary)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &lipsum{
				dictionary: dictionary,
				idxCtr:     int64(-1),
				dictLen:    len(dictionary),
			}
			got := l.SentenceN(tt.N)
			parts := strings.Split(got, " ")
			if len(parts) != tt.expectCt {
				t.Errorf("Length of sentence less than 10 words")
			}
			for _, x := range parts {
				if !tt.expectWords[x] {
					t.Errorf("Did not expect word %s", x)
				}
			}
		})
	}
}

func Test_lipsumParagraph(t *testing.T) {
	tests := []struct {
		name  string
		count int
	}{
		{name: "paragraph", count: 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &lipsum{
				dictionary: dictionary,
				idxCtr:     int64(-1),
				dictLen:    len(dictionary),
			}
			got := l.Paragraph()
			periodCt := strings.Count(got, ".")
			if periodCt < tt.count {
				t.Errorf("Expected atleast %d sentences, got %d sentences", tt.count, periodCt)
			}
		})
	}
}

func Test_lipsumParagraphN(t *testing.T) {
	tests := []struct {
		name  string
		N     int
		count int
	}{
		{name: "single", N: 1, count: 0},
		{name: "five", N: 5, count: 4},
		{name: "hundred", N: 100, count: 99},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &lipsum{
				dictionary: dictionary,
				idxCtr:     int64(-1),
				dictLen:    len(dictionary),
			}
			got := l.ParagraphN(tt.N)
			periodCt := strings.Count(got, ".")
			if periodCt != tt.count {
				t.Errorf("Expected %d sentences, got %d sentences", tt.count, periodCt)
			}
		})
	}
}
