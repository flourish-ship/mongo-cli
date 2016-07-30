package mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type querySet struct {
	session    *mgo.Session
	collection *mgo.Collection
	iter       *mgo.Iter
	query      *mgo.Query
	cond       *Condition
	sort       []string
	selector   map[string]bool
	limit      int
}

type Pagination struct {
	PageIndex    int64       `json:"pageIndex"`
	PageSize     int64       `json:"pageSize"`
	TotalRecords int64       `json:"totalRecords"`
	Records      interface{} `json:"records"`
}

var _ QuerySeter = new(querySet)

func (o *querySet) Filter(field string, opt Operator, value interface{}) QuerySeter {
	if o.cond == nil {
		o.cond = NewCondition()
	}

	o.cond.addCond(field, opt, value)

	return o
}

func (o *querySet) SetCond(newCond *Condition) QuerySeter {
	if o.cond == nil {
		o.cond = NewCondition()
	}

	o.cond.mergeCond(newCond)

	return o
}

func (o *querySet) Raw(doc Map) QuerySeter {
	if o.cond == nil {
		o.cond = NewCondition()
	}

	o.cond.rawCondition = doc
	o.cond.isRawQuery = true

	return o

}

func initQuery(o *querySet) {
	if o.query == nil {
		o.query = o.collection.Find(o.cond.getConditions())
	}

	if len(o.sort) > 0 {
		o.query.Sort(o.sort...)
	}

	if len(o.selector) > 0 {
		o.query.Select(o.selector)
	}

	if o.limit > 0 {
		o.query.Limit(o.limit)
	}

}

//需要显示的字段
func (o *querySet) Fields(fields ...string) QuerySeter {
	if 0 < len(fields) {
		selectors := make(map[string]bool, len(fields))
		for _, field := range fields {
			selectors[field] = true
		}

		o.selector = selectors
	}

	return o

}

func (o *querySet) Distinct(field string, results interface{}) error {
	defer o.session.Close()

	initQuery(o)
	return o.query.Distinct(field, results)
}

// return QuerySeter execution result number
func (o *querySet) Count() (int, error) {
	defer o.session.Close()

	initQuery(o)
	return o.query.Count()
}

// check result empty or not after QuerySeter executed
func (o *querySet) Exist() bool {
	defer o.session.Close()
	initQuery(o)

	cnt, _ := o.query.Count()
	if cnt > 0 {
		return true
	}
	return false
}

func (o *querySet) Update(docs map[string]interface{}) error {
	defer o.session.Close()
	initQuery(o)

	change := mgo.Change{
		Update:    Map{"$set": docs},
		ReturnNew: false,
	}

	_, err := o.query.Apply(change, Map{})
	return err
}

func (o *querySet) UpdateWithCommand(docs map[string]interface{}, command string) error {
	defer o.session.Close()
	initQuery(o)

	change := mgo.Change{
		Update: Map{command: docs},
		Upsert: true,
	}

	_, err := o.query.Apply(change, Map{})
	return err
}

//复杂更新(Update,UpdateWithCommand可由次函数代替)
//docs格式为: M{"$set": M{"slink": g.g/1"}, "$inc": M{"count": 1}}
func (o *querySet) Upsert(docs map[string]interface{}) error {
	defer o.session.Close()
	initQuery(o)

	_, err := o.collection.Upsert(o.cond.getConditions(), docs)
	return err
}

func (o *querySet) UpdateAllWithCommand(docs map[string]interface{}, command string) error {
	defer o.session.Close()

	change := bson.M{
		command: docs,
	}
	_, err := o.collection.UpdateAll(o.cond.getConditions(), change)
	return err
}

//和UpdateAllWithCommand重复,后续考虑移除 和UpdateAllWithCommand重复
func (o *querySet) UpdateAll(docs map[string]interface{}) error {
	defer o.session.Close()

	_, err := o.collection.UpdateAll(o.cond.getConditions(), docs)
	return err
}

func (o *querySet) DeleteAll() error {
	defer o.session.Close()
	initQuery(o)

	_, err := o.collection.RemoveAll(o.cond.getConditions())
	return err
}

func (o *querySet) Delete() error {
	defer o.session.Close()

	initQuery(o)

	change := mgo.Change{
		Remove: true,
	}

	_, err := o.query.Apply(change, Map{})
	return err

}

// query all data and map to containers.
// cols means the columns when querying.
func (o *querySet) All(models interface{}) error {
	defer o.session.Close()

	initQuery(o)

	if o.iter == nil {
		o.iter = o.query.Iter()
	}

	return o.iter.All(models)
}

func (o *querySet) Sort(sortFields map[string]int, keys []string) QuerySeter {
	sortStr := []string{}
	if keys != nil {
		for _, key := range keys {
			val := sortFields[key]
			if 0 < val {
				sortStr = append(sortStr, key)
			} else {
				sortStr = append(sortStr, "-"+key)
			}
		}
	} else {
		for key, val := range sortFields {
			if 0 < val {
				sortStr = append(sortStr, key)
			} else {
				sortStr = append(sortStr, "-"+key)
			}
		}
	}

	o.sort = sortStr

	return o
}

// query one row data and map to containers.
func (o *querySet) One(model interface{}) error {

	initQuery(o)

	defer o.session.Close()

	return o.query.One(model)

}

//limit restricts the maximum number of documents retrieved to n
func (o *querySet) Limit(n int) QuerySeter {
	o.limit = n
	return o
}

func (o *querySet) Pagination(page *Pagination) error {
	defer o.session.Close()

	initQuery(o)
	total, _ := o.query.Count()
	page.TotalRecords = int64(total)
	var models []*interface{}
	if page.Records == nil {
		page.Records = &models
	}
	return o.query.Skip(int((page.PageIndex - 1) * page.PageSize)).Limit(int(page.PageSize)).All(page.Records)

}
