package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type station struct {
	sum float64
	n   int
	min float64
	max float64
}

func parse(path string, w io.Writer) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	m := make(map[string]*station)

	s := bufio.NewScanner(f)
	for s.Scan() {
		p := strings.Split(s.Text(), ";")
		if len(p) != 2 {
			panic("unexpected row")
		}

		name, value := p[0], p[1]

		k, err := strconv.ParseFloat(value, 64)
		if err != nil {
			panic(err)
		}

		if k < -99.9 || k > 99.9 {
			panic("unexpected value")
		}

		if i := m[name]; i == nil {
			m[name] = &station{
				sum: k,
				n:   1,
				min: k,
				max: k,
			}
		} else {
			i.sum += k
			i.n++
			i.min = min(i.min, k)
			i.max = max(i.max, k)
		}
	}

	if err := s.Err(); err != nil {
		panic(err)
	}

	names := make([]string, 0, len(m))
	for s := range m {
		names = append(names, s)
	}

	sort.Strings(names)

	for _, name := range names {
		i := m[name]
		avg := i.sum / float64(i.n)
		fmt.Fprintf(w, "%s:%.1f/%.1f/%.1f\n", name, i.min, math.Trunc(avg*10)/10, i.max)
	}
}

func main() {
	if len(os.Args) != 2 {
		panic("usage: 1brc <filename>")
	}

	parse(os.Args[1], os.Stdout)
}
