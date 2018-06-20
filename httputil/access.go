package httputil

import (
	"context"

	"github.com/containerum/cherry"
	"github.com/containerum/cherry/adaptors/gonic"
	kubeModel "github.com/containerum/kube-client/pkg/model"
	"github.com/gin-gonic/gin"
)

const AccessContext = "access-ctx"

const (
	ProjectParam   = "project"
	NamespaceParam = "namespace"
)

type ProjectAccess struct {
	ProjectID          string            `json:"project_id"`
	ProjectLabel       string            `json:"project_label"`
	NamespacesAccesses []NamespaceAccess `json:"namespaces"`
}

type NamespaceAccess struct {
	NamespaceID    string                    `json:"namespace_id"`
	NamespaceLabel string                    `json:"namespace_label"`
	Access         kubeModel.UserGroupAccess `json:"access"`
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
		project := ctx.Param(ProjectParam)
		ns := ctx.Param(NamespaceParam)

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
