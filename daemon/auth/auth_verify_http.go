// Copyright 2022 Elapse and contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package auth

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/elap5e/penguin/pkg/log"
)

func serveHTTPVerifyCaptcha(l net.Listener) (*CaptchaSign, error) {
	return serveHTTP(l, tmplVerifyCaptcha)
}

func serveHTTPVerifySignInWithCodeCaptach(l net.Listener) (*CaptchaSign, error) {
	return serveHTTP(l, tmplVerifySignInWithCodeCaptach)
}

func serveHTTP(l net.Listener, t string) (*CaptchaSign, error) {
	done := make(chan *CaptchaSign, 1)
	mux := http.NewServeMux()
	mux.Handle("/captcha", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				w.WriteHeader(http.StatusOK)
				fmt.Fprintln(w, t)
			case http.MethodPost:
				w.WriteHeader(http.StatusOK)
				id, _ := strconv.ParseUint(r.FormValue("appid"), 10, 64)
				done <- &CaptchaSign{
					Ticket: r.FormValue("ticket"),
					Random: r.FormValue("randstr"),
					Return: r.FormValue("ret"),
					AppID:  id,
				}
			}
		},
	))
	srv := &http.Server{
		Handler: mux,
	}
	go func() {
		err := srv.Serve(l)
		if err != nil && err != http.ErrServerClosed {
			log.Panicln(err)
		}
	}()
	form := <-done
	ctxShutDown, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()
	if err := srv.Shutdown(ctxShutDown); err != nil {
		return nil, err
	}
	return form, nil
}

const tmplVerifyCaptcha = `<!DOCTYPE html>
<html>

<head lang="zh-CN">
    <meta charset="UTF-8" />
    <meta name="renderer" content="webkit" />
    <meta name="viewport"
        content="width=device-width,initial-scale=1,minimum-scale=1,maximum-scale=1,user-scalable=no" />
    <title>验证码</title>
</head>

<body>
    <div id="cap_iframe" style="width: 230px; height: 220px"></div>
    <script type="text/javascript">
        !(function () {
            var e = document.createElement("script");
            e.type = "text/javascript";
            e.src = "http://captcha.qq.com/template/TCapIframeApi.js" + location.search;
            document.getElementsByTagName("head").item(0).appendChild(e);
            e.onload = function () {
                capInit(document.getElementById("cap_iframe"), {
                    callback: function (a) {
                        var xhr = new XMLHttpRequest();
                        xhr.open("POST", "/captcha", true);
                        var d = new FormData();
                        d.append("ticket", a.ticket);
                        d.append("randstr", a.randstr);
                        d.append("appid", a.appid);
                        d.append("ret", a.ret);
                        xhr.onload = function (e) { window.close(); };
                        xhr.send(d);
                    },
                    showHeader: !1,
                });
            };
        })();
    </script>
</body>

</html>`

const tmplVerifySignInWithCodeCaptach = `<!DOCTYPE html>
<html>

<head lang="zh-CN">
    <meta charset="UTF-8" />
    <meta name="renderer" content="webkit" />
    <meta name="viewport"
        content="width=device-width,initial-scale=1,minimum-scale=1,maximum-scale=1,user-scalable=no" />
    <title>验证码</title>
</head>

<body>
    <div id="cap_iframe" style="width: 230px; height: 220px"></div>
    <script type="text/javascript">
        !(function () {
            var e = document.createElement("script");
            e.type = "text/javascript";
            e.src = "http://captcha.qq.com/TCaptcha.js";
            document.getElementsByTagName("head").item(0).appendChild(e);
            e.onload = function () {
                new window.TencentCaptcha(
                    "2081081773",
                    function (a) {
                        var xhr = new XMLHttpRequest();
                        xhr.open("POST", "/captcha", true);
                        var d = new FormData();
                        d.append("ticket", a.ticket);
                        d.append("randstr", a.randstr);
                        d.append("appid", a.appid);
                        d.append("ret", a.ret);
                        xhr.onload = function (e) { window.close(); };
                        xhr.send(d);
                    },
                    { type: "full", showHeader: !1 }
                ).show();
            };
        })();
    </script>
</body>

</html>`
