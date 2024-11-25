package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/goccy/go-json"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	var user User
	result := make(DomainStat)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if err := json.Unmarshal(scanner.Bytes(), &user); err != nil {
			continue
		}

		atIndex := strings.LastIndex(user.Email, "@")
		if atIndex == -1 {
			continue
		}
		secPart := strings.ToLower(user.Email[atIndex+1:])
		if strings.HasSuffix(secPart, domain) {
			result[secPart]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}

	return result, nil
}
