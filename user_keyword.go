package main

import (
	"labix.org/v2/mgo/bson"
)

type UserKeyword struct {
	Id   bson.ObjectId "_id,omitempty"
	Uid  string
	Key  string
	Type string
	Body string
}
