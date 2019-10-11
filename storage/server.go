package storage

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

type ServerManager struct {
	Db *gorm.DB
}

type Server struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Addr      string    `gorm:"unique; not null" json:"addr"`
	Nick      string    `gorm:"unique; not null" json:"nick"`
	VotesUp   int       `gorm:"default:0" json:"votes_up"`
	VotesDown int       `gorm:"default:0" json:"votes_down"`
	Tags      []*Tag    `gorm:"many2many:server_tags;" json:"tags"`
}

type Tag struct {
	gorm.Model `json:"-"`
	Value      string    `gorm:"unique;not null" json:"value"`
	Servers    []*Server `gorm:"many2many:server_tags;" json:"-"`
}

func NewServerManager(db *gorm.DB) *ServerManager {
	return &ServerManager{Db: db}
}

// CreateServer creates a new server if it doesn't exist, returns existing server or new server
func (mgr ServerManager) CreateServer(host string, nick string) (Server, error) {
	server := Server{
		Addr: host,
		Nick: nick,
	}

	mgr.Db.FirstOrCreate(&server, Server{Addr: host})
	if err := mgr.Db.Error; err != nil {
		return server, err
	}

	return server, nil
}

func (mgr ServerManager) DeleteServer(id uint) error {
	mgr.Db.Where("ID = ?", id).Delete(&Server{})
	if err := mgr.Db.Error; err != nil {
		return err
	}
	return nil
}

func (mgr ServerManager) UpdateServer(server *Server) error {
	mgr.Db.Update(server)
	if err := mgr.Db.Error; err != nil {
		return err
	}
	return nil
}

func (mgr ServerManager) GetServerByAddr(addr string) Server {
	server := Server{}
	mgr.Db.Where("addr = ?", addr).First(&server)
	return server
}

// SearchServers searches all tags, and retrurns *limit* servers, with an offset of *offset*
func (mgr ServerManager) SearchServers(term string, offset int, limit int) ([]Server, error) {
	tags := []Tag{}
	servers := []Server{}

	// search all tags that are similar to term
	mgr.Db.Model(&Tag{}).
		Where("similarity(tags.value, ?) > 0.2", term).
		Offset(offset).
		Limit(limit).
		Find(&tags)

	// find assoicated servers & load the Tags into Server.Tags field
	mgr.Db.Model(&tags).
		Preload("Tags").
		Related(&servers, "Servers")

	if err := mgr.Db.Error; err != nil {
		return nil, err
	}
	return servers, nil
}

func (mgr ServerManager) GetServers(offset int, limit int) []Server {
	servers := []Server{}
	mgr.Db.Model(&Server{}).Offset(offset).Limit(limit).Find(&servers)
	return servers
}

// AddServerTags creates non-existant server tags and adds existing or non-existing tags to some server
func (mgr ServerManager) AddServerTags(server Server, tagsS []string) error {
	var tags []Tag
	for _, tag := range tagsS {
		tag := Tag{
			Value: tag,
		}

		mgr.Db.FirstOrCreate(&tag, tag)
		tags = append(tags, tag)
	}

	db := mgr.Db.Model(&server).Association("Tags").Append(tags)
	if err := db.Error; err != nil {
		return err
	}
	return nil
}

func (s Server) String() string {
	return fmt.Sprintf("Server<(%d) addr: %s, nick: %s>", s.ID, s.Addr, s.Nick)
}

func (t Tag) String() string {
	return fmt.Sprintf("Tag<%s>", t.Value)
}
