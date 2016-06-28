package server_test

import (
	. "config_server/server"

	"config_server/store"
	"errors"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"strings"
)

type BadMockStore struct{}

func (store BadMockStore) Get(key string) (string, error) {
	return "", errors.New("")
}

func (store BadMockStore) Put(key string, value string) error {
	return errors.New("")
}

var _ = Describe("RequestHandlerConcrete", func() {
	Describe("Given a nil store", func() {

		Context("creating the requestHandler", func() {
			It("should return an error", func() {
				putReq, _ := http.NewRequest("PUT", "/v1/config/bla", strings.NewReader("{\"value\":\"blabla\"}"))
				putReq.Header.Set("Content-Type", "application/json")
				putRecorder := httptest.NewRecorder()

				requestHandler := NewConcreteRequestHandler(nil)
				requestHandler.HandleRequest(putRecorder, putReq)

				Expect(putRecorder.Code).To(Equal(http.StatusInternalServerError))
				Expect(putRecorder.Body.String()).To(Equal("DB Store is nil"))
			})
		})
	})

	Describe("Given a server with store", func() {

		var requestHandler RequestHandler

		BeforeEach(func() {
			requestHandler = NewConcreteRequestHandler(store.NewMemoryStore())
		})

		Context("when URL path is invalid", func() {

			It("should return 404 Not Found for invalid paths", func() {
				invalidPaths := []string{"/v1/config/test/case", "/v1"}

				for _, path := range invalidPaths {
					req, _ := http.NewRequest("GET", path, nil)
					recorder := httptest.NewRecorder()
					requestHandler.HandleRequest(recorder, req)

					Expect(recorder.Code).To(Equal(http.StatusNotFound))
				}
			})

			It("should return 404 Not Found for other methods", func() {
				invalidMethods := [...]string{"DELETE", "POST", "PATCH"}
				http.NewRequest("PUT", "/v1/config/bla", strings.NewReader("value=blabla"))

				for _, method := range invalidMethods {
					req, _ := http.NewRequest(method, "/v1/config/bla", nil)
					recorder := httptest.NewRecorder()
					requestHandler.HandleRequest(recorder, req)

					Expect(recorder.Code).To(Equal(http.StatusNotFound))
				}
			})

			It("should return 404 Not Found when key is not provided for fetch", func() {
				req, _ := http.NewRequest("GET", "/v1/config/", nil)
				getRecorder := httptest.NewRecorder()
				requestHandler.HandleRequest(getRecorder, req)

				Expect(getRecorder.Code).To(Equal(http.StatusNotFound))
			})

			It("should return 404 Not Found when key is not provided for update", func() {
				req, _ := http.NewRequest("PUT", "/v1/config/", nil)
				getRecorder := httptest.NewRecorder()
				requestHandler.HandleRequest(getRecorder, req)

				Expect(getRecorder.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("when URL path is valid", func() {

			It("should return 200 OK when config key/value is updated", func() {
				req, _ := http.NewRequest("PUT", "/v1/config/bla", strings.NewReader("{\"value\":\"blabla\"}"))
				req.Header.Set("Content-Type", "application/json")
				recorder := httptest.NewRecorder()
				requestHandler.HandleRequest(recorder, req)

				Expect(recorder.Code).To(Equal(http.StatusOK))
			})

			It("should return 200 OK when valid key is retrieved", func() {
				putReq, _ := http.NewRequest("PUT", "/v1/config/bla", strings.NewReader("{\"value\":\"blabla\"}"))
				putReq.Header.Set("Content-Type", "application/json")
				putRecorder := httptest.NewRecorder()
				requestHandler.HandleRequest(putRecorder, putReq)

				getReq, _ := http.NewRequest("GET", "/v1/config/bla/", nil)
				getRecorder := httptest.NewRecorder()
				requestHandler.HandleRequest(getRecorder, getReq)

				Expect(getRecorder.Code).To(Equal(http.StatusOK))
				Expect(getRecorder.Body.String()).To(Equal("{\"path\":\"bla\",\"value\":\"blabla\"}"))
			})

			It("should return 404 Not Found when key is not found", func() {
				req, _ := http.NewRequest("GET", "/v1/config/test", nil)
				getRecorder := httptest.NewRecorder()
				requestHandler.HandleRequest(getRecorder, req)

				Expect(getRecorder.Code).To(Equal(http.StatusNotFound))
			})

			It("should return 400 Bad Request when value is not provided for update", func() {
				req, _ := http.NewRequest("PUT", "/v1/config/key", nil)
				getRecorder := httptest.NewRecorder()
				requestHandler.HandleRequest(getRecorder, req)

				Expect(getRecorder.Code).To(Equal(http.StatusBadRequest))
			})

			It("should return 500 Internal Server Error if an error occurs", func() {
				requestHandler = NewConcreteRequestHandler(BadMockStore{})

				req, _ := http.NewRequest("GET", "/v1/config/key", nil)
				getRecorder := httptest.NewRecorder()
				requestHandler.HandleRequest(getRecorder, req)

				Expect(getRecorder.Code).To(Equal(http.StatusInternalServerError))
			})
		})

		It("should return 200 Status OK when an integer value is added", func() {
			req, _ := http.NewRequest("PUT", "/v1/config/bla", strings.NewReader("{\"value\":1}"))
			req.Header.Set("Content-Type", "application/json")
			putRecorder := httptest.NewRecorder()
			requestHandler.HandleRequest(putRecorder, req)

			Expect(putRecorder.Code).To(Equal(http.StatusOK))
		})

		It("should return 200 Status OK when a string value is added", func() {
			req, _ := http.NewRequest("PUT", "/v1/config/bla", strings.NewReader("{\"value\":\"str\"}"))
			req.Header.Set("Content-Type", "application/json")
			putRecorder := httptest.NewRecorder()
			requestHandler.HandleRequest(putRecorder, req)

			Expect(putRecorder.Code).To(Equal(http.StatusOK))
		})
	})
})
