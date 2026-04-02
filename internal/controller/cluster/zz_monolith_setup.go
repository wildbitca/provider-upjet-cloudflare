// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	rule "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/access/rule"
	account "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/account/account"
	dnssettings "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/account/dnssettings"
	dnssettingsinternalview "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/account/dnssettingsinternalview"
	member "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/account/member"
	subscription "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/account/subscription"
	token "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/account/token"
	addressmap "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/address/addressmap"
	searchinstance "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/ai/searchinstance"
	searchtoken "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/ai/searchtoken"
	shield "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/api/shield"
	shielddiscoveryoperation "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/api/shielddiscoveryoperation"
	shieldoperation "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/api/shieldoperation"
	shieldoperationschemavalidationsettings "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/api/shieldoperationschemavalidationsettings"
	shieldschema "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/api/shieldschema"
	shieldschemavalidationsettings "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/api/shieldschemavalidationsettings"
	tokenapi "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/api/token"
	smartrouting "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/argo/smartrouting"
	tieredcaching "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/argo/tieredcaching"
	originpulls "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/authenticated/originpulls"
	originpullscertificate "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/authenticated/originpullscertificate"
	originpullshostnamecertificate "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/authenticated/originpullshostnamecertificate"
	originpullssettings "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/authenticated/originpullssettings"
	management "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/bot/management"
	ipprefix "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/byo/ipprefix"
	sfuapp "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/calls/sfuapp"
	turnapp "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/calls/turnapp"
	pack "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/certificate/pack"
	connectorrules "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/cloud/connectorrules"
	filter "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/cloudflare/filter"
	healthcheck "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/cloudflare/healthcheck"
	image "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/cloudflare/image"
	list "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/cloudflare/list"
	organization "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/cloudflare/organization"
	queue "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/cloudflare/queue"
	ruleset "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/cloudflare/ruleset"
	snippet "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/cloudflare/snippet"
	snippets "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/cloudflare/snippets"
	stream "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/cloudflare/stream"
	user "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/cloudflare/user"
	worker "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/cloudflare/worker"
	workflow "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/cloudflare/workflow"
	onerequest "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/cloudforce/onerequest"
	onerequestasset "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/cloudforce/onerequestasset"
	onerequestmessage "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/cloudforce/onerequestmessage"
	onerequestpriority "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/cloudforce/onerequestpriority"
	directoryservice "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/connectivity/directoryservice"
	scanning "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/content/scanning"
	scanningexpression "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/content/scanningexpression"
	hostname "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/custom/hostname"
	hostnamefallbackorigin "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/custom/hostnamefallbackorigin"
	pages "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/custom/pages"
	ssl "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/custom/ssl"
	database "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/d1/database"
	firewall "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/dns/firewall"
	record "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/dns/record"
	zonetransfersacl "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/dns/zonetransfersacl"
	zonetransfersincoming "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/dns/zonetransfersincoming"
	zonetransfersoutgoing "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/dns/zonetransfersoutgoing"
	zonetransferspeer "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/dns/zonetransferspeer"
	zonetransferstsig "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/dns/zonetransferstsig"
	routingaddress "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/email/routingaddress"
	routingcatchall "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/email/routingcatchall"
	routingdns "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/email/routingdns"
	routingrule "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/email/routingrule"
	routingsettings "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/email/routingsettings"
	securityblocksender "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/email/securityblocksender"
	securityimpersonationregistry "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/email/securityimpersonationregistry"
	securitytrusteddomains "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/email/securitytrusteddomains"
	rulefirewall "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/firewall/rule"
	tlssetting "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/hostname/tlssetting"
	config "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/hyperdrive/config"
	variant "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/image/variant"
	certificate "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/keyless/certificate"
	credentialcheck "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/leaked/credentialcheck"
	credentialcheckrule "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/leaked/credentialcheckrule"
	item "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/list/item"
	balancer "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/load/balancer"
	balancermonitor "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/load/balancermonitor"
	balancerpool "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/load/balancerpool"
	retention "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/logpull/retention"
	job "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/logpush/job"
	ownershipchallenge "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/logpush/ownershipchallenge"
	networkmonitoringconfiguration "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/magic/networkmonitoringconfiguration"
	networkmonitoringrule "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/magic/networkmonitoringrule"
	transitconnector "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/magic/transitconnector"
	transitsite "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/magic/transitsite"
	transitsiteacl "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/magic/transitsiteacl"
	transitsitelan "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/magic/transitsitelan"
	transitsitewan "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/magic/transitsitewan"
	wangretunnel "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/magic/wangretunnel"
	wanipsectunnel "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/magic/wanipsectunnel"
	wanstaticroute "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/magic/wanstaticroute"
	transforms "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/managed/transforms"
	certificatemtls "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/mtls/certificate"
	policy "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/notification/policy"
	policywebhooks "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/notification/policywebhooks"
	scheduledtest "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/observatory/scheduledtest"
	profile "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/organization/profile"
	cacertificate "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/origin/cacertificate"
	rulepage "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/page/rule"
	shieldpolicy "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/page/shieldpolicy"
	domain "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/pages/domain"
	project "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/pages/project"
	providerconfig "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/providerconfig"
	consumer "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/queue/consumer"
	bucket "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/r2/bucket"
	bucketcors "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/r2/bucketcors"
	bucketeventnotification "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/r2/bucketeventnotification"
	bucketlifecycle "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/r2/bucketlifecycle"
	bucketlock "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/r2/bucketlock"
	bucketsippy "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/r2/bucketsippy"
	customdomain "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/r2/customdomain"
	manageddomain "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/r2/manageddomain"
	limit "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/rate/limit"
	hostnameregional "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/regional/hostname"
	tieredcache "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/regional/tieredcache"
	domainregistrar "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/registrar/domain"
	validationoperationsettings "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/schema/validationoperationsettings"
	validationschemas "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/schema/validationschemas"
	validationsettings "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/schema/validationsettings"
	rules "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/snippet/rules"
	application "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/spectrum/application"
	connector "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/sso/connector"
	audiotrack "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/stream/audiotrack"
	captionlanguage "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/stream/captionlanguage"
	download "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/stream/download"
	key "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/stream/key"
	liveinput "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/stream/liveinput"
	watermark "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/stream/watermark"
	webhook "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/stream/webhook"
	cache "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/tiered/cache"
	validationconfig "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/token/validationconfig"
	validationrules "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/token/validationrules"
	tls "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/total/tls"
	widget "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/turnstile/widget"
	sslsetting "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/universal/sslsetting"
	normalizationsettings "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/url/normalizationsettings"
	agentblockingrule "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/user/agentblockingrule"
	room "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/waiting/room"
	roomevent "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/waiting/roomevent"
	roomrules "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/waiting/roomrules"
	roomsettings "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/waiting/roomsettings"
	analyticsrule "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/web/analyticsrule"
	analyticssite "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/web/analyticssite"
	hostnameweb3 "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/web3/hostname"
	version "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/worker/version"
	crontrigger "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/workers/crontrigger"
	customdomainworkers "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/workers/customdomain"
	deployment "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/workers/deployment"
	forplatformsdispatchnamespace "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/workers/forplatformsdispatchnamespace"
	kv "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/workers/kv"
	kvnamespace "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/workers/kvnamespace"
	route "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/workers/route"
	script "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/workers/script"
	scriptsubdomain "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/workers/scriptsubdomain"
	trustaccessaicontrolsmcpportal "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustaccessaicontrolsmcpportal"
	trustaccessaicontrolsmcpserver "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustaccessaicontrolsmcpserver"
	trustaccessapplication "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustaccessapplication"
	trustaccesscustompage "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustaccesscustompage"
	trustaccessgroup "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustaccessgroup"
	trustaccessidentityprovider "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustaccessidentityprovider"
	trustaccessinfrastructuretarget "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustaccessinfrastructuretarget"
	trustaccesskeyconfiguration "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustaccesskeyconfiguration"
	trustaccessmtlscertificate "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustaccessmtlscertificate"
	trustaccessmtlshostnamesettings "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustaccessmtlshostnamesettings"
	trustaccesspolicy "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustaccesspolicy"
	trustaccessservicetoken "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustaccessservicetoken"
	trustaccessshortlivedcertificate "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustaccessshortlivedcertificate"
	trustaccesstag "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustaccesstag"
	trustdevicecustomprofile "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustdevicecustomprofile"
	trustdevicecustomprofilelocaldomainfallback "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustdevicecustomprofilelocaldomainfallback"
	trustdevicedefaultprofile "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustdevicedefaultprofile"
	trustdevicedefaultprofilecertificates "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustdevicedefaultprofilecertificates"
	trustdevicedefaultprofilelocaldomainfallback "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustdevicedefaultprofilelocaldomainfallback"
	trustdevicemanagednetworks "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustdevicemanagednetworks"
	trustdevicepostureintegration "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustdevicepostureintegration"
	trustdeviceposturerule "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustdeviceposturerule"
	trustdevicesettings "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustdevicesettings"
	trustdextest "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustdextest"
	trustdlpcustomentry "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustdlpcustomentry"
	trustdlpcustomprofile "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustdlpcustomprofile"
	trustdlpdataset "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustdlpdataset"
	trustdlpentry "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustdlpentry"
	trustdlpintegrationentry "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustdlpintegrationentry"
	trustdlppredefinedentry "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustdlppredefinedentry"
	trustdlppredefinedprofile "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustdlppredefinedprofile"
	trustdnslocation "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustdnslocation"
	trustgatewaycertificate "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustgatewaycertificate"
	trustgatewaylogging "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustgatewaylogging"
	trustgatewaypolicy "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustgatewaypolicy"
	trustgatewayproxyendpoint "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustgatewayproxyendpoint"
	trustgatewaysettings "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustgatewaysettings"
	trustlist "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustlist"
	trustnetworkhostnameroute "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustnetworkhostnameroute"
	trustorganization "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustorganization"
	trustriskbehavior "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustriskbehavior"
	trustriskscoringintegration "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trustriskscoringintegration"
	trusttunnelcloudflared "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trusttunnelcloudflared"
	trusttunnelcloudflaredconfig "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trusttunnelcloudflaredconfig"
	trusttunnelcloudflaredroute "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trusttunnelcloudflaredroute"
	trusttunnelcloudflaredvirtualnetwork "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trusttunnelcloudflaredvirtualnetwork"
	trusttunnelwarpconnector "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zero/trusttunnelwarpconnector"
	cachereserve "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zone/cachereserve"
	cachevariants "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zone/cachevariants"
	dnssec "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zone/dnssec"
	dnssettingszone "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zone/dnssettings"
	hold "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zone/hold"
	lockdown "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zone/lockdown"
	setting "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zone/setting"
	subscriptionzone "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zone/subscription"
	zone "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/cluster/zone/zone"
)

