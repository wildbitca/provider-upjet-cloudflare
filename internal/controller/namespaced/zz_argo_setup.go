// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	smartrouting "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/argo/smartrouting"
	tieredcaching "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/argo/tieredcaching"
)

// Setup_argo creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_argo(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		smartrouting.Setup,
		tieredcaching.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_argo creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_argo(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		smartrouting.SetupGated,
		tieredcaching.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupWebhookWithManager_argo registers conversion webhooks for all resource kinds in the group.
func SetupWebhookWithManager_argo(mgr ctrl.Manager) error {
	for _, setup := range []func(ctrl.Manager) error{
		smartrouting.SetupWebhookWithManager,
		tieredcaching.SetupWebhookWithManager,
	} {
		if err := setup(mgr); err != nil {
			return err
		}
	}
	return nil
}
