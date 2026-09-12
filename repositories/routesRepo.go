package repositories

import (
	"fmt"
	"linkedveil/internal/peer"
	"linkedveil/models"

	"gorm.io/gorm"
)


type RoutesRepository interface{ 
	CreateRoute(*peer.Route) error 
	GetRoute(peer.ID) (*models.RoutesModel, error)
	DeleteRoute(peer.ID) error 
}

func NewRoutesRepo(db *gorm.DB) RoutesRepository{ 
	return &routeRepo{
		db: db,
	}
}

type routeRepo struct{ 
	db *gorm.DB
}

func (p *routeRepo) CreateRoute(route *peer.Route) error { 
	var m models.RoutesModel = models.RoutesModel{
		PeerID: string(route.PeerID),
		Network: route.Network.String(),
	}

	return p.db.Create(m).Error
}

func (p *routeRepo) GetRoute(id peer.ID) (*models.RoutesModel, error) { 
	s := string(id)
	var m models.RoutesModel
	if err := p.db.Find(&m, s).Error; err != nil { 
		return nil, fmt.Errorf("can't find the peer: %w", err)
	}
	return &m, nil
}

func (p *routeRepo) DeleteRoute(id peer.ID) error { 
	s := string(id)
	var m models.PeersModel
	if err := p.db.Delete(&m, s).Error; err != nil { 
		return fmt.Errorf("can't delete the peer: %w", err)
	}
	return nil
}

