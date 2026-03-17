package services

import (
	"context"
	"fmt"
	"food-ordering/internal/dto"
	"food-ordering/internal/models"
	"food-ordering/internal/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type CartService struct {
	repo *repository.CartRepository
}

func NewCartService(cartRepository *repository.CartRepository) *CartService {
	return &CartService{repo: cartRepository}
}

func (s *CartService) AddItem(ctx context.Context, userId string, req dto.AddCartItemRequest) error {

	productId, _ := bson.ObjectIDFromHex(req.ProductId)

	cart, _ := s.repo.FindByUser(ctx, userId)
	if cart == nil {
		cart = &models.Cart{
			UserId: userId,
		}
	}
	found := false

	for i, item := range cart.Items {
		if item.ProductId == productId {
			cart.Items[i].Quantity += req.Quantity
			found = true
			break
		}
	}

	if !found {
		cart.Items = append(cart.Items, models.CartItem{
			ProductId: productId,
			Quantity:  req.Quantity,
		})
	}

	return s.repo.Save(ctx, cart)
}

func (s *CartService) GetCart(ctx context.Context, userId string) (*models.Cart, error) {
	return s.repo.FindByUser(ctx, userId)
}

func (s *CartService) UpdateItem(ctx context.Context, userId string, productId string, quantity int) error {
	objId, err := bson.ObjectIDFromHex(productId)
	if err != nil {
		return err
	}

	cart, err := s.repo.FindByUser(ctx, userId)
	if err != nil {
		return err
	}

	found := false
	for i, item := range cart.Items {
		if item.ProductId == objId {
			cart.Items[i].Quantity = quantity
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("item not found in cart")
	}

	return s.repo.Save(ctx, cart)
}

func (s *CartService) RemoveItem(ctx context.Context, userId string, productId string) error {
	objId, err := bson.ObjectIDFromHex(productId)
	if err != nil {
		return err
	}

	cart, err := s.repo.FindByUser(ctx, userId)
	if err != nil {
		return err
	}

	filtered := cart.Items[:0]
	for _, item := range cart.Items {
		if item.ProductId != objId {
			filtered = append(filtered, item)
		}
	}
	cart.Items = filtered

	return s.repo.Save(ctx, cart)
}

func (s *CartService) ClearCart(ctx context.Context, userId string) error {
	return s.repo.Clear(ctx, userId)
}
