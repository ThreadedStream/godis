package main

import (
	"fmt"
)

func (a *App) saveUserToDatabase(username string, password [32]byte) error {
	query := fmt.Sprintf("INSERT INTO users(username, password) VALUES ('%s', '%v')", username, password)

	_, err := a.Conn.Query(query)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) retrieveUserFromDatabase(username string, password [32]byte) (int, error) {

	query := fmt.Sprintf("SELECT COUNT(1) AS count FROM users WHERE username='%s' AND password='%v'", username, password)

	rows, err := a.Conn.Query(query)
	if err != nil {
		return 0, err
	}

	var count int

	for rows.Next() {
		rows.Scan(&count)
	}

	return count, nil
}
