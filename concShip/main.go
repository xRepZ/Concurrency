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

type Ship struct {
	id         int
	cap        int
	currentCap int
	SType      ShipType
}

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
	tunnel   *Tunnel
	shipType ShipType
}

func newDock(t *Tunnel, st ShipType) *Dock {
	return &Dock{
		tunnel:   t,
		shipType: st,
	}
}
// переделать
// func (d *Dock) unloadShip(wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	fmt.Println("Started")
// 	for k := range d.tunnel.ships {
// 		if k.SType == d.shipType {
// 			//fmt.Println("ожидаем разгрузки, ", d.shipType)
// 			time.Sleep(time.Second * 10)

// 			fmt.Println("Выгружен корабль номер: ", k.id)
// 		} else {
// 			fmt.Println("wrong type ", k.id)
// 		}
// 	}
// }
// Генератор -> Распределитель(принимает, отдаёт) -> ДОКИ


func main() {
	tunnel := newTunnel()

	fmt.Println("размер", len(tunnel.ships))
	gen := newShipGen(tunnel)
	d1 := newDock(tunnel, PROVISION)
	d2 := newDock(tunnel, FUEL)
	d3 := newDock(tunnel, ELECTRONICS)

	//==============

	wg := &sync.WaitGroup{}

	wg.Add(4)

	go gen.startShipGen(wg, tunnel.ships)

	go d1.unloadShip(wg)
	go d2.unloadShip(wg)
	go d3.unloadShip(wg)


	wg.Wait()
	

	//==============

}
