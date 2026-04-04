package repository_test

// TestTransactionRepo_CreateWithBonus_WithBonus
// — создать клиента, бизнес, пополнить баланс бизнеса 100000
// — CreateWithBonus(amount=50000, bonusPercent=10.0)
// — ожидаем bonus_accrued=5000
// — проверить client.bonus_balance=5000
// — проверить business.bonus_balance=95000
// — проверить transaction.status=completed

// TestTransactionRepo_CreateWithBonus_ZeroBonus
// — CreateWithBonus(amount=50000, bonusPercent=0)
// — ожидаем bonus_accrued=0
// — client.bonus_balance не изменился

// TestTransactionRepo_CreateWithBonus_InsufficientBalance
// — бизнес с bonus_balance=0
// — CreateWithBonus(amount=50000, bonusPercent=10.0)
// — ожидаем domain.ErrInsufficientBonusBalance
// — проверить что client.bonus_balance НЕ изменился (rollback сработал)

// TestTransactionRepo_GetByID_Success
// TestTransactionRepo_GetByID_NotFound

// TestTransactionRepo_GetByClientID_Pagination
// — создать 5 транзакций для одного клиента
// — GetByClientID(limit=2, offset=0) → 2 транзакции
// — GetByClientID(limit=2, offset=2) → 2 транзакции
// — GetByClientID(limit=2, offset=4) → 1 транзакция
// — проверить порядок (newest first)

// TestTransactionRepo_Cancel_Success
// — создать транзакцию
// — Cancel
// — GetByID → должна вернуть ErrTransactionNotFound (deleted_at IS NULL)

// TestTransactionRepo_Cancel_AlreadyCancelled
// — отменить транзакцию
// — отменить снова
// — ожидаем domain.ErrTransactionCancelled
