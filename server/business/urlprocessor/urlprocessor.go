package urlprocessor

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

var (
	ErrInvalidOperation = errors.New("invalid operation type")
)

type Operation string

const (
	OpCanonical   Operation = "canonical"
	OpRedirection Operation = "redirection"
	OpAll         Operation = "all"
)

type URLProcessor struct{}

func New() *URLProcessor {
	return &URLProcessor{}
}

func (p *URLProcessor) Process(rawURL string, op string) (string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	operation := Operation(op)

	switch operation {
	case OpCanonical:
		return processCanonical(parsed), nil
	case OpRedirection:
		return processRedirection(parsed), nil
	case OpAll:
		canonicalized := processCanonical(parsed)
		cleanedParsed, err := url.Parse(canonicalized)
		if err != nil {
			return "", fmt.Errorf("canonicalized URL parse failed: %w", err)
		}
		return processRedirection(cleanedParsed), nil
	default:
		return "", fmt.Errorf("unsupported operation type")
	}
}

func processCanonical(u *url.URL) string {
	u.RawQuery = ""
	u.Fragment = ""
	u.Path = strings.TrimRight(u.Path, "/")
	return u.String()
}

func processRedirection(u *url.URL) string {
	u.Scheme = strings.ToLower(u.Scheme)

	// Force domain to www.byfood.com
	u.Host = "www.byfood.com"

	// Lowercase path (important to normalize paths like /FOOD)
	u.Path = strings.ToLower(u.Path)

	return u.String()
}
