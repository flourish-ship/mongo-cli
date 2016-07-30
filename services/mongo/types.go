package mongo

import (
	"gopkg.in/mgo.v2/bson"
)

type IModel interface {
	GetId() bson.ObjectId
	GetMgoInfo() (string, string, string) //database,collection,session
}

type Model struct {
}

type Map map[string]interface{}

type Operator string

const (
	Equal         Operator = "=="
	GT            Operator = ">"
	GTE           Operator = ">="
	LT            Operator = "<"
	LTE           Operator = "<="
	NotEqual      Operator = "!="
	IN            Operator = "in"
	NotIn         Operator = "nin"
	Or            Operator = "or"
	And           Operator = "and"
	Not           Operator = "not"
	Nor           Operator = "nor"
	Exists        Operator = "exists"
	Type          Operator = "type"
	Mod           Operator = "mod"
	Like          Operator = "regex"
	Text          Operator = "text"
	Where         Operator = "where"
	GeoWithIn     Operator = "geoWithin"
	GeoIntersects Operator = "geoIntersects"
	Near          Operator = "near"
	NearSphere    Operator = "nearSphere"
	All           Operator = "all"
	ElemMatch     Operator = "elemMatch"
	Size          Operator = "size"
)

type Expres struct {
	Key string
	Opt Operator
	Val []interface{}
}

const SortDesc int = -1
const SortAsc int = 1

type ICondition interface {
	getConditions() Map
	addCond(string, Operator, interface{}) *Condition
	mergeCond(*Condition) *Condition
}

// query seter
type QuerySeter interface {
	Filter(string, Operator, interface{}) QuerySeter
	SetCond(*Condition) QuerySeter
	Raw(Map) QuerySeter
	Limit(int) QuerySeter
	Pagination(*Pagination) error
	Sort(map[string]int, []string) QuerySeter
	Count() (int, error)
	Exist() bool
	Update(map[string]interface{}) error
	Upsert(map[string]interface{}) error
	UpdateWithCommand(map[string]interface{}, string) error
	UpdateAllWithCommand(map[string]interface{}, string) error
	UpdateAll(map[string]interface{}) error
	Delete() error
	DeleteAll() error
	Fields(...string) QuerySeter
	All(interface{}) error
	One(interface{}) error
	Distinct(string, interface{}) error
}
