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
		c := model.NewCard()
		if c.LoadMoney(0) == nil {
			t.Error("c.LoadBalance(0) nil; want error")
		}
	})
	t.Run("available balance cannot become greater than math.MaxUint64", func(t *testing.T) {
		c := mustCard(t, 1, 0)
		if c.LoadMoney(math.MaxUint64) == nil {
			t.Error("c.LoadMoney(math.MaxUint64) nil; want error")
		}
	})
	t.Run("available balance can reach math.MaxUint64", func(t *testing.T) {
		c := mustCard(t, math.MaxUint64, 0)
		assertCardBalance(t, c, math.MaxUint64, 0)
	})
}

func TestCard_BlockMoney(t *testing.T) {
	t.Run("amount must be greater than zero", func(t *testing.T) {
		c := model.NewCard()
		if c.BlockMoney(0) == nil {
			t.Error("c.BlockMoney(0) nil; want error")
		}
	})
	t.Run("available balance cannot become negative", func(t *testing.T) {
		c := mustCard(t, 13, 0)
		if c.BlockMoney(14) == nil {
			t.Error("c.BlockMoney(14) nil; want error")
		}
	})
	t.Run("blocked balance cannot be more than math.MaxUint64", func(t *testing.T) {
		c := mustCard(t, math.MaxUint64, math.MaxUint64)
		if c.BlockMoney(1) == nil {
			t.Error("c.BlockMoney(1) nil; want error")
		}
	})
	t.Run("available balance can become zero", func(t *testing.T) {
		c := mustCard(t, 13, 0)
		assertCardBalance(t, c, 13, 0)
		if c.BlockMoney(13) != nil {
			t.Fatalf("c.BlockMoney(13) error; want nil")
		}
		assertCardBalance(t, c, 0, 13)
	})
}

func TestCard_ReleaseMoney(t *testing.T) {
	t.Run("cannot release 0", func(t *testing.T) {
		c := model.NewCard()
		if err := c.ReleaseMoney(0); err == nil {
			t.Errorf("c.ReleaseMoney(0) nil; want error")
		}
	})
	t.Run("cannot release more than the blocked balance", func(t *testing.T) {
		c := mustCard(t, 5, 5)
		assertCardBalance(t, c, 0, 5)
		if err := c.ReleaseMoney(6); err == nil {
			t.Errorf("c.ReleaseMoney(6) nil; want error")
		}
	})
	t.Run("cannot release money if available balance becomes more than math.MaxUint64", func(t *testing.T) {
		c := mustCard(t, math.MaxUint64, 1)
		if err := c.LoadMoney(1); err != nil {
			t.Errorf("c.LoadMoney(1) %v; want nil", err)
		}
		assertCardBalance(t, c, math.MaxUint64, 1)
		if err := c.ReleaseMoney(1); err == nil {
			t.Errorf("c.ReleaseMoney(1) nil; want error")
		}
	})
	t.Run("can release multiple times until the blocked balance is reaches 0", func(t *testing.T) {
		c := mustCard(t, 10, 10)
		assertCardBalance(t, c, 0, 10)

		if err := c.ReleaseMoney(2); err != nil {
			t.Errorf("c.ReleaseMoney(2) %v; want nil", err)
		}
		assertCardBalance(t, c, 2, 8)

		if err := c.ReleaseMoney(5); err != nil {
			t.Errorf("c.ReleaseMoney(5) %v; want nil", err)
		}
		assertCardBalance(t, c, 7, 3)

		if err := c.ReleaseMoney(3); err != nil {
			t.Errorf("c.ReleaseMoney(3) %v; want nil", err)
		}
		assertCardBalance(t, c, 10, 0)
	})
}

func TestCard_ChargeMoney(t *testing.T) {
	t.Run("cannot charge 0", func(t *testing.T) {
		c := mustCard(t, 0, 0)
		if c.ChargeMoney(0) == nil {
			t.Error("c.ChargeMoney(0) nil; want error")
		}
	})
	t.Run("cannot charge more than the blocked balance", func(t *testing.T) {
		c := mustCard(t, 2, 2)
		if c.ChargeMoney(3) == nil {
			t.Error("c.ChargeMoney(3) nil; want error")
		}
	})
	t.Run("blocked balance can become zero", func(t *testing.T) {
		c := mustCard(t, 5, 4)
		assertCardBalance(t, c, 1, 4)
		if c.ChargeMoney(4) != nil {
			t.Fatalf("c.ChargeMoney(4) error; want nil")
		}
		assertCardBalance(t, c, 1, 0)
	})
}

