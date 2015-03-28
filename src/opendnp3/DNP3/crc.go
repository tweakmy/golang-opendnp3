package DNP3

import (
	crc "github.com/howeyc/crc16"
)

var DNPTABLE = crc.MakeTable(0xA6BC)

