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
	stat, err := calculateStat(r, domain)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return stat, nil
}

func calculateStat(r io.Reader, domain string) (result DomainStat, err error) {
	reader := bufio.NewReader(r)

	var line []byte
	dotDomain := "." + domain
	result = make(DomainStat)
	for {
		var user User

		line, _, err = reader.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return result, nil
			}

			return nil, fmt.Errorf("error while reading: %w", err)
		}

		if err = user.UnmarshalJSON(line); err != nil {
			return nil, fmt.Errorf("unmarshal json error: %w", err)
		}

		emailDomain, err := getEmailDomain(user.Email)
		if err != nil {
			return nil, fmt.Errorf("parsing email error: %w", err)
		}
		ok := strings.HasSuffix(emailDomain, dotDomain)
		if ok {
			result[emailDomain]++
		}
	}
}

func getEmailDomain(email string) (string, error) {
	parts := strings.SplitN(email, "@", 2)
	if len(parts) < 2 {
		return "", fmt.Errorf("incorrect email format")
	}
	emailDomain := strings.ToLower(parts[1])

	return emailDomain, nil
}
