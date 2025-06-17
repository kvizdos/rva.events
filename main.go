package main

import (
	"fmt"
	"html/template"
	"net/url"
	"slices"
	"time"

	"github.com/kvizdos/easyblog/builder"
	"github.com/kvizdos/easyblog/entrypoint"
)

func rangeEvents(posts builder.PostList) builder.PostList {
	slices.Reverse(posts)

	now := time.Now().UTC().Truncate(24 * time.Hour)

	out := builder.PostList{}
	for _, post := range posts {
		date, err := time.Parse("01/02/2006", post.Date)
		if err != nil {
			panic(err)
		}

		if date.Before(now) {
			continue
		}

		out = append(out, post)
	}
	return out
}

func getSoonEvents(posts builder.PostList) builder.PostList {
	slices.Reverse(posts)

	now := time.Now().UTC().Truncate(24 * time.Hour)
	weekFromNow := now.AddDate(0, 0, 7)

	out := builder.PostList{}
	for _, post := range posts {
		date, err := time.Parse("01/02/2006", post.Date)
		if err != nil {
			panic(err)
		}

		if date.Before(now) || date.After(weekFromNow) {
			continue
		}

		out = append(out, post)
	}
	return out
}

func main() {
	entrypoint.Start(entrypoint.EasyblogOpts{
		CustomFuncs: template.FuncMap{
			"RangeEvents":   rangeEvents,
			"GetSoonEvents": getSoonEvents,
			"GetExternalURL": func(rawURL string) string {
				url, err := url.Parse(rawURL)

				if err != nil {
					panic(err)
				}

				return url.Hostname()
			},
			"GetTime": func(metadata map[string]any) string {
				if time, ok := metadata["Time"].(string); ok {
					return time
				}
				return "Unknown"
			},
			"GetLocation": func(metadata map[string]any) string {
				if location, ok := metadata["Location"].(string); ok {
					return location
				}
				return "Unknown"
			},
			"GetPrice": func(metadata map[string]any) string {
				if price, ok := metadata["Price"].(string); ok {
					return price
				}
				return "Unknown"
			},
			"GetDate": func(metadata map[string]any) string {
				if rawDate, ok := metadata["Date"].(string); ok {
					return rawDate
				}
				return "Unknown"
			},
			"GetRelativeDate": func(metadata map[string]any) string {
				if rawDate, ok := metadata["Date"].(string); ok {
					d, err := time.Parse("01/02/2006", rawDate)
					if err != nil {
						panic(err)
					}

					today := time.Now().Truncate(24 * time.Hour)
					diff := d.Sub(today).Hours() / 24
					days := int(diff)

					switch {
					case days == 0:
						return "TODAY"
					case days == 1:
						return "TOMORROW"
					case days > 1 && days <= 7:
						return fmt.Sprintf("%d DAYS", days)
					default:
						return d.Format("Jan 2, 2006")
					}
				}
				return "Unknown"
			},
		},
	})
}
