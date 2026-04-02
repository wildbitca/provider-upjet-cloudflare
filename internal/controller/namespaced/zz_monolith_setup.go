// SPDX-FileCopyrightText: 2024 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/crossplane/upjet/v2/pkg/controller"

	rule "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/access/rule"
	account "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/account/account"
	dnssettings "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/account/dnssettings"
	dnssettingsinternalview "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/account/dnssettingsinternalview"
	member "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/account/member"
	subscription "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/account/subscription"
	token "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/account/token"
	addressmap "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/address/addressmap"
	searchinstance "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/ai/searchinstance"
	searchtoken "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/ai/searchtoken"
	shield "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/api/shield"
	shielddiscoveryoperation "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/api/shielddiscoveryoperation"
	shieldoperation "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/api/shieldoperation"
	shieldoperationschemavalidationsettings "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/api/shieldoperationschemavalidationsettings"
	shieldschema "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/api/shieldschema"
	shieldschemavalidationsettings "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/api/shieldschemavalidationsettings"
	tokenapi "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/api/token"
	smartrouting "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/argo/smartrouting"
	tieredcaching "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/argo/tieredcaching"
	originpulls "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/authenticated/originpulls"
	originpullscertificate "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/authenticated/originpullscertificate"
	originpullshostnamecertificate "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/authenticated/originpullshostnamecertificate"
	originpullssettings "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/authenticated/originpullssettings"
	management "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/bot/management"
	ipprefix "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/byo/ipprefix"
	sfuapp "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/calls/sfuapp"
	turnapp "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/calls/turnapp"
	pack "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/certificate/pack"
	connectorrules "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/cloud/connectorrules"
	filter "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/cloudflare/filter"
	healthcheck "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/cloudflare/healthcheck"
	image "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/cloudflare/image"
	list "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/cloudflare/list"
	organization "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/cloudflare/organization"
	queue "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/cloudflare/queue"
	ruleset "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/cloudflare/ruleset"
	snippet "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/cloudflare/snippet"
	snippets "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/cloudflare/snippets"
	stream "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/cloudflare/stream"
	user "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/cloudflare/user"
	worker "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/cloudflare/worker"
	workflow "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/cloudflare/workflow"
	onerequest "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/cloudforce/onerequest"
	onerequestasset "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/cloudforce/onerequestasset"
	onerequestmessage "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/cloudforce/onerequestmessage"
	onerequestpriority "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/cloudforce/onerequestpriority"
	directoryservice "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/connectivity/directoryservice"
	scanning "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/content/scanning"
	scanningexpression "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/content/scanningexpression"
	hostname "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/custom/hostname"
	hostnamefallbackorigin "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/custom/hostnamefallbackorigin"
	pages "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/custom/pages"
	ssl "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/custom/ssl"
	database "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/d1/database"
	firewall "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/dns/firewall"
	record "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/dns/record"
	zonetransfersacl "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/dns/zonetransfersacl"
	zonetransfersincoming "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/dns/zonetransfersincoming"
	zonetransfersoutgoing "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/dns/zonetransfersoutgoing"
	zonetransferspeer "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/dns/zonetransferspeer"
	zonetransferstsig "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/dns/zonetransferstsig"
	routingaddress "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/email/routingaddress"
	routingcatchall "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/email/routingcatchall"
	routingdns "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/email/routingdns"
	routingrule "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/email/routingrule"
	routingsettings "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/email/routingsettings"
	securityblocksender "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/email/securityblocksender"
	securityimpersonationregistry "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/email/securityimpersonationregistry"
	securitytrusteddomains "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/email/securitytrusteddomains"
	rulefirewall "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/firewall/rule"
	tlssetting "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/hostname/tlssetting"
	config "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/hyperdrive/config"
	variant "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/image/variant"
	certificate "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/keyless/certificate"
	credentialcheck "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/leaked/credentialcheck"
	credentialcheckrule "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/leaked/credentialcheckrule"
	item "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/list/item"
	balancer "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/load/balancer"
	balancermonitor "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/load/balancermonitor"
	balancerpool "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/load/balancerpool"
	retention "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/logpull/retention"
	job "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/logpush/job"
	ownershipchallenge "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/logpush/ownershipchallenge"
	networkmonitoringconfiguration "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/magic/networkmonitoringconfiguration"
	networkmonitoringrule "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/magic/networkmonitoringrule"
	transitconnector "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/magic/transitconnector"
	transitsite "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/magic/transitsite"
	transitsiteacl "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/magic/transitsiteacl"
	transitsitelan "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/magic/transitsitelan"
	transitsitewan "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/magic/transitsitewan"
	wangretunnel "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/magic/wangretunnel"
	wanipsectunnel "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/magic/wanipsectunnel"
	wanstaticroute "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/magic/wanstaticroute"
	transforms "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/managed/transforms"
	certificatemtls "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/mtls/certificate"
	policy "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/notification/policy"
	policywebhooks "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/notification/policywebhooks"
	scheduledtest "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/observatory/scheduledtest"
	profile "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/organization/profile"
	cacertificate "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/origin/cacertificate"
	rulepage "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/page/rule"
	shieldpolicy "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/page/shieldpolicy"
	domain "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/pages/domain"
	project "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/pages/project"
	providerconfig "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/providerconfig"
	consumer "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/queue/consumer"
	bucket "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/r2/bucket"
	bucketcors "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/r2/bucketcors"
	bucketeventnotification "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/r2/bucketeventnotification"
	bucketlifecycle "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/r2/bucketlifecycle"
	bucketlock "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/r2/bucketlock"
	bucketsippy "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/r2/bucketsippy"
	customdomain "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/r2/customdomain"
	manageddomain "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/r2/manageddomain"
	limit "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/rate/limit"
	hostnameregional "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/regional/hostname"
	tieredcache "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/regional/tieredcache"
	domainregistrar "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/registrar/domain"
	validationoperationsettings "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/schema/validationoperationsettings"
	validationschemas "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/schema/validationschemas"
	validationsettings "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/schema/validationsettings"
	rules "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/snippet/rules"
	application "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/spectrum/application"
	connector "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/sso/connector"
	audiotrack "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/stream/audiotrack"
	captionlanguage "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/stream/captionlanguage"
	download "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/stream/download"
	key "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/stream/key"
	liveinput "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/stream/liveinput"
	watermark "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/stream/watermark"
	webhook "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/stream/webhook"
	cache "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/tiered/cache"
	validationconfig "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/token/validationconfig"
	validationrules "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/token/validationrules"
	tls "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/total/tls"
	widget "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/turnstile/widget"
	sslsetting "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/universal/sslsetting"
	normalizationsettings "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/url/normalizationsettings"
	agentblockingrule "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/user/agentblockingrule"
	room "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/waiting/room"
	roomevent "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/waiting/roomevent"
	roomrules "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/waiting/roomrules"
	roomsettings "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/waiting/roomsettings"
	analyticsrule "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/web/analyticsrule"
	analyticssite "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/web/analyticssite"
	hostnameweb3 "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/web3/hostname"
	version "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/worker/version"
	crontrigger "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/workers/crontrigger"
	customdomainworkers "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/workers/customdomain"
	deployment "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/workers/deployment"
	forplatformsdispatchnamespace "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/workers/forplatformsdispatchnamespace"
	kv "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/workers/kv"
	kvnamespace "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/workers/kvnamespace"
	route "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/workers/route"
	script "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/workers/script"
	scriptsubdomain "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/workers/scriptsubdomain"
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
	cachereserve "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zone/cachereserve"
	cachevariants "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zone/cachevariants"
	dnssec "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zone/dnssec"
	dnssettingszone "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zone/dnssettings"
	hold "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zone/hold"
	lockdown "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zone/lockdown"
	setting "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zone/setting"
	subscriptionzone "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zone/subscription"
	zone "github.com/wildbitca/provider-upjet-cloudflare/internal/controller/namespaced/zone/zone"
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
