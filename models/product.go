package models

type Product struct {
	ID       int
	Name     string
	Price    float64
	Quantity int
}

var Products = []Product{
	{ID: 1, Name: "Product 1", Price: 19.99, Quantity: 100},
	{ID: 2, Name: "Product 2", Price: 29.99, Quantity: 50},
	{ID: 3, Name: "Product 3", Price: 9.99, Quantity: 200},
}
