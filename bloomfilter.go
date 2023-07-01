package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"github.com/bits-and-blooms/bitset"
)

type BloomFilter struct{
	numElem int
	numHashes int
	b *bitset.BitSet
}

func New(numElem int, numHashes int) *BloomFilter {
	return &BloomFilter{
		numElem, numHashes, bitset.New(uint(numElem)),
	}
}

func (bf *BloomFilter) Add(value string) {
	h :=generateHashValue(value)
	h1 := h >> 32
	for i:=0; i<bf.numHashes; i++{
		hashN := h + i*h1 
		if hashN < 0{
			hashN = -hashN
		}
		bf.b.Set(uint(hashN)% uint(bf.numElem))
	}
}

func (bf *BloomFilter) Contains(value string) bool{
	h :=generateHashValue(value)
	h1 := h >> 32
	for i:=0; i<bf.numHashes; i++{
		hashN := h + i*h1 
		if hashN < 0{
			hashN = -hashN
		}
		if !bf.b.Test(uint(hashN)% uint(bf.numElem)) {
			return false
		}
	}
	return true
}

func generateHashValue(value string) int{
	h := sha1.New()
	bits := h.Sum([]byte(value))
	buffer := bytes.NewBuffer(bits)
	result, _ := binary.ReadVarint(buffer)
	return int(result)
}

//Test scenarios
func main(){
	b := New(1000, 4)
	b.Add("Fruit")
	b.Add("Veggies")
	b.Add("Juice")


	fmt.Printf("%v\n", b.Contains("Fruit"))
	fmt.Printf("%v\n", b.Contains("Beer"))
	fmt.Printf("%v\n", b.Contains("Veggies"))
	fmt.Printf("%v\n", b.Contains("Juice"))
	fmt.Printf("%v\n", b.Contains("Apple"))

}