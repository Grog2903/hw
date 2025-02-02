package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(text string) []string {
	split := strings.Fields(text)
	counterMap := make(map[string]int)
	result := make([]string, 0, 10)

	for _, word := range split {
		counterMap[word]++
	}

	type keyValueStruct struct {
		Key   string
		Value int
	}

	sortedStruct := make([]keyValueStruct, 0, len(counterMap))

	for key, value := range counterMap {
		sortedStruct = append(sortedStruct, keyValueStruct{key, value})
	}

	sort.Slice(sortedStruct, func(i, j int) bool {
		if sortedStruct[i].Value != sortedStruct[j].Value {
			return sortedStruct[i].Value > sortedStruct[j].Value
		}

		return sortedStruct[i].Key < sortedStruct[j].Key
	})

	i := 0
	for _, value := range sortedStruct {
		if i >= 10 {
			break
		}

		result = append(result, value.Key)

		i++
	}

	return result
}
