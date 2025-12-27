package models

type RegisterRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CatalogRequest struct {
	Filter string `json:"filter"`
}

type CatalogResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
}

type UpdateCatalogRequest struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Filter      string `json:"filter"`
}
type CreateCatalogRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Filter      string `json:"filter"`
}

type AddBucketRequest struct {
	Id int `json:"id"`
}

type CreateBuckerUserId struct {
	UserId    int   `json:"user_id"`
	Id_bucket []int `json:"id_bucket"`
	Count     int   `json:"count"`
}
