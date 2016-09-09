package store_test

import (
	. "config_server/store"

	"config_server/store/fakes"
	"database/sql"
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("StoreMysql", func() {

	var (
		fakeDbProvider *fakes.FakeDbProvider
		fakeDb         *fakes.FakeIDb
		fakeRow        *fakes.FakeIRow
		fakeResult     *fakes.FakeResult

		store Store
	)

	BeforeEach(func() {
		fakeDbProvider = &fakes.FakeDbProvider{}
		fakeDb = &fakes.FakeIDb{}
		fakeRow = &fakes.FakeIRow{}
		fakeResult = &fakes.FakeResult{}

		store = NewMysqlStore(fakeDbProvider)
	})

	Describe("Get", func() {

		It("closes db connection on exit", func() {
			fakeDb.QueryRowReturns(&fakes.FakeIRow{})
			fakeDbProvider.DbReturns(fakeDb, nil)

			store.Get("Luke")
			Expect(fakeDb.CloseCallCount()).To(Equal(1))
		})

		It("returns value from db query", func() {
			fakeRow.ScanStub = func(dest ...interface{}) error {
				valuePtr, ok := dest[0].(*string)
				Expect(ok).To(BeTrue())
				*valuePtr = "Skywalker"
				return nil
			}

			fakeDb.QueryRowReturns(fakeRow)
			fakeDbProvider.DbReturns(fakeDb, nil)

			value, err := store.Get("Luke")
			Expect(err).To(BeNil())
			Expect(value).To(Equal("Skywalker"))
		})

		It("returns empty string when no result is found", func() {
			fakeRow.ScanReturns(sql.ErrNoRows)

			fakeDb.QueryRowReturns(fakeRow)
			fakeDbProvider.DbReturns(fakeDb, nil)

			value, err := store.Get("luke")
			Expect(err).To(BeNil())
			Expect(value).To(Equal(""))
		})

		It("returns an error when db provider fails to return db", func() {
			dbError := errors.New("connection failure")
			fakeDbProvider.DbReturns(nil, dbError)

			_, err := store.Get("luke")
			Expect(err).ToNot(BeNil())
			Expect(err).To(Equal(dbError))
		})

		It("returns an error when db query fails", func() {
			scanError := errors.New("query failure")
			fakeRow.ScanReturns(scanError)

			fakeDb.QueryRowReturns(fakeRow)
			fakeDbProvider.DbReturns(fakeDb, nil)

			_, err := store.Get("luke")
			Expect(err).ToNot(BeNil())
			Expect(err).To(Equal(scanError))
		})
	})

	Describe("Put", func() {

		It("closes db connection on exit", func() {
			fakeDbProvider.DbReturns(fakeDb, nil)

			store.Put("Luke", "Skywalker")
			Expect(fakeDb.CloseCallCount()).To(Equal(1))
		})

		It("does an insert when key does not exist in database", func() {
			fakeDbProvider.DbReturns(fakeDb, nil)

            err := store.Put("Luke", "Skywalker")
            Expect(err).To(BeNil())

            Expect(fakeDb.ExecCallCount()).To(Equal(1))

            query, values := fakeDb.ExecArgsForCall(0)
            Expect(query).To(Equal("INSERT INTO config VALUES(?,?)"))

            Expect(values[0]).To(Equal("Luke"))
            Expect(values[1]).To(Equal("Skywalker"))
		})

		It("does an update when key exists in database", func() {
			fakeDb.ExecReturns(nil, errors.New("duplicate"))
			fakeDbProvider.DbReturns(fakeDb, nil)

			store.Put("Luke", "Skywalker")

			Expect(fakeDb.ExecCallCount()).To(Equal(2))

			query, values := fakeDb.ExecArgsForCall(0)
            Expect(query).To(Equal("INSERT INTO config VALUES(?,?)"))
            Expect(values[0]).To(Equal("Luke"))
            Expect(values[1]).To(Equal("Skywalker"))

			query, values = fakeDb.ExecArgsForCall(1)
            Expect(query).To(Equal("UPDATE config SET config.config_value = ? WHERE config.config_key = ?"))
            Expect(values[0]).To(Equal("Skywalker"))
            Expect(values[1]).To(Equal("Luke"))
		})
	})

	Describe("Delete", func() {

		It("closes db connection on exit", func() {
			fakeDbProvider.DbReturns(fakeDb, nil)
			store.Delete("Luke")
			Expect(fakeDb.CloseCallCount()).To(Equal(1))
		})

		Context("Key exists", func() {

			BeforeEach(func() {
				fakeDbProvider.DbReturns(fakeDb, nil)
				fakeDb.ExecReturns(fakeResult, nil)

				fakeResult.RowsAffectedReturns(1, nil)
			})

			It("removes value", func() {
				store.Delete("Luke")

                Expect(fakeDb.ExecCallCount()).To(Equal(1))
                query, value := fakeDb.ExecArgsForCall(0)
                Expect(query).To(Equal("DELETE FROM config WHERE config_key = ?"))
                Expect(value[0]).To(Equal("Luke"))
			})

            It("returns true", func() {
                deleted, err := store.Delete("Luke")

                Expect(deleted).To(BeTrue())
                Expect(err).To(BeNil())
            })
		})

		Context("Key does not exist", func() {

			BeforeEach(func() {
				fakeDbProvider.DbReturns(fakeDb, nil)
				fakeDb.ExecReturns(fakeResult, nil)

				fakeResult.RowsAffectedReturns(0, nil)
			})

			It("returns false", func() {
				deleted, err := store.Delete("key")
				Expect(deleted).To(BeFalse())
				Expect(err).To(BeNil())
			})
		})
	})
})
