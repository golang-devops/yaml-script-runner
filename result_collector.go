package main

import (
	"os/exec"
	"sync"
)

type result struct {
	Cmd        *exec.Cmd
	Successful bool
}

type ResultCollector struct {
	sync.RWMutex

	results []*result
}

func (r *ResultCollector) AppendResult(res *result) {
	r.Lock()
	defer r.Unlock()
	r.results = append(r.results, res)
}

func (r *ResultCollector) TotalCount() int {
	return len(r.results)
}

func (r *ResultCollector) SuccessCount() (cnt int) {
	cnt = 0
	for _, res := range r.results {
		if res.Successful {
			cnt++
		}
	}
	return cnt
}

func (r *ResultCollector) FailedCount() (cnt int) {
	cnt = 0
	for _, res := range r.results {
		if !res.Successful {
			cnt++
		}
	}
	return cnt
}

func (r *ResultCollector) FailedDisplayList() []string {
	list := []string{}
	for _, res := range r.results {
		if !res.Successful {
			list = append(list, execCommandToDisplayString(res.Cmd))
		}
	}
	return list
}
