package store_test

import (
	. "config_server/store"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Utils", func() {

	Describe("SQLReplacer", func() {

		It("returns query string with bind parameter ? for mysql", func() {
			adapter := "mysql"
			queryString := "INSERT INTO config VALUES(?, ?)"

			mySqlQueryString := SQLReplacer(adapter, queryString)
			Expect(mySqlQueryString).Should(Equal("INSERT INTO config VALUES(?, ?)"))
		})

 		It("returns query string with bind parameter ${n} for postgres", func() {
			adapter := "postgres"
			queryString := "INSERT INTO config VALUES(?, ?)"

			mySqlQueryString := SQLReplacer(adapter, queryString)
			Expect(mySqlQueryString).Should(Equal("INSERT INTO config VALUES($1, $2)"))
		})
	})
})
