package config

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	wantDBType := "mysql"
	t.Setenv("DB_TYPE", fmt.Sprintf("%s", wantDBType))
	wantDBHost := "mysql-db"
	t.Setenv("DB_HOST", fmt.Sprintf("%s", wantDBHost))
	wantDBPort := 3306
	t.Setenv("DB_PORT", fmt.Sprintf("%d", wantDBPort))
	wantDBName := "todo_db"
	t.Setenv("DB_NAME", fmt.Sprintf("%s", wantDBName))
	wantDBUser := "root"
	t.Setenv("DB_USER", fmt.Sprintf("%s", wantDBUser))
	wantDBPassword := "password"
	t.Setenv("DB_PASSWORD", fmt.Sprintf("%s", wantDBPassword))

	got, err := New()
	if err != nil {
		t.Fatalf("cannot create config: %v", err)
	}
	if got.DBType != wantDBType {
		t.Errorf("want %s, but %s", wantDBType, got.DBType)
	}
	if got.DBHost != wantDBHost {
		t.Errorf("want %s, but %s", wantDBHost, got.DBHost)
	}
	if got.DBPort != wantDBPort {
		t.Errorf("want %d, but %d", wantDBPort, got.DBPort)
	}
	if got.DBName != wantDBName {
		t.Errorf("want %s, but %s", wantDBName, got.DBName)
	}
	if got.DBUser != wantDBUser {
		t.Errorf("want %s, but %s", wantDBUser, got.DBUser)
	}
	if got.DBPassword != wantDBPassword {
		t.Errorf("want %s, but %s", wantDBPassword, got.DBPassword)
	}
}
