package main

import (
	"log"
	"path"
)

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkIsPathAbs(filepath string) {
	if !path.IsAbs(filepath) {
		log.Fatalf("'%s' is %s\n", filepath, path.ErrBadPattern)
	}
}

// TODO: Write tests
func Filter[T any](items []T, fn func(item T) bool) []T {
	result := make([]T, 0, len(items))
	for _, item := range items {
		if fn(item) {
			result = append(result, item)
		}
	}

	return result
}

func Unique[T comparable](items []T) []T {
	slice := make([]T, 0, len(items))
	set := map[T]struct{}{}
	for _, item := range items {
		_, ok := set[item]
		if !ok {
			set[item] = struct{}{}
			slice = append(slice, item)
		}
	}

	return slice
}

// TODO: Write tests
func Compare[T comparable](items1, items2 []T) bool {
	a, b := len(items1), len(items2)
	if a != b {
		return false
	}

	for i := range items1 {
		if !(items1[i] == items2[i]) {
			return false
		}
	}

	return true
}

func In[T comparable](item T, items []T) bool {
	for _, it := range items {
		if it == item {
			return true
		}
	}
	return false
}

func NotIn[T comparable](item T, items []T) bool {
	return !In(item, items)
}
