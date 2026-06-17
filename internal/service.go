package internal

type recordService struct {
	repo *recordRepository
}

func NewRecordService(r *recordRepository) *recordService {
	return &recordService{
		repo: r,
	}
}

func (s *recordService) GetAll() ([]Record, error) {
	return s.repo.GetAll()
}

func (s *recordService) GetByID(id int) (Record, error) {
	return s.repo.GetByID(id)
}

func (s *recordService) Create(record Record) (Record, error) {
	return s.repo.Create(record)
}

func (s *recordService) Update(id int, record Record) (Record, error) {
	return s.repo.Update(id, record)
}
