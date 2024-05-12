package shamir

import "testing"

func TestSealUnseal(t *testing.T) {
	target := make([][]byte, 256)
	for i := 0; i < 256; i++ {
		target[i] = []byte{byte(i)}
	}

	for i := 0; i < 256; i++ {
		for j := 0; j < 256; j++ {
			fn := func(l int) []byte {
				tmp := make([]byte, l)
				for k := 0; k < l; k++ {
					tmp[k] = byte(j)
				}
				return tmp
			}
			ss, err := seal(target[i], 3, 3, fn)
			if err != nil {
				t.FailNow()
			}
			ans, _ := Unseal(ss)
			for k := 0; k < len(target[i]); k++ {
				if ans[k] != target[i][k] {
					t.FailNow()
				}
			}
		}
	}

	for i := 0; i < 256; i++ {
		ss, err := Seal(target[i], 3, 3)
		if err != nil {
			t.FailNow()
		}

		ans, _ := Unseal(ss)
		for k := 0; k < len(target[i]); k++ {
			if ans[k] != target[i][k] {
				t.FailNow()
			}
		}
	}
}
