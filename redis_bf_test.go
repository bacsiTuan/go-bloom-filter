package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"testing"
)

func TestRedisBloomFilter(t *testing.T) {
	ctx := context.Background()

	// Connect to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address
		DB:   0,                // Redis database number
	})
	defer redisClient.Close()

	// Parameters for the Bloom Filter
	n := 20        // Number of items to add
	fpProb := 0.05 // False positive probability
	redisKey := "test_bloom_filter"

	// Clear the Redis key before testing
	err := redisClient.Del(ctx, redisKey).Err()
	if err != nil {
		t.Fatalf("Failed to clear Redis key: %v", err)
	}

	// Initialize the Redis-backed Bloom Filter
	bloomFilter := NewRedisBloomFilter(redisClient, redisKey, n, fpProb)

	t.Logf("Size of bit array: %d", bloomFilter.size)
	t.Logf("False positive Probability: %.2f", fpProb)
	t.Logf("Number of hash functions: %d", bloomFilter.hashCount)

	// Words to be added
	wordsPresent := []string{
		"abound", "abounds", "abundance", "abundant", "accessible",
		"bloom", "blossom", "bolster", "bonny", "bonus",
		"bonuses", "coherent", "cohesive", "colorful", "comely",
		"comfort", "gems", "generosity", "generous", "generously", "genial",
	}

	// Words not added
	wordsAbsent := []string{
		"bluff", "cheater", "hate", "war", "humanity",
		"racism", "hurt", "nuke", "gloomy", "facebook",
		"geeksforgeeks", "twitter",
	}

	// Add words to the Bloom Filter
	for _, word := range wordsPresent {
		bloomFilter.Add(word)
	}

	// Shuffle test words
	shuffle(wordsPresent)
	shuffle(wordsAbsent)

	// Create test set
	testWords := append(wordsPresent[:10], wordsAbsent...)
	shuffle(testWords)

	// Check test words
	for _, word := range testWords {
		if bloomFilter.Check(word) {
			if contains(wordsAbsent, word) {
				t.Logf("'%s' is a false positive!", word)
			} else {
				t.Logf("'%s' is probably present!", word)
			}
		} else {
			if contains(wordsPresent, word) {
				t.Errorf("'%s' should be present but was not found!", word)
			} else {
				t.Logf("'%s' is definitely not present!", word)
			}
		}
	}

	// Cleanup Redis key after test
	err = redisClient.Del(ctx, redisKey).Err()
	if err != nil {
		t.Errorf("Failed to delete Redis key after test: %v", err)
	}
}
