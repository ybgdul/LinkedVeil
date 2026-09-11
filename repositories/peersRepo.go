package repositories

import (
	"fmt"
	"linkedveil/internal/peer"
	"linkedveil/models"

	"gorm.io/gorm"
)

type PeersRepository interface{ 
	CreatePeer(*peer.Peer) error 
	GetPeer(peer.ID) (*models.PeersModel, error)
	DeletePeer(peer.ID) error 
}

func NewPeersRepo(db *gorm.DB) PeersRepository{ 
	return &peerRepo{
		db: db,
	}
}

type peerRepo struct{ 
	db *gorm.DB
}

func (p *peerRepo) CreatePeer(peer *peer.Peer) error { 
	var m models.PeersModel = models.PeersModel{
		ID: string(peer.ID),
		VirtualIP: peer.VirtualIP.String(),
	}

	return p.db.Create(m).Error
}

func (p *peerRepo) GetPeer(id peer.ID) (*models.PeersModel, error) { 
	s := string(id)
	var m models.PeersModel
	if err := p.db.Find(&m, s).Error; err != nil { 
		return nil, fmt.Errorf("can't find the peer: %w", err)
	}
	return &m, nil
}

func (p *peerRepo) DeletePeer(id peer.ID) error { 
	s := string(id)
	var m models.PeersModel
	if err := p.db.Delete(&m, s).Error; err != nil { 
		return fmt.Errorf("can't delete the peer: %w", err)
	}
	return nil
}

