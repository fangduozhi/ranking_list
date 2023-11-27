package ranking_list

import (
	"container/list"
	"time"
)

type PlayerScore struct {
	UserID     string
	Name       string
	Score      int
	ScoreTime  time.Time
	Level      int
	CurrentPos int
}

// 分数到链表节点的映射
type ScoreLinkedLists map[int]*list.List

// 用户ID到用户分数信息的映射
type PlayerInformation map[string]*list.Element
