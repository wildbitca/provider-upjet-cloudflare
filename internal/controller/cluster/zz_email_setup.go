// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	routingaddress "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/email/routingaddress"
	routingcatchall "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/email/routingcatchall"
	routingdns "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/email/routingdns"
	routingrule "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/email/routingrule"
	routingsettings "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/email/routingsettings"
	securityblocksender "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/email/securityblocksender"
	securityimpersonationregistry "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/email/securityimpersonationregistry"
	securitytrusteddomains "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/email/securitytrusteddomains"
)

// Setup_email creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_email(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		routingaddress.Setup,
		routingcatchall.Setup,
		routingdns.Setup,
		routingrule.Setup,
		routingsettings.Setup,
		securityblocksender.Setup,
		securityimpersonationregistry.Setup,
		securitytrusteddomains.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_email creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_email(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		routingaddress.SetupGated,
		routingcatchall.SetupGated,
		routingdns.SetupGated,
		routingrule.SetupGated,
		routingsettings.SetupGated,
		securityblocksender.SetupGated,
		securityimpersonationregistry.SetupGated,
		securitytrusteddomains.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
