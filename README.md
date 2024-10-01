# Purpose
許多密室逃脫及實境解謎遊戲使用 LineBot 平台作為解謎遊戲的輔助工具，其中最常見的方式是業者透過編輯遊戲劇情，讓玩家與 LineBot 對話，並以輸入關鍵字的方式進行闖關。

目前市面上有一些付費工具可以幫助使用者編輯劇本並實現 LineBot 的後端操作，這些工具功能強大，但如果只是製作遊戲劇本，並不需要這麼多的功能。而且這些工具通常需要每月訂閱才能持續使用，對於一些小型工作室來說，這是一筆不小的負擔。

受到某密室逃脫業者的委託，我希望能開發一款簡單易用的 LineBot 劇本製作工具，後續稱為 LineBot Script Creator (LSC)，並搭配製作相應的 LineBot 後端程式，希望能降低使用 LineBot 機器人的成本。

# Objectives
## LineBot Script Creator(LSC)
1. 建立編輯劇本的頁面。
2. 簡易的編輯界面。
3. 可依劇本彈性編輯且自由搭配LineBot中的功能(Message, QuickReply...)。
4. 將編寫的劇本輸入Database中供後續LineBot使用。

## LineBot Backend Program
1. 抓取Database中的劇本。
2. 根據劇本及使用者的狀況回應對應的資訊。
3. 記錄目前使用者的遊玩階段。

# Tool
Golang, Gin, Gorm, PostgreSQL, HTML, CSS, Javascript, GoJS, Docker

# Function introduction
## Data Structure
為了實現彈性編輯遊戲劇本及關卡設置，程式中採用靈活的資料結構——鏈結串列（Linked List）。使用鏈結串列的好處包括：
1. 動態大小：鏈結串列可以根據需要動態增減，不必提前分配固定的記憶體空間，這使得管理不同數量的劇本和關卡變得更加便利。
2. 插入與刪除效率高：在鏈結串列中，插入和刪除節點的操作時間複雜度為 O(1)，這比起陣列在中間位置插入或刪除需要 O(n) 的時間更具優勢。
3. 靈活性：鏈結串列的結構使得在編輯劇本時，可以輕鬆地重新排列或重組不同的節點，方便玩家體驗多樣化的遊戲內容。

## Node結構體
以下是Node結構體的定義，這是遊戲劇本編輯系統中用來表示每個劇本節點的重要數據結構
```go
type Node struct {
	ID           int           `gorm:"primaryKey;autoIncrement"` // 節點的唯一識別碼，自動增量
	Title        string        `gorm:"size:255;not null"`        // 節點的標題
	Type         string        `gorm:"size:255;not null"`        // 節點的類型
	Range        IntArray      `gorm:"type:jsonb;not null"`      // 節點的可用範圍，用來記錄目前Node的資料放在其他Table中的哪幾個ID
	PreviousNode int           `gorm:"index"`                    // 前一個節點的 ID，用於鏈接節點
	NextNode     int           `gorm:"index"`                    // 下一個節點的 ID，用於鏈接節點
	LocX         int           `gorm:"size:50;default:0" json:"locX"` // 節點在畫布上的 X 座標
	LocY         int           `gorm:"size:50;default:0" json:"locY"` // 節點在畫布上的 Y 座標
	Messages     []Message     `gorm:"foreignKey:NodeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // 與節點相關的訊息
	QuickReplies []QuickReply  `gorm:"foreignKey:NodeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // 快速回覆選項
	KeyDecisions []KeyDecision `gorm:"foreignKey:NodeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // 關鍵決策
	TagDecisions []TagDecision `gorm:"foreignKey:NodeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // 標籤決策
	Randoms      []Random      `gorm:"foreignKey:NodeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // 隨機選項
}
```
結構體字段說明：
- ID: 節點的唯一識別碼，採用自動增量的方式確保唯一性。
- Title: 節點的標題，用於描述節點的基本信息。
- Type: 節點的類型，描述節點的具體用途。
- Range: 節點的可用範圍，用來記錄目前Node的資料放在其他Table中的哪幾個ID。比方說，目前的Node的Type是Message，Range是1, 2, 8，則此Node的資料就存在Message table中index是1, 2, 8中。
- PreviousNode 和 NextNode: 用於鏈接節點之間的前後關係，使劇本能夠形成有序的流程。
- LocX 和 LocY: 節點在畫布中的位置，便於在可視化設計工具中進行管理。
- Messages、QuickReplies、KeyDecisions、TagDecisions、Randoms: 這些屬性用於管理與節點相關的各種互動元素，包括訊息、快速回覆、關鍵決策等。
這個結構的設計使得劇本中的每個節點都可以被靈活地定義和控制，支持不同遊戲情境下的需求。

### Node 流程圖示
下圖為 Node 在遊戲劇本中的實際應用示意圖，展示了各個節點如何進行鏈接和交互：
![nodeFlow](images/nodeFlow.png)

## UserSession 結構體
除了節點結構，還需要對每個使用者進行管理，以記錄他們在遊戲中的狀態。以下是 UserSession 的定義：
```
type UserSession struct {
	Index     int    `gorm:"primaryKey;autoIncrement"` // 使用者紀錄的唯一索引，自動增量
	UserID    string `gorm:"size:255;not null"`        // 使用者的唯一識別碼
	CurrentID int    `gorm:"not null"`                 // 使用者目前所在的節點 ID
	Time      time.Time                                // 記錄上次互動發生的時間
}
```
結構體字段說明：
- Index: 此字段作為使用者紀錄的唯一索引，並設置為自動增量。
- UserID: 用於標識使用者，確保每個使用者的數據是唯一且可追溯的。
- CurrentID: 記錄使用者目前所在的節點位置，便於追蹤使用者的進度。
- Time: 互動發生的時間，用於追溯和分析使用者的行為。
UserSession 結構可以幫助管理每個使用者的進度，確保玩家在遊戲過程中的體驗能夠被持續記錄和查找。

### User 行為流程圖示
![userActionFlow](images/userActionFlow.png)
執行流程為：
1. 使用者輸入資訊或執行動作。
2. 到User Tabel中依據UserID查詢目前使用者所在Node。
3. 再到對應的Node中查詢。
4. 依據當前Node的Type以及Range來取得對應的Response。