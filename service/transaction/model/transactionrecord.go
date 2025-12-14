package model

import (
	"context"
	"fmt"
	"time"
)

func (m *defaultTransactionRecordModel) FindRange(ctx context.Context, account string, startTime, endTime time.Time) ([]*TransactionRecord, error) {
	query := fmt.Sprintf("select %s from %s where (`account_from` = ? or `account_to` = ?) and `created_at` between ? and ? order by `created_at` desc",
		transactionRecordRows, m.table)
	var resp []*TransactionRecord
	err := m.conn.QueryRowsCtx(ctx, &resp, query, account, account, startTime, endTime)
	return resp, err
}
