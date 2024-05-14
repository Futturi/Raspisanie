package script

import (
	"encoding/json"
	"fmt"
	"github.com/Futturi/Raspisanie/pkg"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"github.com/tealeg/xlsx"
	"log/slog"
	"strings"
	"time"
)

type Raspisanie struct {
	Pair1 string `json:"pair1"`
	Pair2 string `json:"pair2"`
	Pair3 string `json:"pair3"`
	Pair4 string `json:"pair4"`
	Pair5 string `json:"pair5"`
	Pair6 string `json:"pair6"`
	Pair7 string `json:"pair7"`
	Aud1  string `json:"aud1"`
	Aud2  string `json:"aud2"`
	Aud3  string `json:"aud3"`
	Aud4  string `json:"aud4"`
	Aud5  string `json:"aud5"`
	Aud6  string `json:"aud6"`
	Aud7  string `json:"aud7"`
	Prep1 string `json:"prep1"`
	Prep2 string `json:"prep2"`
	Prep3 string `json:"prep3"`
	Prep4 string `json:"prep4"`
	Prep5 string `json:"prep5"`
	Prep6 string `json:"prep6"`
	Prep7 string `json:"prep7"`
	Vid1  string `json:"vid1"`
	Vid2  string `json:"vid2"`
	Vid3  string `json:"vid3"`
	Vid4  string `json:"vid4"`
	Vid5  string `json:"vid5"`
	Vid6  string `json:"vid6"`
	Vid7  string `json:"vid7"`
}

func InitScript() {
	time.Sleep(time.Second)
	pkfg := pkg.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		Dbname:   viper.GetString("db.namedb"),
		Sslmode:  viper.GetString("db.sslmode"),
	}

	db, err := pkg.InitPostgres(pkfg)
	if err != nil {
		slog.Error("error with db", slog.Any("error", err))
	}

	file, err := xlsx.OpenFile("rasp.xlsx")
	if err != nil {
		slog.Error("error with opening file", slog.Any("error", err))
	}
	arra := make([][]string, 375)
	for _, sheet := range file.Sheets {
		for _, row := range sheet.Rows {
			for i := 4; i < len(row.Cells); i++ {
				arra[i] = append(arra[i], row.Cells[i].String())
			}
		}
	}
	newMap := make(map[string][]string)
	mapfaut := make(map[string][]string)
	mapfprep := make(map[string][]string)
	mapfvid := make(map[string][]string)
	cur := ""
	for ind, v := range arra {
		if len(v) < 1 {
			continue
		}
		if cu(v[1]) {
			newMap[v[1]] = arra[ind][3:]
			cur = v[1]
		}
		if v[2] == "ФИО преподавателя" {
			mapfprep[cur] = arra[ind][3:]
		}
		if v[2] == `№ 
ауд.` {
			mapfaut[cur] = arra[ind][3:]
		}
		if v[2] == `Вид
занятий` {
			mapfvid[cur] = arra[ind][3:]
		}
	}

	for k := range newMap {
		if len(newMap[k]) != len(mapfaut[k]) {
			newMap[k] = newMap[k][:len(mapfaut[k])-1]
		}
	}
	if err := Migrat(viper.GetString("db.host")); err != nil {
		slog.Error("error with migratedb", slog.Any("error", err))
	}
	for grou, v := range newMap {
		for week := 0; week < 17; week++ {
			for day := 0; day < 8; day++ {
				k := newPars(v, week, mapfaut[grou], mapfprep[grou], mapfvid[grou])
				b, err := json.Marshal(k)
				if err != nil {
					slog.Error("error with marshalling", slog.Any("error", err))
				}
				query := fmt.Sprint("INSERT INTO pairs(groups, week, day, pairs) VALUES($1, $2, $3, $4)")
				_, err = db.Exec(query, grou, week, day, string(b))
				if err != nil {
					slog.Error("error with inserting to db", slog.Any("error", err))
				}
			}
		}
	}
	slog.Info("all ok")
}

type T struct {
	Pair1 string `json:"pair1"`
	Pair2 string `json:"pair2"`
	Pair3 string `json:"pair3"`
	Pair4 string `json:"pair4"`
	Pair5 string `json:"pair5"`
	Pair6 string `json:"pair6"`
	Pair7 string `json:"pair7"`
}

func cu(a string) bool {
	c := 0
	for _, v := range a {
		if v == '-' {
			c++
		}
	}
	return c == 2
}

func isSlovo(s string) bool {
	if len(s) == 0 {
		return false
	}
	if strings.Contains(s, "кр. 1") {
		return false
	}
	return strings.ContainsAny(s, "0123456789")
}

