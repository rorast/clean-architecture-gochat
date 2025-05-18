# GoChat 錯誤處理與緩存機制集成文檔

本文檔說明了GoChat項目中新的錯誤處理系統與Redis緩存機制的集成方法。

## 錯誤處理系統架構

新的錯誤處理系統採用分層設計，主要包括以下幾個部分：

### 1. 基礎錯誤處理（pkg/errors）

- **errors.go**: 定義了核心錯誤接口和實現
- **http.go**: 提供HTTP相關的錯誤處理功能
- **log.go**: 提供錯誤日誌記錄功能

### 2. 錯誤碼定義（internal/common/enum）

- **errorcode.go**: 定義了系統中所有錯誤碼和對應的錯誤信息
- **appstatus.go**: 定義了應用狀態碼和對應信息

### 3. 應用錯誤工廠（internal/errors）

- **errors.go**: 提供便捷的應用錯誤創建函數

## Redis緩存機制集成

我們在專案中加入了帶有新錯誤處理機制的Redis緩存實現：

### 1. 增強版Redis緩存（infrastructure/redis）

- **enhanced_message_cache.go**: 使用新的錯誤處理系統的消息緩存實現
- **redis.go**: 改進的Redis連接管理，使用新的錯誤處理系統

### 2. 測試套件（test/）

- **enhanced_redis_cache_test.go**: 測試改進的Redis緩存功能和錯誤處理

### 3. 加強版用例（internal/usecases/message）

- **enhanced_message_usecase.go**: 展示如何在業務邏輯中使用新的錯誤處理系統和緩存機制

## 使用示例

以下是展示如何使用新錯誤處理系統的代碼片段：

```go
// 使用新的錯誤處理系統創建錯誤
if message == nil {
    return appErrors.New(enum.ErrInvalidInput, "消息不能為空")
}

// 封裝現有錯誤
if err := messageRepo.Create(ctx, message); err != nil {
    return appErrors.Wrap(err, enum.ErrMessageSendFailed, map[string]interface{}{
        "userId":   message.UserId,
        "targetId": message.TargetId,
        "content":  message.Content,
    })
}

// 僅記錄錯誤但不中斷流程
if err := messageCache.StorePrivateMessage(ctx, message); err != nil {
    appErrors.WithDevMessage(err, "緩存私人消息失敗").LogError()
}
```

## 異常處理流程

1. **創建異常**：使用 `appErrors.New()` 或 `appErrors.Wrap()` 創建應用錯誤
2. **日誌記錄**：通過 `LogError()` 或設置自定義 logger 進行日誌記錄
3. **錯誤轉換**：使用 `ToResponse()` 將錯誤轉換為HTTP響應

## 未來擴展

未來可以考慮以下擴展：

1. 添加分佈式跟踪ID，關聯相關錯誤
2. 與監控系統集成，提供錯誤報警功能
3. 完善國際化信息支持，提供多語言錯誤提示

## 最佳實踐

1. **統一錯誤碼**：始終使用 `enum` 包中定義的錯誤碼
2. **詳細上下文**：錯誤信息中包含足夠的上下文信息以便排查
3. **區分用戶/開發者信息**：用戶友好信息與開發者調試信息分離
4. **緩存錯誤處理**：緩存操作失敗不應影響主業務流程
5. **非阻塞寫入**：使用goroutine進行非關鍵的緩存寫入操作 