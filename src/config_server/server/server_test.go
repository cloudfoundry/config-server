package server_test

import (
	. "config_server/server"
	"net/http"
	"net/http/httptest"
	"strings"
	"config_server/store"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"errors"
)

type BadMockStore struct {}

func (store BadMockStore) Get(key string) (string, error) {
	return "", errors.New("")
}

func (store BadMockStore) Put(key string, value string) (error) {
	return errors.New("")
}

var _ = Describe("Server", func() {

	Describe("Given a server with no store", func() {

		var configServer ConfigServer

		BeforeEach(func() {
			configServer = NewServer(nil)
		})

		Context("starting the server", func() {
			It("should return an error", func() {
				err := configServer.Start(9000, "/fake/cert/path", "/fake/key/path")
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("DataStore can not be nil"))
			})
		})
	})

	Describe("Given a server with store", func() {

		var configServer ConfigServer

		BeforeEach(func() {
			configServer = NewServer(store.NewMemoryStore())
		})

		Context("when URL path is invalid", func() {

			It("should return 404 Not Found for invalid paths", func() {
				invalidPaths := []string{"/v1/config/test/case", "/v1"}

				for _, path := range invalidPaths {
					req, _ := http.NewRequest("GET", path, nil)
					recorder := httptest.NewRecorder()
					configServer.HandleRequest(recorder, req)

					Expect(recorder.Code).To(Equal(http.StatusNotFound))
				}
			})

			It("should return 404 Not Found for other methods", func() {
				invalidMethods := [...]string{"DELETE", "POST", "PATCH"}
				http.NewRequest("PUT", "/v1/config/bla", strings.NewReader("value=blabla"))

				for _, method := range invalidMethods {
					req, _ := http.NewRequest(method, "/v1/config/bla", nil)
					recorder := httptest.NewRecorder()
					configServer.HandleRequest(recorder, req)

					Expect(recorder.Code).To(Equal(http.StatusNotFound))
				}
			})

			It("should return 404 Not Found when key is not provided for fetch", func() {
				req, _ := http.NewRequest("GET", "/v1/config/", nil)
				getRecorder := httptest.NewRecorder()
				configServer.HandleRequest(getRecorder, req)

				Expect(getRecorder.Code).To(Equal(http.StatusNotFound))
			})

			It("should return 404 Not Found when key is not provided for update", func() {
				req, _ := http.NewRequest("PUT", "/v1/config/", nil)
				getRecorder := httptest.NewRecorder()
				configServer.HandleRequest(getRecorder, req)

				Expect(getRecorder.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("when URL path is valid", func() {

			It("should return 200 OK when config key/value is updated", func() {
				req, _ := http.NewRequest("PUT", "/v1/config/bla", strings.NewReader("{\"value\":\"blabla\"}"))
				req.Header.Set("Content-Type", "application/json")
				recorder := httptest.NewRecorder()
				configServer.HandleRequest(recorder, req)

				Expect(recorder.Code).To(Equal(http.StatusOK))
			})

			It("should return 200 OK when valid key is retrieved", func() {
				putReq, _ := http.NewRequest("PUT", "/v1/config/bla", strings.NewReader("{\"value\":\"blabla\"}"))
				putReq.Header.Set("Content-Type", "application/json")
				putRecorder := httptest.NewRecorder()
				configServer.HandleRequest(putRecorder, putReq)

				getReq, _ := http.NewRequest("GET", "/v1/config/bla/", nil)
				getRecorder := httptest.NewRecorder()
				configServer.HandleRequest(getRecorder, getReq)

				Expect(getRecorder.Code).To(Equal(http.StatusOK))
				Expect(getRecorder.Body.String()).To(Equal("{\"path\":\"bla\",\"value\":\"blabla\"}"))
			})

			It("should return 404 Not Found when key is not found", func() {
				req, _ := http.NewRequest("GET", "/v1/config/test", nil)
				getRecorder := httptest.NewRecorder()
				configServer.HandleRequest(getRecorder, req)

				Expect(getRecorder.Code).To(Equal(http.StatusNotFound))
			})

			It("should return 400 Bad Request when value is not provided for update", func() {
				req, _ := http.NewRequest("PUT", "/v1/config/key", nil)
				getRecorder := httptest.NewRecorder()
				configServer.HandleRequest(getRecorder, req)

				Expect(getRecorder.Code).To(Equal(http.StatusBadRequest))
			})

			It("should return 500 Internal Server Error if an error occurs", func() {
				configServer = NewServer(BadMockStore{})

				req, _ := http.NewRequest("GET", "/v1/config/key", nil)
				getRecorder := httptest.NewRecorder()
				configServer.HandleRequest(getRecorder, req)

				Expect(getRecorder.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})
})
