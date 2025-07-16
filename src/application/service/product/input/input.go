package input

type CreateProductInput struct {
	Name        string `validate:"required,max=200"`
	Description string `validate:"max=500"`
}
