/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"errors"
	"os"
	"testing"
	"time"

	ecc "github.com/ernestio/ernest-config-client"
	"github.com/nats-io/nats"

	. "github.com/smartystreets/goconvey/convey"
)

func wait(ch chan bool) error {
	return waitTime(ch, 500*time.Millisecond)
}

func waitTime(ch chan bool, timeout time.Duration) error {
	select {
	case <-ch:
		return nil
	case <-time.After(timeout):
	}
	return errors.New("timeout")
}

func TestBasicRedirections(t *testing.T) {

	Convey("Given this service is fully set up", t, func() {
		n := ecc.NewConfig(os.Getenv("NATS_URI")).Nats()
		chfak := make(chan bool)
		cherr := make(chan bool)
		chvcl := make(chan bool)

		n.Subscribe("config.get.connectors", func(msg *nats.Msg) {
			n.Publish(msg.Reply, []byte(`{"instances":["fake","vcloud","vcloud-fake","aws","aws-fake"]}`))
		})

		setup()

		n.Subscribe("instance.create.fake", func(msg *nats.Msg) {
			chfak <- true
		})
		n.Subscribe("instance.create.error", func(msg *nats.Msg) {
			cherr <- true
		})
		n.Subscribe("instance.create.vcloud", func(msg *nats.Msg) {
			chvcl <- true
		})
		Convey("When it receives an invalid fake message", func() {
			n.Publish("instance.create", []byte(`{"service":"aaa"}`))
			Convey("Then it should redirect it to instance error creation", func() {
				So(wait(cherr), ShouldBeNil)
			})
		})
		Convey("When it receives a valid fake message", func() {
			n.Publish("instance.create", []byte(`{"service":"aaa","datacenter_type":"fake"}`))
			Convey("Then it should redirect it to a fake connector", func() {
				So(wait(chfak), ShouldBeNil)
			})
		})
		Convey("When it receives a valid vcloud message", func() {
			n.Publish("instance.create", []byte(`{"service":"aaa","datacenter_type":"vcloud"}`))
			Convey("Then it should redirect it to a fake connector", func() {
				So(wait(chvcl), ShouldBeNil)
			})
		})
	})
}
