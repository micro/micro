package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/micro/micro/config/db"
	proto "github.com/micro/micro/config/proto/config"
)

var (
	changeQ = map[string]string{
		"read": `SELECT id, author, comment, timestamp, changeset_timestamp, changeset_checksum, changeset_data, changeset_format, changeset_source 
				from %s.%s where id = ? limit 1`,
		"create": `INSERT INTO %s.%s (id, author, comment, timestamp, changeset_timestamp, changeset_checksum, changeset_data, changeset_format, changeset_source) 
				values(?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"update": `UPDATE %s.%s SET author = ?, comment = ?, timestamp = ?, changeset_timestamp = ?, changeset_checksum = ?, changeset_data = ?, changeset_format = ?,
				changeset_source = ? where id = ? limit 1`,
		"delete": `DELETE from %s.%s where id = ? limit 1`,

		"search": `SELECT id, author, comment, timestamp, changeset_timestamp, changeset_checksum, changeset_data, changeset_format, changeset_source 
				from %s.%s limit ? offset ?`,
		"searchId": `SELECT id, author, comment, timestamp, changeset_timestamp, changeset_checksum, changeset_data, changeset_format, changeset_source 
				from %s.%s where id = ? limit ? offset ?`,
		"searchAuthor": `SELECT id, author, comment, timestamp, changeset_timestamp, changeset_checksum, changeset_data, changeset_format, changeset_source 
				from %s.%s where author = ? limit ? offset ?`,
		"searchIdAndAuthor": `SELECT id, author, comment, timestamp, changeset_timestamp, changeset_checksum, changeset_data, changeset_format, changeset_source 
				from %s.%s where id = ? and author = ? limit ? offset ?`,
	}

	changeLogQ = map[string]string{
		"createLog": `INSERT INTO %s.%s (pid, action, id, path, author, comment, timestamp, changeset_timestamp, changeset_checksum, changeset_data, changeset_format, changeset_source) 
				values(null, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		"readLog": `SELECT pid, action, id, path, author, comment, timestamp, changeset_timestamp, changeset_checksum, changeset_data, changeset_format, changeset_source 
				from %s.%s limit ? offset ?`,
		"readBetween": `SELECT pid, action, id, path, author, comment, timestamp, changeset_timestamp, changeset_checksum, changeset_data, changeset_format, changeset_source 
				from %s.%s where timestamp >= ? and timestamp <= ? limit ? offset ?`,
		"readLogDesc": `SELECT pid, action, id, path, author, comment, timestamp, changeset_timestamp, changeset_checksum, changeset_data, changeset_format, changeset_source 
				from %s.%s order by pid desc limit ? offset ?`,
		"readBetweenDesc": `SELECT pid, action, id, path, author, comment, timestamp, changeset_timestamp, changeset_checksum, changeset_data, changeset_format, changeset_source 
				from %s.%s where timestamp >= ? and timestamp <= ? order by pid desc limit ? offset ?`,
	}

	st = map[string]*sql.Stmt{}
)

type mysql struct {
	db *sql.DB
}

func init() {
	db.Register(new(mysql))
}

func (m *mysql) Init(opts db.Options) error {
	var d *sql.DB
	var err error

	parts := strings.Split(opts.Url, "/")
	if len(parts) != 2 {
		return errors.New("Invalid database url ")
	}

	if len(parts[1]) == 0 {
		return errors.New("Invalid database name ")
	}

	var paramParts []string
	if strings.Contains(opts.Url, "?") {
		paramParts = strings.Split(parts[1], "?")
		parts[1] = paramParts[0]
		paramParts = paramParts[1:]
	}

	url := parts[0]
	database := "`" + parts[1] + "`"

	if d, err = sql.Open("mysql", url+"/"); err != nil {
		return err
	}
	if _, err := d.Exec("CREATE DATABASE IF NOT EXISTS " + database); err != nil {
		return err
	}
	d.Close()
	if d, err = sql.Open("mysql", opts.Url); err != nil {
		return err
	}
	if _, err = d.Exec(changeSchema); err != nil {
		return err
	}
	if _, err = d.Exec(changeLogSchema); err != nil {
		return err
	}

	for query, statement := range changeQ {
		prepared, err := d.Prepare(fmt.Sprintf(statement, database, "configs"))
		if err != nil {
			return err
		}
		st[query] = prepared
	}

	for query, statement := range changeLogQ {
		prepared, err := d.Prepare(fmt.Sprintf(statement, database, "change_log"))
		if err != nil {
			return err
		}
		st[query] = prepared
	}

	m.db = d

	return nil
}

func (m *mysql) Create(change *proto.Change) error {
	// create change entry
	_, err := st["create"].Exec(
		change.Id,
		change.Author,
		change.Comment,
		change.Timestamp,
		change.ChangeSet.Timestamp,
		change.ChangeSet.Checksum,
		change.ChangeSet.Data,
		change.ChangeSet.Format,
		change.ChangeSet.Source,
	)
	if err != nil {
		return err
	}

	// create log entry
	_, err = st["createLog"].Exec(
		"create",
		change.Id,
		change.Path,
		change.Author,
		change.Comment,
		change.Timestamp,
		change.ChangeSet.Timestamp,
		change.ChangeSet.Checksum,
		change.ChangeSet.Data,
		change.ChangeSet.Format,
		change.ChangeSet.Source,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *mysql) Read(id string) (*proto.Change, error) {
	if len(id) == 0 {
		return nil, errors.New("Invalid trace id")
	}

	change := &proto.Change{
		ChangeSet: &proto.ChangeSet{},
	}

	r := st["read"].QueryRow(id)
	if err := r.Scan(
		&change.Id,
		&change.Author,
		&change.Comment,
		&change.Timestamp,
		&change.ChangeSet.Timestamp,
		&change.ChangeSet.Checksum,
		&change.ChangeSet.Data,
		&change.ChangeSet.Format,
		&change.ChangeSet.Source,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("not found")
		}
		return nil, err
	}

	return change, nil
}

func (m *mysql) Delete(change *proto.Change) error {
	_, err := st["delete"].Exec(change.Id)
	if err != nil {
		return err
	}

	// create log entry
	_, err = st["createLog"].Exec(
		"delete",
		change.Id,
		change.Path,
		change.Author,
		change.Comment,
		change.Timestamp,
		change.ChangeSet.Timestamp,
		change.ChangeSet.Checksum,
		change.ChangeSet.Data,
		change.ChangeSet.Format,
		change.ChangeSet.Source,
	)
	if err != nil {
		return err
	}
	return err
}

func (m *mysql) Update(change *proto.Change) error {
	_, err := st["update"].Exec(
		change.Author,
		change.Comment,
		change.Timestamp,
		change.ChangeSet.Timestamp,
		change.ChangeSet.Checksum,
		change.ChangeSet.Data,
		change.ChangeSet.Format,
		change.ChangeSet.Source,
		change.Id,
	)
	if err != nil {
		return err
	}

	// create log entry
	_, err = st["createLog"].Exec(
		"update",
		change.Id,
		change.Path,
		change.Author,
		change.Comment,
		change.Timestamp,
		change.ChangeSet.Timestamp,
		change.ChangeSet.Checksum,
		change.ChangeSet.Data,
		change.ChangeSet.Format,
		change.ChangeSet.Source,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *mysql) Search(id, author string, limit, offset int64) ([]*proto.Change, error) {
	var r *sql.Rows
	var err error

	if len(id) > 0 && len(author) > 0 {
		r, err = st["searchIdAndAuthor"].Query(id, author, limit, offset)
	} else if len(id) > 0 {
		r, err = st["searchId"].Query(id, limit, offset)
	} else if len(author) > 0 {
		r, err = st["searchAuthor"].Query(author, limit, offset)
	} else {
		r, err = st["search"].Query(limit, offset)
	}

	if err != nil {
		return nil, err
	}
	defer r.Close()

	var changes []*proto.Change

	for r.Next() {
		change := &proto.Change{
			ChangeSet: &proto.ChangeSet{},
		}
		if err := r.Scan(
			&change.Id,
			&change.Author,
			&change.Comment,
			&change.Timestamp,
			&change.ChangeSet.Timestamp,
			&change.ChangeSet.Checksum,
			&change.ChangeSet.Data,
			&change.ChangeSet.Format,
			&change.ChangeSet.Source,
		); err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("not found")
			}
			return nil, err
		}
		changes = append(changes, change)

	}
	if r.Err() != nil {
		return nil, err
	}

	return changes, nil
}

