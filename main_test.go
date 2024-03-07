package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
)

func testParse(n int, t *testing.T) {
	suffix := fmt.Sprintf(".%d.txt", n)

	p := bytes.Buffer{}
	parse("testdata/measurements"+suffix, &p)

	result := "testdata/result" + suffix
	expect := "testdata/expect" + suffix

	if err := os.WriteFile(result, p.Bytes(), 0664); err != nil {
		t.Fatal(err)
	}

	f, err := os.ReadFile(expect)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(f, p.Bytes()) {
		t.Fatalf("diff %s %s", expect, result)
	}
}

func TestParse1K(t *testing.T) { testParse(1000, t) }
func TestParse1M(t *testing.T) { testParse(1000000, t) }
func TestParse1G(t *testing.T) { testParse(1000000000, t) }

func benchmark(n int, b *testing.B) {
	path := fmt.Sprintf("testdata/measurements.%d.txt", n)
	for n := 0; n < b.N; n++ {
		parse(path, io.Discard)
	}
}

func BenchmarkParse1K(b *testing.B) { benchmark(1000, b) }
func BenchmarkParse1M(b *testing.B) { benchmark(1000000, b) }
func BenchmarkParse1G(b *testing.B) { benchmark(1000000000, b) }
