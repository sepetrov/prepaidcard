package model_test

import (
	"math"
	"testing"
	"time"

	"github.com/satori/go.uuid"
	"github.com/sepetrov/prepaidcard/pkg/internal/model"
)

func TestCard_LoadMoney(t *testing.T) {
	t.Run("amount must be greater than zero", func(t *testing.T) {
		t.Parallel()
		c := model.NewCard()
		if c.LoadMoney(0) == nil {
			t.Error("c.LoadBalance(0) nil; want error")
		}
	})
	t.Run("available balance cannot become greater than math.MaxUint64", func(t *testing.T) {
		t.Parallel()
		c := model.NewCard()
		if err := c.LoadMoney(1); err != nil {
			t.Fatalf("c.LoadMoney(1) err; want nil; %v", err)
		}
		if c.LoadMoney(math.MaxUint64) == nil {
			t.Error("c.LoadMoney(math.MaxUint64) nil; want error")
		}
	})
	t.Run("available balance can reach math.MaxUint64", func(t *testing.T) {
		t.Parallel()
		c := model.NewCard()
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
		c := model.NewCard()
		if c.BlockMoney(0) == nil {
			t.Error("c.BlockMoney(0) nil; want error")
		}
	})
	t.Run("available balance cannot become negative", func(t *testing.T) {
		t.Parallel()
		c := model.NewCard()
		if err := c.LoadMoney(13); err != nil {
			t.Fatalf("c.LoadMoney(13) err; want nil; %v", err)
		}
		if c.BlockMoney(14) == nil {
			t.Error("c.LoadMoney(math.MaxUint64) nil; want error")
		}
	})
	t.Run("available balance can become zero", func(t *testing.T) {
		t.Parallel()
		c := model.NewCard()
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
		c := model.NewCard()
		if c.ChargeMoney(0) == nil {
			t.Error("c.ChargeMoney(0) nil; want error")
		}
	})
	t.Run("cannot charge more than the blocked balance", func(t *testing.T) {
		t.Parallel()
		c := model.NewCard()
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
		c := model.NewCard()
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

func TestNewAuthorizationRequest(t *testing.T) {
	t.Run("cannot block 0", func(t *testing.T) {
		_, err := model.NewAuthorizationRequest(*model.NewCard(), uuid.NewV4(), 0)
		if err == nil {
			t.Errorf("NewAuthorizationRequest() = AuthorizationRequest{}, nil; want AuthorizationRequest{}, error")
		}
	})
	t.Run("success", func(t *testing.T) {
		c := model.NewCard()
		c.LoadMoney(100)
		m := uuid.NewV4()
		a := uint64(100)
		b := time.Now()
		req, err := model.NewAuthorizationRequest(*c, m, a)
		if err != nil {
			t.Errorf("NewAuthorizationRequest() = %+v, %v; want nil", req, err)
		}
		if req.CardUUID() != c.UUID() {
			t.Errorf("req.CardUUID() = %v; want %v", req.CardUUID(), c.UUID())
		}
		if req.MerchantUUID() != m {
			t.Errorf("req.MerchantUUID() = %v; want %v", req.MerchantUUID(), m)
		}
		if req.BlockedAmount() != a {
			t.Errorf("req.BlockedAmount() = %v; want %v", req.BlockedAmount(), a)
		}
		if req.CapturedAmount() != 0 {
			t.Errorf("req.CapturedAmount() = %v; want 0", req.CapturedAmount())
		}
		if req.RefundedAmount() != 0 {
			t.Errorf("req.RefundedAmount() = %v; want 0", req.RefundedAmount())
		}
		if len(req.History()) != 1 {
			t.Errorf("len(req.History()) = %v; want 1", len(req.History()))
		}
		s := req.History()[0]
		if s.BlockedAmount() != req.BlockedAmount() {
			t.Errorf("s.BlockedAmount() != %v; want %v", s.BlockedAmount(), req.BlockedAmount())
		}
		if s.CapturedAmount() != req.CapturedAmount() {
			t.Errorf("s.CapturedAmount() != %v; want %v", s.CapturedAmount(), req.CapturedAmount())
		}
		if s.RefundedAmount() != req.RefundedAmount() {
			t.Errorf("s.CapturedAmount() != %s; want %s", s.CapturedAmount(), req.CapturedAmount())
		}
		if s.CreatedAt().Before(b) || s.CreatedAt().After(time.Now()) {
			t.Errorf("s.CreatedAt() = %v; want <= time.Now()", s.CreatedAt())
		}
	})

}
}