func TestNewAuthorizationRequest(t *testing.T) {
	t.Run("cannot block 0", func(t *testing.T) {
		_, err := model.NewAuthorizationRequest(model.NewCard(), uuid.NewV4(), 0)
		if err == nil {
			t.Errorf("NewAuthorizationRequest() = AuthorizationRequest{}, nil; want AuthorizationRequest{}, error")
		}
	})
	t.Run("success", func(t *testing.T) {
		c := model.NewCard()
		c.LoadMoney(100)
		m := uuid.NewV4()
		b := time.Now()
		req, err := model.NewAuthorizationRequest(c, m, 70)
		if err != nil {
			t.Fatalf("NewAuthorizationRequest() = %+v, %v; want nil", req, err)
		}
		if req.CardUUID() != c.UUID() {
			t.Errorf("req.CardUUID() = %v; want %v", req.CardUUID(), c.UUID())
		}
		if req.MerchantUUID() != m {
			t.Errorf("req.MerchantUUID() = %v; want %v", req.MerchantUUID(), m)
		}
		if len(req.History()) != 1 {
			t.Errorf("len(req.History()) = %v; want 1", len(req.History()))
		}
		s := req.History()[0]
		if s.CreatedAt().Before(b) || s.CreatedAt().After(time.Now()) {
			t.Errorf("s.CreatedAt() = %v; want <= time.Now()", s.CreatedAt())
		}
		assertAuthorizationRequestBalance(t, req, 70, 0, 0)
		assertCardBalance(t, c, 30, 70)
	})

}

func TestAuthorizationRequest_Reverse(t *testing.T) {
	t.Run("cannot reverse from different card", func(t *testing.T) {
		c1 := mustCard(t, 1, 0)
		c2 := mustCard(t, 10, 10)
		if c1.UUID() == c2.UUID() {
			t.Fatalf("c1.UUID() == c2.UUID(); want %v != %v", c1.UUID(), c2.UUID())
		}
		req, err := model.NewAuthorizationRequest(c1, uuid.NewV4(), 1)
		if err != nil {
			t.Errorf("NewAuthorizationRequest() = AuthorizationRequest{}, %v; want AuthorizationRequest{}, nil", err)
		}
		if err := req.Reverse(c2, 1); err == nil {
			t.Error("req.Reverse(c2, 0) = nil; want error")
		}
	})
	t.Run("cannot reverse 0", func(t *testing.T) {
		c := model.NewCard()
		c.LoadMoney(100)
		req, err := model.NewAuthorizationRequest(c, uuid.NewV4(), 100)
		if err != nil {
			t.Errorf("NewAuthorizationRequest() = AuthorizationRequest{}, %v; want AuthorizationRequest{}, nil", err)
		}
		if err := req.Reverse(c, 0); err == nil {
			t.Error("req.Reverse(0) = nil; want error")
		}
	})
	t.Run("cannot reverse more than the blocked amount", func(t *testing.T) {
		c := model.NewCard()
		c.LoadMoney(100)
		req, err := model.NewAuthorizationRequest(c, uuid.NewV4(), 50)
		if err != nil {
			t.Errorf("NewAuthorizationRequest() = AuthorizationRequest{}, %v; want AuthorizationRequest{}, nil", err)
		}
		if err := req.Reverse(c, 51); err == nil {
			t.Error("req.Reverse(51) = nil; want error")
		}
	})
	t.Run("cannot reverse money if available balance becomes more than math.MaxUint64", func(t *testing.T) {
		c := mustCard(t, math.MaxUint64, 0)
		req, err := model.NewAuthorizationRequest(c, uuid.NewV4(), 1)
		if err != nil {
			t.Errorf("NewAuthorizationRequest() = AuthorizationRequest{}, %v; want AuthorizationRequest{}, nil", err)
		}
		assertCardBalance(t, c, math.MaxUint64-1, 1)

		if err := c.LoadMoney(1); err != nil {
			t.Fatalf("c.LoadMoney(1) %v; want nil", err)
		}
		assertCardBalance(t, c, math.MaxUint64, 1)

		if err := req.Reverse(c, 1); err == nil {
			t.Fatal("req.Reverse(c, 1) nil; want error")
		}
	})
	t.Run("can reverse multiple times until the blocked amount reaches 0", func(t *testing.T) {
		c := model.NewCard()
		c.LoadMoney(50)
		assertCardBalance(t, c, 50, 0)
		req, err := model.NewAuthorizationRequest(c, uuid.NewV4(), 50)
		if err != nil {
			t.Errorf("NewAuthorizationRequest() = AuthorizationRequest{}, %v; want AuthorizationRequest{}, nil", err)
		}
		assertAuthorizationRequestBalance(t, req, 50, 0, 0)
		assertCardBalance(t, c, 0, 50)

		if err := req.Reverse(c, 10); err != nil {
			t.Errorf("req.Reverse(10) = %v; want nil", err)
		}
		assertAuthorizationRequestBalance(t, req, 40, 0, 0)
		assertCardBalance(t, c, 10, 40)

		if err := req.Reverse(c, 15); err != nil {
			t.Errorf("req.Reverse(15) = %v; want nil", err)
		}
		assertAuthorizationRequestBalance(t, req, 25, 0, 0)
		assertCardBalance(t, c, 25, 25)

		if err := req.Reverse(c, 25); err != nil {
			t.Errorf("req.Reverse(25) = %v; want nil", err)
		}
		assertAuthorizationRequestBalance(t, req, 0, 0, 0)
		assertCardBalance(t, c, 50, 0)
	})
}

