package producer

import (
	"fmt"
	"net/http"
	"sort"
	"time"
)

type SuperV interface {
	Get() string
}

type Handler func(http.Response, *http.Request) error

type Offer struct {
	ID       int
	CreateAt time.Time
}

type OfferList []Offer
func (o OfferList) Len() int {
	return len(o)
}
func (o OfferList) Swap(i, j int) {
	o[i], o[j] = o[j], o[i]
}
func (o OfferList) Less(i, j int) bool {
	return o[i].CreateAt.After(o[j].CreateAt)
}

type WannaV struct {
	name string
}

func (w WannaV) Get() string {
	oList := make([]Offer, 0)
	ol := OfferList(oList)
	sort.Sort(&ol)
	return w.name
}

func Hi(v SuperV) string {
	value := v.(*WannaV) // ext: WannaOne
	Say(value)
	return fmt.Sprintf("Hi, folks! this is I %d", toInt(value.name))
}

func Say(w *WannaV) {
	w.Get()
}

func Next(h Handler) {
}
