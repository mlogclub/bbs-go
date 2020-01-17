package sitemap

import (
	"time"
)

const (
	ChangefreqAlways  = "always"
	ChangefreqHourly  = "hourly"
	ChangefreqDaily   = "daily"
	ChangefreqWeekly  = "weekly"
	ChangefreqMonthly = "monthly"
	ChangefreqYearly  = "yearly"
	ChangefreqNever   = "never"
)

type IndexURL struct {
	Loc     string
	Lastmod time.Time
}

func (u IndexURL) String() string {
	return "<sitemap>" +
		"<loc>" + u.Loc + "</loc>" +
		"<lastmod>" + u.Lastmod.Format(time.RFC3339) + "</lastmod>" +
		"</sitemap>"
}

type URL struct {
	Loc        string
	Lastmod    time.Time
	Changefreq string
	Priority   string
}

func (u URL) String() string {
	return "<url>" +
		"<loc>" + u.Loc + "</loc>" +
		"<lastmod>" + u.Lastmod.Format(time.RFC3339) + "</lastmod>" +
		"<changefreq>" + u.Changefreq + "</changefreq>" +
		"<priority>" + u.Priority + "</priority>" +
		"</url>"
}
