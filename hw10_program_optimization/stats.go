package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	reader := bufio.NewReader(r)

	var i int
	var line []byte
	for {
		var user User

		line, _, err = reader.ReadLine()
		if errors.Is(err, io.EOF) {
			return result, nil
		}
		if err != nil {
			return result, fmt.Errorf("error while reading: %w", err)
		}

		if err = user.UnmarshalJSON(line); err != nil {
			return
		}
		result[i] = user
		i++
	}
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	s := "." + domain
	for _, user := range u {
		ok := strings.HasSuffix(user.Email, s)
		if ok {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}
	return result, nil
}
