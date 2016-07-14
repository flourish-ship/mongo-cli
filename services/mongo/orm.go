package mongo

import (
	"errors"
	"gopkg.in/mgo.v2"
	"reflect"
)

var (
	ErrMultiRows    = errors.New("<QuerySeter> return multi rows")
	ErrNoRows       = errors.New("<QuerySeter> no row found")
	ErrArgs         = errors.New("<Ormer> args error may be empty")
	ErrNotImplement = errors.New("have not implement")
)

type Ormer interface {
	Insert(IModel) error
	InsertMulti(interface{}) error
	Read(IModel) error
	Delete(IModel) error
	Find(IModel) QuerySeter
	Update(IModel, map[string]interface{}) error
	UpdateWithCommand(IModel, map[string]interface{}, string) error
}

type orm struct {
}

// create new orm
func NewOrm() Ormer {
	o := new(orm)
	return o
}

func (o *orm) Insert(md IModel) error {

	d, c, s := md.GetMgoInfo()
	fn := func(c *mgo.Collection) error {
		return c.Insert(md)
	}
	return WithCollection(d, c, s, fn)

}
func (o *orm) InsertMulti(mds interface{}) error {

	if reflect.TypeOf(mds).Kind() != reflect.Slice {
		return ErrArgs
	}
	arrValue := reflect.ValueOf(mds)

	models := make([]interface{}, arrValue.Len())
	cnt := arrValue.Len()
	for k := 0; k < cnt; k++ {
		models[k] = arrValue.Index(k).Interface()
	}

	md := models[0]
	model := md.(IModel)
	d, c, s := model.GetMgoInfo()

	fn := func(c *mgo.Collection) error {
		if cnt == 1 {
			return c.Insert(model)
		} else {
			return c.Insert(models...)
		}

	}
	return WithCollection(d, c, s, fn)
}

func (o *orm) Update(model IModel, attributes map[string]interface{}) error {
	d, c, s := model.GetMgoInfo()

	fn := func(c *mgo.Collection) error {
		if nil == attributes {
			return c.UpdateId(model.GetId(), model)
		} else {
			return c.UpdateId(model.GetId(), Map{"$set": attributes})
		}
	}

	return WithCollection(d, c, s, fn)

}

func (o *orm) UpdateWithCommand(model IModel, attributes map[string]interface{}, command string) error {
	d, c, s := model.GetMgoInfo()

	fn := func(c *mgo.Collection) error {
		return c.UpdateId(model.GetId(), Map{command: attributes})
	}

	return WithCollection(d, c, s, fn)
}

// delete model in database
func (o *orm) Delete(model IModel) error {
	d, c, s := model.GetMgoInfo()

	fn := func(c *mgo.Collection) error {
		return c.RemoveId(model.GetId())
	}

	return WithCollection(d, c, s, fn)
}

func (o *orm) Read(model IModel) error {
	var md IModel
	d, c, s := model.GetMgoInfo()
	fn := func(c *mgo.Collection) error {
		return c.FindId(model.GetId()).One(md)
	}

	err := WithCollection(d, c, s, fn)
	if err == nil {
		model = md
	}

	return err

}

func (o *orm) Find(model IModel) (qs QuerySeter) {

	d, c, s := model.GetMgoInfo()
	session, err := CopySession(s)
	if err != nil {
		return nil
	}

	collection := session.DB(d).C(c)

	query := &querySet{}
	query.collection = collection
	query.session = session
	query.cond = &Condition{}
	return query
}
