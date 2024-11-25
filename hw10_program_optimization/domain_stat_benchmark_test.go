package hw10programoptimization

import (
	"bytes"
	"testing"
)

var testData = []byte(`{"email":"user1@example.com"}
{"email":"user2@example.com"}
{"email":"user3@sub.example.com"}
{"email":"user4@example.org"}
{"email":"user5@sub.example.com"}
{"email":"user6@another.com"}
{"email":"user7@example.com"}`)

func BenchmarkGetDomainStat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(testData)
		_, err := GetDomainStat(r, "example.com")
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}
