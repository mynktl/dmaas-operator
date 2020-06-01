package prebackupaction

import "github.com/mayadata-io/dmaas-operator/pkg/apis/mayadata.io/v1alpha1"

type apiResource struct {
	name       string
	kind       string
	apiVersion string
	namespace  string
}

type Request struct {
	*v1alpha1.PreBackupAction

	LabeledList   map[apiResource]bool
	AnnotatedList map[apiResource]bool
}
