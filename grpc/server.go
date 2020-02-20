package grpc

import (
	"context"

	"go.opencensus.io/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"wiliam.dev/product"
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

	// TODO: Preallocate our initial slice for pagination values.
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

	span.SetStatus(trace.Status{
		Code:    trace.StatusCodeOK,
		Message: "Ok",
	})

	response := &grpcv1.ListProductsResponse{Products: p}
	return response, nil
}

// CreateProduct ...
func (s *ProductAPIServer) CreateProduct(ctx context.Context, r *grpcv1.CreateProductRequest) (*grpcv1.CreateProductResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateProduct not implemented")
}

// NewProductAPIServer create a grpc product API server
func NewProductAPIServer(useCase product.UseCase) *ProductAPIServer {
	return &ProductAPIServer{UseCase: useCase}
}
