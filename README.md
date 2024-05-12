# toy-shamir-secret-sharing
Toy implementation of Shamir's Secret Sharing Scheme

# What is this?
- This implementation is shamir's secret sharing scheme on 8-bit Galois Field.
    - Please click [here](https://en.wikipedia.org/wiki/Shamir%27s_secret_sharing) for details of shamir's secret sharing.
- No special speed-up implementation, such as using SIMD operations, is used. 
- This is a toy implementation and should not be used in a production environment. We do not take any responsibility for the use of this implementation. 

# Go Implementation
## How-to-use
### Seal
```
$ cat example/sample.txt
Hello world!!
こんにちは世界!!
$ go run main.go -filename example/sample.txt -mode seal > example/seal.txt
$ cat example/seal.txt
[{"Pt":1,"Val":"gMdvXurgdzw+z8Ar5T2G9nKzra6Jl+nIV8msO1q3SRAepKE6yQ=="},{"Pt":2,"Val":"wzpqCH67d8nqMTc1smQpb0pD3Ok3rS+1NnF97l5CQYEO99YX6g=="}]
```
### Unseal
```
$ go run main.go -filename example/seal.txt
Hello world!!
こんにちは世界!!
```