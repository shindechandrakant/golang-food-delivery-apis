package services

import (
	"context"
	"fmt"
	"food-ordering/internal/dto"
	"food-ordering/internal/models"
	"food-ordering/internal/repository"
)

type CartService struct {
	repo repository.CartRepository
}

func NewCartService(repo repository.CartRepository) *CartService {
	return &CartService{repo: repo}
}

func (s *CartService) AddItem(ctx context.Context, userId string, req dto.AddCartItemRequest) error {
	cart, err := s.repo.FindByUser(ctx, userId)
	if err != nil {
		return err
	}

	for i, item := range cart.Items {
		if item.ProductId == req.ProductId {
			cart.Items[i].Quantity += req.Quantity
			return s.repo.Save(ctx, cart)
		}
	}

	cart.Items = append(cart.Items, models.CartItem{
		ProductId: req.ProductId,
		Quantity:  req.Quantity,
	})

	return s.repo.Save(ctx, cart)
}

func (s *CartService) GetCart(ctx context.Context, userId string) (*models.Cart, error) {
	return s.repo.FindByUser(ctx, userId)
}

func (s *CartService) UpdateItem(ctx context.Context, userId string, productId string, quantity int) error {
	cart, err := s.repo.FindByUser(ctx, userId)
	if err != nil {
		return err
	}

	for i, item := range cart.Items {
		if item.ProductId == productId {
			cart.Items[i].Quantity = quantity
			return s.repo.Save(ctx, cart)
		}
	}

	return fmt.Errorf("item not found in cart")
}

func (s *CartService) RemoveItem(ctx context.Context, userId string, productId string) error {
	cart, err := s.repo.FindByUser(ctx, userId)
	if err != nil {
		return err
	}

	filtered := cart.Items[:0]
	for _, item := range cart.Items {
		if item.ProductId != productId {
			filtered = append(filtered, item)
		}
	}
	cart.Items = filtered

	return s.repo.Save(ctx, cart)
}

func (s *CartService) ClearCart(ctx context.Context, userId string) error {
	return s.repo.Clear(ctx, userId)
}
