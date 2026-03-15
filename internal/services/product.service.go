package services

import (
	"context"
	"food-ordering/internal/dto"
	"food-ordering/internal/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetProducts(ctx context.Context) ([]dto.ProductResponse, error) {

	products, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var response []dto.ProductResponse

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

		response = append(response, dto.ProductResponse{
			Id:          p.Id.Hex(),
			Name:        p.Name,
			Cuisines:    p.Cuisines,
			Category:    p.Category,
			Price:       p.Price,
			Description: p.Description,
			Rating:      p.Rating,
			Image:       images,
		})
	}

	return response, nil
}

func (s *ProductService) GetProduct(ctx context.Context, id string) (*dto.ProductResponse, error) {

	ObjectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	product, err := s.repo.FindById(ctx, ObjectId)
	if err != nil {
		return nil, err
	}

	var images []dto.Image
	for _, img := range product.Image {
		images = append(images, dto.Image{
			Thumbnail: img.Thumbnail,
			Mobile:    img.Mobile,
			Tablet:    img.Tablet,
			Desktop:   img.Desktop,
		})
	}

	return &dto.ProductResponse{
		Id:          product.Id.Hex(),
		Name:        product.Name,
		Cuisines:    product.Cuisines,
		Category:    product.Category,
		Price:       product.Price,
		Description: product.Description,
		Rating:      product.Rating,
		Image:       images,
	}, nil
}
