package stringops

import (
	"unicode/utf8"
)

func ChunkSize(items []string, size int) (chunks [][]string) {
	if len(items) == 0 || size <= 0 {
		return nil
	}

	currentChunk := make([]string, 0, size)
	currentChunkSize := 0

	for _, item := range items {
		itemSize := len(item)
		if currentChunkSize+itemSize > size {
			chunks = append(chunks, currentChunk)
			currentChunk = make([]string, 0, size)
			currentChunkSize = 0
		}
		currentChunk = append(currentChunk, item)
		currentChunkSize += itemSize
	}

	if len(currentChunk) > 0 {
		chunks = append(chunks, currentChunk)
	}

	return chunks
}

func listLen(items []string) (sum int) {
	for _, i := range items {
		sum += utf8.RuneCountInString(i)
	}
	return sum
}
