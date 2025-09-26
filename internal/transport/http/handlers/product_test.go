package handlers

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/xxthunderblastxx/ase-challenge/internal/domain/product"
	apperrors "github.com/xxthunderblastxx/ase-challenge/internal/pkg/errors"
)

// mockProductService implements the product.Service interface for testing
type mockProductService struct {
	products []product.Product
	getError error
}

func (m *mockProductService) CreateProduct(*product.Product) error {
	return nil
}

func (m *mockProductService) GetAllProducts() ([]product.Product, error) {
	if m.getError != nil {
		return nil, m.getError
	}
	return m.products, nil
}

func (m *mockProductService) GetProductByID(string) (*product.Product, error) {
	return nil, nil
}

func (m *mockProductService) UpdateProduct(string, *product.Product) error {
	return nil
}

func (m *mockProductService) DeleteProduct(string) error {
	return nil
}

func (m *mockProductService) IncermentStock(string, int) error {
	return nil
}

func (m *mockProductService) DecrementStock(string, int) error {
	return nil
}

func TestProductHandler_GetAllProducts(t *testing.T) {
	// Setup test products
	testProducts := []product.Product{
		{
			Name:             "High Stock Product",
			Description:      "Product with high stock",
			StockQuantity:    50,
			LowStockThresold: 10,
		},
		{
			Name:             "Low Stock Product 1",
			Description:      "Product with low stock",
			StockQuantity:    5,
			LowStockThresold: 10,
		},
		{
			Name:             "Low Stock Product 2",
			Description:      "Product at threshold",
			StockQuantity:    15,
			LowStockThresold: 15,
		},
		{
			Name:             "Normal Stock Product",
			Description:      "Product with normal stock",
			StockQuantity:    25,
			LowStockThresold: 5,
		},
	}

	t.Run("returns all products when no filter applied", func(t *testing.T) {
		mockService := &mockProductService{
			products: testProducts,
		}
		handler := NewProductHandler(mockService)

		app := fiber.New()
		app.Get("/products", handler.GetAllProducts())

		req := httptest.NewRequest("GET", "/products", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != 200 {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		var response struct {
			Data []product.Product `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		if len(response.Data) != 4 {
			t.Errorf("expected 4 products, got %d", len(response.Data))
		}
	})

	t.Run("returns only low stock products when low-stock=true", func(t *testing.T) {
		mockService := &mockProductService{
			products: testProducts,
		}
		handler := NewProductHandler(mockService)

		app := fiber.New()
		app.Get("/products", handler.GetAllProducts())

		req := httptest.NewRequest("GET", "/products?low-stock=true", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != 200 {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		var response struct {
			Data []product.Product `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		// Should return 2 products:
		// - "Low Stock Product 1" (stock=5, threshold=10)
		// - "Low Stock Product 2" (stock=15, threshold=15)
		if len(response.Data) != 2 {
			t.Errorf("expected 2 low stock products, got %d", len(response.Data))
		}

		// Verify the correct products are returned
		expectedNames := map[string]bool{
			"Low Stock Product 1": true,
			"Low Stock Product 2": true,
		}
		for _, p := range response.Data {
			if !expectedNames[p.Name] {
				t.Errorf("unexpected product in low stock results: %s", p.Name)
			}
			// Verify each product is actually low stock
			if p.StockQuantity > p.LowStockThresold {
				t.Errorf("product %s should not be in low stock results (stock=%d, threshold=%d)",
					p.Name, p.StockQuantity, p.LowStockThresold)
			}
		}
	})

	t.Run("returns empty array when no products are low stock", func(t *testing.T) {
		highStockProducts := []product.Product{
			{
				Name:             "High Stock Product 1",
				StockQuantity:    100,
				LowStockThresold: 10,
			},
			{
				Name:             "High Stock Product 2",
				StockQuantity:    50,
				LowStockThresold: 5,
			},
		}

		mockService := &mockProductService{
			products: highStockProducts,
		}
		handler := NewProductHandler(mockService)

		app := fiber.New()
		app.Get("/products", handler.GetAllProducts())

		req := httptest.NewRequest("GET", "/products?low-stock=true", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != 200 {
			t.Errorf("expected status 200, got %d", resp.StatusCode)
		}

		var response struct {
			Data []product.Product `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		if len(response.Data) != 0 {
			t.Errorf("expected 0 low stock products, got %d", len(response.Data))
		}
	})

	t.Run("ignores low-stock filter when value is not 'true'", func(t *testing.T) {
		mockService := &mockProductService{
			products: testProducts,
		}
		handler := NewProductHandler(mockService)

		app := fiber.New()
		app.Get("/products", handler.GetAllProducts())

		// Test with low-stock=false
		req := httptest.NewRequest("GET", "/products?low-stock=false", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}

		var response struct {
			Data []product.Product `json:"data"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			t.Fatal(err)
		}

		// Should return all products since filter is not "true"
		if len(response.Data) != 4 {
			t.Errorf("expected 4 products, got %d", len(response.Data))
		}
	})

	t.Run("handles service error correctly", func(t *testing.T) {
		mockService := &mockProductService{
			getError: apperrors.NewDatabaseError("database connection failed"),
		}
		handler := NewProductHandler(mockService)

		app := fiber.New()
		app.Get("/products", handler.GetAllProducts())

		req := httptest.NewRequest("GET", "/products", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != 500 {
			t.Errorf("expected status 500, got %d", resp.StatusCode)
		}
	})

	t.Run("handles service error correctly with low-stock filter", func(t *testing.T) {
		mockService := &mockProductService{
			getError: apperrors.NewDatabaseError("database connection failed"),
		}
		handler := NewProductHandler(mockService)

		app := fiber.New()
		app.Get("/products", handler.GetAllProducts())

		req := httptest.NewRequest("GET", "/products?low-stock=true", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatal(err)
		}

		// Should return error before filtering is applied
		if resp.StatusCode != 500 {
			t.Errorf("expected status 500, got %d", resp.StatusCode)
		}
	})
}
