package sample

import (
	"github.com/jimyag/grpc-go/pb"
)

//
// NewKeyboard
//  @Description: 生成键盘
//  @return *pb.KeyBoard
//
func NewKeyboard() *pb.KeyBoard {
	keyboard := &pb.KeyBoard{
		Layout:  randomKeyboardLayout(),
		Backlit: randomBool(),
	}
	return keyboard
}

//
// NewCPU
//  @Description: 生成 CPU 信息
//  @return *pb.CPU
//
func NewCPU() *pb.CPU {
	brand := randomCPUBand()
	name := randomCPUName(brand)
	numberCores := randomInt(2, 8)
	numberThreads := randomInt(numberCores, 14)
	minGhz := randomFloat(1.0, 4.5)
	maxGhz := randomFloat(minGhz, minGhz*3)
	cpu := &pb.CPU{
		Brand:         brand,
		Name:          name,
		NumberCores:   uint32(numberCores),
		NumberThreads: uint32(numberThreads),
		MinGhz:        float32(minGhz),
		MaxGhz:        float32(maxGhz),
	}
	return cpu
}

//
// NewGPU
//  @Description: 生成 GPU 信息
//  @return *pb.GPU
//
func NewGPU() *pb.GPU {
	brand := randomGPUBrand()
	gpuName := randomGpuName(brand)
	minGhz := randomFloat(1.0, 1.5)
	maxGhz := randomFloat(minGhz, minGhz*2)
	memory := &pb.Memory{
		Value: uint64(randomInt(2, 6)),
		Unit:  pb.Memory_GB,
	}
	gpu := &pb.GPU{
		Brand:  brand,
		Name:   gpuName,
		MinGhz: float32(minGhz),
		MaxGhz: float32(maxGhz),
		Memory: memory,
	}
	return gpu
}

//
// NewRAM
//  @Description: 生成 RAM
//  @return *pb.Memory
//
func NewRAM() *pb.Memory {
	memory := &pb.Memory{
		Value: uint64(randomInt(8, 64)),
		Unit:  pb.Memory_GB,
	}
	return memory
}

//
// NewSSD
//  @Description: 生成 SSD
//  @return *pb.Storage
//
func NewSSD() *pb.Storage {
	memory := &pb.Memory{
		Value: uint64(randomInt(256, 1024)),
		Unit:  pb.Memory_GB,
	}
	ssd := &pb.Storage{
		Driver: pb.Storage_SSD,
		Memory: memory,
	}
	return ssd
}

//
// NewScreen
//  @Description: 生成屏幕
//  @return *pb.Screen
//
func NewScreen() *pb.Screen {
	sizeInch := randomFloat(14.1, 22.0)
	screen := &pb.Screen{
		SizeInch:   float32(sizeInch),
		Resolution: randomResolution(),
		Panel:      randomPanel(),
		MultiTouch: randomBool(),
	}
	return screen
}

//
// NewLaptop
//  @Description: 电脑
//  @return *pb.Laptop
//
func NewLaptop() *pb.Laptop {
	brand := randomLapTopBrand()
	name := randomLapTopName(brand)
	laptop := &pb.Laptop{
		Id:          randomId(),
		Brand:       brand,
		Name:        name,
		Cpu:         NewCPU(),
		Memory:      NewRAM(),
		Gpu:         []*pb.GPU{NewGPU()},
		Storage:     []*pb.Storage{NewSSD(), NewSSD()},
		Screen:      NewScreen(),
		Keyboard:    NewKeyboard(),
		Weight:      &pb.Laptop_WeightKg{WeightKg: randomFloat(1.0, 3.0)},
		PriceUsd:    randomFloat(600.0, 1500.1),
		ReleaseYear: uint32(randomInt(2015, 2022)),
	}
	return laptop
}

func NewRandomLaptopScore() float64 {
	return float64(randomInt(1, 10))
}
