import "./main.tsp";
using Http;

namespace AppAPI;

model User {
  id: int64;
  name: string;
  email: string;
  is_active: boolean;
  created_at: utcDateTime;
  updated_at: utcDateTime;
}

model CreateUserInput {
  name: string;
  email: string;
}

model UpdateUserInput {
  name: string;
  email: string;
  is_active?: boolean;
}

model CountResponse {
  count: int64;
}

model LoginRequest {
  email: string;
  password: string;
}

model TokenResponse {
  token: string;
}

@tag("Auth")
@route("/auth")
interface Auth {
  @post
  @route("/login")
  login(@body body: LoginRequest): TokenResponse;
}

@tag("Users")
@route("/api/v1/users")
interface Users {
  @post
  create(@body body: CreateUserInput): User;

  @get
  list(): User[];

  @get
  @route("/{id}")
  getById(@path id: int64): User;

  @put
  @route("/{id}")
  update(@path id: int64, @body body: UpdateUserInput): User;

  @delete
  @route("/{id}")
  delete(@path id: int64): void;

  @get
  @route("/active/{active}")
  getByActive(@path active: boolean): User[];

  @get
  @route("/count")
  getCount(): CountResponse;
}
