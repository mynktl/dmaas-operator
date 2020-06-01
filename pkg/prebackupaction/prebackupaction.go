package prebackupaction

import (
	"github.com/mayadata-io/dmaas-operator/pkg/apis/mayadata.io/v1alpha1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

type preBackupAction struct {
	kubeClient    kubernetes.Interface
	dynamicClient dynamic.Interface
	parser        map[string]parser
}

func NewPreBackupAction(
	kubeClient kubernetes.Interface,
	dynamicClient dynamic.Interface,
) PreBackupActioner {
	return &preBackupAction{
		dynamicClient: dynamicClient,
		kubeClient:    kubeClient,
		parser:        getSupportedParser(),
	}
}

func (p *preBackupAction) Action(obj *v1alpha1.PreBackupAction) error {
	return nil
}

func getSupportedParser() map[string]parser {
	return map[string]parser{}
}
