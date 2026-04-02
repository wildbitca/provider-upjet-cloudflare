// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	shield "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/api/shield"
	shielddiscoveryoperation "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/api/shielddiscoveryoperation"
	shieldoperation "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/api/shieldoperation"
	shieldoperationschemavalidationsettings "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/api/shieldoperationschemavalidationsettings"
	shieldschema "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/api/shieldschema"
	shieldschemavalidationsettings "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/api/shieldschemavalidationsettings"
	token "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/api/token"
)

// Setup_api creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_api(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		shield.Setup,
		shielddiscoveryoperation.Setup,
		shieldoperation.Setup,
		shieldoperationschemavalidationsettings.Setup,
		shieldschema.Setup,
		shieldschemavalidationsettings.Setup,
		token.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_api creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_api(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		shield.SetupGated,
		shielddiscoveryoperation.SetupGated,
		shieldoperation.SetupGated,
		shieldoperationschemavalidationsettings.SetupGated,
		shieldschema.SetupGated,
		shieldschemavalidationsettings.SetupGated,
		token.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
