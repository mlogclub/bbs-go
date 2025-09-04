package urls

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

func (builder *UrlBuilder) AddQuery(name, value string) *UrlBuilder {
	builder.query.Add(name, value)
	return builder
}

func (builder *UrlBuilder) AddQueries(queries map[string]string) *UrlBuilder {
	for name, value := range queries {
		builder.AddQuery(name, value)
	}
	return builder
}

func (builder *UrlBuilder) GetQuery() url.Values {
	return builder.query
}

func (builder *UrlBuilder) GetURL() *url.URL {
	return builder.u
}

func (builder *UrlBuilder) Build() *url.URL {
	builder.u.RawQuery = builder.query.Encode()
	return builder.u
}

func (builder *UrlBuilder) BuildStr() string {
	return builder.Build().String()
}
