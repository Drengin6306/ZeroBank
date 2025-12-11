package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ TransactionRecordModel = (*customTransactionRecordModel)(nil)

type (
	// TransactionRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customTransactionRecordModel.
	TransactionRecordModel interface {
		transactionRecordModel
		withSession(session sqlx.Session) TransactionRecordModel
	}

	customTransactionRecordModel struct {
		*defaultTransactionRecordModel
	}
)

// NewTransactionRecordModel returns a model for the database table.
func NewTransactionRecordModel(conn sqlx.SqlConn) TransactionRecordModel {
	return &customTransactionRecordModel{
		defaultTransactionRecordModel: newTransactionRecordModel(conn),
	}
}

func (m *customTransactionRecordModel) withSession(session sqlx.Session) TransactionRecordModel {
	return NewTransactionRecordModel(sqlx.NewSqlConnFromSession(session))
}
