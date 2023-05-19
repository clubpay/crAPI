/*
 * Licensed under the Apache License, Version 2.0 (the “License”);
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package router

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"crapi.proj/goservice/api/config"
	"crapi.proj/goservice/api/controllers"
	"crapi.proj/goservice/api/middlewares"
	"github.com/gorilla/mux"
)

type Server config.Server

var controller = controllers.Server{}

// initializeRoutes initialize routes of url with Authentication or without Authentication
func (server *Server) InitializeRoutes() *mux.Router {

	controller.DB = server.DB

	controller.Client = server.Client

	server.Router.Use(middlewares.AccessControlMiddleware)
	// Post Route
	server.Router.HandleFunc("/community/api/v2/community/posts/recent", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(controller.GetPost, server.DB))).Methods("GET", "OPTIONS")

	server.Router.HandleFunc("/community/api/v2/community/posts/{postID}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(controller.GetPostByID, server.DB))).Methods("GET", "OPTIONS")

	server.Router.HandleFunc("/community/api/v2/community/posts", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(controller.AddNewPost, server.DB))).Methods("POST", "OPTIONS")

	server.Router.HandleFunc("/community/api/v2/community/posts/{postID}/comment", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(controller.Comment, server.DB))).Methods("POST", "OPTIONS")

	//Coupon Route
	server.Router.HandleFunc("/community/api/v2/coupon/new-coupon", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(controller.AddNewCoupon, server.DB))).Methods("POST", "OPTIONS")

	server.Router.HandleFunc("/community/api/v2/coupon/validate-coupon", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(controller.ValidateCoupon, server.DB))).Methods("POST", "OPTIONS")

	//Health
	server.Router.HandleFunc("/community/home", middlewares.SetMiddlewareJSON(controller.Home)).Methods("GET")
	return server.Router
}

func isTrue(a string) bool {
	a = strings.ToLower(a)
	true_list := []string{"true", "1", "t", "y", "yes"}
	for _, b := range true_list {
		if b == a {
			return true
		}
	}
	return false
}

func (server *Server) Run(addr string) {
	fmt.Println("Listening to port " + os.Getenv("SERVER_PORT"))
	tls_enabled, is_tls := os.LookupEnv("TLS_ENABLED")
	if is_tls && isTrue(tls_enabled) {
		// Check if env variable TLS_CERTIFICATE is set then use it as certificate else default to certs/server.crt
		certificate, is_cert := os.LookupEnv("TLS_CERTIFICATE")
		if !is_cert || certificate == "" {
			certificate = "certs/server.crt"
		}
		// Check if env variable TLS_KEY is set then use it as key else default to certs/server.key
		key, is_key := os.LookupEnv("TLS_KEY")
		if !is_key || key == "" {
			key = "certs/server.key"
		}
		err := http.ListenAndServeTLS(addr, certificate, key, server.Router)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		err := http.ListenAndServe(addr, server.Router)
		if err != nil {
			fmt.Println(err)
		}
	}
}
