package property

type PropertyDTO struct {
	Division     string
	Station      string
	Sector       string
	Group        string
	FlatNo       string
	ReversePrice int
	EMD          int
}

type Tree struct {
	Text     string `json:"name"`
	SsgId    int    `json:"ssgId"`
	Children []Tree `json:"children"`
}
type SSG struct {
	SSGId   int    `json:"ssgId"`
	Station string `json:"station"`
	Sector  string `json:"sector"`
	Pgroup  string `json:"pgroup"`
}

type FRE struct {
	FREId        int    `json:"freId"`
	FlatNo       string `json:"flatNo"`
	ReservePrice string `json:"reservePrice"`
	EMD          string `json:"emd"`
}

type DivisionView struct {
	DivisionID int    `json:"divisionId"`
	Name       string `json:"name"`
}
