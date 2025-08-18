package metrics

import (
    "container/heap"
    "context"
    "sync"

    "github.com/gatij/goUrlShortener/internal/model"
)

// DomainHeapItem represents an item in our priority queue
type DomainHeapItem struct {
    domain       string
    shortenCount int
    index        int // index in the heap
}

// DomainMaxHeap implements a max heap of domains by shorten count
type DomainMaxHeap []*DomainHeapItem

// Required heap interface methods
func (h DomainMaxHeap) Len() int           { return len(h) }
func (h DomainMaxHeap) Less(i, j int) bool { return h[i].shortenCount > h[j].shortenCount } // Max heap
func (h DomainMaxHeap) Swap(i, j int) {
    h[i], h[j] = h[j], h[i]
    h[i].index = i
    h[j].index = j
}

func (h *DomainMaxHeap) Push(x interface{}) {
    n := len(*h)
    item := x.(*DomainHeapItem)
    item.index = n
    *h = append(*h, item)
}

func (h *DomainMaxHeap) Pop() interface{} {
    old := *h
    n := len(old)
    item := old[n-1]
    old[n-1] = nil     // avoid memory leak
    item.index = -1    // for safety
    *h = old[0 : n-1]
    return item
}

// MemoryStorage implements the metrics Storage interface with optimized data structures
type MemoryStorage struct {
    domains     map[string]model.DomainMetrics // Maps domain name to metrics
    domainHeap  *DomainMaxHeap                 // Max heap for quick access to top domains
    domainItems map[string]*DomainHeapItem     // Maps domain name to heap item for quick updates
    mu          sync.RWMutex                   // Protects the data structures
}

// NewMemoryStorage creates a new in-memory metrics storage
func NewMemoryStorage() *MemoryStorage {
    h := &DomainMaxHeap{}
    heap.Init(h)
    
    return &MemoryStorage{
        domains:     make(map[string]model.DomainMetrics),
        domainHeap:  h,
        domainItems: make(map[string]*DomainHeapItem),
    }
}

// SaveDomainMetrics stores metrics for a domain
func (s *MemoryStorage) SaveDomainMetrics(ctx context.Context, metrics model.DomainMetrics) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    // Store domain metrics in map
    s.domains[metrics.Domain] = metrics
    
    // Update or insert in the heap
    if item, exists := s.domainItems[metrics.Domain]; exists {
        // Update existing item
        item.shortenCount = metrics.ShortenCount
        heap.Fix(s.domainHeap, item.index)
    } else {
        // Create new item
        item := &DomainHeapItem{
            domain:       metrics.Domain,
            shortenCount: metrics.ShortenCount,
        }
        heap.Push(s.domainHeap, item)
        s.domainItems[metrics.Domain] = item
    }
    
    return nil
}

// GetTopDomains retrieves the top N domains based on shorten count
func (s *MemoryStorage) GetTopDomains(ctx context.Context, limit int) ([]model.DomainMetrics, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    // Create a copy of the heap to avoid modifying the original
    h := make(DomainMaxHeap, len(*s.domainHeap))
    copy(h, *s.domainHeap)
    
    // Extract top N domains
    result := make([]model.DomainMetrics, 0, limit)
    count := 0
    
    for h.Len() > 0 && count < limit {
        item := heap.Pop(&h).(*DomainHeapItem)
        metrics := s.domains[item.domain]
        result = append(result, metrics)
        count++
    }
    
    return result, nil
}

func (s *MemoryStorage) GetDomainMetrics(ctx context.Context, domain string) (model.DomainMetrics, bool, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    metrics, exists := s.domains[domain]
    if !exists {
        return model.DomainMetrics{}, false, nil
    }
    
    return metrics, true, nil
}