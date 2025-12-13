package mysql

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LimitAmountModel = (*customLimitAmountModel)(nil)

type (
	// LimitAmountModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLimitAmountModel.
	LimitAmountModel interface {
		limitAmountModel
	}

	customLimitAmountModel struct {
		*defaultLimitAmountModel
	}
)

// NewLimitAmountModel returns a model for the database table.
func NewLimitAmountModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) LimitAmountModel {
	return &customLimitAmountModel{
		defaultLimitAmountModel: newLimitAmountModel(conn, c, opts...),
	}
}
