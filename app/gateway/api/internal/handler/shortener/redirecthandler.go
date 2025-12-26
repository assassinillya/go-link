// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package shortener

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-link/app/gateway/api/internal/logic/shortener"
	"go-link/app/gateway/api/internal/svc"
	"go-link/app/gateway/api/internal/types"
)

func RedirectHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetLinkReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := shortener.NewRedirectLogic(r.Context(), svcCtx)
		resp, err := l.Redirect(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			//httpx.OkJsonCtx(r.Context(), w, resp)
			http.Redirect(w, r, resp.Url, http.StatusFound) // 不返回json而是302跳转
		}
	}
}
