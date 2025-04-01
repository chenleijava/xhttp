## Define a unified return message based on the go-zero template, customize the template

- the goctl tool version: 1.8.1
- [go-zero Template customization](https://go-zero.dev/docs/tutorials/customization/template)
- The template is generated based on the goctl tool. For api generation, the tpl/api/handler.tpl file is referenced
- Detailed changes are as follows:

```text
package handler

import (
    "github.com/zeromicro/go-zero/rest/httpx"
    "net/http"
    "backend/render" 
    {{.ImportPackages}}
)

func {{.HandlerName}}(svcCtx *svc.ServiceContext) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        {{if .HasRequest}}var req types.{{.RequestType}}
        if err := httpx.Parse(r, &req); err != nil {
            httpx.Error(w, err)
            return
        }{{end}}

        l := {{.LogicName}}.New{{.LogicType}}(r.Context(), svcCtx)
        {{if .HasResp}}resp, {{end}}err := l.{{.Call}}({{if .HasRequest}}&req{{end}})
        {{if .HasResp}}render.ResponseJson(w, resp, err){{else}}render.ResponseJson(w, nil, err){{end}}
    }
}
```

```makefile
API_DIR:=./backend/etcdwebtool.api # api文件所在目录
TARGET_DIR= ./backend #代码输出目录
goctl-api-go:
	goctl api go -api ${API_DIR} -dir ${TARGET_DIR} --style go_zero --home=./tpl
```

---

### gen code such as:

```go
package handler

import (
	"backend/internal/logic/auth"
	"backend/internal/svc"
	"backend/internal/types"
	render "github.com/chenleijava/xhttp"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := auth.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(&req)
		render.ResponseJson(w, resp, err)
	}
}
```

Enjoy coding !