func newPars(arr []string, k int, aut, prep, vid []string) []Raspisanie {
	result := make([]Raspisanie, 0)
	arr1 := make([]string, 7) // это для пар, нужно сделать для всего остального такие же
	arra := make([]string, 7)
	arrp := make([]string, 7)
	arrv := make([]string, 7)
	if k%2 == 0 {
		c := 0
		for ind, v := range arr {
			if ind%2 == 0 {
				if isSlovo(v) && !(strings.Contains(v, fmt.Sprint(k-1)) && !strings.Contains(v, "1"+fmt.Sprint(k-1))) {
					arr1[c] = ""
					c++
					if c == 7 {
						result = append(result, Raspisanie{
							Pair1: arr1[0],
							Pair2: arr1[1],
							Pair3: arr1[2],
							Pair4: arr1[3],
							Pair5: arr1[4],
							Pair6: arr1[5],
							Pair7: arr1[6],
							Aud1:  arra[0],
							Aud2:  arra[1],
							Aud3:  arra[2],
							Aud4:  arra[3],
							Aud5:  arra[4],
							Aud6:  arra[5],
							Aud7:  arra[6],
							Prep1: arrp[0],
							Prep2: arrp[1],
							Prep3: arrp[2],
							Prep4: arrp[3],
							Prep5: arrp[4],
							Prep6: arrp[5],
							Prep7: arrp[6],
							Vid1:  arrv[0],
							Vid2:  arrv[1],
							Vid3:  arrv[2],
							Vid4:  arrv[3],
							Vid5:  arrv[4],
							Vid6:  arrv[5],
							Vid7:  arrv[6],
						})
						c = 0
					}
				} else {
					arr1[c] = v
					arrv[c] = vid[ind]
					arrp[c] = prep[ind]
					arra[c] = aut[ind]
					c++
					if c == 7 {
						result = append(result, Raspisanie{
							Pair1: arr1[0],
							Pair2: arr1[1],
							Pair3: arr1[2],
							Pair4: arr1[3],
							Pair5: arr1[4],
							Pair6: arr1[5],
							Pair7: arr1[6],
							Aud1:  arra[0],
							Aud2:  arra[1],
							Aud3:  arra[2],
							Aud4:  arra[3],
							Aud5:  arra[4],
							Aud6:  arra[5],
							Aud7:  arra[6],
							Prep1: arrp[0],
							Prep2: arrp[1],
							Prep3: arrp[2],
							Prep4: arrp[3],
							Prep5: arrp[4],
							Prep6: arrp[5],
							Prep7: arrp[6],
							Vid1:  arrv[0],
							Vid2:  arrv[1],
							Vid3:  arrv[2],
							Vid4:  arrv[3],
							Vid5:  arrv[4],
							Vid6:  arrv[5],
							Vid7:  arrv[6],
						})
						c = 0
					}
				}
			}
		}
	}
	if k%2 != 0 {
		c := 0
		for ind, v := range arr {
			if ind%2 != 0 {
				if isSlovo(v) && !(strings.Contains(v, fmt.Sprint(k-1)) && !strings.Contains(v, "1"+fmt.Sprint(k-1))) {
					arr1[c] = ""
					c++
					if c == 7 {
						result = append(result, Raspisanie{
							Pair1: arr1[0],
							Pair2: arr1[1],
							Pair3: arr1[2],
							Pair4: arr1[3],
							Pair5: arr1[4],
							Pair6: arr1[5],
							Pair7: arr1[6],
							Aud1:  arra[0],
							Aud2:  arra[1],
							Aud3:  arra[2],
							Aud4:  arra[3],
							Aud5:  arra[4],
							Aud6:  arra[5],
							Aud7:  arra[6],
							Prep1: arrp[0],
							Prep2: arrp[1],
							Prep3: arrp[2],
							Prep4: arrp[3],
							Prep5: arrp[4],
							Prep6: arrp[5],
							Prep7: arrp[6],
							Vid1:  arrv[0],
							Vid2:  arrv[1],
							Vid3:  arrv[2],
							Vid4:  arrv[3],
							Vid5:  arrv[4],
							Vid6:  arrv[5],
							Vid7:  arrv[6],
						})
						c = 0
					}
				} else {
					arr1[c] = v
					arrv[c] = vid[ind]
					arrp[c] = prep[ind]
					arra[c] = aut[ind]
					c++
					if c == 7 {
						result = append(result, Raspisanie{
							Pair1: arr1[0],
							Pair2: arr1[1],
							Pair3: arr1[2],
							Pair4: arr1[3],
							Pair5: arr1[4],
							Pair6: arr1[5],
							Pair7: arr1[6],
							Aud1:  arra[0],
							Aud2:  arra[1],
							Aud3:  arra[2],
							Aud4:  arra[3],
							Aud5:  arra[4],
							Aud6:  arra[5],
							Aud7:  arra[6],
							Prep1: arrp[0],
							Prep2: arrp[1],
							Prep3: arrp[2],
							Prep4: arrp[3],
							Prep5: arrp[4],
							Prep6: arrp[5],
							Prep7: arrp[6],
							Vid1:  arrv[0],
							Vid2:  arrv[1],
							Vid3:  arrv[2],
							Vid4:  arrv[3],
							Vid5:  arrv[4],
							Vid6:  arrv[5],
							Vid7:  arrv[6],
						})
						c = 0
					}
				}
			}
		}
	}
	return result
}

func Migrat(host string) error {
	time.Sleep(time.Second)
	m, err := migrate.New(
		"file://migrate", fmt.Sprintf(
			"postgres://root:root@%s:5432/root?sslmode=disable", host))
	if err != nil {
		return err
	}
	return m.Up()
}
