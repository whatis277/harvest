package subscription

import (
	"testing"

	"harvest/bean/internal/usecases/interfaces"

	"harvest/bean/internal/driver/postgres/postgrestest"
)

var (
	userWithSubsID   = "00000000-0000-0000-0002-000000000001"
	userWithNoSubsID = "00000000-0000-0000-0002-000000000002"

	methodID = "00000000-0000-0000-0001-000000000001"

	subID     = "00000000-0000-0000-0000-000000000001"
	missingID = "11111111-1111-1111-1111-111111111111"
)

func TestDataSouce(t *testing.T) {
	db := postgrestest.DBTest(t)
	ds := New(db)

	t.Run("create", func(t *testing.T) {
		create(t, ds)
	})

	t.Run("find_by_id", func(t *testing.T) {
		findByID(t, ds)
	})

	t.Run("find_by_user_id", func(t *testing.T) {
		findByUserID(t, ds)
	})

	t.Run("delete", func(t *testing.T) {
		delete(t, ds)
	})
}

func create(t *testing.T, ds interfaces.SubscriptionDataSource) {
	sub, err := ds.Create(userWithSubsID, methodID, "action-create", "bean", 1000, 1, "month")
	if err != nil {
		t.Fatalf("failed to create subscription: %s", err)
	}

	if sub.ID == "" {
		t.Error("expected subscription ID, got empty string")
	}

	if sub.UserID != userWithSubsID {
		t.Errorf("expected user ID: %s, got: %s", userWithSubsID, sub.UserID)
	}

	if sub.PaymentMethodID != methodID {
		t.Errorf("expected payment method ID: %s, got: %s", methodID, sub.PaymentMethodID)
	}

	if sub.Label != "action-create" {
		t.Errorf("expected label: %s, got: %s", "action-create", sub.Label)
	}

	if sub.Provider != "bean" {
		t.Errorf("expected provider: %s, got: %s", "bean", sub.Provider)
	}

	if sub.Amount != 1000 {
		t.Errorf("expected amount: %d, got: %d", 1000, sub.Amount)
	}

	if sub.Interval != 1 {
		t.Errorf("expected interval: %d, got: %d", 1, sub.Interval)
	}

	if sub.Period != "month" {
		t.Errorf("expected period: %s, got: %s", "month", sub.Period)
	}

	if err := ds.Delete(userWithSubsID, sub.ID); err != nil {
		t.Fatalf("failed to cleanup subscription: %s", err)
	}
}

func findByID(t *testing.T, ds interfaces.SubscriptionDataSource) {
	t.Run("existing_subscription", func(t *testing.T) {
		sub, err := ds.FindByID(userWithSubsID, subID)
		if err != nil {
			t.Fatalf("expected nil error, got: %s", err)
		}

		if sub == nil {
			t.Error("expected subscription, got nil")
		}
	})

	t.Run("not_owned_subscription", func(t *testing.T) {
		sub, err := ds.FindByID(userWithNoSubsID, subID)
		if err != nil {
			t.Fatalf("expected nil error, got: %s", err)
		}

		if sub != nil {
			t.Errorf("expected nil subscription, got: %v", sub)
		}
	})

	t.Run("missing_subscription", func(t *testing.T) {
		sub, err := ds.FindByID(userWithSubsID, missingID)
		if err != nil {
			t.Fatalf("expected nil error, got: %s", err)
		}

		if sub != nil {
			t.Errorf("expected nil subscription, got: %v", sub)
		}
	})
}

func findByUserID(t *testing.T, ds interfaces.SubscriptionDataSource) {
	t.Run("has_subscriptions", func(t *testing.T) {
		subs, err := ds.FindByUserID(userWithSubsID)
		if err != nil {
			t.Fatalf("failed to find subscriptions by user id: %s", err)
		}

		if len(subs) != 2 {
			t.Errorf("expected %d subscriptions, got: %d", 2, len(subs))
		}

		for _, sub := range subs {
			if sub.UserID != userWithSubsID {
				t.Errorf("expected user ID: %s, got: %s", userWithSubsID, sub.UserID)
			}
		}
	})

	t.Run("no_subscriptions", func(t *testing.T) {
		subs, err := ds.FindByUserID(userWithNoSubsID)
		if err != nil {
			t.Fatalf("failed to find subscriptions by user id: %s", err)
		}

		if len(subs) != 0 {
			t.Errorf("expected %d subscriptions, got: %d", 0, len(subs))
		}
	})
}

func delete(t *testing.T, ds interfaces.SubscriptionDataSource) {
	t.Run("existing_subscription", func(t *testing.T) {
		sub, err := ds.Create(userWithSubsID, methodID, "action-delete", "bean", 1000, 1, "month")
		if err != nil {
			t.Fatalf("failed to create subscription: %s", err)
		}

		if err = ds.Delete(userWithSubsID, sub.ID); err != nil {
			t.Fatalf("failed to delete subscription: %s", err)
		}

		sub, err = ds.FindByID(userWithSubsID, sub.ID)
		if err != nil {
			t.Fatalf("failed to find subscription: %s", err)
		}

		if sub != nil {
			t.Errorf("expected nil subscription, got: %v", sub)
		}
	})

	t.Run("not_owned_subscription", func(t *testing.T) {
		if err := ds.Delete(userWithNoSubsID, subID); err != nil {
			t.Errorf("failed to delete subscription: %s", err)
		}

		sub, err := ds.FindByID(userWithSubsID, subID)
		if err != nil {
			t.Fatalf("failed to find subscription: %s", err)
		}

		if sub == nil {
			t.Error("expected subscription, got nil")
		}
	})

	t.Run("missing_subscription", func(t *testing.T) {
		if err := ds.Delete(userWithSubsID, missingID); err != nil {
			t.Errorf("failed to delete subscription: %s", err)
		}
	})
}
