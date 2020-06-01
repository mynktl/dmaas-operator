/*
Copyright 2020 The MayaData Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    https://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"time"

	v1alpha1 "github.com/mayadata-io/dmaas-operator/pkg/apis/mayadata.io/v1alpha1"
	clientset "github.com/mayadata-io/dmaas-operator/pkg/generated/clientset/versioned"
	informers "github.com/mayadata-io/dmaas-operator/pkg/generated/informers/externalversions/mayadata.io/v1alpha1"
	dmaaslister "github.com/mayadata-io/dmaas-operator/pkg/generated/listers/mayadata.io/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	prebackupaction "github.com/mayadata-io/dmaas-operator/pkg/prebackupaction"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/sirupsen/logrus"

	apimachineryclock "k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

var (
	preBackupActionSyncPeriod = 30 * time.Second
)

type preBackupActionController struct {
	*controller
	namespace     string
	kubeClient    kubernetes.Interface
	dmaasClient   clientset.Interface
	dynamicClient dynamic.Interface
	lister        dmaaslister.PreBackupActionLister
	prebackup     prebackupaction.PreBackupActioner
	clock         apimachineryclock.Clock
}

// NewPreBackupActionController returns controller for prebackupaction resource
func NewPreBackupActionController(
	namespace string,
	kubeClient kubernetes.Interface,
	dmaasClient clientset.Interface,
	dynamicClient dynamic.Interface,
	preBackupActionInformer informers.PreBackupActionInformer,
	logger logrus.FieldLogger,
	clock apimachineryclock.Clock,
	numWorker int,
) Controller {
	c := &preBackupActionController{
		controller:    newController("prebackupaction", logger, numWorker),
		namespace:     namespace,
		kubeClient:    kubeClient,
		lister:        preBackupActionInformer.Lister(),
		dmaasClient:   dmaasClient,
		dynamicClient: dynamicClient,
		clock:         clock,
	}

	c.reconcile = c.processPreBackupAction
	c.syncPeriod = preBackupActionSyncPeriod

	preBackupActionInformer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				c.enqueue(obj)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				_ = oldObj
				c.enqueue(newObj)
			},
			DeleteFunc: func(obj interface{}) {
				c.enqueue(obj)
			},
		},
	)
	return c
}

func (p *preBackupActionController) processPreBackupAction(key string) error {
	log := p.logger.WithField("key", key)

	log.Debug("processing prebackupaction")

	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		log.WithError(err).Errorf("failed to split key")
		return nil
	}

	original, err := p.lister.PreBackupActions(ns).Get(name)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}
		return errors.Wrapf(err, "failed to get prebackupaction")
	}

	isActionOneTime := isPreBackupActionOneTime(original)

	// if prebackupaction has OneTimeAction set and phase is completed
	// then we don't need to process it
	if isActionOneTime && original.Status.Phase == v1alpha1.PreBackupActionPhaseCompleted {
		return nil
	}

	// validate and update status to InProgress or return error

	preBackupObj := original.DeepCopy()

	// we are processing preBackupObj so update failure time to empty
	preBackupObj.Status.LastFailureTimestamp = nil

	if err = p.prebackup.Action(preBackupObj); err != nil {
		preBackupObj.Status.LastFailureTimestamp = &metav1.Time{Time: p.clock.Now()}
	} else {
		preBackupObj.Status.LastSuccessfulTimestamp = &metav1.Time{Time: p.clock.Now()}
	}

	if isActionOneTime && preBackupObj.Status.LastFailureTimestamp == nil {
		preBackupObj.Status.Phase = v1alpha1.PreBackupActionPhaseCompleted
	}

	// update the object
	return nil
}

func isPreBackupActionOneTime(preBackupActionObj *v1alpha1.PreBackupAction) bool {
	if preBackupActionObj.Spec.OneTimeAction != nil &&
		*preBackupActionObj.Spec.OneTimeAction == true {
		return true
	}
	return false
}
