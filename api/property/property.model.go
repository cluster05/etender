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
	Name             string `json:"name"`
	Text             string `json:"text"`
	SsgId            int    `json:"ssgId"`
	Children         []Tree `json:"children"`
	Level            string `json:"level"`
	ReservePrice     int    `json:"reservePrice"`
	EMD              int    `json:"emd"`
	EditNodeDisabled bool   `json:"editNodeDisabled"`
}
type SSG struct {
	SSGId        int    `json:"ssgId"`
	Station      string `json:"station"`
	Sector       string `json:"sector"`
	Pgroup       string `json:"pgroup"`
	ReservePrice string `json:"reservePrice"`
	EMD          string `json:"emd"`
}

type FRE struct {
	FREId  int    `json:"freId"`
	FlatNo string `json:"flatNo"`
}

type DivisionView struct {
	DivisionID int    `json:"divisionId"`
	Name       string `json:"name"`
}
