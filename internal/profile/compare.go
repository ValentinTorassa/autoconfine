package profile

import (
	"fmt"
	"sort"

	"github.com/ValentinTorassa/autoconfine/internal/traceio"
)

// CompareResult describe las diferencias entre dos conjuntos.
type CompareResult struct {
	Added   []string `json:"added"`
	Removed []string `json:"removed"`
	Common  []string `json:"common"`
}

// CompareTraces compara syscalls entre dos trazas.
func CompareTraces(pathA, pathB string) (*CompareResult, error) {
	a, err := traceio.ReadSyscalls(pathA)
	if err != nil {
		return nil, fmt.Errorf("leer traza A: %w", err)
	}
	b, err := traceio.ReadSyscalls(pathB)
	if err != nil {
		return nil, fmt.Errorf("leer traza B: %w", err)
	}
	return compareMaps(a, b), nil
}

// CompareProfiles compara syscalls permitidas entre dos perfiles seccomp.
func CompareProfiles(pathA, pathB string) (*CompareResult, error) {
	a, err := AllowedSyscalls(pathA)
	if err != nil {
		return nil, fmt.Errorf("leer perfil A: %w", err)
	}
	b, err := AllowedSyscalls(pathB)
	if err != nil {
		return nil, fmt.Errorf("leer perfil B: %w", err)
	}
	return compareSlices(a, b), nil
}

func compareMaps(a, b map[string]int) *CompareResult {
	setA := make(map[string]struct{}, len(a))
	for k := range a {
		setA[k] = struct{}{}
	}
	setB := make(map[string]struct{}, len(b))
	for k := range b {
		setB[k] = struct{}{}
	}
	return diffSets(setA, setB)
}

func compareSlices(a, b []string) *CompareResult {
	setA := sliceToSet(a)
	setB := sliceToSet(b)
	return diffSets(setA, setB)
}

func sliceToSet(s []string) map[string]struct{} {
	set := make(map[string]struct{}, len(s))
	for _, v := range s {
		set[v] = struct{}{}
	}
	return set
}

func diffSets(a, b map[string]struct{}) *CompareResult {
	res := &CompareResult{}
	for k := range a {
		if _, ok := b[k]; ok {
			res.Common = append(res.Common, k)
		} else {
			res.Removed = append(res.Removed, k)
		}
	}
	for k := range b {
		if _, ok := a[k]; !ok {
			res.Added = append(res.Added, k)
		}
	}
	sort.Strings(res.Added)
	sort.Strings(res.Removed)
	sort.Strings(res.Common)
	return res
}
