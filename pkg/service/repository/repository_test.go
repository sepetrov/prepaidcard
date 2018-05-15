// +build integration

package repository_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/satori/go.uuid"

	"github.com/sepetrov/prepaidcard/pkg/internal/model"
	"github.com/sepetrov/prepaidcard/pkg/service/repository"

	_ "github.com/go-sql-driver/mysql"
)

const sqlInsertCard = "INSERT INTO card (uuid, available_balance, blocked_balance) VALUES (?, ?, ?)"
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

func TestGetCard(t *testing.T) {
	db := db(t)
	defer db.Close()
	t.Run("returns ErrNotFound", func(t *testing.T) {
		card := model.NewCard()
		repo := repository.New(db)
		stmt, err := db.Prepare(sqlInsertCard)
		if err != nil {
			t.Fatalf("cannot prepare statement to save card: %v", err)
		}
		defer stmt.Close()
		if _, err := stmt.Exec(card.UUID(), card.AvailableBalance(), card.BlockedBalance()); err != nil {
			t.Fatal(err)
		}
		defer func() {
			if _, err := db.Exec(sqlDeleteCard); err != nil {
				t.Fatalf("cannot delete test card: %v", err)
			}
		}()

		card, err = repo.GetCard(uuid.NewV4())
		if err != repository.ErrNotFound {
			t.Fatalf("got error %v, want ErrNotFound", err)
		}
	})
	t.Run("returns card with UUID", func(t *testing.T) {
		repo := repository.New(db)
		stmt, err := db.Prepare(sqlInsertCard)
		if err != nil {
			t.Fatalf("cannot prepare statement to save card: %v", err)
		}
		defer stmt.Close()

		// insert first card
		card := model.NewCard()
		if _, err := stmt.Exec(card.UUID(), card.AvailableBalance(), card.BlockedBalance()); err != nil {
			t.Fatal(err)
		}
		defer func() {
			if _, err := db.Exec(sqlDeleteCard); err != nil {
				t.Fatalf("cannot delete test card: %v", err)
			}
		}()

		// insert second card
		card = model.NewCard()
		if _, err := stmt.Exec(card.UUID(), card.AvailableBalance(), card.BlockedBalance()); err != nil {
			t.Fatal(err)
		}

		res, err := repo.GetCard(card.UUID())
		if err != nil {
			t.Fatalf("got error %v, want nil", err)
		}
		if res.UUID() != card.UUID() {
			t.Errorf("got uuid %q, want %q", res.UUID().String(), card.UUID().String())
		}
		if res.AvailableBalance() != card.AvailableBalance() {
			t.Errorf("got available_balance %d, want %d", res.AvailableBalance(), card.AvailableBalance())
		}
		if res.BlockedBalance() != card.BlockedBalance() {
			t.Errorf("got blocked_balance %d, want %d", res.BlockedBalance(), card.BlockedBalance())
		}
	})
}

func db(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatal(err)
	}
	return db
}
