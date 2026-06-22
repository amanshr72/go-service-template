import "./main.tsp";

using Http;

namespace AppAPI;

model Product {
  id: int64;
  name: string;
  price: float64;
}

model CreateProductRequest {
  name: string;
  price: float64;
}

model ErrorResponse {
  error: string;
}

@tag("Products")
@route("/api/v1/products")
interface Products {
  @post
  create(@body body: CreateProductRequest): Product;

  @get
  getAll(): Product[];

  @get
  @route("/{id}")
  getById(@path id: int64): Product;
}
