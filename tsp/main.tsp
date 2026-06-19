import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";
import "./product.tsp";

using Http;
using Rest;

@service(#{ title: "Go CRUD2 Product API" })
namespace ProductAPI;