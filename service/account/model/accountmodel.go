package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AccountModel = (*customAccountModel)(nil)

type (
	// AccountModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAccountModel.
	AccountModel interface {
		accountModel
		WithSession(session sqlx.Session) AccountModel
		FindOneForUpdate(ctx context.Context, session sqlx.Session, accountId string) (*Account, error)
		Trans(ctx context.Context, fn func(context.Context, sqlx.Session) error) error
	}

	customAccountModel struct {
		*defaultAccountModel
	}
)

// NewAccountModel returns a model for the database table.
func NewAccountModel(conn sqlx.SqlConn) AccountModel {
	return &customAccountModel{
		defaultAccountModel: newAccountModel(conn),
	}
}

func (m *customAccountModel) WithSession(session sqlx.Session) AccountModel {
	return NewAccountModel(sqlx.NewSqlConnFromSession(session))
}

// FindOneForUpdate 使用悲观锁查询账户信息
func (m *customAccountModel) FindOneForUpdate(ctx context.Context, session sqlx.Session, accountId string) (*Account, error) {
	query := fmt.Sprintf("select %s from %s where `account_id` = ? limit 1 for update", accountRows, m.tableName())
	var resp Account
	err := session.QueryRowCtx(ctx, &resp, query, accountId)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// Trans 执行事务
func (m *customAccountModel) Trans(ctx context.Context, fn func(context.Context, sqlx.Session) error) error {
	return m.conn.TransactCtx(ctx, fn)
}
