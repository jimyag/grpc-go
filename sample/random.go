package sample

import (
	"github.com/google/uuid"
	"github.com/jimyag/grpc-go/pb"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func randomKeyboardLayout() pb.KeyBoard_Layout {
	switch rand.Intn(3) {
	case 1:
		return pb.KeyBoard_QWERTY
	case 2:
		return pb.KeyBoard_QWERTZ
	default:
		return pb.KeyBoard_AZERTY

	}
}

func randomBool() bool {
	return rand.Intn(2) == 1
}

func randomCPUBand() string {
	return randomStringFormat("Intel", "AMD")
}
func randomCPUName(brand string) string {
	if brand == "Intel" {
		return randomStringFormat(
			"Core i9-9800k",
			"Core i7-9750H",
			"Core i7-8750H")
	} else {
		return randomStringFormat("Ryzen 7 pro 2700U",
			"Ryzen 5 pro 3700U",
			"Ryzen 3 pro 3200GE")
	}
}

func randomInt(min int, max int) int {
	return min + rand.Intn(max-min+1)
}

func randomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
func randomStringFormat(a ...string) string {
	n := len(a)
	if n == 0 {
		return ""
	}
	return a[rand.Intn(n)]
}

func randomGPUBrand() string {
	return randomStringFormat("NVIDIA", "AMD")
}

func randomGpuName(brand string) string {
	if brand == "NVIDIA" {
		return randomStringFormat(
			"RTX 2060",
			"RTX 2070", "GTX 3060",
			"GTX 3070")
	} else {
		return randomStringFormat(
			"RX 590",
			"RX 580",
			"RX 5700-XT")
	}
}

func randomPanel() pb.Screen_Panel {
	if rand.Intn(2) == 1 {
		return pb.Screen_IPS
	}
	return pb.Screen_OLED
}

func randomResolution() *pb.Screen_Resolution {
	height := randomInt(1080, 4230)
	width := height * 16 / 9
	res := &pb.Screen_Resolution{
		Width:  uint32(width),
		Height: uint32(height),
	}
	return res
}

func randomId() string {
	return uuid.New().String()
}

func randomLapTopBrand() string {
	return randomStringFormat("Apple", "Dell", "Lenovo")
}
func randomLapTopName(brand string) string {
	switch brand {
	case "Apple":
		return randomStringFormat("Macbook air", "Macbook pro")
	case "Dell":
		return randomStringFormat("XPS", "Vostro")
	default:
		return randomStringFormat("Thinkpad X1", "Thinkpad X3", "Thinkpad T2")

	}
}
