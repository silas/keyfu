package main

import (
	"labix.org/v2/mgo/bson"
)

type Catalog struct {
	Id    bson.ObjectId "_id,omitempty"
	Name  string
	Root  string
	Tags  string
	Value string
}
