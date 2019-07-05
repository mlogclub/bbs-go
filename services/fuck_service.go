package services

var FuckService = newFuckService()

func newFuckService() *fuckService {
	return &fuckService{
		count: 0,
	}
}

type fuckService struct {
	count int
}

func (this *fuckService) Incr() int {
	this.count++
	return this.count
}

func (this *fuckService) Get() int {
	return this.count
}
