package repositories

type BoSeRepository interface {
}

type BoSeRepositoryImpl struct {
}

func NewBoSeRepository() BoSeRepository {
	return &BoSeRepositoryImpl{}
}

func (r *BoSeRepositoryImpl) CreateSession() {

}