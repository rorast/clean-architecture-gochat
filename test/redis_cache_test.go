package test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

type Message struct {
	ID        string    `json:"id"`
	UserID    int64     `json:"user_id"`
	TargetID  int64     `json:"target_id"`
	Content   string    `json:"content"`
	Type      int       `json:"type"`
	Media     int       `json:"media"`
	CreatedAt time.Time `json:"created_at"`
}

func TestRedisMessageCache(t *testing.T) {
	// 連接到 Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6378", // 使用 Docker 映射的端口
		Password: "",               // 如果有設置密碼，請在這裡填寫
		DB:       0,
	})

	ctx := context.Background()

	// 測試 1: 測試連接
	t.Run("Test Redis Connection", func(t *testing.T) {
		_, err := rdb.Ping(ctx).Result()
		assert.NoError(t, err, "Redis 連接應該成功")
	})

	// 測試 2: 測試消息緩存
	t.Run("Test Message Cache", func(t *testing.T) {
		// 創建測試消息
		msg := Message{
			ID:        "test_msg_1",
			UserID:    1,
			TargetID:  2,
			Content:   "Hello, this is a test message",
			Type:      1,
			Media:     1,
			CreatedAt: time.Now(),
		}

		// 將消息轉換為 JSON
		msgJSON, err := json.Marshal(msg)
		assert.NoError(t, err, "消息序列化應該成功")

		// 生成緩存 key
		cacheKey := fmt.Sprintf("chat:msg:%d:%d", msg.UserID, msg.TargetID)

		// 存儲消息到 Redis
		err = rdb.LPush(ctx, cacheKey, msgJSON).Err()
		assert.NoError(t, err, "消息應該成功存儲到 Redis")

		// 從 Redis 讀取消息
		result, err := rdb.LRange(ctx, cacheKey, 0, 0).Result()
		assert.NoError(t, err, "應該能夠從 Redis 讀取消息")
		assert.NotEmpty(t, result, "讀取的結果不應該為空")

		// 解析讀取的消息
		var retrievedMsg Message
		err = json.Unmarshal([]byte(result[0]), &retrievedMsg)
		assert.NoError(t, err, "應該能夠解析消息")

		// 驗證消息內容
		assert.Equal(t, msg.ID, retrievedMsg.ID, "消息 ID 應該匹配")
		assert.Equal(t, msg.Content, retrievedMsg.Content, "消息內容應該匹配")
	})

	// 測試 3: 測試消息過期
	t.Run("Test Message Expiration", func(t *testing.T) {
		cacheKey := "chat:msg:test:expiration"

		// 設置帶過期時間的消息
		err := rdb.Set(ctx, cacheKey, "test message", 1*time.Second).Err()
		assert.NoError(t, err, "應該能夠設置帶過期時間的消息")

		// 立即讀取
		val, err := rdb.Get(ctx, cacheKey).Result()
		assert.NoError(t, err, "應該能夠讀取消息")
		assert.Equal(t, "test message", val, "消息內容應該匹配")

		// 等待過期
		time.Sleep(2 * time.Second)

		// 再次讀取
		_, err = rdb.Get(ctx, cacheKey).Result()
		assert.Error(t, err, "消息應該已經過期")
		assert.Equal(t, redis.Nil, err, "應該返回 redis.Nil 錯誤")
	})

	// 測試 4: 測試批量消息操作
	t.Run("Test Batch Message Operations", func(t *testing.T) {
		cacheKey := "chat:msg:batch:test"

		// 清除可能存在的舊數據
		rdb.Del(ctx, cacheKey)

		// 批量添加消息
		messages := []Message{
			{
				ID:        "msg1",
				UserID:    1,
				TargetID:  2,
				Content:   "Message 1",
				Type:      1,
				Media:     1,
				CreatedAt: time.Now(),
			},
			{
				ID:        "msg2",
				UserID:    1,
				TargetID:  2,
				Content:   "Message 2",
				Type:      1,
				Media:     1,
				CreatedAt: time.Now(),
			},
		}

		// 將消息添加到列表
		for _, msg := range messages {
			msgJSON, err := json.Marshal(msg)
			assert.NoError(t, err, "消息序列化應該成功")
			err = rdb.LPush(ctx, cacheKey, msgJSON).Err()
			assert.NoError(t, err, "應該能夠添加消息到列表")
		}

		// 獲取列表長度
		length, err := rdb.LLen(ctx, cacheKey).Result()
		assert.NoError(t, err, "應該能夠獲取列表長度")
		assert.Equal(t, int64(len(messages)), length, "列表長度應該匹配")

		// 分頁獲取消息
		result, err := rdb.LRange(ctx, cacheKey, 0, -1).Result()
		assert.NoError(t, err, "應該能夠獲取所有消息")
		assert.Equal(t, len(messages), len(result), "獲取的消息數量應該匹配")

		// 清理測試數據
		rdb.Del(ctx, cacheKey)
	})

	// 測試 5: 查看和操作緩存資料
	t.Run("Test View and Manipulate Cache Data", func(t *testing.T) {
		// 清理之前的測試數據
		pattern := "chat:msg:*"
		iter := rdb.Scan(ctx, 0, pattern, 0).Iterator()
		for iter.Next(ctx) {
			rdb.Del(ctx, iter.Val())
		}

		// 1. 插入一些測試數據
		testMessages := []Message{
			{
				ID:        "test1",
				UserID:    1,
				TargetID:  2,
				Content:   "Hello from user 1 to 2",
				Type:      1,
				Media:     1,
				CreatedAt: time.Now(),
			},
			{
				ID:        "test2",
				UserID:    2,
				TargetID:  1,
				Content:   "Reply from user 2 to 1",
				Type:      1,
				Media:     1,
				CreatedAt: time.Now(),
			},
			{
				ID:        "test3",
				UserID:    1,
				TargetID:  3,
				Content:   "Hello from user 1 to 3",
				Type:      1,
				Media:     1,
				CreatedAt: time.Now(),
			},
		}

		// 2. 將消息存入不同的聊天緩存中
		for _, msg := range testMessages {
			// 為每對用戶創建一個唯一的緩存鍵
			cacheKey := fmt.Sprintf("chat:msg:%d:%d", msg.UserID, msg.TargetID)
			msgJSON, err := json.Marshal(msg)
			assert.NoError(t, err, "消息序列化應該成功")
			err = rdb.LPush(ctx, cacheKey, msgJSON).Err()
			assert.NoError(t, err, "消息應該成功存儲到 Redis")
		}

		// 3. 列出所有聊天緩存鍵
		fmt.Println("\n=== 所有聊天緩存鍵 ===")
		iter = rdb.Scan(ctx, 0, pattern, 0).Iterator()
		for iter.Next(ctx) {
			key := iter.Val()
			fmt.Printf("找到緩存鍵: %s\n", key)

			// 4. 獲取每個緩存的消息數量
			length, err := rdb.LLen(ctx, key).Result()
			assert.NoError(t, err)
			fmt.Printf("緩存 %s 中的消息數量: %d\n", key, length)

			// 5. 讀取該緩存中的所有消息
			messages, err := rdb.LRange(ctx, key, 0, -1).Result()
			assert.NoError(t, err)

			fmt.Printf("緩存 %s 中的消息內容:\n", key)
			for i, msgJSON := range messages {
				var msg Message
				err = json.Unmarshal([]byte(msgJSON), &msg)
				assert.NoError(t, err)
				fmt.Printf("  %d) 發送者: %d, 接收者: %d, 內容: %s\n",
					i+1, msg.UserID, msg.TargetID, msg.Content)
			}
			fmt.Println()
		}

		// 6. 測試特定用戶對話的緩存
		testUserID := int64(1)
		testTargetID := int64(2)
		cacheKey := fmt.Sprintf("chat:msg:%d:%d", testUserID, testTargetID)

		// 獲取特定對話的消息
		messages, err := rdb.LRange(ctx, cacheKey, 0, -1).Result()
		assert.NoError(t, err)

		fmt.Printf("\n=== 用戶 %d 和用戶 %d 之間的對話 ===\n", testUserID, testTargetID)
		for i, msgJSON := range messages {
			var msg Message
			err = json.Unmarshal([]byte(msgJSON), &msg)
			assert.NoError(t, err)
			fmt.Printf("%d) %s\n", i+1, msg.Content)
		}

		// 7. 測試分頁功能
		fmt.Println("\n=== 測試分頁功能 ===")
		pageSize := int64(2)
		totalMessages, err := rdb.LLen(ctx, cacheKey).Result()
		assert.NoError(t, err)

		for page := int64(0); page*pageSize < totalMessages; page++ {
			start := page * pageSize
			end := start + pageSize - 1

			messages, err := rdb.LRange(ctx, cacheKey, start, end).Result()
			assert.NoError(t, err)

			fmt.Printf("第 %d 頁消息:\n", page+1)
			for i, msgJSON := range messages {
				var msg Message
				err = json.Unmarshal([]byte(msgJSON), &msg)
				assert.NoError(t, err)
				fmt.Printf("  %d) %s\n", start+int64(i)+1, msg.Content)
			}
		}

		// 8. 測試消息更新
		fmt.Println("\n=== 測試消息更新 ===")
		// 獲取第一條消息
		firstMsg, err := rdb.LIndex(ctx, cacheKey, 0).Result()
		assert.NoError(t, err)

		var msgToUpdate Message
		err = json.Unmarshal([]byte(firstMsg), &msgToUpdate)
		assert.NoError(t, err)

		// 更新消息內容
		msgToUpdate.Content = "Updated: " + msgToUpdate.Content
		updatedMsgJSON, err := json.Marshal(msgToUpdate)
		assert.NoError(t, err)

		// 使用 LSET 更新消息
		err = rdb.LSet(ctx, cacheKey, 0, updatedMsgJSON).Err()
		assert.NoError(t, err)

		// 驗證更新
		updatedMsg, err := rdb.LIndex(ctx, cacheKey, 0).Result()
		assert.NoError(t, err)
		var verifyMsg Message
		err = json.Unmarshal([]byte(updatedMsg), &verifyMsg)
		assert.NoError(t, err)
		fmt.Printf("更新後的消息: %s\n", verifyMsg.Content)
	})

	// 清理所有測試數據
	t.Cleanup(func() {
		pattern := "chat:msg:*"
		iter := rdb.Scan(ctx, 0, pattern, 0).Iterator()
		for iter.Next(ctx) {
			rdb.Del(ctx, iter.Val())
		}
		assert.NoError(t, iter.Err(), "清理測試數據時不應該出錯")
	})
}
