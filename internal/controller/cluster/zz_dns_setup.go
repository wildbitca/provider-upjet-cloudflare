// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	firewall "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/dns/firewall"
	record "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/dns/record"
	zonetransfersacl "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/dns/zonetransfersacl"
	zonetransfersincoming "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/dns/zonetransfersincoming"
	zonetransfersoutgoing "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/dns/zonetransfersoutgoing"
	zonetransferspeer "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/dns/zonetransferspeer"
	zonetransferstsig "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/dns/zonetransferstsig"
)

// Setup_dns creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_dns(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		firewall.Setup,
		record.Setup,
		zonetransfersacl.Setup,
		zonetransfersincoming.Setup,
		zonetransfersoutgoing.Setup,
		zonetransferspeer.Setup,
		zonetransferstsig.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_dns creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_dns(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		firewall.SetupGated,
		record.SetupGated,
		zonetransfersacl.SetupGated,
		zonetransfersincoming.SetupGated,
		zonetransfersoutgoing.SetupGated,
		zonetransferspeer.SetupGated,
		zonetransferstsig.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
