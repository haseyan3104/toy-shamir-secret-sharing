package shamir

import (
	"encoding/base64"
	"fmt"
	"sync"

	"secret.sharing/gf256"
)

type shamirShare struct {
	point  byte
	values []byte
}

type ShamirShare struct {
	Pt  byte
	Val string
}

func toBase64(share shamirShare) ShamirShare {
	return ShamirShare{share.point, base64.StdEncoding.EncodeToString(share.values)}
}

func fromBase64(share ShamirShare) (shamirShare, error) {
	values, err := base64.StdEncoding.DecodeString(share.Val)
	if err != nil {
		fmt.Println(err)
		return shamirShare{}, err
	}
	return shamirShare{share.Pt, values}, nil
}

func Seal(target []byte, k, n int) ([]ShamirShare, error) {
	return seal(target, k, n, gf256.NewRandBytes)
}

func seal(target []byte, k, n int, fn func(int) []byte) ([]ShamirShare, error) {
	if k-1 <= 0 {
		return nil, fmt.Errorf("must be k-1 > 0")
	}
	if n < k {
		return nil, fmt.Errorf("must be n >= k")
	}
	result := make([]shamirShare, n)
	for i := 0; i < n; i++ {
		result[i] = shamirShare{byte(i + 1), append(target[:0:0], target...)}
	}

	rands := make([][]byte, k-1)
	for i := 0; i < k-1; i++ {
		rands[i] = fn(len(target))
	}
	var wg sync.WaitGroup
	for j := 0; j < n; j++ {
		wg.Add(1)
		go func(j int) {
			x := result[j].point
			for t, length := 0, len(target); t < length; t++ {
				for s := 0; s < k-1; s++ {
					result[j].values[t] = gf256.Add(result[j].values[t], gf256.Mul(gf256.Pow(x, s+1), rands[s][t]))
				}
			}
			wg.Done()
		}(j)
	}
	wg.Wait()

	ans := make([]ShamirShare, n)
	for i := 0; i < n; i++ {
		ans[i] = toBase64(result[i])
	}

	return ans, nil
}

func Unseal(data []ShamirShare) ([]byte, error) {
	shares := make([]shamirShare, len(data))
	{
		dedup := make(map[byte]struct{})
		for i, l := 0, len(data); i < l; i++ {
			_, ok := dedup[data[i].Pt]
			if ok {
				return nil, fmt.Errorf("avoid divide zero")
			} else {
				dedup[data[i].Pt] = struct{}{}
			}
		}
	}
	for i, l := 0, len(data); i < l; i++ {
		tmp, err := fromBase64(data[i])
		if err != nil {
			return nil, err
		}
		shares[i] = tmp
	}
	length := len(shares[0].values)
	for i, l := 1, len(data); i < l; i++ {
		if len(shares[i].values) != length {
			return nil, fmt.Errorf("must be same length for shares")
		}
	}
	result := make([]byte, length)

	for t := 0; t < length; t++ {
		for j, l := 0, len(shares); j < l; j++ {
			s := byte(1)
			for m := 0; m < l; m++ {
				if j != m {
					s = gf256.Mul(s, gf256.Div(shares[m].point, gf256.Sub(shares[m].point, shares[j].point)))
				}
			}
			result[t] = gf256.Add(result[t], gf256.Mul(s, shares[j].values[t]))
		}
	}

	return result, nil
}
