package model_test

import (
	"math"
	"testing"
	"time"

	"github.com/gofrs/uuid"

	"github.com/sepetrov/prepaidcard/pkg/internal/model"
	h "github.com/sepetrov/prepaidcard/pkg/internal/testing"
)

func TestCard_LoadMoney(t *testing.T) {
	t.Run("amount must be greater than zero", func(t *testing.T) {
		c, err := model.NewCard()
		h.MustNotErr(t, err, "%v")
		h.MustErr(t, c.LoadMoney(0), "c.LoadBalance(0) nil; want error")
	})
	t.Run("available balance cannot become greater than math.MaxUint64", func(t *testing.T) {
		c := mustCard(t, 1, 0)
		h.MustErr(t, c.LoadMoney(math.MaxUint64), "c.LoadMoney(math.MaxUint64) nil; want error")
	})
	t.Run("available balance can reach math.MaxUint64", func(t *testing.T) {
		c := mustCard(t, math.MaxUint64, 0)
		assertCardBalance(t, c, math.MaxUint64, 0)
	})
}

func TestNewAuthorizationRequest(t *testing.T) {
	t.Run("cannot block 0", func(t *testing.T) {
		c, err := model.NewCard()
		h.MustNotErr(t, err, "%v")
		_, err = model.NewAuthorizationRequest(c, uuid.Must(uuid.NewV4()), 0)
		h.MustErr(t, err, "NewAuthorizationRequest() = AuthorizationRequest{}, nil; want AuthorizationRequest{}, error")
	})
	t.Run("success", func(t *testing.T) {
		c, err := model.NewCard()
		h.MustNotErr(t, err, "%v")
		c.LoadMoney(100)
		m := uuid.Must(uuid.NewV4())
		b := time.Now()
		req, err := model.NewAuthorizationRequest(c, m, 70)
		h.MustNotErr(t, err, "NewAuthorizationRequest() = %+v, %v; want nil", req)
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
		req := mustAuthorizationRequest(t, c1, 1)
		h.MustErr(t, req.Reverse(c2, 1), "req.Reverse(c2, 0) = nil; want error")
	})
	t.Run("cannot reverse 0", func(t *testing.T) {
		c, req := mustCardWithAuthorizationRequest(t, 100, 100)
		h.MustErr(t, req.Reverse(c, 0), "req.Reverse(c, 0) = nil; want error")
	})
	t.Run("cannot reverse more than the blocked amount", func(t *testing.T) {
		c, req := mustCardWithAuthorizationRequest(t, 100, 50)
		h.MustErr(t, req.Reverse(c, 51), "req.Reverse(51) = nil; want error")
	})
	t.Run("cannot reverse money if available balance becomes more than math.MaxUint64", func(t *testing.T) {
		c, req := mustCardWithAuthorizationRequest(t, math.MaxUint64, 1)
		h.MustNotErr(t, c.LoadMoney(1), "c.LoadMoney(1) %v; want nil")
		assertCardBalance(t, c, math.MaxUint64, 1)
		h.MustErr(t, req.Reverse(c, 1), "req.Reverse(c, 1) nil; want error")
	})
	t.Run("can reverse multiple times until the blocked amount reaches 0", func(t *testing.T) {
		c, req := mustCardWithAuthorizationRequest(t, 50, 50)
		assertAuthorizationRequestBalance(t, req, 50, 0, 0)
		assertCardBalance(t, c, 0, 50)

		h.MustNotErr(t, req.Reverse(c, 10), "req.Reverse(10) = %v; want nil")
		assertAuthorizationRequestBalance(t, req, 40, 0, 0)
		assertCardBalance(t, c, 10, 40)

		h.MustNotErr(t, req.Reverse(c, 15), "req.Reverse(15) = %v; want nil")
		assertAuthorizationRequestBalance(t, req, 25, 0, 0)
		assertCardBalance(t, c, 25, 25)

		h.MustNotErr(t, req.Reverse(c, 25), "req.Reverse(25) = %v; want nil")
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
		t.Errorf("req.RefundedAmount() = %v; want %v", req.RefundedAmount(), r)
	}
	assertAuthorizationRequestSnapshot(t, req)
}

func assertAuthorizationRequestSnapshot(t *testing.T, req *model.AuthorizationRequest) {
	t.Helper()

	if len(req.History()) > 1 {
		last := req.History()[len(req.History())-1]
		beforeLast := req.History()[len(req.History())-2]
		if last.CreatedAt().Before(beforeLast.CreatedAt()) {
			t.Fatalf("incorrect change log order %v", req.History())
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
	c, err := model.NewCard()
	h.MustNotErr(t, err, "%v")
	if l > 0 {
		h.MustNotErr(t, c.LoadMoney(l), "Card.LoadMoney(%v) %v; want nil; mustCard", l)
	}
	if b > 0 {
		mustAuthorizationRequest(t, c, b)
	}
	return c
}

func mustAuthorizationRequest(t *testing.T, c *model.Card, b uint64) *model.AuthorizationRequest {
	t.Helper()
	req, err := model.NewAuthorizationRequest(c, uuid.Must(uuid.NewV4()), b)
	h.MustNotErr(t, err, "NewAuthorizationRequest(c, uuid.Must(uuid.NewV4()), %v) %v; want nil; mustAuthorizationRequest", b)
	return req
}

func mustCardWithAuthorizationRequest(t *testing.T, l, b uint64) (*model.Card, *model.AuthorizationRequest) {
	t.Helper()
	c := mustCard(t, l, 0)
	req := mustAuthorizationRequest(t, c, b)
	return c, req
}
