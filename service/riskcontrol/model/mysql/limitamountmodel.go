package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ LimitAmountModel = (*customLimitAmountModel)(nil)

type (
	// LimitAmountModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLimitAmountModel.
	LimitAmountModel interface {
		limitAmountModel
		withSession(session sqlx.Session) LimitAmountModel
	}

	customLimitAmountModel struct {
		*defaultLimitAmountModel
	}
)

// NewLimitAmountModel returns a model for the database table.
func NewLimitAmountModel(conn sqlx.SqlConn) LimitAmountModel {
	return &customLimitAmountModel{
		defaultLimitAmountModel: newLimitAmountModel(conn),
	}
}

func (m *customLimitAmountModel) withSession(session sqlx.Session) LimitAmountModel {
	return NewLimitAmountModel(sqlx.NewSqlConnFromSession(session))
}
