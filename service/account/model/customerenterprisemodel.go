package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ CustomerEnterpriseModel = (*customCustomerEnterpriseModel)(nil)

type (
	// CustomerEnterpriseModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCustomerEnterpriseModel.
	CustomerEnterpriseModel interface {
		customerEnterpriseModel
		withSession(session sqlx.Session) CustomerEnterpriseModel
	}

	customCustomerEnterpriseModel struct {
		*defaultCustomerEnterpriseModel
	}
)

// NewCustomerEnterpriseModel returns a model for the database table.
func NewCustomerEnterpriseModel(conn sqlx.SqlConn) CustomerEnterpriseModel {
	return &customCustomerEnterpriseModel{
		defaultCustomerEnterpriseModel: newCustomerEnterpriseModel(conn),
	}
}

func (m *customCustomerEnterpriseModel) withSession(session sqlx.Session) CustomerEnterpriseModel {
	return NewCustomerEnterpriseModel(sqlx.NewSqlConnFromSession(session))
}
