package data

import (
	"context"
	"database/sql"
	"time"
)


type Permissions []string

type PermissionModel struct {
	DB *sql.DB
}

func (p Permissions) Include(code string) bool {
	for i := range p {
		if p[i] == code {
			return true
		}
	}
	return false
}

func (p *PermissionModel) GetAllForUser(userId int64) (*Permissions,error) {

	var sql string = `
		SELECT p.code
		FROM permissions as p 
		INNER JOIN users_permission as up 
		ON p.id = up.permission_id
		INNER JOIN users as u 
		ON up.user_id = u.id
		WHERE u.id=$1
		`

	args := []interface{}{
		userId,
	}

	permissions := Permissions{}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	
	defer cancel()

	rows, err := p.DB.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var permission string 
		err = rows.Scan(&permission)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	if err =rows.Err() ; err != nil {
		return nil, err
	}

	return &permissions, nil
}
