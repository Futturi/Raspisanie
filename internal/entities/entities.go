package entities

type Group struct {
	Name string `json:"name" db:"groups"`
	Week int    `json:"week" db:"week"`
	Day  int    `json:"day" db:"day"`
}

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

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Group    string `json:"group"`
}
