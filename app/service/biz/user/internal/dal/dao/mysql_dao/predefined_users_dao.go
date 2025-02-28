/*
 * WARNING! All changes made in this file will be lost!
 *   Created from by 'dalgen'
 *
 * Copyright (c) 2021-present,  Teamgram Studio (https://teamgram.io).
 *  All rights reserved.
 *
 * Author: teamgramio (teamgram.io@gmail.com)
 */

package mysql_dao

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/teamgram/marmota/pkg/stores/sqlx"
	"github.com/teamgram/teamgram-server/app/service/biz/user/internal/dal/dataobject"

	"github.com/zeromicro/go-zero/core/logx"
)

var _ *sql.Result

type PredefinedUsersDAO struct {
	db *sqlx.DB
}

func NewPredefinedUsersDAO(db *sqlx.DB) *PredefinedUsersDAO {
	return &PredefinedUsersDAO{db}
}

// Insert
// insert into predefined_users(first_name, last_name, username, phone, code, verified) values (:first_name, :last_name, :username, :phone, :code, :verified)
// TODO(@benqi): sqlmap
func (dao *PredefinedUsersDAO) Insert(ctx context.Context, do *dataobject.PredefinedUsersDO) (lastInsertId, rowsAffected int64, err error) {
	var (
		query = "insert into predefined_users(first_name, last_name, username, phone, code, verified) values (:first_name, :last_name, :username, :phone, :code, :verified)"
		r     sql.Result
	)

	r, err = dao.db.NamedExec(ctx, query, do)
	if err != nil {
		logx.WithContext(ctx).Errorf("namedExec in Insert(%v), error: %v", do, err)
		return
	}

	lastInsertId, err = r.LastInsertId()
	if err != nil {
		logx.WithContext(ctx).Errorf("lastInsertId in Insert(%v)_error: %v", do, err)
		return
	}
	rowsAffected, err = r.RowsAffected()
	if err != nil {
		logx.WithContext(ctx).Errorf("rowsAffected in Insert(%v)_error: %v", do, err)
	}

	return
}

// InsertTx
// insert into predefined_users(first_name, last_name, username, phone, code, verified) values (:first_name, :last_name, :username, :phone, :code, :verified)
// TODO(@benqi): sqlmap
func (dao *PredefinedUsersDAO) InsertTx(tx *sqlx.Tx, do *dataobject.PredefinedUsersDO) (lastInsertId, rowsAffected int64, err error) {
	var (
		query = "insert into predefined_users(first_name, last_name, username, phone, code, verified) values (:first_name, :last_name, :username, :phone, :code, :verified)"
		r     sql.Result
	)

	r, err = tx.NamedExec(query, do)
	if err != nil {
		logx.WithContext(tx.Context()).Errorf("namedExec in Insert(%v), error: %v", do, err)
		return
	}

	lastInsertId, err = r.LastInsertId()
	if err != nil {
		logx.WithContext(tx.Context()).Errorf("lastInsertId in Insert(%v)_error: %v", do, err)
		return
	}
	rowsAffected, err = r.RowsAffected()
	if err != nil {
		logx.WithContext(tx.Context()).Errorf("rowsAffected in Insert(%v)_error: %v", do, err)
	}

	return
}

// SelectByPhone
// select id, phone, first_name, last_name, username, code, verified, registered_user_id from predefined_users where phone = :phone and deleted = 0 limit 1
// TODO(@benqi): sqlmap
func (dao *PredefinedUsersDAO) SelectByPhone(ctx context.Context, phone string) (rValue *dataobject.PredefinedUsersDO, err error) {
	var (
		query = "select id, phone, first_name, last_name, username, code, verified, registered_user_id from predefined_users where phone = ? and deleted = 0 limit 1"
		rows  *sqlx.Rows
	)
	rows, err = dao.db.Query(ctx, query, phone)

	if err != nil {
		logx.WithContext(ctx).Errorf("queryx in SelectByPhone(_), error: %v", err)
		return
	}

	defer rows.Close()

	do := &dataobject.PredefinedUsersDO{}
	if rows.Next() {
		// TODO(@benqi): not use reflect
		err = rows.StructScan(do)
		if err != nil {
			logx.WithContext(ctx).Errorf("structScan in SelectByPhone(_), error: %v", err)
			return
		} else {
			rValue = do
		}
	}

	return
}

// SelectPredefinedUsersAll
// select id, phone, first_name, last_name, username, code, verified, registered_user_id from predefined_users where deleted = 0 order by username asc
// TODO(@benqi): sqlmap
func (dao *PredefinedUsersDAO) SelectPredefinedUsersAll(ctx context.Context) (rList []dataobject.PredefinedUsersDO, err error) {
	var (
		query = "select id, phone, first_name, last_name, username, code, verified, registered_user_id from predefined_users where deleted = 0 order by username asc"
		rows  *sqlx.Rows
	)
	rows, err = dao.db.Query(ctx, query)

	if err != nil {
		logx.WithContext(ctx).Errorf("queryx in SelectPredefinedUsersAll(_), error: %v", err)
		return
	}

	defer rows.Close()

	var values []dataobject.PredefinedUsersDO
	for rows.Next() {
		v := dataobject.PredefinedUsersDO{}

		// TODO(@benqi): not use reflect
		err = rows.StructScan(&v)
		if err != nil {
			logx.WithContext(ctx).Errorf("structScan in SelectPredefinedUsersAll(_), error: %v", err)
			return
		}
		values = append(values, v)
	}
	rList = values

	return
}

