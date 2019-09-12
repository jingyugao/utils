package main

import (
	"reflect"
	"sort"
)

func CheckCircle(m map[string][]string) (allPath [][]string) {
	for k := range m {
		deepSearch(k, m, nil, &allPath)
	}

	for i, path := range allPath {
		sort.Strings(path)
		allPath[i] = append(path, path[0])
	}
	Unique(&allPath)
	return

}

func Unique(slc interface{}) {
	rval := reflect.ValueOf(slc)
	if rval.Kind() != reflect.Ptr {
		panic("param not ptr")
	}

	rval = reflect.Indirect(rval)

	num := rval.Len()
	result := reflect.MakeSlice(rval.Type(), 0, num)

	for i := 0; i < num; i++ {
		flag := true
		for j := 0; j < result.Len(); j++ {
			if reflect.DeepEqual(rval.Index(i).Interface(), result.Index(j).Interface()) {
				flag = false
				break
			}
		}
		if flag {
			result = reflect.Append(result, rval.Index(i))
		}
	}
	rval.Set(result)
	return
}

func inSlice(target string, slice []string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}

func deepSearch(curNode string, tree map[string][]string, curPath []string, allPath *[][]string) {
	if inSlice(curNode, curPath) {
		*allPath = append(*allPath, curPath)
		return
	}
	nextPath := make([]string, len(curPath)+1)
	copy(nextPath, curPath)
	nextPath[len(nextPath)-1] = curNode

	children, hasChildren := tree[curNode]
	if !hasChildren {
		return
	}

	for _, child := range children {
		deepSearch(child, tree, nextPath, allPath)
	}

}
