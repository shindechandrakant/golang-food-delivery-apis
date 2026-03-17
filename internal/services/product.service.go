package services

import (
	"context"
	"food-ordering/internal/dto"
	"food-ordering/internal/models"
	"food-ordering/internal/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetProducts(ctx context.Context, filter dto.ProductFilter) ([]dto.ProductResponse, error) {
	products, err := s.repo.FindWithFilter(ctx, filter)
	if err != nil {
		return nil, err
	}

	response := make([]dto.ProductResponse, 0, len(products))
	for _, p := range products {
		response = append(response, mapProductToDTO(p))
	}
	return response, nil
}

func (s *ProductService) GetProduct(ctx context.Context, id string) (*dto.ProductResponse, error) {
	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	product, err := s.repo.FindById(ctx, objectId)
	if err != nil {
		return nil, err
	}

	resp := mapProductToDTO(*product)
	return &resp, nil
}

func mapProductToDTO(p models.Product) dto.ProductResponse {
	images := make([]dto.Image, 0, len(p.Image))
	for _, img := range p.Image {
		images = append(images, dto.Image{
			Thumbnail: img.Thumbnail,
			Mobile:    img.Mobile,
			Tablet:    img.Tablet,
			Desktop:   img.Desktop,
		})
	}
	return dto.ProductResponse{
		Id:          p.Id.Hex(),
		Name:        p.Name,
		Cuisines:    p.Cuisines,
		Category:    p.Category,
		Price:       p.Price,
		Description: p.Description,
		Rating:      p.Rating,
		Image:       images,
	}
}
