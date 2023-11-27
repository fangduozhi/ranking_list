package ranking_list

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

// Leaderboard 代表排行榜
type Leaderboard struct {
	scoreLinkedLists  ScoreLinkedLists
	playerInformation PlayerInformation
	Lock              sync.Mutex
}

// 新建排行榜
func NewLeaderboard() *Leaderboard {
	return &Leaderboard{
		scoreLinkedLists:  make(ScoreLinkedLists),
		playerInformation: make(PlayerInformation),
	}
}

// 更新或者新增玩家得分
func (lb *Leaderboard) UpdatePlayerScore(userID, name string, score int, level int) {
	currentTime := time.Now()
	lb.Lock.Lock()
	defer lb.Lock.Unlock()

	// 如果用户已经存在，移除旧分数对应的节点
	if playerElem, exists := lb.playerInformation[userID]; exists {
		lb.scoreLinkedLists[playerElem.Value.(*PlayerScore).Score].Remove(playerElem)
	}

	// 创建新的玩家分数实例
	newPlayerScore := &PlayerScore{
		UserID:    userID,
		Name:      name,
		Score:     score,
		ScoreTime: currentTime,
		Level:     level,
	}

	// 获取或创建分数对应的链表
	scoreList, exists := lb.scoreLinkedLists[score]
	if !exists {
		scoreList = list.New()
		lb.scoreLinkedLists[score] = scoreList
	}

	var laterUserCount int
	for sc := 0; sc <= score; sc-- {
		if scList, ex := lb.scoreLinkedLists[sc]; ex {
			laterUserCount += scList.Len()
		}
	}
	// 将用户添加到分数对应的链表中
	playerElem := scoreList.PushBack(newPlayerScore)
	// 当前用户的排名
	newPlayerScore.CurrentPos = laterUserCount + scoreList.Len()

	// 更新该用户的信息
	lb.playerInformation[userID] = playerElem
}

// 输出排行榜
func (lb *Leaderboard) PrintLeaderboard() {
	// 由于分数有范围，我们按分数从高到低遍历
	for score := 10000; score >= 0; score-- {
		if scoreList, exists := lb.scoreLinkedLists[score]; exists {
			for elem := scoreList.Front(); elem != nil; elem = elem.Next() {
				playerScore := elem.Value.(*PlayerScore)
				fmt.Printf("UserID: %s, Name: %s, Score: %d, Time: %s, Level: %d\n",
					playerScore.UserID, playerScore.Name, playerScore.Score, playerScore.ScoreTime.Format(time.RFC3339), playerScore.Level)
			}
		}
	}
}

func (lb *Leaderboard) GetUserRangePlayer(userID string, scope int) []*PlayerScore {
	userScore, ok := lb.playerInformation[userID]
	if !ok {
		return nil
	}
	userScore.Prev()

	return nil
}

func main() {
	leaderboard := NewLeaderboard()

	leaderboard.UpdatePlayerScore("1", "Alice", 9500, 5)
	leaderboard.UpdatePlayerScore("2", "Bob", 9400, 4)
	leaderboard.UpdatePlayerScore("3", "Charlie", 9500, 7)

	leaderboard.PrintLeaderboard()
}
