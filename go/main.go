package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"secret.sharing/shamir"
)

func inputFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)

	result := make([]byte, 0)
	tmp := make([]byte, 1024)

	for {
		size, err := reader.Read(tmp)
		if err == io.EOF || size == 0 {
			break
		}
		result = append(result, tmp[0:size]...)

	}
	return result, nil
}

func seal(data []byte, k, n int, verbose, verify bool) {
	var s time.Time
	if verbose {
		fmt.Printf("input file size: %d bytes\n", len(data))
		s = time.Now()
	}
	shares, err := shamir.Seal(data, k, n)
	if err != nil {
		panic(err)
	}

	if verbose {
		fmt.Printf("seal time: %s\n", time.Since(s))
	}
	if verify {
		if verbose {
			s = time.Now()
		}
		unsealData, err := shamir.Unseal(shares)
		if err != nil {
			panic(err)
		}
		if verbose {
			fmt.Printf("unseal time: %s\n", time.Since(s))
			s = time.Now()
		}
		count := 0
		for i, n := 0, len(unsealData); i < n; i++ {
			if data[i] != unsealData[i] {
				count++
			}
		}
		if verbose {
			fmt.Printf("verify time: %s\nnumber of diff bytes %d\n", time.Since(s), count)
		}
	}

	out, err := json.Marshal(shares)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", out)
}

func unseal(data []byte, verbose bool) {
	var shares []shamir.ShamirShare
	err := json.Unmarshal(data, &shares)
	if err != nil {
		panic(err)
	}
	var s time.Time

	if verbose {
		s = time.Now()
	}
	out, err := shamir.Unseal(shares)

	if err != nil {
		panic(err)
	}
	if verbose {
		fmt.Printf("unseal time: %s\n", time.Since(s))
	}

	fmt.Printf("%s\n", out)

}

func main() {
	mode := flag.String("mode", "seal", "seal or unseal")
	name := flag.String("filename", "", "source file")
	k := flag.Int("k", 2, "k is k-out-of-n")
	n := flag.Int("n", 2, "n is k-out-of-n")
	verbose := flag.Bool("verbose", false, "verbose")
	verify := flag.Bool("verify", true, "verification")
	flag.Parse()

	data, err := inputFile(*name)
	if err != nil {
		panic(err)
	}

	if *mode == "seal" {
		seal(data, *k, *n, *verbose, *verify)
	} else if *mode == "unseal" {
		unseal(data, *verbose)
	} else {
		panic(fmt.Errorf("not support mode %s", *mode))
	}
}
