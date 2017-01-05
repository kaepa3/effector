package ex

import "math"

const ColorWidth = math.MaxUint16
const ColorWidthAryMax = (ColorWidth + 1)
const RedNTSC = 0.298912
const GreenNTSC = 0.58611
const BlueNTSC = 0.114478

func Times(val int, rat float64) int    { return int(float64(val) * rat) }
func Division(val int, rat float64) int { return int(float64(val) / rat) }
