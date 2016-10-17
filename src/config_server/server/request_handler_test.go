package server_test

import (
	. "config_server/server"
	. "config_server/server/fakes"
	. "config_server/store/fakes"
	. "config_server/types/fakes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"config_server/config"
	"config_server/store"
	"config_server/types"
	"encoding/json"
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
				_, err := NewRequestHandler(nil, types.NewValueGeneratorConcrete(config.ServerConfig{}))
				Expect(err.Error()).To(Equal("Data store must be set"))
			})
		})
	})

	Describe("Given a server with store", func() {

		var requestHandler http.Handler
		var mockTokenValidator *FakeTokenValidator
		var mockStore *FakeStore
		var mockValueGeneratorFactory *FakeValueGeneratorFactory
		var mockValueGenerator *FakeValueGenerator

		BeforeEach(func() {
			mockTokenValidator = &FakeTokenValidator{}
			mockStore = &FakeStore{}
			mockValueGeneratorFactory = &FakeValueGeneratorFactory{}
			mockValueGenerator = &FakeValueGenerator{}
			requestHandler, _ = NewRequestHandler(mockStore, mockValueGeneratorFactory)
		})

		Context("when URL path is invalid", func() {
			It("should return 400 Bad Request", func() {
				invalidPaths := []string{"/v1", "/v1/", "/v1/data", "/v1/data/"}
				validMethods := []string{"GET", "PUT", "POST", "DELETE"}

				for _, method := range validMethods {
					for _, path := range invalidPaths {
						req, _ := http.NewRequest(method, path, nil)
						recorder := httptest.NewRecorder()
						requestHandler.ServeHTTP(recorder, req)

						Expect(recorder.Code).To(Equal(http.StatusBadRequest))
					}
				}
			})

			Context("when key path parameter is missing", func() {
				It("should return 400 Bad Request", func() {
					validMethods := []string{"GET", "PUT", "POST", "DELETE"}
					for _, method := range validMethods {
						req, _ := http.NewRequest(method, "/v1/data/", nil)
						getRecorder := httptest.NewRecorder()
						requestHandler.ServeHTTP(getRecorder, req)

						Expect(getRecorder.Code).To(Equal(http.StatusBadRequest))
					}
				})
			})

			Context("when key format is invalid", func() {
				It("should return 400 Bad Request", func() {
					inValidUrlPaths := []string{
						"/v1/data/key//path//",
						"/v1/data/key/{/*",
						"/v1/data/key/@?/",
					}

					validMethods := []string{"GET", "PUT", "POST", "DELETE"}

					for _, method := range validMethods {
						for _, path := range inValidUrlPaths {
							req, _ := http.NewRequest(method, path, nil)
							recorder := httptest.NewRecorder()
							requestHandler.ServeHTTP(recorder, req)

							Expect(recorder.Code).To(Equal(http.StatusBadRequest))
							Expect(recorder.Body).To(ContainSubstring("Key must consist of alphanumeric, underscores, dashes, and forward slashes"))
						}
					}
				})
			})
		})

		Context("when URL path is valid", func() {
			Context("when http method is not supported", func() {
				It("should return 405 Method Not Allowed", func() {
					invalidMethods := []string{"PATCH"}
					http.NewRequest("PUT", "/v1/data/bla", strings.NewReader("value=blabla"))

					for _, method := range invalidMethods {
						req, _ := http.NewRequest(method, "/v1/data/bla", nil)
						recorder := httptest.NewRecorder()
						requestHandler.ServeHTTP(recorder, req)

						Expect(recorder.Code).To(Equal(http.StatusMethodNotAllowed))
					}
				})
			})

			Context("when http method is supported", func() {
				validUrlPaths := map[string]string{
					"/v1/data/smurf":                          "smurf",
					"/v1/data/smurf/color":                    "smurf/color",
					"/v1/data/smurf/color/darkness":           "smurf/color/darkness",
					"/v1/data/smurf/color/dark_ness/name-tag": "smurf/color/dark_ness/name-tag",
				}

				Describe("GET", func() {

					It("can handle all types of validkeys", func() {
						respValue := store.Configuration{
							Value: `{"value":"common value"}`,
						}
						mockStore.GetReturns(respValue, nil)
						var counter int = 0
						for path, extractedKey := range validUrlPaths {
							getReq, _ := http.NewRequest("GET", path, nil)
							getRecorder := httptest.NewRecorder()

							requestHandler.ServeHTTP(getRecorder, getReq)
							key := mockStore.GetArgsForCall(counter)

							Expect(key).To(Equal(extractedKey))

							Expect(getRecorder.Code).To(Equal(http.StatusOK))
							counter = counter + 1
						}
					})

					Context("when key exists", func() {
						It("returns value in the store", func() {
							values := []string{
								`123`,
								`"blabla"`,
								`{"key":"blabla"}`,
							}

							for _, value := range values {
								respValue := store.Configuration{
									Value: `{"value":` + value + "}",
									Key:   "bla",
									Id:    "some_id",
								}
								mockStore.GetReturns(respValue, nil)

								getReq, _ := http.NewRequest("GET", "/v1/data/bla/", nil)
								getRecorder := httptest.NewRecorder()
								requestHandler.ServeHTTP(getRecorder, getReq)

								Expect(getRecorder.Code).To(Equal(http.StatusOK))
								expectedString := `{"id":"some_id","path":"bla","value":` + value + "}"
								Expect(getRecorder.Body.String()).To(Equal(expectedString))
							}
						})
					})

					Context("when key does not exist", func() {
						It("should return 404 Not Found", func() {
							req, _ := http.NewRequest("GET", "/v1/data/test", nil)
							getRecorder := httptest.NewRecorder()
							requestHandler.ServeHTTP(getRecorder, req)

							Expect(getRecorder.Code).To(Equal(http.StatusNotFound))
						})
					})

					Context("when store errors", func() {
						It("returns 500 Internal Server Error", func() {
							mockStore.GetReturns(store.Configuration{}, errors.New("Kaboom!"))

							getReq, _ := http.NewRequest("GET", "/v1/data/bla/", nil)
							getRecorder := httptest.NewRecorder()
							requestHandler.ServeHTTP(getRecorder, getReq)

							Expect(getRecorder.Code).To(Equal(http.StatusInternalServerError))
						})
					})
				})

				Describe("PUT", func() {
					It("can handle all types of validkeys", func() {
						config := store.Configuration{
							Key: "bla",
							Value: `{"value":"burpees"}`,
							Id: "1",
						}
						mockStore.GetReturns(config, nil)

						var counter int = 0
						for path, extractedKey := range validUrlPaths {
							req, _ := http.NewRequest("PUT", path, strings.NewReader(`{"value":"str"}`))
							putRecorder := httptest.NewRecorder()
							requestHandler.ServeHTTP(putRecorder, req)

							key, _ := mockStore.PutArgsForCall(counter)

							Expect(key).To(Equal(extractedKey))
							Expect(putRecorder.Code).To(Equal(http.StatusOK))
							counter++
						}
					})

					Context("when request body is NOT in the specified format", func() {
						Context("when body is empty", func() {
							It("should return 400 Bad Request", func() {
								req, _ := http.NewRequest("PUT", "/v1/data/key", nil)
								putRecorder := httptest.NewRecorder()
								requestHandler.ServeHTTP(putRecorder, req)

								Expect(putRecorder.Body.String()).To(ContainSubstring("Request can't be empty"))
								Expect(putRecorder.Code).To(Equal(http.StatusBadRequest))
							})
						})

						Context("when body is NOT JSON string", func() {
							It("should return 400 Bad Request", func() {
								req, _ := http.NewRequest("PUT", "/v1/data/key", strings.NewReader(`smurf`))
								putRecorder := httptest.NewRecorder()
								requestHandler.ServeHTTP(putRecorder, req)

								Expect(putRecorder.Body.String()).To(ContainSubstring("Request Body should be JSON string"))
								Expect(putRecorder.Code).To(Equal(http.StatusBadRequest))
							})
						})

						Context("when body is JSON string but NOT as expected", func() {
							It("should return 400 Bad Request", func() {
								req, _ := http.NewRequest("PUT", "/v1/data/key", strings.NewReader(`{"smurf":"blue"}`))
								putRecorder := httptest.NewRecorder()
								requestHandler.ServeHTTP(putRecorder, req)

								Expect(putRecorder.Body.String()).To(ContainSubstring("JSON request body shoud contain the key 'value'"))
								Expect(putRecorder.Code).To(Equal(http.StatusBadRequest))
							})
						})
					})

					Context("when request body is in the specified format", func() {

						BeforeEach(func() {
							config := store.Configuration{
								Key: "bla",
								Value: `{"value":"burpees"}`,
								Id: "1",
							}
							mockStore.GetReturns(config, nil)
						})


						It("returns value, path and id in the response", func() {
							req, _ := http.NewRequest("PUT", "/v1/data/bla", strings.NewReader(`{"value":"str"}`))
							putRecorder := httptest.NewRecorder()
							requestHandler.ServeHTTP(putRecorder, req)

							key := mockStore.GetArgsForCall(0)
							Expect(key).To(Equal("bla"))

							Expect(putRecorder.Body.String()).To(Equal(`{"id":"1","path":"bla","value":"burpees"}`))
						})

						Context("when value is a string ", func() {
							It("should store value in a specific JSON format and respond with 204 StatusNoContent", func() {
								req, _ := http.NewRequest("PUT", "/v1/data/bla", strings.NewReader(`{"value":"str"}`))
								putRecorder := httptest.NewRecorder()
								requestHandler.ServeHTTP(putRecorder, req)

								Expect(mockStore.PutCallCount()).To(Equal(1))
								key, value := mockStore.PutArgsForCall(0)

								Expect(key).To(Equal("bla"))
								Expect(value).To(Equal(`{"value":"str"}`))
								Expect(putRecorder.Code).To(Equal(http.StatusOK))
							})
						})

						Context("when value is a number", func() {
							It("should store value in a specific JSON format and respond with 204 StatusNoContent", func() {
								req, _ := http.NewRequest("PUT", "/v1/data/bla", strings.NewReader(`{"value":123}`))
								putRecorder := httptest.NewRecorder()
								requestHandler.ServeHTTP(putRecorder, req)

								Expect(mockStore.PutCallCount()).To(Equal(1))
								key, value := mockStore.PutArgsForCall(0)

								Expect(key).To(Equal("bla"))
								Expect(value).To(Equal(`{"value":123}`))
								Expect(putRecorder.Code).To(Equal(http.StatusOK))
							})
						})

						Context("when value is a JSON hash", func() {
							It("should store value in a specific JSON format and respond with 204 StatusNoContent", func() {
								requestBody := `{"value":{"age":10,"color":"red"}}`
								valueToStore := `{"value":{"age":10,"color":"red"}}`

								req, _ := http.NewRequest("PUT", "/v1/data/bla", strings.NewReader(requestBody))
								putRecorder := httptest.NewRecorder()
								requestHandler.ServeHTTP(putRecorder, req)

								Expect(mockStore.PutCallCount()).To(Equal(1))
								key, value := mockStore.PutArgsForCall(0)

								Expect(key).To(Equal("bla"))
								Expect(value).To(Equal(valueToStore))
								Expect(putRecorder.Code).To(Equal(http.StatusOK))
							})
						})
					})
				})

				Describe("POST", func() {
					It("can handle all types of validkeys", func() {
						respValue := store.Configuration{
							Value: `{"value":"anything"}`,
						}
						mockStore.GetReturns(respValue, nil)

						var counter int = 0
						for path, extractedKey := range validUrlPaths {
							req, _ := http.NewRequest("POST", path, strings.NewReader(`{"type":"password","parameters":{}}`))
							recorder := httptest.NewRecorder()
							requestHandler.ServeHTTP(recorder, req)

							key := mockStore.GetArgsForCall(counter)

							Expect(key).To(Equal(extractedKey))
							counter++
						}
					})

					Context("when request body is NOT in the specified format", func() {
						Context("when body is empty", func() {
							It("should return 400 Bad Request", func() {
								req, _ := http.NewRequest("POST", "/v1/data/key", nil)
								recorder := httptest.NewRecorder()
								requestHandler.ServeHTTP(recorder, req)

								Expect(recorder.Body.String()).To(ContainSubstring("Request can't be empty"))
								Expect(recorder.Code).To(Equal(http.StatusBadRequest))
							})
						})

						Context("when body is NOT JSON string", func() {
							It("should return 400 Bad Request", func() {
								req, _ := http.NewRequest("POST", "/v1/data/key", strings.NewReader("smurf"))
								recorder := httptest.NewRecorder()
								requestHandler.ServeHTTP(recorder, req)

								Expect(recorder.Body.String()).To(ContainSubstring("Request Body should be JSON string"))
								Expect(recorder.Code).To(Equal(http.StatusBadRequest))
							})
						})

						Context("when body is JSON string but NOT as expected", func() {
							It("should return 400 Bad Request", func() {
								req, _ := http.NewRequest("POST", "/v1/data/key", strings.NewReader(`{"smurf":"blue"}`))
								recorder := httptest.NewRecorder()
								requestHandler.ServeHTTP(recorder, req)

								Expect(recorder.Body.String()).To(ContainSubstring("JSON request body shoud contain the key 'type'"))
								Expect(recorder.Code).To(Equal(http.StatusBadRequest))
							})
						})
					})

					Context("when request body is in the specified format", func() {

						Describe("Password generation", func() {
							Context("when value already exists", func() {
								It("should not generate a password", func() {
									mockStore.GetStub = func(key string) (store.Configuration, error) {
										respValue := store.Configuration{
											Value: `{"value":"smurf"}`,
											Id:    "some_id",
											Key:   "bla",
										}
										return respValue, nil
									}

									postReq, _ := http.NewRequest("POST", "/v1/data/bla/", strings.NewReader(`{"type":"password","parameters":{}}`))

									recorder := httptest.NewRecorder()
									requestHandler.ServeHTTP(recorder, postReq)

									Expect(recorder.Code).To(Equal(http.StatusOK))
									Expect(recorder.Body.String()).To(Equal(`{"id":"some_id","path":"bla","value":"smurf"}`))
									Expect(mockValueGeneratorFactory.GetGeneratorCallCount()).To(Equal(0))
								})
							})

							Context("when value does NOT exist", func() {
								It("should return generated password", func() {
									requestHandler, _ = NewRequestHandler(store.NewMemoryStore(), types.NewValueGeneratorConcrete(config.ServerConfig{}))

									postReq, _ := http.NewRequest("POST", "/v1/data/bla/", strings.NewReader(`{"type":"password","parameters":{}}`))

									recorder := httptest.NewRecorder()
									requestHandler.ServeHTTP(recorder, postReq)

									Expect(recorder.Code).To(Equal(http.StatusCreated))

									var data map[string]string
									json.Unmarshal(recorder.Body.Bytes(), &data)

									Expect(data["path"]).To(Equal("bla"))
									Expect(data["value"]).Should(MatchRegexp("[a-z0-9]{20}"))
								})
							})
						})

						Describe("Certificate generation", func() {
							Context("when value already exists", func() {
								It("should not generate certificates", func() {

									mockStore.GetStub = func(key string) (store.Configuration, error) {
										respValue := store.Configuration{
											Value: `{"value":"smurf"}`,
											Id:    "some_id",
											Key:   "bla",
										}
										return respValue, nil
									}

									postReq, _ := http.NewRequest("POST", "/v1/data/bla/", strings.NewReader(`{"type":"certificate","parameters":{}}`))

									recorder := httptest.NewRecorder()
									requestHandler.ServeHTTP(recorder, postReq)

									Expect(recorder.Code).To(Equal(http.StatusOK))
									Expect(recorder.Body.String()).To(Equal(`{"id":"some_id","path":"bla","value":"smurf"}`))
									Expect(mockValueGeneratorFactory.GetGeneratorCallCount()).To(Equal(0))
								})
							})

							Context("when value does NOT exist", func() {
								It("should return generated certificate, its private key and root certificate used to sign the generated certificate", func() {
									requestHandler, _ = NewRequestHandler(store.NewMemoryStore(), mockValueGeneratorFactory)
									mockValueGeneratorFactory.GetGeneratorReturns(mockValueGenerator, nil)

									mockValueGenerator.GenerateReturns(types.CertResponse{
										Certificate: "fake-certificate",
										PrivateKey:  "fake-private-key",
										CA:          "fake-ca",
									}, nil)

									postReq, _ := http.NewRequest("POST", "/v1/data/bla/", strings.NewReader(`{"type":"certificate","parameters":{"common_name": "asdf", "alternative_names":["nam1", "name2"]}}`))

									recorder := httptest.NewRecorder()
									requestHandler.ServeHTTP(recorder, postReq)

									Expect(recorder.Code).To(Equal(http.StatusCreated))

									var data map[string]interface{}
									json.Unmarshal(recorder.Body.Bytes(), &data)

									Expect(data["path"]).To(Equal("bla"))

									value := data["value"].(map[string]interface{})
									Expect(value["certificate"]).To(Equal("fake-certificate"))
									Expect(value["private_key"]).To(Equal("fake-private-key"))
									Expect(value["ca"]).To(Equal("fake-ca"))
								})
							})
						})
					})
				})

				Describe("DELETE", func() {
					It("can handle all types of validkeys", func() {
						mockStore.DeleteReturns(true, nil)

						var counter int = 0
						for path, extractedKey := range validUrlPaths {
							req, _ := http.NewRequest("DELETE", path, nil)

							recorder := httptest.NewRecorder()
							requestHandler.ServeHTTP(recorder, req)

							Expect(mockStore.DeleteArgsForCall(counter)).To(Equal(extractedKey))
							counter++
						}
					})

					Context("Key exists", func() {
						BeforeEach(func() {
							mockStore.DeleteReturns(true, nil)
						})

						It("should delete value", func() {
							req, _ := http.NewRequest("DELETE", "/v1/data/bla", nil)

							putRecorder := httptest.NewRecorder()
							requestHandler.ServeHTTP(putRecorder, req)

							Expect(mockStore.DeleteCallCount()).To(Equal(1))
							Expect(mockStore.DeleteArgsForCall(0)).To(Equal("bla"))
						})

						It("should return 204 Status No Content", func() {
							req, _ := http.NewRequest("DELETE", "/v1/data/bla", nil)
							req.Header.Set("Authorization", "bearer fake-auth-header")

							putRecorder := httptest.NewRecorder()
							requestHandler.ServeHTTP(putRecorder, req)

							Expect(putRecorder.Code).To(Equal(http.StatusNoContent))
						})
					})

					Context("Key does not exist", func() {
						It("should return 404 Status Not Found", func() {
							req, _ := http.NewRequest("DELETE", "/v1/data/bla", nil)
							req.Header.Set("Authorization", "bearer fake-auth-header")

							putRecorder := httptest.NewRecorder()
							requestHandler.ServeHTTP(putRecorder, req)

							Expect(putRecorder.Code).To(Equal(http.StatusNotFound))
							Expect(mockStore.DeleteCallCount()).To(Equal(1))
							Expect(mockStore.DeleteArgsForCall(0)).To(Equal("bla"))
						})
					})
				})
			})
		})
	})
})
