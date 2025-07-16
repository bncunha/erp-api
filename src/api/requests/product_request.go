package request

type CreateProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}