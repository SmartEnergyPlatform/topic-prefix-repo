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
	"net/http"
	"github.com/SmartEnergyPlatform/util/http/logger"
	"github.com/SmartEnergyPlatform/util/http/response"

	"github.com/julienschmidt/httprouter"
)

func StartRest() {
	log.Println("start server on port: ", Config.ServerPort)
	httpHandler := getRoutes()
	logger := logger.New(httpHandler, Config.LogLevel)
	log.Println(http.ListenAndServe(":"+Config.ServerPort, logger))
}

func getRoutes() (router *httprouter.Router) {
	router = httprouter.New()

	router.GET("/get/routes", func(res http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		routes, err := GetAllRoutes()
		if err == nil {
			response.To(res).Json(routes)
		} else {
			log.Println("error on GetAllRoutes(): ", err)
			response.To(res).DefaultError("serverside error", 500)
		}
	})

	router.GET("/get/routes/:topic/:device/:service", func(res http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		topic := ps.ByName("topic")
		device := ps.ByName("device")
		service := ps.ByName("service")
		routes, err := GetRoutes(topic, device, service)
		result := []string{}
		for _, route := range routes {
			result = append(result, route.Target)
		}
		if err == nil {
			response.To(res).Json(result)
		} else {
			log.Println("error on GetRoutes(): ", err)
			response.To(res).DefaultError("serverside error", 500)
		}
	})

	router.GET("/get/routes/:topic/:device", func(res http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		topic := ps.ByName("topic")
		device := ps.ByName("device")
		service := ""
		routes, err := GetRoutes(topic, device, service)
		result := []string{}
		for _, route := range routes {
			result = append(result, route.Target)
		}
		if err == nil {
			response.To(res).Json(result)
		} else {
			log.Println("error on GetRoutes(): ", err)
			response.To(res).DefaultError("serverside error", 500)
		}
	})

	router.GET("/get/routes/:topic", func(res http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		topic := ps.ByName("topic")
		log.Println("DEBUG: find routes without prefix")
		routes, err := GetRoutesWithEmptyPrefix(topic)
		result := []string{}
		for _, route := range routes {
			result = append(result, route.Target)
		}
		if err == nil {
			response.To(res).Json(result)
		} else {
			log.Println("error on GetRoutes(): ", err)
			response.To(res).DefaultError("serverside error", 500)
		}
	})

	//create new route
	router.POST("/add/route/:topic/:device/:service/:target", func(res http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		topic := ps.ByName("topic")
		device := ps.ByName("device")
		service := ps.ByName("service")
		target := ps.ByName("target")
		err := AddRoute(topic, device, service, target)
		if err == nil {
			response.To(res).Text("ok")
		} else {
			log.Println("error on AddRoute(): ", err)
			response.To(res).DefaultError("serverside error", 500)
		}
	})

	//remove route
	router.DELETE("/remove/route/:topic/:device/:service/:target", func(res http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		topic := ps.ByName("topic")
		device := ps.ByName("device")
		service := ps.ByName("service")
		target := ps.ByName("target")
		err := RemoveRoute(topic, device, service, target)
		if err == nil {
			response.To(res).Text("ok")
		} else {
			log.Println("error on RemoveRoute(): ", err)
			response.To(res).DefaultError("serverside error", 500)
		}
	})

	//remove all routes to target with prefix
	router.DELETE("/remove/prefix/:device/:service/:target", func(res http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		device := ps.ByName("device")
		service := ps.ByName("service")
		target := ps.ByName("target")
		err := RemovePrefix(device, service, target)
		if err == nil {
			response.To(res).Text("ok")
		} else {
			log.Println("error on RemovePrefix(): ", err)
			response.To(res).DefaultError("serverside error", 500)
		}
	})

	//remove all routes to target
	router.DELETE("/remove/target/:target", func(res http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		target := ps.ByName("target")
		err := RemoveTarget(target)
		if err == nil {
			response.To(res).Text("ok")
		} else {
			log.Println("error on RemoveTarget(): ", err)
			response.To(res).DefaultError("serverside error", 500)
		}
	})

	router.GET("/topics", func(res http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		topics, err := GetTopics()
		if err == nil {
			response.To(res).Json(topics)
		} else {
			log.Println("error on GetTopics(): ", err)
			response.To(res).DefaultError("serverside error", 500)
		}
	})

	return
}
