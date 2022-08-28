package statistics

import (
	"sync"

	"github.com/Eliad-S/Permutation_web_service/db"
)

type Statistics struct {
	TotalWords          uint32  `json:"totalWords"`
	TotalRequests       uint32  `json:"totalRequests"`
	AvgProcessingTimeNs float32 `json:"avgProcessingTimeNs"`
}

const Stat_table = "Statistics"

var stats Statistics
var mu sync.RWMutex // guards balance

func GetStats() Statistics {
	mu.RLock() // readers lock
	defer mu.RUnlock()
	return stats
}

func Init() {
	stats = Statistics{0, 0, 0}
	Set_TotalWords()
}

// func Import_stat_from_db() {
// 	stats = db.Import_stats(Stat_table)
// }

func Set_TotalWords() {
	//mutex
	mu.Lock()
	defer mu.Unlock()
	stats.TotalWords = db.Get_total_words()
}

func Get_TotalWords() uint32 {
	//mutex
	mu.RLock() // readers lock
	defer mu.RUnlock()
	return stats.TotalWords
}

func Inc_requests(requestProcessingTime int64) {
	//mutex
	mu.Lock()
	defer mu.Unlock()
	curTotalRequests := stats.TotalRequests
	curAvgProcessingTimeNs := stats.AvgProcessingTimeNs
	stats.TotalRequests += 1
	stats.AvgProcessingTimeNs = (curAvgProcessingTimeNs*float32(curTotalRequests) + float32(requestProcessingTime)) / float32(stats.TotalRequests)
}
