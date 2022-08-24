package main

import (
	"log"
	"net/http"

	"github.com/blackwind-code/BWOJ/pkg/util"
	"github.com/didip/tollbooth"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

var CONF gjson.Result

// Made by Heeyong Yoon
type MultipleDirStaticServe struct{}

func (h MultipleDirStaticServe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["project_id"]
	// DB에서 path 찾는다
	path := util.PathSanitize("/storage/" + id)
	http.FileServer(http.Dir(path)).ServeHTTP(w, r)
}

func limiter(f func(http.ResponseWriter, *http.Request)) http.Handler {
	if CONF.Get("http_request_limiter.use").Bool() {
		max_req_option := CONF.Get("http_request_limiter.max_request_per_second")
		max_req := float64(10) // maximum number of reqests to limit per second
		if max_req_option.Exists() {
			max_req = max_req_option.Float()
		}

		return tollbooth.LimitFuncHandler(
			tollbooth.NewLimiter(max_req, nil).
				SetStatusCode(http.StatusTooManyRequests).
				SetMessageContentType("text/plain; charset=utf-8").
				SetMessage("Your request rate is too high. Are you a bot?").
				SetOnLimitReached(func(w http.ResponseWriter, r *http.Request) {
					logrus.WithFields(logrus.Fields{
						"HTTP status":     http.StatusTooManyRequests,
						"X-Forwarded-For": r.Header.Get("X-Forwarded-For"),
						"X-Real-IP":       r.Header.Get("X-Real-IP"),
					}).Info("Detect high request rate. Block temporarily.")
				}), f)
	} else {
		return http.HandlerFunc(f)
	}
}

func Q(w http.ResponseWriter, r *http.Request) {
	out := ""
	out, _ = sjson.Set(out, "success", true)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(out))
}

func A(w http.ResponseWriter, r *http.Request) {
	content := []byte{}
	if _, e := r.Body.Read(content); e != nil {
		log.Println(e)
		return
	}

	in := gjson.ParseBytes(content)
	log.Println(in.Get("answer").String())

	out := ""
	out, _ = sjson.Set(out, "success", true)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(out))
}

func init() {
	CONF = gjson.Parse(util.File2Str("./config.json"))
	init_logrus(CONF.Get("log.path").String(), CONF.Get("log.level").String(), CONF.Get("log.format").String())
}

func main() {
	r := mux.NewRouter()

	r.Methods(http.MethodGet).Path("/q").Handler(limiter(Q))
	r.Methods(http.MethodPost).Path("/a").Handler(limiter(A))
	r.Methods(http.MethodGet).Path("/resource/").Handler(MultipleDirStaticServe{})

	r.Use(mux.CORSMethodMiddleware(r))
	http.ListenAndServe(":"+CONF.Get("server_port").String(), r)
}
