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

type SSG struct {
	SSGId   int
	Station string
	Sector  string
	Pgroup  string
}
type FRE struct {
	FREId        int
	FlatNo       string
	ReservePrice string
	EMD          string
}

type DivisionView struct {
	DivisionID int
	Name       string
}
