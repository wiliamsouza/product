package usecase

import (
	"context"
	"log"
	"strings"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"go.opencensus.io/trace"
	"google.golang.org/grpc/metadata"

	"wiliam.dev/product"
	"wiliam.dev/product/entity"
	"wiliam.dev/product/grpc/client/promotion/v1alpha1"
)

// Ensure PromotionUseCase implements product.UseCase interface.
var _ product.UseCase = &PromotionUseCase{}

//PromotionUseCase implements and extend product.UseCase interface using decorator design pattern.
type PromotionUseCase struct {
	Promotion v1alpha1.PromotionAPIClient
	Product   product.UseCase
}

// List products add promotion informations for products.
func (u *PromotionUseCase) List(ctx context.Context) ([]*entity.Product, error) {
	ctx, span := trace.StartSpan(ctx, "usecase.PromotionUseCase.List")
	defer span.End()
	products, err := u.Product.List(ctx)
	if err != nil {
		wrappedErr := errors.Wrapf(
			err,
			"struct=usecase.PromotionUseCase, method=List, error=product_list_failed",
		)
		return nil, wrappedErr
	}

	userID := ""
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if uID, ok := md["x-user-id"]; ok {
			userID = strings.Join(uID, ",")
		}
	}

	for _, p := range products {
		request := v1alpha1.RetrievePromotionRequest{
			UserId:    uuid.NewV4().String(),
			ProductId: p.ID,
		}
		_, err := u.Promotion.RetrievePromotion(ctx, &request)
		if err != nil {
			log.Printf("struct=usecase.PromotionUseCase, method=List, error=%s", err)
		}
	}
	return products, nil
}

// Create product just call product use case.
func (u *PromotionUseCase) Create(ctx context.Context, p *entity.Product) (*entity.Product, error) {
	return u.Product.Create(ctx, p)
}

//NewPromotionUseCase create a product use case instance.
func NewPromotionUseCase(product product.UseCase, promotion v1alpha1.PromotionAPIClient) *PromotionUseCase {
	return &PromotionUseCase{
		Promotion: promotion,
		Product:   product,
	}
}
