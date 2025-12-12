package mysql

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ RiskRecordModel = (*customRiskRecordModel)(nil)

type (
	// RiskRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRiskRecordModel.
	RiskRecordModel interface {
		riskRecordModel
		withSession(session sqlx.Session) RiskRecordModel
	}

	customRiskRecordModel struct {
		*defaultRiskRecordModel
	}
)

// NewRiskRecordModel returns a model for the database table.
func NewRiskRecordModel(conn sqlx.SqlConn) RiskRecordModel {
	return &customRiskRecordModel{
		defaultRiskRecordModel: newRiskRecordModel(conn),
	}
}

func (m *customRiskRecordModel) withSession(session sqlx.Session) RiskRecordModel {
	return NewRiskRecordModel(sqlx.NewSqlConnFromSession(session))
}
