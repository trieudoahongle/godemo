package cars

import (
	"fmt"
)

// This struct is lowercased. Packages that import `cars` will not be able
// to instantiate a literal car.
type car struct {
	// This will not be exported.
	speed int
	// This will be exported.
	Wheels int
}

// This func is uppercased. Packages that import `cars` can instantiate a
// new car using this method. This allows us more control around car creation.
func NewCar() car {
	return car{}
}

// Lowercased method name. It is not visible outside this package.
func (car car) secretHonk() {
	fmt.Println("Shhh! h o n k  h o n k")
}

// Uppercased method name. It is visible outside this package.
func (car car) HonkTheHorn() {
	car.secretHonk()

	fmt.Println("HONK HONK! I am a car!")
}
