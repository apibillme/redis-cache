package main

import (
	"testing"
	"time"

	"github.com/apibillme/restly"
	"github.com/apibillme/stubby"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSpec(t *testing.T) {

	Convey("Test Proxy", t, func() {

		go server(128, time.Duration(200), "", "8000")

		Convey("Success - should get key value from cache when found", func() {
			stub1 := stubby.StubFunc(&redisNewClient, nil)
			defer stub1.Reset()
			stub2 := stubby.StubFunc(&cachedGet, "555", true)
			defer stub2.Reset()
			req := restly.New()
			res, statusCode, err := restly.GetJSON(req, "http://localhost:8000/", `?key=123`)
			So(err, ShouldBeNil)
			So(statusCode, ShouldResemble, 200)
			value := res.Get("value").String()
			So(value, ShouldResemble, "555")
		})

		Convey("Success - should save key value into cache when not found in cache", func() {
			stub1 := stubby.StubFunc(&redisNewClient, nil)
			defer stub1.Reset()
			stub2 := stubby.StubFunc(&cachedGet, "", false)
			defer stub2.Reset()
			stub3 := stubby.StubFunc(&redisGet, "555")
			defer stub3.Reset()
			stub4 := stubby.StubFunc(&cachedSet, true)
			defer stub4.Reset()
			req := restly.New()
			res, statusCode, err := restly.GetJSON(req, "http://localhost:8000/", `?key=123`)
			So(err, ShouldBeNil)
			So(statusCode, ShouldResemble, 200)
			value := res.Get("value").String()
			So(value, ShouldResemble, "555")
		})

		Convey("Failure - cannot save into cache when not found in cache", func() {
			stub1 := stubby.StubFunc(&redisNewClient, nil)
			defer stub1.Reset()
			stub2 := stubby.StubFunc(&cachedGet, "", false)
			defer stub2.Reset()
			stub3 := stubby.StubFunc(&redisGet, "555")
			defer stub3.Reset()
			stub4 := stubby.StubFunc(&cachedSet, false)
			defer stub4.Reset()
			req := restly.New()
			res, statusCode, err := restly.GetJSON(req, "http://localhost:8000/", `?key=123`)
			So(err, ShouldBeNil)
			So(statusCode, ShouldResemble, 500)
			error := res.Get("error").String()
			So(error, ShouldResemble, "cannot set new value in cache")
		})
	})
}
