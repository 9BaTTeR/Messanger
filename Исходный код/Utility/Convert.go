package Utility

func BoolToUInt(b bool) uint8 {
	if !b {
		return 0
	} else {
		return 1
	}
}

func UintToBool(u uint64) bool {
	if u == 0 {
		return false
	} else {
		return true
	}
}
