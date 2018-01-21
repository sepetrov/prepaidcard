package model

import (
	"math"
	"testing"

	"github.com/satori/go.uuid"
)

func TestCard_LoadMoney(t *testing.T) {
	t.Run("amount must be greater than zero", func(t *testing.T) {
		t.Parallel()
		c := &Card{UUID: uuid.NewV4()}
		if c.LoadMoney(0) == nil {
			t.Error("c.LoadBalance(0) nil; want error")
		}
	})
	t.Run("available balance cannot become greater than math.MaxUint64", func(t *testing.T) {
		t.Parallel()
		c := &Card{UUID: uuid.NewV4(), AvailableBalance: 1}
		if c.LoadMoney(math.MaxUint64) == nil {
			t.Error("c.LoadMoney(math.MaxUint64) nil; want error")
		}
	})
	t.Run("available balance can reach math.MaxUint64", func(t *testing.T) {
		t.Parallel()
		c := &Card{UUID: uuid.NewV4()}
		if err := c.LoadMoney(math.MaxUint64); err != nil {
			t.Errorf("c.LoadMoney(math.MaxUint64) %v; want nil", err)
		}
		if c.AvailableBalance != math.MaxUint64 {
			t.Errorf("c.AvailableBalance == %+v; want math.MaxUint64", c.AvailableBalance)
		}
		if c.BlockedBalance != 0 {
			t.Errorf("c.BlockedBalance == %+v; want 0", c.BlockedBalance)
		}
	})
}
	})
}
