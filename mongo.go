/*
 * Copyright 2018 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"log"
	"sync"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var instance *mgo.Session
var once sync.Once

func getDb() *mgo.Session {
	once.Do(func() {
		session, err := mgo.Dial(Config.MongoUrl)
		if err != nil {
			log.Fatal("error on connection to mongodb: ", err)
		}
		session.SetMode(mgo.Monotonic, true)
		instance = session
	})
	return instance.Copy()
}

func getRoutesCollection() (session *mgo.Session, collection *mgo.Collection) {
	session = getDb()
	collection = session.DB(Config.MongoTable).C(Config.RoutesCollection)
	err := collection.EnsureIndexKey("topic")
	if err != nil {
		log.Fatal("error on db topic index: ", err)
	}
	err = collection.EnsureIndexKey("prefix")
	if err != nil {
		log.Fatal("error on db prefix index: ", err)
	}
	err = collection.EnsureIndexKey("target")
	if err != nil {
		log.Fatal("error on db target index: ", err)
	}
	return
}

type Route struct {
	Topic  string `bson:"topic,omitempty"`
	Target string `bson:"target,omitempty"`
	Device string  `json:"device,omitempty" bson:"device,omitempty"`
	Service	string `json:"service,omitempty" bson:"service,omitempty"`
}

func GetAllRoutes() (result []Route, err error) {
	session, collection := getRoutesCollection()
	defer session.Close()
	err = collection.Find(nil).All(&result)
	return
}

func GetRoutes(topic string, device string, service string) (result []Route, err error) {
	session, collection := getRoutesCollection()
	defer session.Close()
	err = collection.Find(bson.M{"topic": topic, "device": bson.M{"$in": []string{device, "*"}},"service": bson.M{"$in": []string{service, "*"}}}).All(&result)
	return
}

func GetRoutesWithEmptyPrefix(topic string) (result []Route, err error) {
	session, collection := getRoutesCollection()
	defer session.Close()
	err = collection.Find(bson.M{"topic": topic, "$where": "(!this.service || this.service.length == 0) && (!this.device || this.device.length == 0)"}).All(&result)
	return
}

func AddRoute(topic string, device string, service string, target string) (err error) {
	session, collection := getRoutesCollection()
	defer session.Close()
	err = collection.Insert(Route{Topic: topic, Device: device, Service: service, Target: target})
	return
}

func RemoveTarget(target string) (err error) {
	session, collection := getRoutesCollection()
	defer session.Close()
	_, err = collection.RemoveAll(Route{Target: target})
	return
}

func RemovePrefix(device string, service string, target string) (err error) {
	session, collection := getRoutesCollection()
	defer session.Close()
	_, err = collection.RemoveAll(Route{Device: device, Service: service, Target: target})
	return
}

func RemoveRoute(topic string, device string, service string, target string) (err error) {
	session, collection := getRoutesCollection()
	defer session.Close()
	_, err = collection.RemoveAll(Route{Topic: topic, Device: device, Service: service, Target: target})
	return
}

func GetTopics() (topics []string, err error) {
	session, collection := getRoutesCollection()
	defer session.Close()
	err = collection.Find(nil).Distinct("topic", &topics)
	return
}
