package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ CustomerIndividualModel = (*customCustomerIndividualModel)(nil)

type (
	// CustomerIndividualModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCustomerIndividualModel.
	CustomerIndividualModel interface {
		customerIndividualModel
		withSession(session sqlx.Session) CustomerIndividualModel
	}

	customCustomerIndividualModel struct {
		*defaultCustomerIndividualModel
	}
)

// NewCustomerIndividualModel returns a model for the database table.
func NewCustomerIndividualModel(conn sqlx.SqlConn) CustomerIndividualModel {
	return &customCustomerIndividualModel{
		defaultCustomerIndividualModel: newCustomerIndividualModel(conn),
	}
}

func (m *customCustomerIndividualModel) withSession(session sqlx.Session) CustomerIndividualModel {
	return NewCustomerIndividualModel(sqlx.NewSqlConnFromSession(session))
}
