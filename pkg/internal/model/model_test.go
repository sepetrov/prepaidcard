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

func TestCard_BlockMoney(t *testing.T) {
	t.Run("amount must be greater than zero", func(t *testing.T) {
		t.Parallel()
		c := &Card{UUID: uuid.NewV4()}
		if c.BlockMoney(0) == nil {
			t.Error("c.BlockMoney(0) nil; want error")
		}
	})
	t.Run("available balance cannot become negative", func(t *testing.T) {
		t.Parallel()
		c := &Card{UUID: uuid.NewV4(), AvailableBalance: 13}
		if c.BlockMoney(14) == nil {
			t.Error("c.LoadMoney(math.MaxUint64) nil; want error")
		}
	})
	t.Run("available balance can become zero", func(t *testing.T) {
		t.Parallel()
		c := &Card{UUID: uuid.NewV4(), AvailableBalance: 13}
		if c.BlockMoney(13) != nil {
			t.Fatalf("c.LoadMoney(math.MaxUint64) error; want nil")
		}
		if c.AvailableBalance != 0 {
			t.Errorf("c.AvailableBalance == %+v; want 0", c.AvailableBalance)
		}

		if c.BlockedBalance != 13 {
			t.Errorf("c.BlockedBalance == %+v; want 13", c.BlockedBalance)
		}
	})
}

func TestCard_ChargeMoney(t *testing.T) {
	t.Run("cannot charge 0", func(t *testing.T) {
		t.Parallel()
		c := &Card{UUID: uuid.NewV4()}
		if c.ChargeMoney(0) == nil {
			t.Error("c.ChargeMoney(0) nil; want error")
		}
	})
	t.Run("cannot charge more than the blocked balance", func(t *testing.T) {
		t.Parallel()
		c := &Card{UUID: uuid.NewV4(), BlockedBalance: 2}
		if c.ChargeMoney(3) == nil {
			t.Error("c.ChargeMoney(3) nil; want error")
		}
	})
	t.Run("blocked balance can become zero", func(t *testing.T) {
		t.Parallel()
		c := &Card{UUID: uuid.NewV4(), AvailableBalance: 3, BlockedBalance: 4}
		if c.ChargeMoney(4) != nil {
			t.Fatalf("c.ChargeMoney(13) error; want nil")
		}
		if c.AvailableBalance != 3 {
			t.Errorf("c.AvailableBalance == %+v; want 0", c.AvailableBalance)
		}

		if c.BlockedBalance != 0 {
			t.Errorf("c.BlockedBalance == %+v; want 13", c.BlockedBalance)
		}
	})
}
