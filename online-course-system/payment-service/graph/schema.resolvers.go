package graph

import (
	"payment-service/prisma/db"
)

// Resolver struct should contain the Prisma client
type Resolver struct {
	Prisma *db.PrismaClient
}

// Mutation returns the mutation resolver
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

// Query returns the query resolver
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
