package partner

type ProductCollectionRelationshipService struct{}

func (s *ProductCollectionRelationshipService) Clone() *ProductCollectionRelationshipService {
	res := *s
	return &res
}