func (m *mysql) AuditLog(from, to, limit, offset int64, reverse bool) ([]*proto.ChangeLog, error) {
	var r *sql.Rows
	var err error

	if from == 0 && to == 0 {
		q := "readLog"
		if reverse {
			q += "Desc"
		}
		r, err = st[q].Query(limit, offset)
	} else {
		q := "readBetween"
		if reverse {
			q += "Desc"
		}
		r, err = st[q].Query(from, to, limit, offset)
	}

	if err != nil {
		return nil, err
	}
	defer r.Close()

	var logs []*proto.ChangeLog

	for r.Next() {
		var id int

		log := &proto.ChangeLog{
			Change: &proto.Change{
				ChangeSet: &proto.ChangeSet{},
			},
		}
		if err := r.Scan(
			&id,
			&log.Action,
			&log.Change.Id,
			&log.Change.Path,
			&log.Change.Author,
			&log.Change.Comment,
			&log.Change.Timestamp,
			&log.Change.ChangeSet.Timestamp,
			&log.Change.ChangeSet.Checksum,
			&log.Change.ChangeSet.Data,
			&log.Change.ChangeSet.Format,
			&log.Change.ChangeSet.Source,
		); err != nil {
			if err == sql.ErrNoRows {
				return nil, errors.New("not found")
			}
			return nil, err
		}
		logs = append(logs, log)

	}
	if r.Err() != nil {
		return nil, err
	}

	return logs, nil
}

func (m *mysql) String() string {
	return "mysql"
}
