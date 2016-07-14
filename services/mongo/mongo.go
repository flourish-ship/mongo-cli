// MongoDB Session 管理器
// @Author Heroic Yang

package mongo

import (
	"fmt"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
)

var (
	singleton mongoManager
)

type (
	mongoSession struct {
		mongoDialInfo *mgo.DialInfo
		mongoSession  *mgo.Session
	}

	mongoManager struct {
		sessions map[string]mongoSession
	}
)

// 连接 MongoDB，并初始化名为 sessionName 的 Session
func Startup(sessionName, hosts, dbName, username, password string) error {
	if singleton.sessions == nil {
		singleton = mongoManager{
			sessions: make(map[string]mongoSession),
		}
	}

	if _, ok := singleton.sessions[sessionName]; ok {
		return nil
	}

	Hosts := strings.Split(hosts, ",")
	return CreateSession(sessionName, Hosts, dbName, username, password)
}

// 断开连接并关闭所有的 Session
func Shutdown() {
	for _, session := range singleton.sessions {
		CloseSession(session.mongoSession)
	}
}

// 创建 Session
func CreateSession(sessionName string, hosts []string, dbName, username, password string) error {
	mongoSession := mongoSession{
		mongoDialInfo: &mgo.DialInfo{
			Addrs:    hosts,
			Timeout:  60 * time.Second,
			Database: dbName,
			Username: username,
			Password: password,
		},
	}

	var err error
	mongoSession.mongoSession, err = mgo.DialWithInfo(mongoSession.mongoDialInfo)
	if err != nil {
		return err
	}

	mongoSession.mongoSession.SetSafe(&mgo.Safe{})
	singleton.sessions[sessionName] = mongoSession

	return nil
}

// 复制 Session
func CopySession(sessionName string) (*mgo.Session, error) {
	session := singleton.sessions[sessionName]

	if session.mongoSession == nil {
		return nil, fmt.Errorf("Unable To Locate Session %s", sessionName)
	}

	return session.mongoSession.Copy(), nil
}

// 关闭 Session
func CloseSession(mongoSession *mgo.Session) {
	mongoSession.Close()
}

// 根据 sessionName 获取 Session 对应的数据库
func GetSessionDatabase(sessionName string) (string, error) {
	session := singleton.sessions[sessionName]

	if session.mongoSession == nil {
		return "", fmt.Errorf("Unable To Locate Session %s", sessionName)
	}

	return session.mongoDialInfo.Database, nil
}

func WithCollection(databaseName, collection, sessionName string, fn func(*mgo.Collection) error) error {
	session, err := CopySession(sessionName)
	if err != nil {
		return err
	}
	defer session.Close()
	c := session.DB(databaseName).C(collection)
	return fn(c)
}
