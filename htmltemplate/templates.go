package htmltemplate

const (
	ProviderFilterErr string = `<html><h2>Unauthorized!</h2>
	<p>Error while applying additional provider specific authorization filter:</p>
	<p><b>%v</b></p></html>`
	GenericFilterErr string = `<html><h2>Unauthorized!</h2>
	<p>Error while applying generic authorization filter:</p>
	<p><b>%v</b></p></html>`
	JwtIssueErr string = `<html><h2>Unauthorized!</h2>
	<p>Error while trying to issue a new authentication token (JWT):</p>
	<p><b>%v</b></p></html>`
	NoOriginCookieErr string = `<html><h2>Error while processing request!</h2>
	<p>After successful authentication at the Oauth provider, this CF SSO service will ` +
		`try to redirect back to the application. The cookie providing information about ` +
		`the apps domain/url seems not present:</p>
	<p><b>%v</b></p></html>`
	NoValidOriginErr string = `<html><h2>Error while processing request!</h2>
	<p>After successful authentication at the Oauth provider, this CF SSO service will ` +
		`try to redirect back to the application. The cookie providing information about ` +
		`the apps domain/url is present, but is does not seem to contain a valid domain/url:</p>
	<p><b>%v</b></p></html>`
	NoForwardUrlErr string = "<html><h2>CF did not send Fordward URL in Header 'X-Cf-Forwarded-Url'</h2></html>"
	NotSameApexErr  string = `<html><h2>The SSO Service and the App are not hosted on the same apex domain!</h2>
	<p>Apex domain for SSO Service is: <b>%v</b></p></html>`
)
