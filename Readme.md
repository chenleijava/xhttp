## 基于go-zero 模板定义统一返回消息，自定义模板

- 当前 goctl 生成工具版本是1.8.1
- [go-zero 模板定制](https://go-zero.dev/docs/tutorials/customization/template)
- 模板是基于goctl工具生成，对于api生成，会引用tpl/api/handler.tpl文件 
- 详细修改如下：

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