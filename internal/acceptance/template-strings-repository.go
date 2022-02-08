package acceptance

const (
	TemplateStringProxyRepository = TemplateStringNameOnline +
		TemplateStringCleanup +
		TemplateStringHTTPClient +
		TemplateStringNegativeCache +
		TemplateStringProxy +
		TemplateStringRoutingRule +
		TemplateStringStorage +
		TemplateStringEnd
	TemplateStringHostedRepository = TemplateStringNameOnline +
		TemplateStringCleanup +
		TemplateStringComponent +
		TemplateStringStorageHosted +
		TemplateStringEnd
	TemplateStringGroupRepository = TemplateStringNameOnline +
		TemplateStringGroup +
		TemplateStringStorage +
		TemplateStringEnd
	TemplateStringGroupDeployRepository = TemplateStringNameOnline +
		TemplateStringGroupDeploy +
		TemplateStringStorage +
		TemplateStringEnd

	TemplateStringNameOnline = `
	name   = "{{ .Name }}"
	online = {{ .Online }}
`

	TemplateStringCleanup = `
{{ if .Cleanup }}
	cleanup {
		policy_names = [
		{{- range $val := .Cleanup.PolicyNames }}
			"{{ $val }}",
		{{ end -}}
		]
	}
{{ end -}}
`

	TemplateStringComponent = `
{{ if .Component }}
	component {
		proprietary_components = {{ .Component.ProprietaryComponents }}
	}
{{ end -}}
`

	TemplateStringHTTPClient = `
{{ if .HTTPClient }}
	http_client {
		auto_block = {{ .HTTPClient.AutoBlock }}
		blocked    = {{ .HTTPClient.Blocked }}

		{{ if .HTTPClient.Authentication -}}
		authentication {
			ntlm_domain = "{{ .HTTPClient.Authentication.NTLMDomain }}"
			ntlm_host   = "{{ .HTTPClient.Authentication.NTLMHost }}"
			{{ if .HTTPClient.Authentication.Password -}}
			password    = "{{ .HTTPClient.Authentication.Password }}"
			{{ end -}}
			type        = "{{ .HTTPClient.Authentication.Type }}"
			{{ if .HTTPClient.Authentication.Username -}}
			username    = "{{ .HTTPClient.Authentication.Username }}"
			{{ end -}}
		}
		{{ end -}}

		{{ if .HTTPClient.Connection -}}
		connection {
			{{ if .HTTPClient.Connection.EnableCircularRedirects -}}
			enable_circular_redirects = {{ .HTTPClient.Connection.EnableCircularRedirects }}
			{{ end -}}
			{{ if .HTTPClient.Connection.EnableCookies -}}
			enable_cookies = {{ .HTTPClient.Connection.EnableCookies }}
			{{ end -}}
			{{ if .HTTPClient.Connection.Retries -}}
			retries = {{ .HTTPClient.Connection.Retries }}
			{{ end -}}
			{{ if .HTTPClient.Connection.Timeout -}}
			timeout = {{ .HTTPClient.Connection.Timeout }}
			{{ end -}}
			{{ if .HTTPClient.Connection.UserAgentSuffix -}}
			user_agent_suffix = "{{ .HTTPClient.Connection.UserAgentSuffix }}"
			{{ end -}}
			{{ if .HTTPClient.Connection.UseTrustStore -}}
			use_trust_store = {{ .HTTPClient.Connection.UseTrustStore }}
			{{ end -}}
		}
		{{ end -}}
	}
{{ end -}}
`

	TemplateStringGroup = `
	group {
		member_names = [
		{{- range $val := .Group.MemberNames }}
			"{{ $val }}",
		{{ end -}}
		]
	}
`

	TemplateStringGroupDeploy = `
group {
	member_names = [
	{{- range $val := .Group.MemberNames }}
		"{{ $val }}",
	{{ end -}}
	]
{{- if .Group.WritableMember }}
	writable_member = "{{ .Group.WritableMember }}"
{{- end }}
}
`

	TemplateStringNegativeCache = `
{{ if .NegativeCache }}
	negative_cache {
		{{ if .NegativeCache.Enabled }}
		enabled = {{ .NegativeCache.Enabled }}
		{{ end -}}
		{{ if .NegativeCache.TTL }}
		ttl = {{ .NegativeCache.TTL }}
		{{ end }}
	}
{{ end -}}
`

	TemplateStringProxy = `
	proxy {
		remote_url = "{{ .Proxy.RemoteURL }}"
		{{ if .Proxy.ContentMaxAge }}
			content_max_age = "{{ .Proxy.ContentMaxAge }}"
		{{ end -}}
		{{ if .Proxy.MetadataMaxAge }}
			metadata_max_age = "{{ .Proxy.MetadataMaxAge }}"
		{{ end -}}
	}
`

	TemplateStringRoutingRule = `
	{{ if .RoutingRule }}
		routing_rule = nexus_routing_rule.acceptance.name
	{{ end -}}
`

	TemplateStringStorage = `
	storage {
		blob_store_name                = "{{ .Storage.BlobStoreName }}"
		strict_content_type_validation = {{ .Storage.StrictContentTypeValidation }}
	}
`

	TemplateStringStorageHosted = `
	storage {
		blob_store_name                = "{{ .Storage.BlobStoreName }}"
		strict_content_type_validation = {{ .Storage.StrictContentTypeValidation }}
		{{- if .Storage.WritePolicy }}
		write_policy                   = "{{ .Storage.WritePolicy }}"
		{{- end }}
	}
`

	TemplateStringEnd = `
}
`
)
