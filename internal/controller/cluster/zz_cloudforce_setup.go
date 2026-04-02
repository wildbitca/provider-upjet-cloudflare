// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	onerequest "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/cloudforce/onerequest"
	onerequestasset "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/cloudforce/onerequestasset"
	onerequestmessage "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/cloudforce/onerequestmessage"
	onerequestpriority "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/cloudforce/onerequestpriority"
)

// Setup_cloudforce creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_cloudforce(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		onerequest.Setup,
		onerequestasset.Setup,
		onerequestmessage.Setup,
		onerequestpriority.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_cloudforce creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_cloudforce(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		onerequest.SetupGated,
		onerequestasset.SetupGated,
		onerequestmessage.SetupGated,
		onerequestpriority.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
