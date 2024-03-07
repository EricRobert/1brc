TESTS=1K|1M|1G

test: testdata/.done
	go test -v -run '$(TESTS)'

bench: testdata/.done
	go test -v -run no -bench '$(TESTS)' -test.benchmem

profile: testdata/.done
	go test -v -run no -bench 1G -cpuprofile cpu.out -memprofile mem.out

testdata/.done:
	git clone --depth 1 https://github.com/gunnarmorling/1brc
	cd 1brc/src/main/python && python3 create_measurements.py 1_000_000_000
	mkdir -p testdata
	cat 1brc/data/measurements.txt | head -1000 >testdata/measurements.1000.txt
	cat 1brc/data/measurements.txt | head -1000000 >testdata/measurements.1000000.txt
	mv 1brc/data/measurements.txt testdata/measurements.1000000000.txt
	rm -rf 1brc
	touch testdata/.done
