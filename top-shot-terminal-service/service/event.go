package service

import (
	"net/http"
)


type Event struct {
	id int
	playId int
	play string
	setId int
	setName int
	serialNumber int
	price float32
	uri string
	lowAsk float32
}

func hydrateMetadata(e *Event) {

}


