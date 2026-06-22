import "./main.tsp";

using Http;

namespace AppAPI;

model Product {
  id: int64;
  name: string;
  price: float64;
  description: string;
  category: string;
  stock: int32;
}

model CreateProductRequest {
  @minLength(1) name: string;
  @minValue(0) price: float64;
  @minLength(1) description: string;
  @minLength(1) category: string;
  @minValue(0) stock: int32;
}

model ErrorResponse {
  error: string;
  fields?: Record<string>;
}

@tag("Products")
@route("/api/v1/products")
interface Products {
  @post
  create(@body body: CreateProductRequest):
    | {
        @statusCode code: 201;
        @body body: Product;
      }
    | {
        @statusCode code: 400;
        @body body: ErrorResponse;
      };

  @get
  getAll(): {
    @statusCode code: 200;
    @body body: Product[];
  };

  @get
  @route("/{id}")
  getById(@path id: int64):
    | {
        @statusCode code: 200;
        @body body: Product;
      }
    | {
        @statusCode code: 404;
        @body body: ErrorResponse;
      };
}
