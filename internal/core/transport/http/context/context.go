package core_http_context

import "context"

type contextKey string

const accountIDKey contextKey = "account_id"

func WithAccountID(ctx context.Context, accountID int) context.Context {
	return context.WithValue(ctx, accountIDKey, accountID)
}

func AccountIDFromContext(ctx context.Context) (int, bool) {
	accountID, ok := ctx.Value(accountIDKey).(int)
	return accountID, ok
}
