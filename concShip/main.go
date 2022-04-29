package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

type ShipType int

const numOfShips = 60
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
	ships chan *Ship
}

func newTunnel() *Tunnel {

	shipsIn := make(chan *Ship, 5)
	return &Tunnel{
		ships: shipsIn,
	}
}

func (t *Tunnel) addShipToTunnel(sh *Ship) {
	//fmt.Println("положили корабль type.... ", sh.SType)
	t.ships <- sh

}

type Dock struct {
	SType   ShipType
	distrib *TunnelDistributor
}

func newDock(d *TunnelDistributor, st ShipType) *Dock {
	return &Dock{
		SType:   st,
		distrib: d,
	}
}

type TunnelDistributor struct {
	tunnel *Tunnel

	distribChan map[ShipType]chan *Ship
}

func newDistributor(t *Tunnel) *TunnelDistributor {
	fchan := make(chan *Ship)
	echan := make(chan *Ship)
	pchan := make(chan *Ship)
	dChan := make(map[ShipType]chan *Ship)
	dChan[FUEL] = fchan
	dChan[ELECTRONICS] = echan
	dChan[PROVISION] = pchan

	return &TunnelDistributor{
		tunnel:      t,
		distribChan: dChan,
	}
}

func (d *TunnelDistributor) getShips(s ShipType) chan *Ship {
	ch := d.distribChan[s]
	return ch
}

func (dist *TunnelDistributor) distributor(wg *sync.WaitGroup) {
	defer wg.Done()
	for k := range dist.distribChan {
		defer close(dist.distribChan[k])
	}
	for sh := range dist.tunnel.ships {
		dChan, ok := dist.distribChan[sh.SType]
		if !ok {
			log.Println("invalid type")
			continue
		}
		dChan <- sh

	}

}

// переделать
func (d *Dock) unloadShip(wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("Started")
	for k := range d.distrib.distribChan[d.SType] {

		//fmt.Println("ожидаем разгрузки, ", d.shipType)
		for k.currentCap > 0 {
			time.Sleep(time.Microsecond)
			k.currentCap -= 10
		}
		fmt.Println("Выгружен корабль номер: ", k.id)
	}
}

// Генератор -> Распределитель(принимает, отдаёт) -> ДОКИ

func main() {
	tunnel := newTunnel()

	fmt.Println("размер", len(tunnel.ships))
	gen := newShipGen(tunnel)
	dist := newDistributor(tunnel)
	d1 := newDock(dist, PROVISION)
	d2 := newDock(dist, FUEL)
	d3 := newDock(dist, ELECTRONICS)

	//==============
	// osSigCh := make(chan os.Signal, 1)
	// signal.Notify(osSigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// go func() {
	// 	for {
	// 		s := <-osSigCh
	// 		log.Println("Received exit signal!")
	// 		tunnel.ships.Close()
	// 	}
	// }()

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