func assertAuthorizationRequestBalance(t *testing.T, req *model.AuthorizationRequest, b, c, r uint64) {
	t.Helper()
	if req.BlockedAmount() != b {
		t.Errorf("req.BlockedAmount() = %v; want %v", req.BlockedAmount(), b)
	}
	if req.CapturedAmount() != c {
		t.Errorf("req.CapturedAmount() = %v; want %c", req.CapturedAmount(), c)
	}
	if req.RefundedAmount() != r {
		t.Errorf("req.RefundedAmount() = %v; want %r", req.RefundedAmount(), r)
	}
	assertAuthorizationRequestSnapshot(t, req)
}

func assertAuthorizationRequestSnapshot(t *testing.T, req *model.AuthorizationRequest) {
	t.Helper()

	if len(req.History()) > 1 {
		last := req.History()[len(req.History())-1]
		beforeLast := req.History()[len(req.History())-2]
		if last.CreatedAt().Before(beforeLast.CreatedAt()) {
			t.Fatalf("incorrect change log order %s", req.History())
		}
	}

	s := req.History()[len(req.History())-1]
	if s.BlockedAmount() != req.BlockedAmount() {
		t.Errorf("s.BlockedAmount() != req.BlockedAmount(); want %v == %v", s.BlockedAmount(), req.BlockedAmount())
	}
	if s.CapturedAmount() != req.CapturedAmount() {
		t.Errorf("s.CapturedAmount() != req.CapturedAmount(); want %v == %v", s.CapturedAmount(), req.CapturedAmount())
	}
	if s.RefundedAmount() != req.RefundedAmount() {
		t.Errorf("s.RefundedAmount() != req.RefundedAmount(); want %v == %v", s.RefundedAmount(), req.RefundedAmount())
	}
}

func assertCardBalance(t *testing.T, c *model.Card, a, b uint64) {
	t.Helper()
	if c.AvailableBalance() != a {
		t.Errorf("c.AvailableBalance() = %v; want %v", c.AvailableBalance(), a)
	}
	if c.BlockedBalance() != b {
		t.Errorf("c.BlockedBalance() = %v; want %v", c.BlockedBalance(), b)
	}
}

func mustCard(t *testing.T, l, b uint64) *model.Card {
	t.Helper()
	c := model.NewCard()
	if l > 0 {
		err := c.LoadMoney(l)
		if err != nil {
			t.Fatalf("Card.LoadMoney(%v) %v; want nil", l, err)
		}
	}
	if b > 0 {
		err := c.BlockMoney(b)
		if err != nil {
			t.Fatalf("Card.BlockMoney(%v) %v; want nil", b, err)
		}
	}
	return c
}
