package main

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"syscall"
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

	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}

	data, err := syscall.Mmap(int(f.Fd()), 0, int(fi.Size()), syscall.PROT_READ, syscall.MAP_SHARED)
	if err != nil {
		panic(err)
	}

	defer syscall.Munmap(data)

	m := make(map[string]*station)
	a := 0

	for {
		b := bytes.IndexByte(data[a:], '\n')
		if b < 0 {
			break
		}

		p := string(data[a : a+b])
		a += b + 1

		i := strings.LastIndex(p, ";")
		if i < 0 {
			panic("unexpected row")
		}

		name, value := p[0:i], p[i+1:]

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
