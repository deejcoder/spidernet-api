package storage

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type ServerManager struct {
	Db *gorm.DB
}

type Server struct {
	gorm.Model
	Addr      string `gorm:"unique; not null"`
	Nick      string `gorm:"unique; not null"`
	VotesUp   int    `gorm:"default:0"`
	VotesDown int    `gorm:"default:0"`
	Tags      []*Tag `gorm:"many2many:server_tags;"`
}

type Tag struct {
	gorm.Model
	Value   string    `gorm:"unique;not null"`
	Servers []*Server `gorm:"many2many:server_tags;"`
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

func (mgr ServerManager) DeleteServer(id uint) {
	mgr.Db.Where("ID = ?", id).Delete(&Server{})
}

func (mgr ServerManager) UpdateServer(server *Server) {
	mgr.Db.Update(server)
}

func (mgr ServerManager) GetServerByAddr(addr string) Server {
	server := Server{}
	mgr.Db.Where("addr = ?", addr).First(&server)
	return server
}

func (mgr ServerManager) SearchServers(term string, offset int, limit int) []Server {
	tags := []Tag{}
	servers := []Server{}

	// we want to search for all tags which are similar to the term, and then find the assoicated servers
	mgr.Db.Model(&Tag{}).Where("similarity(tags.value, ?) > 0.2", term).Offset(offset).Limit(limit).Find(&tags)
	mgr.Db.Model(&tags).Related(&servers, "Servers")
	return servers
}

func (s Server) String() string {
	return fmt.Sprintf("Server<(%d) %s --> %s>", s.ID, s.Addr, s.Nick)
}
