package httputil

import (
	"context"
	"fmt"

	"github.com/containerum/cherry"
	"github.com/containerum/cherry/adaptors/gonic"
	kubeModel "github.com/containerum/kube-client/pkg/model"
	"github.com/gin-gonic/gin"
)

const AccessContext = "access-ctx"

type ProjectAccess struct {
	ProjectID          string
	ProjectLabel       string
	NamespacesAccesses []NamespaceAccess
}

type NamespaceAccess struct {
	NamespaceID    string
	NamespaceLabel string
	Access         kubeModel.UserGroupAccess
}

type Permissions interface {
	GetAllAccesses(ctx context.Context) ([]ProjectAccess, error)
	GetNamespaceAccess(ctx context.Context, projectID, namespaceID string) (NamespaceAccess, error)
}

type AccessChecker struct {
	PermissionsClient Permissions
	AccessError       cherry.ErrConstruct
	NotFoundError     cherry.ErrConstruct
}

func (a *AccessChecker) CheckAccess(requiredAccess string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requiredAccessParsed, err := kubeModel.UserGroupAccessString(requiredAccess)
		if err != nil {
			panic(err)
		}
		ns := ctx.Param("namespace")
		project := ctx.Param("project")
		fmt.Println(ns, project)

		namespaceAccess, err := a.PermissionsClient.GetNamespaceAccess(ctx, project, ns)
		if err != nil {
			gonic.Gonic(a.AccessError(), ctx)
			return
		}

		if namespaceAccess.Access < requiredAccessParsed {
			gonic.Gonic(a.NotFoundError(), ctx)
			return
		}

		rctx := context.WithValue(ctx.Request.Context(), AccessContext, namespaceAccess)
		ctx.Request = ctx.Request.WithContext(rctx)

	}
}