// Setup_monolith creates all controllers with the supplied logger and adds them to
// the supplied manager.
func Setup_monolith(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		rule.Setup,
		account.Setup,
		dnssettings.Setup,
		dnssettingsinternalview.Setup,
		member.Setup,
		subscription.Setup,
		token.Setup,
		addressmap.Setup,
		searchinstance.Setup,
		searchtoken.Setup,
		shield.Setup,
		shielddiscoveryoperation.Setup,
		shieldoperation.Setup,
		shieldoperationschemavalidationsettings.Setup,
		shieldschema.Setup,
		shieldschemavalidationsettings.Setup,
		tokenapi.Setup,
		smartrouting.Setup,
		tieredcaching.Setup,
		originpulls.Setup,
		originpullscertificate.Setup,
		originpullshostnamecertificate.Setup,
		originpullssettings.Setup,
		management.Setup,
		ipprefix.Setup,
		sfuapp.Setup,
		turnapp.Setup,
		pack.Setup,
		connectorrules.Setup,
		filter.Setup,
		healthcheck.Setup,
		image.Setup,
		list.Setup,
		organization.Setup,
		queue.Setup,
		ruleset.Setup,
		snippet.Setup,
		snippets.Setup,
		stream.Setup,
		user.Setup,
		worker.Setup,
		workflow.Setup,
		onerequest.Setup,
		onerequestasset.Setup,
		onerequestmessage.Setup,
		onerequestpriority.Setup,
		directoryservice.Setup,
		scanning.Setup,
		scanningexpression.Setup,
		hostname.Setup,
		hostnamefallbackorigin.Setup,
		pages.Setup,
		ssl.Setup,
		database.Setup,
		firewall.Setup,
		record.Setup,
		zonetransfersacl.Setup,
		zonetransfersincoming.Setup,
		zonetransfersoutgoing.Setup,
		zonetransferspeer.Setup,
		zonetransferstsig.Setup,
		routingaddress.Setup,
		routingcatchall.Setup,
		routingdns.Setup,
		routingrule.Setup,
		routingsettings.Setup,
		securityblocksender.Setup,
		securityimpersonationregistry.Setup,
		securitytrusteddomains.Setup,
		rulefirewall.Setup,
		tlssetting.Setup,
		config.Setup,
		variant.Setup,
		certificate.Setup,
		credentialcheck.Setup,
		credentialcheckrule.Setup,
		item.Setup,
		balancer.Setup,
		balancermonitor.Setup,
		balancerpool.Setup,
		retention.Setup,
		job.Setup,
		ownershipchallenge.Setup,
		networkmonitoringconfiguration.Setup,
		networkmonitoringrule.Setup,
		transitconnector.Setup,
		transitsite.Setup,
		transitsiteacl.Setup,
		transitsitelan.Setup,
		transitsitewan.Setup,
		wangretunnel.Setup,
		wanipsectunnel.Setup,
		wanstaticroute.Setup,
		transforms.Setup,
		certificatemtls.Setup,
		policy.Setup,
		policywebhooks.Setup,
		scheduledtest.Setup,
		profile.Setup,
		cacertificate.Setup,
		rulepage.Setup,
		shieldpolicy.Setup,
		domain.Setup,
		project.Setup,
		providerconfig.Setup,
		consumer.Setup,
		bucket.Setup,
		bucketcors.Setup,
		bucketeventnotification.Setup,
		bucketlifecycle.Setup,
		bucketlock.Setup,
		bucketsippy.Setup,
		customdomain.Setup,
		manageddomain.Setup,
		limit.Setup,
		hostnameregional.Setup,
		tieredcache.Setup,
		domainregistrar.Setup,
		validationoperationsettings.Setup,
		validationschemas.Setup,
		validationsettings.Setup,
		rules.Setup,
		application.Setup,
		connector.Setup,
		audiotrack.Setup,
		captionlanguage.Setup,
		download.Setup,
		key.Setup,
		liveinput.Setup,
		watermark.Setup,
		webhook.Setup,
		cache.Setup,
		validationconfig.Setup,
		validationrules.Setup,
		tls.Setup,
		widget.Setup,
		sslsetting.Setup,
		normalizationsettings.Setup,
		agentblockingrule.Setup,
		room.Setup,
		roomevent.Setup,
		roomrules.Setup,
		roomsettings.Setup,
		analyticsrule.Setup,
		analyticssite.Setup,
		hostnameweb3.Setup,
		version.Setup,
		crontrigger.Setup,
		customdomainworkers.Setup,
		deployment.Setup,
		forplatformsdispatchnamespace.Setup,
		kv.Setup,
		kvnamespace.Setup,
		route.Setup,
		script.Setup,
		scriptsubdomain.Setup,
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
		cachereserve.Setup,
		cachevariants.Setup,
		dnssec.Setup,
		dnssettingszone.Setup,
		hold.Setup,
		lockdown.Setup,
		setting.Setup,
		subscriptionzone.Setup,
		zone.Setup,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}

