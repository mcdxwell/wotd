package wotd

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

func cachePath() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("could not get user cache directory: %w", err)
	}
	return path.Join(cacheDir, "wotd", "wotd.json"), nil
}

func saveInfo(date, word string) error {
	wInfo := Word{Date: date, Word: word}

	// Get the current list of words (handling errors if the cache read fails)
	info, err := searchCache()
	if err != nil {
		// If we fail to read the cache, start with a new list, logging the failure.
		info = []Word{}
	}

	info = append(info, wInfo)
	return saveWord(info)
}

func saveWord(w []Word) error {
	wordData, err := json.MarshalIndent(w, "", " \t")
	if err != nil {
		return fmt.Errorf("failed to marshal word data to JSON: %w", err)
	}

	cachePath, err := cachePath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(path.Dir(cachePath), 0755); err != nil {
		return fmt.Errorf("failed to create cache directory structure: %w", err)
	}

	if err := os.WriteFile(cachePath, wordData, 0644); err != nil {
		return fmt.Errorf("failed to write data to cache file %s: %w", cachePath, err)
	}
	return nil
}

func searchCache() ([]Word, error) {
	cachePath, err := cachePath()
	if err != nil {
		return nil, err
	}

	wordData, err := os.ReadFile(cachePath)
	if err != nil {
		// If the file is not found, treat it as a fresh start and return an empty slice.
		if os.IsNotExist(err) {
			return []Word{}, nil
		}
		return nil, fmt.Errorf("failed to read cache file %s: %w", cachePath, err)
	}

	var wotd []Word
	if err = json.Unmarshal(wordData, &wotd); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON data from %s: %w", cachePath, err)
	}
	return wotd, nil
}
