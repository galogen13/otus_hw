package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

//easyjson:json
type User struct {
	Email string
}

type DomainStat map[string]int

var (
	ErrEmptyDomain = errors.New("empty domain")
)

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	var (
		result = DomainStat{}
	)

	if domain == "" {
		return nil, ErrEmptyDomain
	}

	reader := bufio.NewReader(r)

	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			if len(line) == 0 {
				return result, nil
			}
			result, err = getDomainStatInLine(line, domain, result)
			if err != nil {
				return nil, err
			}
			return result, nil
		}
		if err != nil {
			return nil, err
		}

		result, err = getDomainStatInLine(line, domain, result)
		if err != nil {
			return nil, err
		}
	}
}

func getDomainStatInLine(line []byte, domain string, result DomainStat) (DomainStat, error) {
	user := User{}
	err := user.UnmarshalJSON(line)
	if err != nil {
		return nil, err
	}

	if user.Email == "" {
		return result, nil
	}

	if contain := strings.Contains(user.Email, "."+domain); contain {
		result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
	}
	return result, nil
}