// SetupGated_monolith creates all controllers with the supplied logger and adds them to
// the supplied manager gated.
func SetupGated_monolith(mgr ctrl.Manager, o controller.Options) error {
	for _, setup := range []func(ctrl.Manager, controller.Options) error{
		rule.SetupGated,
		account.SetupGated,
		dnssettings.SetupGated,
		dnssettingsinternalview.SetupGated,
		member.SetupGated,
		subscription.SetupGated,
		token.SetupGated,
		addressmap.SetupGated,
		searchinstance.SetupGated,
		searchtoken.SetupGated,
		shield.SetupGated,
		shielddiscoveryoperation.SetupGated,
		shieldoperation.SetupGated,
		shieldoperationschemavalidationsettings.SetupGated,
		shieldschema.SetupGated,
		shieldschemavalidationsettings.SetupGated,
		tokenapi.SetupGated,
		smartrouting.SetupGated,
		tieredcaching.SetupGated,
		originpulls.SetupGated,
		originpullscertificate.SetupGated,
		originpullshostnamecertificate.SetupGated,
		originpullssettings.SetupGated,
		management.SetupGated,
		ipprefix.SetupGated,
		sfuapp.SetupGated,
		turnapp.SetupGated,
		pack.SetupGated,
		connectorrules.SetupGated,
		filter.SetupGated,
		healthcheck.SetupGated,
		image.SetupGated,
		list.SetupGated,
		organization.SetupGated,
		queue.SetupGated,
		ruleset.SetupGated,
		snippet.SetupGated,
		snippets.SetupGated,
		stream.SetupGated,
		user.SetupGated,
		worker.SetupGated,
		workflow.SetupGated,
		onerequest.SetupGated,
		onerequestasset.SetupGated,
		onerequestmessage.SetupGated,
		onerequestpriority.SetupGated,
		directoryservice.SetupGated,
		scanning.SetupGated,
		scanningexpression.SetupGated,
		hostname.SetupGated,
		hostnamefallbackorigin.SetupGated,
		pages.SetupGated,
		ssl.SetupGated,
		database.SetupGated,
		firewall.SetupGated,
		record.SetupGated,
		zonetransfersacl.SetupGated,
		zonetransfersincoming.SetupGated,
		zonetransfersoutgoing.SetupGated,
		zonetransferspeer.SetupGated,
		zonetransferstsig.SetupGated,
		routingaddress.SetupGated,
		routingcatchall.SetupGated,
		routingdns.SetupGated,
		routingrule.SetupGated,
		routingsettings.SetupGated,
		securityblocksender.SetupGated,
		securityimpersonationregistry.SetupGated,
		securitytrusteddomains.SetupGated,
		rulefirewall.SetupGated,
		tlssetting.SetupGated,
		config.SetupGated,
		variant.SetupGated,
		certificate.SetupGated,
		credentialcheck.SetupGated,
		credentialcheckrule.SetupGated,
		item.SetupGated,
		balancer.SetupGated,
		balancermonitor.SetupGated,
		balancerpool.SetupGated,
		retention.SetupGated,
		job.SetupGated,
		ownershipchallenge.SetupGated,
		networkmonitoringconfiguration.SetupGated,
		networkmonitoringrule.SetupGated,
		transitconnector.SetupGated,
		transitsite.SetupGated,
		transitsiteacl.SetupGated,
		transitsitelan.SetupGated,
		transitsitewan.SetupGated,
		wangretunnel.SetupGated,
		wanipsectunnel.SetupGated,
		wanstaticroute.SetupGated,
		transforms.SetupGated,
		certificatemtls.SetupGated,
		policy.SetupGated,
		policywebhooks.SetupGated,
		scheduledtest.SetupGated,
		profile.SetupGated,
		cacertificate.SetupGated,
		rulepage.SetupGated,
		shieldpolicy.SetupGated,
		domain.SetupGated,
		project.SetupGated,
		providerconfig.SetupGated,
		consumer.SetupGated,
		bucket.SetupGated,
		bucketcors.SetupGated,
		bucketeventnotification.SetupGated,
		bucketlifecycle.SetupGated,
		bucketlock.SetupGated,
		bucketsippy.SetupGated,
		customdomain.SetupGated,
		manageddomain.SetupGated,
		limit.SetupGated,
		hostnameregional.SetupGated,
		tieredcache.SetupGated,
		domainregistrar.SetupGated,
		validationoperationsettings.SetupGated,
		validationschemas.SetupGated,
		validationsettings.SetupGated,
		rules.SetupGated,
		application.SetupGated,
		connector.SetupGated,
		audiotrack.SetupGated,
		captionlanguage.SetupGated,
		download.SetupGated,
		key.SetupGated,
		liveinput.SetupGated,
		watermark.SetupGated,
		webhook.SetupGated,
		cache.SetupGated,
		validationconfig.SetupGated,
		validationrules.SetupGated,
		tls.SetupGated,
		widget.SetupGated,
		sslsetting.SetupGated,
		normalizationsettings.SetupGated,
		agentblockingrule.SetupGated,
		room.SetupGated,
		roomevent.SetupGated,
		roomrules.SetupGated,
		roomsettings.SetupGated,
		analyticsrule.SetupGated,
		analyticssite.SetupGated,
		hostnameweb3.SetupGated,
		version.SetupGated,
		crontrigger.SetupGated,
		customdomainworkers.SetupGated,
		deployment.SetupGated,
		forplatformsdispatchnamespace.SetupGated,
		kv.SetupGated,
		kvnamespace.SetupGated,
		route.SetupGated,
		script.SetupGated,
		scriptsubdomain.SetupGated,
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
		cachereserve.SetupGated,
		cachevariants.SetupGated,
		dnssec.SetupGated,
		dnssettingszone.SetupGated,
		hold.SetupGated,
		lockdown.SetupGated,
		setting.SetupGated,
		subscriptionzone.SetupGated,
		zone.SetupGated,
	} {
		if err := setup(mgr, o); err != nil {
			return err
		}
	}
	return nil
}
