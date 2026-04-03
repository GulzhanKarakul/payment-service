package domain

import "errors"

var (
	ErrNotFound                 = errors.New("not found")
	ErrInvalidInput             = errors.New("invalid input")
	ErrInternalError            = errors.New("internal error")
	ErrClientNotFound           = errors.New("client not found")
	ErrClientIsNotActive        = errors.New("client is not active")
	ErrClientAlreadyExist       = errors.New("client already exists")
	ErrBusinessNotFound         = errors.New("business not found")
	ErrBusinessIsNotActive      = errors.New("business is not active")
	ErrBusinessAlreadyExist     = errors.New("business already exists")
	ErrTransactionNotFound      = errors.New("transaction not found")
	ErrTransactionCancelled     = errors.New("transaction already cancelled")
	ErrTransactionProcessed     = errors.New("transaction already processed")
	ErrInsufficientBonusBalance = errors.New("insufficient bonus balance")
	ErrBonusSettingsNotFound    = errors.New("bonus settings not found")
	ErrBonusSettingsIsNotActive = errors.New("bonus settings is not active")
)
