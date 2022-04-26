package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

//var wg = &sync.WaitGroup{}
//var mu = &sync.Mutex{}

type ShipType int

const numOfShips = 10
const (
	big   = 100
	med   = 50
	small = 10
)

const (
	PROVISION ShipType = iota
	FUEL
	ELECTRONICS
)

// type Ship interface {
// 	GetType()
// }

type Ship struct {
	id         int
	cap        int
	currentCap int
	SType      ShipType
}

// type FuelShip struct {
// 	id         int
// 	cap        int
// 	currentCap int
// 	SType      ShipType
// }

// type ElectShip struct {
// 	id         int
// 	cap        int
// 	currentCap int
// 	SType      ShipType
// }

// type ProvShip struct {
// 	id         int
// 	cap        int
// 	currentCap int
// 	SType      ShipType
// }

func createShips(id int) *Ship {
	rand.Seed(time.Now().UnixNano())
	capacity := make(map[int]int)
	capacity[0] = big
	capacity[1] = small
	capacity[2] = med
	cap := capacity[rand.Intn(3)]
	shipType := make(map[int]ShipType)
	shipType[0] = PROVISION
	shipType[1] = FUEL
	shipType[2] = ELECTRONICS
	sType := shipType[rand.Intn(3)]
	return &Ship{
		id:         id,
		cap:        cap,
		currentCap: cap,
		SType:      sType,
	}
}

type ShipGen struct {
	tunnel     *Tunnel
	numOfShips int
}

func newShipGen(t *Tunnel) *ShipGen {
	return &ShipGen{
		tunnel:     t,
		numOfShips: numOfShips,
	}
}

func (s *ShipGen) startShipGen(wg *sync.WaitGroup, t chan *Ship) {
	defer wg.Done()
	defer close(t)
	for i := 0; i < numOfShips; i++ {
		s.tunnel.addShipToTunnel(createShips(i))
		log.Println("Корабль сгенерирован")
		//time.Sleep(time.Second)

	}

	fmt.Println("OK")

	//s.tunnel.showSheepsIn()

}

type Tunnel struct {
	ships      chan *Ship
	tunnelSize int
	//isFilled   bool
}

func newTunnel() *Tunnel {

	shipsIn := make(chan *Ship, 5)
	return &Tunnel{
		ships:      shipsIn,
		tunnelSize: 3,
		//isFilled:   false,
	}
}

func (t *Tunnel) addShipToTunnel(sh *Ship) {
	fmt.Println("положили корабль type.... ", sh.SType)
	t.ships <- sh

}

type Dock struct {
	ships    chan *Ship
	shipType ShipType
}

func newDock(c chan *Ship, st ShipType) *Dock {

	return &Dock{
		ships:    c,
		shipType: st,
	}
}

type TunnelDistributor struct {
	tunnel    *Tunnel
	FuelChan  chan *Ship
	ElectChan chan *Ship
	ProvChan  chan *Ship
}

func newDistributor(t *Tunnel) *TunnelDistributor {
	fchan := make(chan *Ship)
	echan := make(chan *Ship)
	pchan := make(chan *Ship)
	return &TunnelDistributor{
		tunnel:    t,
		FuelChan:  fchan,
		ElectChan: echan,
		ProvChan:  pchan,
	}
}

func (dist *TunnelDistributor) distributor(wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(dist.FuelChan)
	defer close(dist.ProvChan)
	defer close(dist.ElectChan)
	for sh := range dist.tunnel.ships {
		switch sh.SType {
		case FUEL:
			dist.FuelChan <- sh
		case PROVISION:
			dist.ProvChan <- sh
		case ELECTRONICS:
			dist.ElectChan <- sh
		}
	}

}

// переделать
func (d *Dock) unloadShip(wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Started")
	for k := range d.ships {
		if k.SType == d.shipType {
			//fmt.Println("ожидаем разгрузки, ", d.shipType)
			time.Sleep(time.Second * 1)

			fmt.Println("Выгружен корабль номер: ", k.id)
		} else {
			fmt.Println("wrong type ", k.id)
		}
	}
}

// Генератор -> Распределитель(принимает, отдаёт) -> ДОКИ

func main() {
	tunnel := newTunnel()

	fmt.Println("размер", len(tunnel.ships))
	gen := newShipGen(tunnel)
	dist := newDistributor(tunnel)
	d1 := newDock(dist.ProvChan, PROVISION)
	d2 := newDock(dist.FuelChan, FUEL)
	d3 := newDock(dist.ElectChan, ELECTRONICS)

	//==============

	wg := &sync.WaitGroup{}

	wg.Add(5)

	go gen.startShipGen(wg, tunnel.ships)
	go dist.distributor(wg)

	go d1.unloadShip(wg)
	go d2.unloadShip(wg)
	go d3.unloadShip(wg)

	wg.Wait()

	//==============

}
