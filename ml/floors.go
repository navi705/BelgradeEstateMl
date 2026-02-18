package main

import "fmt"

func GetFloorLabel(f float32) string {
	switch f {
	case -3.0:
		return "Suteren"
	case -2.0:
		return "Prisuteren"
	case -0.5:
		return "Niskoprizemlje"
	case 0.0:
		return "Prizemlje"
	case 0.5:
		return "Visokoprizemlje"
	case 1000.0:
		return "Potkrovlje"
	case -5.0:
		return "Nepoznato"
	default:
		if f < 0 && f > -1 {
			return "Suteren/NPR"
		}
		return fmt.Sprintf("%g", f)
	}
}

func NormalizeFloorValue(f float32, total float32) float32 {
	if f == 1000.0 {
		if total > 0 {
			return total
		}
		return 1.0
	}
	if f == -5.0 || f == -100.0 {
		return 0.0
	}
	return f
}
