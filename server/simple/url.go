package simple

import "net/url"

type UrlBuilder struct {
	u     *url.URL
	query url.Values
}

func ParseUrl(rawUrl string) *UrlBuilder {
	ub := &UrlBuilder{}
	ub.u, _ = url.Parse(rawUrl)
	ub.query = ub.u.Query()
	return ub
}

func (this *UrlBuilder) AddQuery(name, value string) *UrlBuilder {
	this.query.Add(name, value)
	return this
}

func (this *UrlBuilder) GetQuery() url.Values {
	return this.query
}

func (this *UrlBuilder) GetURL() *url.URL {
	return this.u
}

func (this *UrlBuilder) Build() *url.URL {
	this.u.RawQuery = this.query.Encode()
	return this.u
}

func (this *UrlBuilder) BuildStr() string {
	return this.Build().String()
}
