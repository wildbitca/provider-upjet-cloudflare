// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	scanning "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/content/scanning"
	scanningexpression "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/content/scanningexpression"
)

// Setup_content creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_content(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		scanning.Setup,
		scanningexpression.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_content creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_content(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		scanning.SetupGated,
		scanningexpression.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupWebhookWithManager_content registers conversion webhooks for all resource kinds in the group.
func SetupWebhookWithManager_content(mgr ctrl.Manager) error {
	for _, setup := range []func(ctrl.Manager) error{
		scanning.SetupWebhookWithManager,
		scanningexpression.SetupWebhookWithManager,
	} {
		if err := setup(mgr); err != nil {
			return err
		}
	}
	return nil
}
