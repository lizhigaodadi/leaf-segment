package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/EslRain/leaf-segment/model"
	"github.com/go-sql-driver/mysql"
	"time"
)

func (d *Dao) CreateLeaf(ctx context.Context, leaf *model.Leaf) error {
	now := time.Now().Unix()
	query := fmt.Sprintf("INSERT INTO %s (biz_tag, max_id, step, update_time) VALUES (?,?,?,?,?)", model.LeafTableName())

	res, err := d.sql.ExecContext(ctx, query, leaf.BizTag, leaf.MaxID, leaf.Step, uint64(now))
	if err != nil {
		fmt.Printf("insert leaf failed; leaf: %v; err: %v", leaf, err)
		return err
	}

	_, err = res.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}

func (d *Dao) NextSegment(ctx context.Context, bizTag string) (*model.Leaf, error) {
	//开启事务
	tx, err := d.sql.Begin()
	defer func() {
		if err != nil {
			err = tx.Rollback()
			if err != sql.ErrTxDone && err != nil {
				fmt.Println("rollback error")
			}
		}
	}()

	if err = d.checkError(err); err != nil {
		return nil, err
	}

	err = d.UpdateMaxID(ctx, bizTag, tx)
	if err = d.checkError(err); err != nil {
		return nil, err
	}

	leaf, err := d.Get(ctx, bizTag, tx)
	if err = d.checkError(err); err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err = d.checkError(err); err != nil {
		return nil, err
	}
	return leaf, nil
}

func (d *Dao) checkError(err error) error {
	if err == nil {
		return nil
	}
	if message, ok := err.(*mysql.MySQLError); ok {
		fmt.Println("it's sql error; str:%v", message.Message)
	}
	return errors.New("db error, " + err.Error())
}

func (d *Dao) Get(ctx context.Context, bizTag string, tx *sql.Tx) (*model.Leaf, error) {
	query := fmt.Sprintf("SELECT id, biz_tag, max_id, step, update_time FROM %s WHERE biz_tag=?", model.LeafTableName())
	var leaf model.Leaf
	var err error
	if tx != nil {
		err = tx.QueryRowContext(ctx, query, bizTag).Scan(&leaf.ID, &leaf.BizTag, &leaf.MaxID, &leaf.Step, &leaf.UpdateTime)
	} else {
		err = d.sql.QueryRowContext(ctx, query, bizTag).Scan(&leaf.ID, &leaf.BizTag, &leaf.MaxID,
			&leaf.Step, &leaf.UpdateTime)
	}

	if err != nil {
		fmt.Printf("get leaf failed; biz_tag:%s;err: %v", bizTag, err)
		return nil, err
	}
	return &leaf, nil
}

func (d *Dao) UpdateMaxID(ctx context.Context, bizTag string, tx *sql.Tx) error {
	query := fmt.Sprintf("UPDATE %s SET max_id = max_id + step, update_time = ? WHERE biz_tag = ?", model.LeafTableName())
	var err error
	var res sql.Result
	now := uint64(time.Now().Unix())
	if tx != nil {
		res, err = tx.ExecContext(ctx, query, now, bizTag)
	} else {
		res, err = d.sql.ExecContext(ctx, query, now, bizTag)
	}

	if err != nil {
		fmt.Printf("update max_id failed; bizTag: %s; err: %v", bizTag, err)
		return err
	}

	rowsID, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsID == 0 {
		return errors.New("no update")
	}
	return nil
}
