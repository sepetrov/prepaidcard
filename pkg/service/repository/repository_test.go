// +build integration

package repository_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/sepetrov/prepaidcard/pkg/internal/model"
	"github.com/sepetrov/prepaidcard/pkg/service/repository"

	_ "github.com/go-sql-driver/mysql"
)

const sqlSelectCardWithUUID = "SELECT uuid, available_balance, blocked_balance FROM card WHERE uuid = ?"
const sqlDeleteCard = "DELETE FROM card"

var dsn = fmt.Sprintf(
	"%s:%s@tcp(%s:%s)/%s",
	os.Getenv("TEST_DB_USER"),
	os.Getenv("TEST_DB_PASSWORD"),
	os.Getenv("TEST_DB_HOST"),
	os.Getenv("TEST_DB_PORT"),
	os.Getenv("TEST_DB_NAME"),
)

func TestSaveCard(t *testing.T) {
	db := db(t)
	defer db.Close()

	card := model.NewCard()
	repo := repository.New(db)
	if err := repo.SaveCard(card); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if _, err := db.Exec(sqlDeleteCard); err != nil {
			t.Fatalf("cannot delete test card: %v", err)
		}
	}()

	res := struct {
		uuid             string
		availableBalance uint64
		blockedBalance   uint64
	}{}
	row := db.QueryRow(sqlSelectCardWithUUID, card.UUID())
	if err := row.Scan(&res.uuid, &res.availableBalance, &res.blockedBalance); err != nil {
		t.Fatalf("got error, want one row: %v", err)
	}
	if res.uuid != card.UUID().String() {
		t.Errorf("got uuid %q, want %q", res.uuid, card.UUID().String())
	}
	if res.availableBalance != card.AvailableBalance() {
		t.Errorf("got available_balance %d, want %d", res.availableBalance, card.AvailableBalance())
	}
	if res.blockedBalance != card.BlockedBalance() {
		t.Errorf("got blocked_balance %d, want %d", res.blockedBalance, card.BlockedBalance())
	}
}

func db(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatal(err)
	}
	return db
}
