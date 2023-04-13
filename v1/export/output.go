package export

import (
	"github.com/globalsign/mgo"
	"github.com/pojiang20/musical-funicular/v1/util"
	"io"
)

type SinkWriter interface {
	io.Closer

	WriteMap(data map[string]interface{}) error
	Flush() error
}

type MongoWriter struct {
	config   *util.MongodbInputConfig
	session  *mgo.Session
	collName string
}

func NewMongoWriter(config *util.MongodbInputConfig) (SinkWriter, error) {

}

func (m MongoWriter) Close() error {
	//TODO implement me
	panic("implement me")
}

func (m MongoWriter) WriteMap(data map[string]interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (m MongoWriter) Flush() error {
	//TODO implement me
	panic("implement me")
}
