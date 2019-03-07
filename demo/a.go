package demo

import (
	"database/sql"
)

type A struct {
	Db0 *sql.DB `di:"db"`  	//单例模式
	Db1 *sql.DB `di:"db"` 	//单例模式
	B0  *B      `di:"b,prototype"` //工厂模式
	B1  *B      `di:"b,prototype"` //工厂模式
}

func NewA() *A {
	return &A{}
}

func (p *A) Version() (string, error) {
	rows, err := p.Db0.Query("SELECT VERSION() as version")
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var version string
	if rows.Next() {
		if err := rows.Scan(&version); err != nil {
			return "", err
		}
	}
	if err := rows.Err(); err != nil {
		return "", err
	}
	return version, nil
}
