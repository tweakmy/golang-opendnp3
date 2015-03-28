package DNP3

import (

)

type GrpVar struct {
	Grp uint
	Var uint
}

var (
	Grp02Var02 = GrpVar{2,2}
	Grp32Var02 = GrpVar{32,2}
)