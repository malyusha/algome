package sourceproviders

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/malyusha/algome/genreadme/generator"
)

const (
	leetcodeURL = "https://leetcode.com"
)

type leetcodeProvider struct {
	httpClient *http.Client
}

type leetcodeRt struct {
	base http.RoundTripper
}

func (l *leetcodeRt) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
	return l.base.RoundTrip(req)
}

// NewLeetcodeProvider returns new instance of leetcode problems provider.
func NewLeetcodeProvider() (*leetcodeProvider, error) {
	tp := &http.Transport{
		DisableKeepAlives:  false,
		DisableCompression: false,
	}
	provider := &leetcodeProvider{
		httpClient: &http.Client{
			Transport: &leetcodeRt{base: tp},
			Timeout:   time.Second * 15,
		},
	}

	return provider, nil
}

type LeetcodeProblemsResponse struct {
	Pairs []LeetcodeStatPair `json:"stat_status_pairs"`
}

type LeetcodeStatPair struct {
	Stat       LeetcodeStat `json:"stat"`
	Difficulty struct {
		Level int64 `json:"level"`
	} `json:"difficulty"`
}
type LeetcodeStat struct {
	QuestionId        int64  `json:"question_id"`
	QuestionTitle     string `json:"question__title"`
	QuestionTitleSlug string `json:"question__title_slug"`
}

func (l *leetcodeProvider) newRequest(ctx context.Context, path string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/%s", leetcodeURL, path), http.NoBody)
	return req.WithContext(ctx)
}

func (l *leetcodeProvider) GetAllProblems(ctx context.Context) ([]generator.Problem, error) {
	req := l.newRequest(ctx, "problems/all")
	resp, err := l.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("bad status code '%d'", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	dst := LeetcodeProblemsResponse{}
	if err := json.Unmarshal(data, &dst); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	out := make([]generator.Problem, len(dst.Pairs))
	for i, pair := range dst.Pairs {
		out[i] = problemFromLeetcode(pair)
	}
	return out, nil
}

func problemFromLeetcode(pair LeetcodeStatPair) generator.Problem {
	problem := generator.Problem{
		ID:    pair.Stat.QuestionId,
		URL:   fmt.Sprintf("%s/problems/%s", leetcodeURL, pair.Stat.QuestionTitleSlug),
		Slug:  pair.Stat.QuestionTitleSlug,
		Title: pair.Stat.QuestionTitle,
	}

	switch pair.Difficulty.Level {
	case 1:
		problem.Difficulty = generator.DifficultyEasy
	case 2:
		problem.Difficulty = generator.DifficultyMedium
	case 3:
		problem.Difficulty = generator.DifficultyHard
	}

	return problem
}
