package models

import (

)


type PeersModel struct{ 
	ID string `json:"id" gorm:"primary_key"`
	VirtualIP string `json:"virtual_ip"`
}

type RoutesModel struct{ 
	ID string `json:"id" gorm:"primary_key"`
	Network string `json:"network"`
	PeerID string `json:"peer_id"`
}