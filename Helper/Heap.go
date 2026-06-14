package helper

import (
	"container/heap"
	"sync"

	model "github.com/S-Unknown047/LoadBalancer/Model"
)

type ServerHeap []*model.Server

func (h ServerHeap) Len() int {
	return len(h)
}

// min heap
func (h ServerHeap) Less(i, j int) bool {
	return h[i].Connection < h[j].Connection
}

// max heap
// func (h IntHeap) Less(i, j int) bool {
// 	return h[i] > h[j]
// }

func (h *ServerHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
	(*h)[i].Index = i
	(*h)[j].Index = j
}

func (h *ServerHeap) Push(x interface{}) {
	srv := x.(*model.Server)
	srv.Index = len(*h)
	*h = append(*h, srv)
}

func (h *ServerHeap) Pop() interface{} {
	old := *h
	n := len(old)
	srv := old[n-1]
	srv.Index = -1
	*h = old[:n-1]
	return srv
}

var (
	H      = &ServerHeap{}
	HeapMu sync.Mutex
)

func EnterData(server *[]model.Server) {
	HeapMu.Lock()
	defer HeapMu.Unlock()
	*H = make(ServerHeap, len(*server))
	for i := range *server {
		(*server)[i].Index = i
		(*H)[i] = &(*server)[i]
	}

	heap.Init(H)
}

func GetServer() *model.Server {

	HeapMu.Lock()
	defer HeapMu.Unlock()
	srv := (*H)[0]
	srv.Connection++
	heap.Fix(H, 0)
	return srv

}

func HeapFix(srv *model.Server) {
	HeapMu.Lock()
	defer HeapMu.Unlock()
	heap.Fix(H, srv.Index)
}

// func UpdateHeap() {
// 	model.Server = h.Pop()
// }

// func main() {
// 	h := &IntHeap{5, 3, 5}
// 	heap.Init(h)

// 	heap.Push(h, 2)

// 	fmt.Println(h.Pop())n
// }
