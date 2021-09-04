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

type Ssg struct {
	SsgId   int
	Station string
	Sector  string
	Pgroup  string
}
type Fre struct {
	FreId        int
	FlatNo       string
	ReservePrice string
	EMD          string
}

type Division struct {
	DivisionID int
	Name       string
}