// SelectPredefinedUsersAllWithCB
// select id, phone, first_name, last_name, username, code, verified, registered_user_id from predefined_users where deleted = 0 order by username asc
// TODO(@benqi): sqlmap
func (dao *PredefinedUsersDAO) SelectPredefinedUsersAllWithCB(ctx context.Context, cb func(i int, v *dataobject.PredefinedUsersDO)) (rList []dataobject.PredefinedUsersDO, err error) {
	var (
		query = "select id, phone, first_name, last_name, username, code, verified, registered_user_id from predefined_users where deleted = 0 order by username asc"
		rows  *sqlx.Rows
	)
	rows, err = dao.db.Query(ctx, query)

	if err != nil {
		logx.WithContext(ctx).Errorf("queryx in SelectPredefinedUsersAll(_), error: %v", err)
		return
	}

	defer func() {
		rows.Close()
		if err == nil && cb != nil {
			for i := 0; i < len(rList); i++ {
				cb(i, &rList[i])
			}
		}
	}()

	var values []dataobject.PredefinedUsersDO
	for rows.Next() {
		v := dataobject.PredefinedUsersDO{}

		// TODO(@benqi): not use reflect
		err = rows.StructScan(&v)
		if err != nil {
			logx.WithContext(ctx).Errorf("structScan in SelectPredefinedUsersAll(_), error: %v", err)
			return
		}
		values = append(values, v)
	}
	rList = values

	return
}

// Delete
// update predefined_users set deleted = 0 where phone = :phone
// TODO(@benqi): sqlmap
func (dao *PredefinedUsersDAO) Delete(ctx context.Context, phone string) (rowsAffected int64, err error) {
	var (
		query   = "update predefined_users set deleted = 0 where phone = ?"
		rResult sql.Result
	)
	rResult, err = dao.db.Exec(ctx, query, phone)

	if err != nil {
		logx.WithContext(ctx).Errorf("exec in Delete(_), error: %v", err)
		return
	}

	rowsAffected, err = rResult.RowsAffected()
	if err != nil {
		logx.WithContext(ctx).Errorf("rowsAffected in Delete(_), error: %v", err)
	}

	return
}

// update predefined_users set deleted = 0 where phone = :phone
// DeleteTx
// TODO(@benqi): sqlmap
func (dao *PredefinedUsersDAO) DeleteTx(tx *sqlx.Tx, phone string) (rowsAffected int64, err error) {
	var (
		query   = "update predefined_users set deleted = 0 where phone = ?"
		rResult sql.Result
	)
	rResult, err = tx.Exec(query, phone)

	if err != nil {
		logx.WithContext(tx.Context()).Errorf("exec in Delete(_), error: %v", err)
		return
	}

	rowsAffected, err = rResult.RowsAffected()
	if err != nil {
		logx.WithContext(tx.Context()).Errorf("rowsAffected in Delete(_), error: %v", err)
	}

	return
}

// Update
// update predefined_users set %s where phone = :phone
// TODO(@benqi): sqlmap
func (dao *PredefinedUsersDAO) Update(ctx context.Context, cMap map[string]interface{}, phone string) (rowsAffected int64, err error) {
	names := make([]string, 0, len(cMap))
	aValues := make([]interface{}, 0, len(cMap))
	for k, v := range cMap {
		names = append(names, k+" = ?")
		aValues = append(aValues, v)
	}

	var (
		query   = fmt.Sprintf("update predefined_users set %s where phone = ?", strings.Join(names, ", "))
		rResult sql.Result
	)

	aValues = append(aValues, phone)

	rResult, err = dao.db.Exec(ctx, query, aValues...)

	if err != nil {
		logx.WithContext(ctx).Errorf("exec in Update(_), error: %v", err)
		return
	}

	rowsAffected, err = rResult.RowsAffected()
	if err != nil {
		logx.WithContext(ctx).Errorf("rowsAffected in Update(_), error: %v", err)
	}

	return
}

// UpdateTx
// update predefined_users set %s where phone = :phone
// TODO(@benqi): sqlmap
func (dao *PredefinedUsersDAO) UpdateTx(tx *sqlx.Tx, cMap map[string]interface{}, phone string) (rowsAffected int64, err error) {
	names := make([]string, 0, len(cMap))
	aValues := make([]interface{}, 0, len(cMap))
	for k, v := range cMap {
		names = append(names, k+" = ?")
		aValues = append(aValues, v)
	}

	var (
		query   = fmt.Sprintf("update predefined_users set %s where phone = ?", strings.Join(names, ", "))
		rResult sql.Result
	)

	aValues = append(aValues, phone)

	rResult, err = tx.Exec(query, aValues...)

	if err != nil {
		logx.WithContext(tx.Context()).Errorf("exec in Update(_), error: %v", err)
		return
	}

	rowsAffected, err = rResult.RowsAffected()
	if err != nil {
		logx.WithContext(tx.Context()).Errorf("rowsAffected in Update(_), error: %v", err)
	}

	return
}
