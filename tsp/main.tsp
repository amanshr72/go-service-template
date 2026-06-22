import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";
import "./product.tsp";
import "./user.tsp";

using Http;
using Rest;

@service(#{ title: "Go CRUD2 User/Product API" })
namespace AppAPI;