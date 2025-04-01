package graph

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

// Query returns the query resolver
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

// QueryResolver handles queries like `getPayment` and `listPayments`
type queryResolver struct{ *Resolver }

// MutationResolver handles mutations like `createPayment`
type mutationResolver struct{ *Resolver }
