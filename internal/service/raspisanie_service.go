package service

import (
	"github.com/Futturi/Raspisanie/internal/entities"
	"github.com/Futturi/Raspisanie/internal/repository"
)

type RaspService struct {
	repo repository.Raspisanie
}

func NewRaspService(repo repository.Raspisanie) *RaspService {
	return &RaspService{repo: repo}
}

//func zapisMap(files ...*xlsx.File) <-chan map[int][]string {
//	r := make(chan map[int][]string)
//	go func() {
//		result := make(map[int][]string)
//		mu := sync.Mutex{}
//		wg := &sync.WaitGroup{}
//		for _, file := range files{
//			for _, sheet := range file.Sheets {
//				for _, row := range sheet.Rows {
//					for i := 4; i < len(row.Cells); i++ {
//						i := i
//						row := row
//						wg.Add(1)
//						go func() {
//							defer wg.Done()
//							mu.Lock()
//							result[i] = append(result[i], row.Cells[i].String())
//							mu.Unlock()
//						}()
//					}
//				}
//			}
//		}
//		wg.Wait()
//		r <- result
//	}()
//	return r
//}
//
//func nMap(mapas chan map[int][]string) <-chan map[string][]string {
//	out := make(chan map[string][]string)
//	go func() {
//		defer close(out)
//		for mapa := range mapas{
//			newMap := make(map[string][]string)
//			for key, v := range mapa{
//				if cu(v[1]){
//					newMap[v[1]] = mapa[key][3:]
//				}
//			}
//			out <- newMap
//		}
//	}()
//	return out
//}
//
//func pari(mapas chan map[string][]string) chan entities.Raspisanie{
//	out := make(chan )
//}

func (r *RaspService) GetRasp(group entities.Group, gr string) (entities.Raspisanie, error) {
	group.Week = group.Week + 1
	group.Day = group.Day - 1
	return r.repo.GetRasp(group, gr)
}
