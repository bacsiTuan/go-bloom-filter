package main

import (
	"hash/fnv"
	"math"
	"sync"
)

type BloomFilter struct {
	bitArray  []bool
	size      int
	hashCount int
	lock      sync.Mutex
}

func NewBloomFilter(itemsCount int, fpProb float64) *BloomFilter {
	size := getSize(itemsCount, fpProb)
	hashCount := getHashCount(size, itemsCount)

	return &BloomFilter{
		bitArray:  make([]bool, size),
		size:      size,
		hashCount: hashCount,
	}
}

func (bf *BloomFilter) Add(item string) {
	digests := bf.getDigests(item)
	bf.lock.Lock()
	defer bf.lock.Unlock()

	for _, digest := range digests {
		bf.bitArray[digest] = true
	}
}

func (bf *BloomFilter) Check(item string) bool {
	digests := bf.getDigests(item)
	bf.lock.Lock()
	defer bf.lock.Unlock()

	for _, digest := range digests {
		if !bf.bitArray[digest] {
			return false
		}
	}
	return true
}

func (bf *BloomFilter) getDigests(item string) []int {
	digests := make([]int, bf.hashCount)

	for i := 0; i < bf.hashCount; i++ {
		h := fnv.New32a()
		_, err := h.Write([]byte(item))
		if err != nil {
			return nil
		}
		seed := uint32(i)
		digest := int(h.Sum32()+seed) % bf.size
		digests[i] = digest
	}
	return digests
}

func getSize(n int, p float64) int {
	/*
		Return the size of bit array(m) to used using
		following formula
		m = -(n * lg(p)) / (lg(2)^2)
		n : int
			number of items expected to be stored in filter
		p : float
			False Positive probability in decimal
	*/
	return int(-float64(n) * math.Log(p) / (math.Pow(math.Log(2), 2)))
}

func getHashCount(m, n int) int {
	/*
		Return the hash function(k) to be used with following formula
		k = (m/n) * lg(2)

		m : int
			size of bit array
		n : int
			number of items expected to be stored in filter
	*/
	return int(float64(m) / float64(n) * math.Log(2))
}

func main() {
	itemsCount := 20
	fpProb := 0.05

	bf := NewBloomFilter(itemsCount, fpProb)
	bf.Add("test-item")

	if bf.Check("test-item") {
		println("Item exists")
	} else {
		println("Item does not exist")
	}
	println(bf.Check("test1-item"))

}
