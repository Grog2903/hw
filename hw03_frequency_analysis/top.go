package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var reg = regexp.MustCompile(`([.,;:!?'])`)

func Top10(text string) []string {
	var split = strings.Fields(text)
	var counterMap = make(map[string]int)
	var result []string

	for _, word := range split {
		if word == "-" {
			continue
		}

		lowerKeyAndReg := reg.ReplaceAllString(strings.ToLower(word), "")

		if lowerKeyAndReg == "" {
			continue
		}

		counterMap[lowerKeyAndReg] = counterMap[lowerKeyAndReg] + 1
	}

	type keyValueStruct struct {
		Key   string
		Value int
	}

	var sortedStruct []keyValueStruct

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
