package export

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/pojiang20/musical-funicular/v1/util"
	"io"
	"time"
)

type SourceReader interface {
	io.Closer

	ReadMap(data *map[string]interface{}) error
	Next() (map[string]interface{}, error)
	Count() (int, error)
}

type MongoReader struct {
	query  bson.M
	fields []string
	count  int

	dbName   string
	collName string
	sess     *mgo.Session
	coll     *mgo.Collection
	iter     *mgo.Iter
}

func NewMongoReader(config *util.MongodbInputConfig, query bson.M, fields []string) (SourceReader, error) {
	tmp := &MongoReader{
		query:    query,
		fields:   fields,
		dbName:   config.DB,
		collName: config.Coll,
	}

	var err error
	tmp.sess, err = mgo.Dial(config.Host)
	if err != nil {
		util.Zlog.Errorf("mongo session init error:%v", err)
		return nil, err
	}
	tmp.coll = tmp.sess.DB(config.DB).C(config.Coll)
	tmp.iter = tmp.coll.Find(query).Iter()
	return tmp, nil
}

func (m *MongoReader) Count() (int, error) {
	if m.count != 0 {
		return m.count, nil
	}

	var err error
	m.count, err = m.coll.Find(m.query).Count()

	if err != nil {
		return 0, err
	}
	return m.count, nil
}

func (m *MongoReader) Close() error {
	if m.sess != nil {
		m.sess.Close()
	}
	return nil
}

func (m *MongoReader) ReadMap(data *map[string]interface{}) (err error) {
	if m.iter == nil {
		return fmt.Errorf("iter is nil")
	}
	defer func() {
		//存在mongo断连，导致panic的情况
		if rerr := recover(); rerr != nil {
			util.Zlog.Errorf("mongo read panic recover error:%v", rerr)
			err = util.ErrInterrupted
		}
	}()

	retry := 0
	for {
		if m.iter.Next(data) {
			return nil
		} else {
			err = m.iter.Err()
			if err == nil {
				return io.EOF
			}

			if retry > 3 {
				util.Zlog.Errorf("retry greater than 3,sleep 1s")
				time.Sleep(time.Second)
			}
			m.reInitIter()
			retry++
		}
	}
}

func (m *MongoReader) Next() (map[string]interface{}, error) {
	var data map[string]interface{}
	if err := m.ReadMap(&data); err != nil {
		return nil, err
	}
	return data, nil
}

// 处理迭代器失效情况
func (m *MongoReader) reInitIter() {
	err := m.iter.Close()
	if err != nil {
		util.Zlog.Errorf("mongo close error: %v", err)
		m.sess.Refresh()
	}
	m.iter = nil

	m.iter = m.coll.Find(m.query).Iter()
}
