package mongo

func NewCondition() *Condition {
	condition := new(Condition)
	condition.conditions = Map{}
	return condition
}

type Condition struct {
	conditions   Map
	isRawQuery   bool //原生查询
	rawCondition Map  //原生查询语句
}

func (self *Condition) And(field string, opt Operator, val interface{}) *Condition {
	cond := parseCondition(opt, val)
	return self.logicQuery(And, field, cond)
}

func (self *Condition) Or(field string, opt Operator, val interface{}) *Condition {
	cond := parseCondition(opt, val)
	return self.logicQuery(Or, field, cond)
}

func (self *Condition) Nor(field string, opt Operator, val interface{}) *Condition {
	cond := parseCondition(opt, val)
	return self.logicQuery(Nor, field, cond)
}

func (self *Condition) ElemMatch(field string, opt Operator, val interface{}) *Condition {
	cond := parseCondition(opt, val)
	return self.logicQuery(ElemMatch, field, cond)
}

func (self *Condition) logicQuery(opt Operator, field string, express Map) *Condition {

	if self.conditions == nil {
		self.conditions = Map{}
	}

	var operator string
	if opt == Nor {
		operator = "$nor"
	} else if opt == Or {
		operator = "$or"
	} else if opt == And {
		operator = "$and"
	} else if opt == ElemMatch {
		operator = "$elemMatch"
	}

	if v, ok := self.conditions[operator]; ok && v != nil {
		if lq, ok := v.([]Map); ok {
			lq = append(lq, Map{field: express})
			self.conditions[operator] = lq
		}
	} else {
		lq := make([]Map, 1)
		lq[0] = Map{field: express}
		self.conditions[operator] = lq
	}

	return self
}

func isLogicQuery(opt Operator) bool {
	if opt == Or || opt == And || opt == Nor || opt == ElemMatch {
		return true
	}

	return false
}

func parseCondition(opt Operator, value interface{}) Map {
	condMap := Map{}
	if opt == Equal {
		condMap = Map{"$eq": value}
	} else if opt == NotEqual {
		condMap = Map{"$ne": value}
	} else if opt == LT {
		condMap = Map{"$lt": value}
	} else if opt == LTE {
		condMap = Map{"$lte": value}
	} else if opt == GT {
		condMap = Map{"$gt": value}
	} else if opt == GTE {
		condMap = Map{"$gte": value}
	} else if opt == IN {
		condMap = Map{"$in": value}
	} else if opt == NotIn {
		condMap = Map{"$nin": value}
	} else if opt == Size {
		condMap = Map{"$size": value}
	} else if opt == All {
		condMap = Map{"$all": value}
	} else if opt == Where {
		condMap = Map{"$where": value}
	} else if opt == Type {
		condMap = Map{"$type": value}
	} else if opt == Exists {
		condMap = Map{"$exists": value}
	} else if opt == ElemMatch {
		condMap = Map{"$elemMatch": value}
	} else if opt == Like {
		condMap = Map{"$regex": value, "$options": "i"}
	} else if opt == Not {
		condMap = Map{"$not": value}
	} else if opt == GeoWithIn {
		condMap = Map{"$geoWithin": value}
	}

	return condMap

}

func (self *Condition) addCond(field string, opt Operator, value interface{}) *Condition {

	if self.conditions == nil {
		self.conditions = Map{}
	}

	parsedCond := parseCondition(opt, value)

	if isLogicQuery(opt) {
		return self
	}

	//check the field exist or not
	if clone, ok := self.conditions[field]; ok {
		delete(self.conditions, field)

		newCond := Map{}
		for k, v := range clone.(Map) {
			newCond[k] = v
		}

		for k, v := range parsedCond {
			newCond[k] = v
		}

		self.conditions[field] = newCond
		return self
	}

	self.conditions[field] = parsedCond
	return self
}

func (self *Condition) mergeCond(cond *Condition) *Condition {
	if self.conditions == nil {
		self.conditions = Map{}
	}
	for k, v := range cond.conditions {
		self.conditions[k] = v
	}

	return self
}

func (self *Condition) getConditions() Map {
	if self.isRawQuery {
		self.conditions = self.rawCondition
	} else if self.conditions == nil {
		self.conditions = Map{}
	}

	return self.conditions
}
