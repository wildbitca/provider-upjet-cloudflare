// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	trustaccessaicontrolsmcpportal "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustaccessaicontrolsmcpportal"
	trustaccessaicontrolsmcpserver "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustaccessaicontrolsmcpserver"
	trustaccessapplication "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustaccessapplication"
	trustaccesscustompage "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustaccesscustompage"
	trustaccessgroup "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustaccessgroup"
	trustaccessidentityprovider "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustaccessidentityprovider"
	trustaccessinfrastructuretarget "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustaccessinfrastructuretarget"
	trustaccesskeyconfiguration "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustaccesskeyconfiguration"
	trustaccessmtlscertificate "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustaccessmtlscertificate"
	trustaccessmtlshostnamesettings "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustaccessmtlshostnamesettings"
	trustaccesspolicy "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustaccesspolicy"
	trustaccessservicetoken "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustaccessservicetoken"
	trustaccessshortlivedcertificate "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustaccessshortlivedcertificate"
	trustaccesstag "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustaccesstag"
	trustdevicecustomprofile "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustdevicecustomprofile"
	trustdevicecustomprofilelocaldomainfallback "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustdevicecustomprofilelocaldomainfallback"
	trustdevicedefaultprofile "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustdevicedefaultprofile"
	trustdevicedefaultprofilecertificates "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustdevicedefaultprofilecertificates"
	trustdevicedefaultprofilelocaldomainfallback "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustdevicedefaultprofilelocaldomainfallback"
	trustdevicemanagednetworks "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustdevicemanagednetworks"
	trustdevicepostureintegration "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustdevicepostureintegration"
	trustdeviceposturerule "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustdeviceposturerule"
	trustdevicesettings "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustdevicesettings"
	trustdextest "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustdextest"
	trustdlpcustomentry "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustdlpcustomentry"
	trustdlpcustomprofile "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustdlpcustomprofile"
	trustdlpdataset "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustdlpdataset"
	trustdlpentry "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustdlpentry"
	trustdlpintegrationentry "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustdlpintegrationentry"
	trustdlppredefinedentry "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustdlppredefinedentry"
	trustdlppredefinedprofile "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustdlppredefinedprofile"
	trustdnslocation "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustdnslocation"
	trustgatewaycertificate "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustgatewaycertificate"
	trustgatewaylogging "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustgatewaylogging"
	trustgatewaypolicy "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustgatewaypolicy"
	trustgatewayproxyendpoint "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustgatewayproxyendpoint"
	trustgatewaysettings "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustgatewaysettings"
	trustlist "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustlist"
	trustnetworkhostnameroute "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustnetworkhostnameroute"
	trustorganization "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustorganization"
	trustriskbehavior "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustriskbehavior"
	trustriskscoringintegration "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trustriskscoringintegration"
	trusttunnelcloudflared "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trusttunnelcloudflared"
	trusttunnelcloudflaredconfig "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trusttunnelcloudflaredconfig"
	trusttunnelcloudflaredroute "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trusttunnelcloudflaredroute"
	trusttunnelcloudflaredvirtualnetwork "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trusttunnelcloudflaredvirtualnetwork"
	trusttunnelwarpconnector "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zero/trusttunnelwarpconnector"
)

