// pkg/graphql/models.go
package graphql

type Account struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Orders []*Order `json:"orders"`
}
