package store_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"config_server/config"
	"config_server/store"
)

var _ = Describe("DatabaseFactory", func() {

	Describe("Given database config", func() {

		It("returns correct connection string for mysql", func() {

			dbConfig := config.DBConfig {
				Adapter: "mysql",
				User: "bosh",
				Password: "bosh-password",
				Host: "host",
				Port: 0,
				Name: "dbconfig",
			}

			connectionString, err := store.ConnectionString(dbConfig)
			Expect(err).To(BeNil())
			Expect(connectionString).To(Equal("bosh:bosh-password@tcp(host:0)/dbconfig"))
		})

		It("returns correct connection string for postgres", func() {
			dbConfig := config.DBConfig {
				Adapter: "postgres",
				User: "bosh",
				Password: "bosh-password",
				Host: "host",
				Port: 0,
				Name: "dbconfig",
			}

			connectionString, err := store.ConnectionString(dbConfig)
			Expect(err).To(BeNil())
			Expect(connectionString).To(Equal("user=bosh password=bosh-password dbname=dbconfig sslmode=disable"))
		})

		It("returns error for unsupported adapater", func() {
			dbConfig := config.DBConfig {
				Adapter: "mongo",
				User: "bosh",
				Password: "bosh-password",
				Host: "host",
				Port: 0,
				Name: "dbconfig",
			}

			_, err := store.ConnectionString(dbConfig)
			Expect(err).ToNot(BeNil())
		})
	})
})
