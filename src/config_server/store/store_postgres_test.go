package store_test

import (
	. "config_server/store"

	"config_server/store/fakes"
	"database/sql"
	"errors"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StorePostgres", func() {

	var (
		fakeDbProvider *fakes.FakeDbProvider
		mysqlStore     Store
	)

	BeforeEach(func() {
		fakeDbProvider = &fakes.FakeDbProvider{}
		mysqlStore = NewPostgresStore(fakeDbProvider)
	})

	Describe("Get", func() {

		It("closes db connection on exit", func() {
			fakeDb := &fakes.FakeIDb{}
			fakeDb.QueryRowReturns(&fakes.FakeIRow{})
			fakeDbProvider.DbReturns(fakeDb, nil)

			mysqlStore.Get("Luke")
			Expect(fakeDb.CloseCallCount()).To(Equal(1))
		})

		It("returns value from db query", func() {
			fakeRow := &fakes.FakeIRow{}
			fakeRow.ScanStub = func(dest ...interface{}) error {
				valuePtr, ok := dest[0].(*string)
				Expect(ok).To(BeTrue())
				*valuePtr = "Skywalker"
				return nil
			}

			fakeDb := &fakes.FakeIDb{}
			fakeDb.QueryRowReturns(fakeRow)
			fakeDbProvider.DbReturns(fakeDb, nil)

			value, err := mysqlStore.Get("Luke")
			Expect(err).To(BeNil())
			Expect(value).To(Equal("Skywalker"))
		})

		It("returns empty string when no result is found", func() {
			fakeRow := &fakes.FakeIRow{}
			fakeRow.ScanReturns(sql.ErrNoRows)

			fakeDb := &fakes.FakeIDb{}
			fakeDb.QueryRowReturns(fakeRow)
			fakeDbProvider.DbReturns(fakeDb, nil)

			value, err := mysqlStore.Get("luke")
			Expect(err).To(BeNil())
			Expect(value).To(Equal(""))
		})

		It("returns an error when db provider fails to return db", func() {
			dbError := errors.New("connection failure")
			fakeDbProvider.DbReturns(nil, dbError)

			_, err := mysqlStore.Get("luke")
			Expect(err).ToNot(BeNil())
			Expect(err).To(Equal(dbError))
		})

		It("returns an error when db query fails", func() {
			scanError := errors.New("query failure")
			fakeRow := &fakes.FakeIRow{}
			fakeRow.ScanReturns(scanError)

			fakeDb := &fakes.FakeIDb{}
			fakeDb.QueryRowReturns(fakeRow)
			fakeDbProvider.DbReturns(fakeDb, nil)

			_, err := mysqlStore.Get("luke")
			Expect(err).ToNot(BeNil())
			Expect(err).To(Equal(scanError))
		})
	})

	Describe("Put", func() {

		It("closes db connection on exit", func() {
			fakeDb := &fakes.FakeIDb{}
			fakeDbProvider.DbReturns(fakeDb, nil)

			mysqlStore.Put("Luke", "Skywalker")
			Expect(fakeDb.CloseCallCount()).To(Equal(1))
		})

		It("does an insert when key does not exist in database", func() {
			fakeDb := &fakes.FakeIDb{}
			fakeDbProvider.DbReturns(fakeDb, nil)

			err := mysqlStore.Put("Luke", "Skywalker")

			Expect(fakeDb.ExecCallCount()).To(Equal(1))

			query, _ := fakeDb.ExecArgsForCall(0)
			query = strings.ToUpper(query)
			Expect(strings.HasPrefix(query, "INSERT")).To(BeTrue())

			Expect(err).To(BeNil())
		})

		It("does an update when key exists in database", func() {
			fakeDb := &fakes.FakeIDb{}
			fakeDb.ExecReturns(nil, errors.New("duplicate"))
			fakeDbProvider.DbReturns(fakeDb, nil)

			mysqlStore.Put("Luke", "Skywalker")

			Expect(fakeDb.ExecCallCount()).To(Equal(2))

			query, _ := fakeDb.ExecArgsForCall(0)
			query = strings.ToUpper(query)
			Expect(strings.HasPrefix(query, "INSERT")).To(BeTrue())

			query, _ = fakeDb.ExecArgsForCall(1)
			query = strings.ToUpper(query)
			Expect(strings.HasPrefix(query, "UPDATE")).To(BeTrue())
		})
	})
})
