package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sort"
	"syscall"
)

type station struct {
	sum int64
	n   int
	min int
	max int
}

func num2str(k int) string {
	u, d := k/10, k%10

	if k >= 0 {
		return fmt.Sprintf("%d.%d", u, d)
	}

	u = -u
	d = -d
	return fmt.Sprintf("-%d.%d", u, d)
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

		p := data[a : a+b]
		a += b + 1

		k := 0
		i := b - 1
		for d := 1; p[i] != ';'; i-- {
			if p[i] == '-' {
				k = -k
				continue
			}

			if p[i] != '.' {
				k += (int(p[i]) - 0x30) * d
				d *= 10
			}
		}

		if k < -999 || k > 999 {
			panic("unexpected value")
		}

		name := string(p[0:i])

		if i := m[name]; i == nil {
			m[name] = &station{
				sum: int64(k),
				n:   1,
				min: k,
				max: k,
			}
		} else {
			i.sum += int64(k)
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
		avg := int(i.sum*10/int64(i.n)) / 10
		fmt.Fprintf(w, "%s:%s/%s/%s\n", name, num2str(i.min), num2str(avg), num2str(i.max))
	}
}

func main() {
	if len(os.Args) != 2 {
		panic("usage: 1brc <filename>")
	}

	parse(os.Args[1], os.Stdout)
}
