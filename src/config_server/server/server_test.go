package server_test

import (
	. "config_server/server"
	"net/http"
	"net/http/httptest"
	"strings"
	"config_server/server"
	"config_server/store"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {

	Describe("Given a server with no store", func() {

		var configServer ConfigServer

		BeforeEach(func() {
			configServer = server.NewServer(nil)
		})

		Context("starting the server", func() {
			It("should return an error", func() {
				err := configServer.Start(9000)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("DataStore can not be nil"))
			})
		})
	})

	Describe("Given a server with store", func() {

		var configServer ConfigServer

		BeforeEach(func() {
			configServer = server.NewServer(store.NewMemoryStore())
		})

		It("should return 200 OK for PUT", func() {
			req, _ := http.NewRequest("PUT", "/v1/config/bla", strings.NewReader("value=blabla"))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded;")
			recorder := httptest.NewRecorder()
			configServer.HandleRequest(recorder, req)

			Expect(recorder.Code).To(Equal(http.StatusOK))
		})

		It("should return 404 when path is not found for PUT", func() {
			req, _ := http.NewRequest("PUT", "/v1/bla", strings.NewReader("value=blabla"))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded;")
			recorder := httptest.NewRecorder()
			configServer.HandleRequest(recorder, req)

			Expect(recorder.Code).To(Equal(http.StatusNotFound))
		})

		It("should return 404 when path is not found for GET", func() {
			req, _ := http.NewRequest("GET", "/v1/bla", nil)
			recorder := httptest.NewRecorder()
			configServer.HandleRequest(recorder, req)

			Expect(recorder.Code).To(Equal(http.StatusNotFound))
		})

		It("should return 200 OK for PUT when valid key is retrieved", func() {
			putReq, _ := http.NewRequest("PUT", "/v1/config/bla", strings.NewReader("value=blabla"))
			putReq.Header.Set("Content-Type", "application/x-www-form-urlencoded;")
			putRecorder := httptest.NewRecorder()
			configServer.HandleRequest(putRecorder, putReq)

			getReq, _ := http.NewRequest("GET", "/v1/config/bla", nil)
			getRecorder := httptest.NewRecorder()
			configServer.HandleRequest(getRecorder, getReq)

			Expect(getRecorder.Code).To(Equal(http.StatusOK))
			Expect(getRecorder.Body.String()).To(Equal("{\"path\":\"bla\",\"value\":\"blabla\"}"))
		})
	})
})
