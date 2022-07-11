package lfiscanner

import (
	"bytes"
	"errors"
	"net/url"
	"strings"
)

type Client interface {
	GetResponseBodyFrom(url string) ([]byte, error)
}

type Scanner struct {
	client Client
	cfg    *Config
}

func New(cfg *Config, cl Client) *Scanner {
	return &Scanner{
		client: cl,
		cfg:    cfg,
	}
}

// ScanURL scans a url by replacing query parameters with variations of Target
func (sc *Scanner) ScanUrl(urlStr string) ([]string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	if u.RawQuery == "" {
		return nil, errors.New("nothing to scan")
	}
	results := make([]string, 0, len(sc.cfg.Targets))
	for target, confirmStr := range sc.cfg.Targets {
		finalTarget := buildFinalPath(sc.cfg.LevelUpAttempts, target)
		backIndex := len(finalTarget) - len(target) + 1
		for i := 0; i < sc.cfg.LevelUpAttempts; i++ {
			// currentTarget is target with level ups added
			currentTarget := finalTarget[backIndex-i*3:]
			if i == 0 {
				currentTarget = "/" + currentTarget
			}
			sc.replaceQueryParams(u, currentTarget)
			respBody, err := sc.client.GetResponseBodyFrom(u.String())
			if err != nil && bytes.Contains(respBody, []byte(confirmStr)) {
				results = append(results, u.String())
				break
			}
		}
	}
	return results, nil
}

// replaceQueryParams replaces query params with newParam
func (sc *Scanner) replaceQueryParams(u *url.URL, newParam string) {
	qs := u.Query()
	for param := range qs {
		qs.Set(param, newParam)
	}
	u.RawQuery = qs.Encode()
}

// buildFinalPath builds path to target file with prepended "../" 'attempts' times
func buildFinalPath(attempts int, initialPath string) string {
	maxLen := 3*attempts + len(initialPath) - 1
	var str strings.Builder
	str.Grow(maxLen)
	for i := 0; i < attempts; i++ {
		str.WriteString("../")
	}
	str.WriteString(initialPath[1:])
	return str.String()
}
