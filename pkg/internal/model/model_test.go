package model

import (
	"math"
	"testing"
)

func TestCard_LoadMoney(t *testing.T) {
	t.Run("amount must be greater than zero", func(t *testing.T) {
		t.Parallel()
		c := NewCard()
		if c.LoadMoney(0) == nil {
			t.Error("c.LoadBalance(0) nil; want error")
		}
	})
	t.Run("available balance cannot become greater than math.MaxUint64", func(t *testing.T) {
		t.Parallel()
		c := NewCard()
		if err := c.LoadMoney(1); err != nil {
			t.Fatalf("c.LoadMoney(1) err; want nil; %v", err)
		}
		if c.LoadMoney(math.MaxUint64) == nil {
			t.Error("c.LoadMoney(math.MaxUint64) nil; want error")
		}
	})
	t.Run("available balance can reach math.MaxUint64", func(t *testing.T) {
		t.Parallel()
		c := NewCard()
		if err := c.LoadMoney(math.MaxUint64); err != nil {
			t.Errorf("c.LoadMoney(math.MaxUint64) %v; want nil", err)
		}
		if c.AvailableBalance() != math.MaxUint64 {
			t.Errorf("c.AvailableBalance() == %+v; want math.MaxUint64", c.AvailableBalance())
		}
		if c.BlockedBalance() != 0 {
			t.Errorf("c.BlockedBalance() == %+v; want 0", c.BlockedBalance())
		}
	})
}

func TestCard_BlockMoney(t *testing.T) {
	t.Run("amount must be greater than zero", func(t *testing.T) {
		t.Parallel()
		c := NewCard()
		if c.BlockMoney(0) == nil {
			t.Error("c.BlockMoney(0) nil; want error")
		}
	})
	t.Run("available balance cannot become negative", func(t *testing.T) {
		t.Parallel()
		c := NewCard()
		if err := c.LoadMoney(13); err != nil {
			t.Fatalf("c.LoadMoney(13) err; want nil; %v", err)
		}
		if c.BlockMoney(14) == nil {
			t.Error("c.LoadMoney(math.MaxUint64) nil; want error")
		}
	})
	t.Run("available balance can become zero", func(t *testing.T) {
		t.Parallel()
		c := NewCard()
		if err := c.LoadMoney(13); err != nil {
			t.Fatalf("c.LoadMoney(13) err; want nil; %v", err)
		}
		if c.BlockMoney(13) != nil {
			t.Fatalf("c.LoadMoney(math.MaxUint64) error; want nil")
		}
		if c.AvailableBalance() != 0 {
			t.Errorf("c.AvailableBalance() == %+v; want 0", c.AvailableBalance())
		}

		if c.BlockedBalance() != 13 {
			t.Errorf("c.BlockedBalance() == %+v; want 13", c.BlockedBalance())
		}
	})
}

func TestCard_ChargeMoney(t *testing.T) {
	t.Run("cannot charge 0", func(t *testing.T) {
		t.Parallel()
		c := NewCard()
		if c.ChargeMoney(0) == nil {
			t.Error("c.ChargeMoney(0) nil; want error")
		}
	})
	t.Run("cannot charge more than the blocked balance", func(t *testing.T) {
		t.Parallel()
		c := NewCard()
		if err := c.LoadMoney(5); err != nil {
			t.Fatalf("c.LoadMoney(5) err; want nil; %v", err)
		}
		if err := c.BlockMoney(2); err != nil {
			t.Fatalf("c.BlockMoney(2) nil; want error; %v", err)
		}
		if c.ChargeMoney(3) == nil {
			t.Error("c.ChargeMoney(3) nil; want error")
		}
	})
	t.Run("blocked balance can become zero", func(t *testing.T) {
		t.Parallel()
		c := NewCard()
		if err := c.LoadMoney(5); err != nil {
			t.Fatalf("c.LoadMoney(5) err; want nil; %v", err)
		}
		if err := c.BlockMoney(4); err != nil {
			t.Fatalf("c.BlockMoney(4) nil; want error; %v", err)
		}
		if c.ChargeMoney(4) != nil {
			t.Fatalf("c.ChargeMoney(4) error; want nil")
		}
		if c.AvailableBalance() != 1 {
			t.Errorf("c.AvailableBalance() == %+v; want 0", c.AvailableBalance())
		}

		if c.BlockedBalance() != 0 {
			t.Errorf("c.BlockedBalance() == %+v; want 13", c.BlockedBalance())
		}
	})
}
