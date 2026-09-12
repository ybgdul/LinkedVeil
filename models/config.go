package models


type PeersModel struct{ 
	ID string `json:"id" gorm:"primary_key"`
	VirtualIP string `json:"virtual_ip"`
}

type RoutesModel struct{ 
	Network string `json:"network" gorm:"primary_key"`
	PeerID string `json:"peer_id" gorm:"not null;index"`
}