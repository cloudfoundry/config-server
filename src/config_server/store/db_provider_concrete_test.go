package store_test

import (
	"config_server/config"
	. "config_server/store"
	fakes "config_server/store/storefakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DbProviderConcrete", func() {

	var fakeDb *fakes.FakeIDb
	var fakeSQL *fakes.FakeISql

	BeforeEach(func() {
		fakeDb = &fakes.FakeIDb{}
		fakeSQL = &fakes.FakeISql{}
		fakeSQL.OpenReturns(fakeDb, nil)
	})

	It("configures max open/idle connections", func() {

		dbConfig := config.DBConfig{
			Adapter:  "mysql",
			User:     "bosh",
			Password: "somethingsafe",
			Host:     "host",
			Port:     0,
			Name:     "dbconfig",
			ConnectionOptions: config.DBConnectionConfig{
				MaxOpenConnections: 12,
				MaxIdleConnections: 6,
			},
		}

		_, err := NewConcreteDbProvider(fakeSQL, dbConfig).Db()
		Expect(err).To(BeNil())
		Expect(fakeSQL.OpenCallCount()).To(Equal(1))

		Expect(fakeDb.SetMaxOpenConnsCallCount()).To(Equal(1))
		Expect(fakeDb.SetMaxOpenConnsArgsForCall(0)).To(Equal(12))

		Expect(fakeDb.SetMaxIdleConnsCallCount()).To(Equal(1))
		Expect(fakeDb.SetMaxIdleConnsArgsForCall(0)).To(Equal(6))
	})

	It("returns correct connection string for mysql", func() {

		dbConfig := config.DBConfig{
			Adapter:  "mysql",
			User:     "bosh",
			Password: "somethingsafe",
			Host:     "host",
			Port:     0,
			Name:     "dbconfig",
		}

		_, err := NewConcreteDbProvider(fakeSQL, dbConfig).Db()
		Expect(err).To(BeNil())
		Expect(fakeSQL.OpenCallCount()).To(Equal(1))

		driverName, dataSourceName, _ := fakeSQL.OpenArgsForCall(0)
		Expect(driverName).To(Equal(dbConfig.Adapter))
		Expect(dataSourceName).To(Equal("bosh:somethingsafe@tcp(host:0)/dbconfig"))
	})

	It("returns correct connection string for postgres", func() {

		dbConfig := config.DBConfig{
			Adapter:  "postgres",
			User:     "bosh",
			Password: "somethingsafe",
			Host:     "host",
			Port:     0,
			Name:     "dbconfig",
		}

		_, err := NewConcreteDbProvider(fakeSQL, dbConfig).Db()
		Expect(err).To(BeNil())
		Expect(fakeSQL.OpenCallCount()).To(Equal(1))

		driverName, dataSourceName, _ := fakeSQL.OpenArgsForCall(0)
		Expect(driverName).To(Equal(dbConfig.Adapter))
		Expect(dataSourceName).To(Equal("user=bosh password=somethingsafe dbname=dbconfig sslmode=disable"))
	})

	It("returns error for unsupported adapater", func() {

		dbConfig := config.DBConfig{
			Adapter:  "mongo",
			User:     "bosh",
			Password: "somethingsafe",
			Host:     "host",
			Port:     0,
			Name:     "dbconfig",
		}

		_, err := NewConcreteDbProvider(fakeSQL, dbConfig).Db()
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(Equal("Failed to generate DB connection string: Unsupported adapter: mongo"))
		Expect(fakeSQL.OpenCallCount()).To(Equal(0))
	})
})