// Setup_zero creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_zero(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		trustaccessaicontrolsmcpportal.Setup,
		trustaccessaicontrolsmcpserver.Setup,
		trustaccessapplication.Setup,
		trustaccesscustompage.Setup,
		trustaccessgroup.Setup,
		trustaccessidentityprovider.Setup,
		trustaccessinfrastructuretarget.Setup,
		trustaccesskeyconfiguration.Setup,
		trustaccessmtlscertificate.Setup,
		trustaccessmtlshostnamesettings.Setup,
		trustaccesspolicy.Setup,
		trustaccessservicetoken.Setup,
		trustaccessshortlivedcertificate.Setup,
		trustaccesstag.Setup,
		trustdevicecustomprofile.Setup,
		trustdevicecustomprofilelocaldomainfallback.Setup,
		trustdevicedefaultprofile.Setup,
		trustdevicedefaultprofilecertificates.Setup,
		trustdevicedefaultprofilelocaldomainfallback.Setup,
		trustdevicemanagednetworks.Setup,
		trustdevicepostureintegration.Setup,
		trustdeviceposturerule.Setup,
		trustdevicesettings.Setup,
		trustdextest.Setup,
		trustdlpcustomentry.Setup,
		trustdlpcustomprofile.Setup,
		trustdlpdataset.Setup,
		trustdlpentry.Setup,
		trustdlpintegrationentry.Setup,
		trustdlppredefinedentry.Setup,
		trustdlppredefinedprofile.Setup,
		trustdnslocation.Setup,
		trustgatewaycertificate.Setup,
		trustgatewaylogging.Setup,
		trustgatewaypolicy.Setup,
		trustgatewayproxyendpoint.Setup,
		trustgatewaysettings.Setup,
		trustlist.Setup,
		trustnetworkhostnameroute.Setup,
		trustorganization.Setup,
		trustriskbehavior.Setup,
		trustriskscoringintegration.Setup,
		trusttunnelcloudflared.Setup,
		trusttunnelcloudflaredconfig.Setup,
		trusttunnelcloudflaredroute.Setup,
		trusttunnelcloudflaredvirtualnetwork.Setup,
		trusttunnelwarpconnector.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_zero creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_zero(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		trustaccessaicontrolsmcpportal.SetupGated,
		trustaccessaicontrolsmcpserver.SetupGated,
		trustaccessapplication.SetupGated,
		trustaccesscustompage.SetupGated,
		trustaccessgroup.SetupGated,
		trustaccessidentityprovider.SetupGated,
		trustaccessinfrastructuretarget.SetupGated,
		trustaccesskeyconfiguration.SetupGated,
		trustaccessmtlscertificate.SetupGated,
		trustaccessmtlshostnamesettings.SetupGated,
		trustaccesspolicy.SetupGated,
		trustaccessservicetoken.SetupGated,
		trustaccessshortlivedcertificate.SetupGated,
		trustaccesstag.SetupGated,
		trustdevicecustomprofile.SetupGated,
		trustdevicecustomprofilelocaldomainfallback.SetupGated,
		trustdevicedefaultprofile.SetupGated,
		trustdevicedefaultprofilecertificates.SetupGated,
		trustdevicedefaultprofilelocaldomainfallback.SetupGated,
		trustdevicemanagednetworks.SetupGated,
		trustdevicepostureintegration.SetupGated,
		trustdeviceposturerule.SetupGated,
		trustdevicesettings.SetupGated,
		trustdextest.SetupGated,
		trustdlpcustomentry.SetupGated,
		trustdlpcustomprofile.SetupGated,
		trustdlpdataset.SetupGated,
		trustdlpentry.SetupGated,
		trustdlpintegrationentry.SetupGated,
		trustdlppredefinedentry.SetupGated,
		trustdlppredefinedprofile.SetupGated,
		trustdnslocation.SetupGated,
		trustgatewaycertificate.SetupGated,
		trustgatewaylogging.SetupGated,
		trustgatewaypolicy.SetupGated,
		trustgatewayproxyendpoint.SetupGated,
		trustgatewaysettings.SetupGated,
		trustlist.SetupGated,
		trustnetworkhostnameroute.SetupGated,
		trustorganization.SetupGated,
		trustriskbehavior.SetupGated,
		trustriskscoringintegration.SetupGated,
		trusttunnelcloudflared.SetupGated,
		trusttunnelcloudflaredconfig.SetupGated,
		trusttunnelcloudflaredroute.SetupGated,
		trusttunnelcloudflaredvirtualnetwork.SetupGated,
		trusttunnelwarpconnector.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
