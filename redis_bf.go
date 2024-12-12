package main

import (
	"context"
	"github.com/go-redis/redis/v8"
)

// RedisBloomFilter is a Redis-backed Bloom Filter using murmur3 hash function.
type RedisBloomFilter struct {
	redisClient *redis.Client
	redisKey    string
	fpProb      float64
	size        int
	hashCount   int
}

// NewRedisBloomFilter initializes a new RedisBloomFilter.
func NewRedisBloomFilter(redisClient *redis.Client, redisKey string, itemsCount int, fpProb float64) *RedisBloomFilter {
	size := getSize(itemsCount, fpProb)
	hashCount := getHashCount(size, itemsCount)
	return &RedisBloomFilter{
		redisClient: redisClient,
		redisKey:    redisKey,
		fpProb:      fpProb,
		size:        size,
		hashCount:   hashCount,
	}
}

// Add inserts an item into the Bloom Filter.
func (bf *RedisBloomFilter) Add(item string) {
	ctx := context.Background()
	for i := 0; i < bf.hashCount; i++ {
		digest := murmur3Hash(item, i) % bf.size
		bf.redisClient.SetBit(ctx, bf.redisKey, int64(digest), 1)
	}
}

// Check verifies if an item might be in the Bloom Filter.
func (bf *RedisBloomFilter) Check(item string) bool {
	ctx := context.Background()
	for i := 0; i < bf.hashCount; i++ {
		digest := murmur3Hash(item, i) % bf.size
		if bf.redisClient.GetBit(ctx, bf.redisKey, int64(digest)).Val() == 0 {
			return false // Definitely not present
		}
	}
	return true // Might be present
}

// murmur3Hash generates a hash value for a given string and seed using Murmur3 hash algorithm.
func murmur3Hash(data string, seed int) int {
	// Using a simple hash implementation as a placeholder; replace with Murmur3 hash library.
	var hash int
	for i := 0; i < len(data); i++ {
		hash = (hash*seed + int(data[i])) % 2147483647
	}
	return hash
}
