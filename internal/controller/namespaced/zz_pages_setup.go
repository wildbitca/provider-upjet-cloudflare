// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	domain "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/pages/domain"
	project "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/pages/project"
)

// Setup_pages creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_pages(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		domain.Setup,
		project.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_pages creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_pages(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		domain.SetupGated,
		project.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
