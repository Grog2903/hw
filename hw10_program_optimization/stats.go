package hw10programoptimization

import (
	"bufio"
	"fmt"
	"github.com/goccy/go-json"
	"io"
	"strings"
)

type User struct {
	Email string
}

type users []User

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

func getUsers(r io.Reader) (result users, err error) {
	var user User
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if err = json.Unmarshal(scanner.Bytes(), &user); err != nil {
			continue
		}

		result = append(result, user)
	}

	if err := scanner.Err(); err != nil {
		return result, fmt.Errorf("scanner error: %w", err)
	}

	return result, nil
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		if strings.Contains(user.Email, domain) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}
	return result, nil
}
