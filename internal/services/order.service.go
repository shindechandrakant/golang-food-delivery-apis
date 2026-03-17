package services

import (
	"context"
	"fmt"
	"food-ordering/internal/dto"
	"food-ordering/internal/models"
	"food-ordering/internal/promo"
	"food-ordering/internal/repository"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type OrderService struct {
	orderRepo   *repository.OrderRepository
	productRepo *repository.ProductRepository
	promo       *promo.Validator
}

func NewOrderService(
	orderRepo *repository.OrderRepository,
	productRepo *repository.ProductRepository,
	promoValidator *promo.Validator,
) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		promo:       promoValidator,
	}
}

func (s *OrderService) PlaceOrder(ctx context.Context, req dto.PlaceOrderRequest) (*dto.OrderResponse, error) {

	// Validate coupon code if provided.
	if req.CouponCode != "" && !s.promo.IsValid(req.CouponCode) {
		return nil, fmt.Errorf("invalid coupon code")
	}

	// Parse and deduplicate product IDs.
	objectIds := make([]bson.ObjectID, 0, len(req.Items))
	seen := make(map[string]struct{})
	for _, item := range req.Items {
		if _, dup := seen[item.ProductId]; dup {
			continue
		}
		seen[item.ProductId] = struct{}{}

		objId, err := bson.ObjectIDFromHex(item.ProductId)
		if err != nil {
			return nil, fmt.Errorf("invalid productId %q", item.ProductId)
		}
		objectIds = append(objectIds, objId)
	}

	// Fetch all products in one query.
	products, err := s.productRepo.FindByIds(ctx, objectIds)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}

	// Verify every requested product exists.
	productMap := make(map[string]models.Product, len(products))
	for _, p := range products {
		productMap[p.Id.Hex()] = p
	}
	for _, item := range req.Items {
		if _, ok := productMap[item.ProductId]; !ok {
			return nil, fmt.Errorf("product %q not found", item.ProductId)
		}
	}

	// Build order items.
	orderItems := make([]models.OrderItem, len(req.Items))
	for i, item := range req.Items {
		orderItems[i] = models.OrderItem{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
		}
	}

	order := &models.Order{
		UUID:  uuid.New().String(),
		Items: orderItems,
	}

	if err := s.orderRepo.Save(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to save order: %w", err)
	}

	// Build response.
	respItems := make([]dto.OrderItemResponse, len(order.Items))
	for i, item := range order.Items {
		respItems[i] = dto.OrderItemResponse{
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
		}
	}

	respProducts := make([]dto.ProductResponse, 0, len(products))
	for _, p := range products {
		var images []dto.Image
		for _, img := range p.Image {
			images = append(images, dto.Image{
				Thumbnail: img.Thumbnail,
				Mobile:    img.Mobile,
				Tablet:    img.Tablet,
				Desktop:   img.Desktop,
			})
		}
		respProducts = append(respProducts, dto.ProductResponse{
			Id:          p.Id.Hex(),
			Name:        p.Name,
			Category:    p.Category,
			Price:       p.Price,
			Cuisines:    p.Cuisines,
			Description: p.Description,
			Rating:      p.Rating,
			Image:       images,
		})
	}

	return &dto.OrderResponse{
		Id:       order.UUID,
		Items:    respItems,
		Products: respProducts,
	}, nil
}
