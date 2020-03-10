package usecase

import (
	"context"
	"log"
	"math"
	"strings"

	"github.com/pkg/errors"
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
		ctxPromotion, spanPromotion := trace.StartSpan(ctx, "usecase.PromotionUseCase.RetrievePromotion")
		defer spanPromotion.End()

		request := v1alpha1.RetrievePromotionRequest{
			UserId:    userID,
			ProductId: p.ID,
		}

		promotion, err := u.Promotion.RetrievePromotion(ctxPromotion, &request)
		if err != nil {
			log.Printf("struct=usecase.PromotionUseCase, method=List, error=%s", err)
			spanPromotion.SetStatus(trace.Status{Code: trace.StatusCodeInternal, Message: err.Error()})

			continue
		}
		d := promotion.Discounts[0]
		valueInCents := ((float32(p.PriceInCents) * d.Pct) / 100)
		// Using math.Ceil here to avoid problem with products
		// with lower price ie: 10 cents where if given 5% of
		// discount will result in 0.5 cents which is not compatible
		// with API definition that required value_in_cents to be int.
		p.Discount = entity.Discount{
			Pct:          d.Pct,
			ValueInCents: int32(math.Ceil(float64(valueInCents))),
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
