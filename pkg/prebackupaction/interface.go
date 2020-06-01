package prebackupaction

import (
	"github.com/mayadata-io/dmaas-operator/pkg/apis/mayadata.io/v1alpha1"
)

// PreBackupActioner returns interface to execute operation on prebackupaction resource
type PreBackupActioner interface {
	// Action to execute prebackupaction
	Action(preBackupActionObj *v1alpha1.PreBackupAction) error
}

// parser returns interface to parse resource
type parser interface {
	// parse to parse the given resource
	parse(resource apiResource) error
}
