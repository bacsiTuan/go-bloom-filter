package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func shuffle(slice []string) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}

func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func TestBloomFilter(t *testing.T) {
	n := 20        // Number of items to add
	fpProb := 0.05 // False positive probability
	bloom := NewBloomFilter(n, fpProb)

	if bloom.size <= 0 {
		t.Errorf("Expected bit array size to be greater than 0, got %d", bloom.size)
	}

	if bloom.hashCount <= 0 {
		t.Errorf("Expected number of hash functions to be greater than 0, got %d", bloom.hashCount)
	}

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

	// Add words to BloomFilter
	for _, word := range wordsPresent {
		bloom.Add(word)
	}

	// Shuffle test words
	shuffle(wordsPresent)
	shuffle(wordsAbsent)

	// Combine present and absent words for testing
	testWords := append(wordsPresent[:10], wordsAbsent...)
	shuffle(testWords)

	// Check test words
	for _, word := range testWords {
		if bloom.Check(word) {
			if contains(wordsAbsent, word) {
				fmt.Printf("'%s' is a false positive!\n", word)
			} else {
				fmt.Printf("'%s' is probably present!\n", word)
			}
		} else {
			if contains(wordsPresent, word) {
				t.Errorf("'%s' should be present but was not found!\n", word)
			} else {
				fmt.Printf("'%s' is definitely not present!\n", word)
			}
		}
	}
}
