package main

import (
	"container/heap"
	"fmt"
	"sort"
)

// (definec c-trace (n :pos) :lop
//   (cond
//    ((== n 1) (cons n nil))
//    ((evenp n) (cons n (c-trace (/ n 2))))
//    (t  (cons n (c-trace (1+ (* 3 n)))))))
//  Write code to find a number, *e-num*, between 1 and 2^33 such that
//  (len (c-trace *e-num*)) is the 11th largest number in the set
//  { n :  n = (len (c-trace i)) ^ 1 <= i <= 2^33 }.

type trace struct {
	i      int
	result int
}

// An TraceHeap is a min-heap of traces.
type TraceHeap []trace

func (h TraceHeap) Len() int           { return len(h) }
func (h TraceHeap) Less(i, j int) bool { return h[i].result < h[j].result }
func (h TraceHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *TraceHeap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(trace))
}

func (h *TraceHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func nth_largest_ctrace_length(minsize int, maxsize int, nlargest int, result chan TraceHeap) {
	slice := make(TraceHeap, nlargest)
	largest_n := &slice
	visited := make(map[int]bool)
	for i := minsize; i <= maxsize; i++ {
		n := i
		result := 0

		for true {
			if n == 1 {
				result += 1
				break
			} else if n%2 == 0 {
				result += 1
				n = n / 2
			} else {
				result += 1
				n = n*3 + 1
			}
		}
		if i%1000000 == 0 {
			fmt.Printf("%v%% %v\n", (float32(i-minsize)/float32(maxsize-minsize))*100, largest_n)
		}
		if _, ok := visited[result]; !ok {
			visited[result] = true
			heap.Push(largest_n, trace{i: i, result: result})
			heap.Pop(largest_n)
		}
	}
	result <- *largest_n
}

func nth_largest_ctrace_length_parallel(maxsize int, nlargest int, threads int) int {
	results := make(chan TraceHeap)

	for i := 0; i < threads; i++ {
		size := maxsize / threads
		fmt.Printf("starting thread %v, %v-%v\n", i, i*size+1, (i+1)*size)
		go nth_largest_ctrace_length(i*size+1, (i+1)*size, nlargest, results)
	}

	done := 0
	threadHeap := make(TraceHeap, 0)
	for done < threads {
		select {
		case val := <-results:
			for _, val := range val {
				heap.Push(&threadHeap, val)
			}
			done += 1
		}
	}
	close(results)
	sort.Sort(&threadHeap)
	nth := threadHeap[len(threadHeap)-(1+nlargest)]
	fmt.Printf("%v\n", nth)
	return nth.result
}

func main() {
	fmt.Println(nth_largest_ctrace_length_parallel(8589934592, 11, 16))
}
