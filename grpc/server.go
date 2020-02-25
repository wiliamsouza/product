package grpc

import (
	"context"

	"go.opencensus.io/trace"
	"wiliam.dev/product"
	"wiliam.dev/product/entity"
	grpcv1 "wiliam.dev/product/grpc/v1beta1"
)

// Ensure ProductAPIServer implements grpcv1.ProductAPIServer.
var _ grpcv1.ProductAPIServer = &ProductAPIServer{}

// ProductAPIServer implements gRCP product API server
type ProductAPIServer struct {
	UseCase product.UseCase
}

// ListProducts returns a products list
func (s *ProductAPIServer) ListProducts(
	ctx context.Context,
	r *grpcv1.ListProductsRequest) (*grpcv1.ListProductsResponse, error) {

	ctx, span := trace.StartSpan(ctx, "grpc.ProductAPIServer.ListProducts")
	defer span.End()

	products, err := s.UseCase.List(ctx)
	if err != nil {
		span.SetStatus(trace.Status{
			Code:    trace.StatusCodeUnknown,
			Message: err.Error(),
		})
		return nil, err
	}

	span.Annotate([]trace.Attribute{
		trace.Int64Attribute("total", int64(len(products))),
	}, "Products listed")

	var p []*grpcv1.CreateProductResponse
	for _, product := range products {
		p = append(p,
			&grpcv1.CreateProductResponse{
				Id:           product.ID,
				Title:        product.Title,
				Description:  product.Description,
				PriceInCents: product.PriceInCents,
			})
	}

	response := &grpcv1.ListProductsResponse{Products: p}

	span.SetStatus(trace.Status{
		Code:    trace.StatusCodeOK,
		Message: "Ok",
	})

	return response, nil
}

// CreateProduct create a new product
func (s *ProductAPIServer) CreateProduct(ctx context.Context, r *grpcv1.CreateProductRequest) (*grpcv1.CreateProductResponse, error) {
	ctx, span := trace.StartSpan(ctx, "grpc.ProductAPIServer.CreateProduct")
	defer span.End()

	p := &entity.Product{
		Title:        r.Title,
		Description:  r.Description,
		PriceInCents: r.PriceInCents,
	}
	product, err := s.UseCase.Create(ctx, p)
	if err != nil {
		span.SetStatus(trace.Status{
			Code:    trace.StatusCodeUnknown,
			Message: err.Error(),
		})
		return nil, err
	}

	response := &grpcv1.CreateProductResponse{
		Id:           product.ID,
		Title:        product.Title,
		Description:  product.Description,
		PriceInCents: product.PriceInCents,
	}

	span.SetStatus(trace.Status{
		Code:    trace.StatusCodeOK,
		Message: "Ok",
	})

	return response, nil
}

// NewProductAPIServer create a grpc product API server
func NewProductAPIServer(useCase product.UseCase) *ProductAPIServer {
	return &ProductAPIServer{UseCase: useCase}
}
