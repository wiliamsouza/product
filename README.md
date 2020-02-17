# Product

Products is a gRPC service that implements products listings.

## Instalation

```bash
make deps
```

## Testing

```bash
make test
```

## Database


First install https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

This command bellow will create `product` table.

```bash
make migrate
```

If you want to populate database with mock data just run:

```bash
make populate
```

To create a new migration use:

```bash
migrate create -ext .sql -dir scripts/migrations/ <migration_name>
```

## Docker

Build:

```bash
docker build -t wiliam.dev/product .
```

Run:

```bash
docker run --rm wiliam.dev/product
```

## FAQ

Why use diferente models for each layer(ie: grpc, domain and database)?


For diferent in grpc and domain https://github.com/golang/protobuf/issues/156


For diference in domain and database:
```
--- FAIL: TestProductRepository (0.00s)
    --- FAIL: TestProductRepository/list_with_one_row_in_database (0.00s)
        product_test.go:46:
                Error Trace:    product_test.go:46
                Error:          Expected nil, but got: struct=repository.Product, method=List, error=select_fail: missing destination name price_in_cents in *[]*product.Product
                Test:           TestProductRepository/list_with_one_row_in_database
```
