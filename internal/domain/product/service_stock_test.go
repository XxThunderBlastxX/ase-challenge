package product

import (
	"testing"

	apperrors "github.com/xxthunderblastxx/ase-challenge/internal/pkg/errors"
	"gorm.io/gorm"
)

// mockRepo is an in-memory implementation of the Repository interface for testing.
type mockRepo struct {
	products map[string]*Product

	// For verifying that update functions were called with expected values.
	lastUpdatedID          string
	lastUpdatedColumn      string
	lastUpdatedColumnValue any
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		products: make(map[string]*Product),
	}
}

func (m *mockRepo) Create(p *Product) error {
	// Not needed for current tests
	return nil
}

func (m *mockRepo) GetAll() ([]Product, error) {
	out := make([]Product, 0, len(m.products))
	for _, p := range m.products {
		out = append(out, *p)
	}
	return out, nil
}

func (m *mockRepo) GetByID(id string) (*Product, error) {
	p, ok := m.products[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return p, nil
}

func (m *mockRepo) UpdateAllColumn(id string, p *Product) error {
	if _, ok := m.products[id]; !ok {
		return gorm.ErrRecordNotFound
	}
	m.products[id] = p
	m.lastUpdatedID = id
	m.lastUpdatedColumn = "ALL"
	return nil
}

func (m *mockRepo) UpdateSingleColumn(id string, column string, value any) error {
	p, ok := m.products[id]
	if !ok {
		return gorm.ErrRecordNotFound
	}
	switch column {
	case "stock_quantity":
		if v, ok := value.(int); ok {
			p.StockQuantity = v
		}
	default:
		// ignore unknown column for test simplicity
	}
	m.lastUpdatedID = id
	m.lastUpdatedColumn = column
	m.lastUpdatedColumnValue = value
	return nil
}

func (m *mockRepo) Delete(id string) error {
	if _, ok := m.products[id]; !ok {
		return gorm.ErrRecordNotFound
	}
	delete(m.products, id)
	return nil
}

// --- Helper assertions ---

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func assertAppErrorCode(t *testing.T, err error, code apperrors.ErrorCode) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected error with code %s, got nil", code)
	}
	appErr, ok := err.(*apperrors.AppError)
	if !ok {
		t.Fatalf("expected *AppError, got %T (%v)", err, err)
	}
	if appErr.Code != code {
		t.Fatalf("expected error code %s, got %s (message: %s)", code, appErr.Code, appErr.Message)
	}
}

// --- Tests ---

func TestService_IncermentStock(t *testing.T) {
	repo := newMockRepo()
	repo.products["p1"] = &Product{
		Name:             "Widget",
		Description:      "Test product",
		StockQuantity:    10,
		LowStockThresold: 5,
	}
	svc := NewService(repo)

	t.Run("successfully increments stock", func(t *testing.T) {
		err := svc.IncermentStock("p1", 5)
		assertNoError(t, err)

		p, _ := repo.GetByID("p1")
		if p.StockQuantity != 15 {
			t.Fatalf("expected stock 15, got %d", p.StockQuantity)
		}
		if repo.lastUpdatedColumn != "stock_quantity" {
			t.Fatalf("expected stock_quantity update, got %s", repo.lastUpdatedColumn)
		}
		if v, ok := repo.lastUpdatedColumnValue.(int); !ok || v != 15 {
			t.Fatalf("expected last updated value 15, got %#v", repo.lastUpdatedColumnValue)
		}
	})

	t.Run("error on zero quantity", func(t *testing.T) {
		err := svc.IncermentStock("p1", 0)
		assertAppErrorCode(t, err, apperrors.InvalidInput)
	})

	t.Run("error on negative quantity", func(t *testing.T) {
		err := svc.IncermentStock("p1", -3)
		assertAppErrorCode(t, err, apperrors.InvalidInput)
	})

	t.Run("error on empty id", func(t *testing.T) {
		err := svc.IncermentStock("", 5)
		assertAppErrorCode(t, err, apperrors.MissingRequiredData)
	})

	t.Run("error when product not found", func(t *testing.T) {
		err := svc.IncermentStock("does-not-exist", 5)
		assertAppErrorCode(t, err, apperrors.ProductNotFound)
	})
}

func TestService_DecrementStock(t *testing.T) {
	repo := newMockRepo()
	repo.products["p1"] = &Product{
		Name:             "Widget",
		Description:      "Test product",
		StockQuantity:    10,
		LowStockThresold: 5,
	}
	svc := NewService(repo)

	t.Run("successfully decrements stock", func(t *testing.T) {
		err := svc.DecrementStock("p1", 3)
		assertNoError(t, err)

		p, _ := repo.GetByID("p1")
		if p.StockQuantity != 7 {
			t.Fatalf("expected stock 7, got %d", p.StockQuantity)
		}
		if repo.lastUpdatedColumn != "stock_quantity" {
			t.Fatalf("expected stock_quantity update, got %s", repo.lastUpdatedColumn)
		}
	})

	t.Run("decrement to exact zero", func(t *testing.T) {
		// Set to known stock
		repo.products["p2"] = &Product{
			Name:             "Gadget",
			StockQuantity:    5,
			LowStockThresold: 3,
		}
		err := svc.DecrementStock("p2", 5)
		assertNoError(t, err)
		p, _ := repo.GetByID("p2")
		if p.StockQuantity != 0 {
			t.Fatalf("expected stock 0, got %d", p.StockQuantity)
		}
	})

	t.Run("error on zero quantity", func(t *testing.T) {
		err := svc.DecrementStock("p1", 0)
		assertAppErrorCode(t, err, apperrors.InvalidInput)
	})

	t.Run("error on negative quantity", func(t *testing.T) {
		err := svc.DecrementStock("p1", -2)
		assertAppErrorCode(t, err, apperrors.InvalidInput)
	})

	t.Run("error on empty id", func(t *testing.T) {
		err := svc.DecrementStock("", 2)
		assertAppErrorCode(t, err, apperrors.MissingRequiredData)
	})

	t.Run("error when product not found", func(t *testing.T) {
		err := svc.DecrementStock("does-not-exist", 1)
		assertAppErrorCode(t, err, apperrors.ProductNotFound)
	})

	t.Run("error when decrement exceeds available stock", func(t *testing.T) {
		// Current p1 stock is 7 from earlier test
		before := repo.products["p1"].StockQuantity
		err := svc.DecrementStock("p1", before+1)
		assertAppErrorCode(t, err, apperrors.InsufficientStock)

		after := repo.products["p1"].StockQuantity
		if after != before {
			t.Fatalf("stock should not change on insufficient stock error; before=%d after=%d", before, after)
		}
	})
}
