package main

import (
	"fmt"
	"html/template"
	"math"
	"net/url"
	"slices"
	"time"

	"github.com/fogleman/gg"
	"github.com/kvizdos/easyblog/builder"
	"github.com/kvizdos/easyblog/entrypoint"
)

func rangeEvents(posts builder.PostList) builder.PostList {
	p := slices.Clone(posts)
	slices.Reverse(p)

	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic(err)
	}

	// Midnight local time
	now := time.Now().In(loc)
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	out := builder.PostList{}
	for _, post := range p {
		date, err := time.ParseInLocation("01/02/2006", post.Date, loc)
		if err != nil {
			panic(err)
		}
		date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc)

		if date.Before(now) {
			continue
		}

		out = append(out, post)
	}
	return out
}

func getSoonEvents(posts builder.PostList) builder.PostList {
	p := slices.Clone(posts)
	slices.Reverse(p)

	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic(err)
	}

	// Ensure 'now' is midnight in local time
	now := time.Now().In(loc)
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)

	weekFromNow := now.AddDate(0, 0, 7)

	out := builder.PostList{}
	for _, post := range p {
		date, err := time.ParseInLocation("01/02/2006", post.Date, loc)
		if err != nil {
			panic(err)
		}

		// Ensure date is also midnight in local time
		date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc)

		// Keep only events today through 7 days out
		if !date.Before(now) && !date.After(weekFromNow) {
			out = append(out, post)
		}
	}
	return out
}

func findMaxFontSize(dc *gg.Context, fontPath, text string, maxWidth, maxHeight float64) float64 {
	minSize := 10.0
	maxSize := 300.0
	var bestSize float64

	for i := 0; i < 10; i++ { // 10 iterations = good enough precision
		mid := (minSize + maxSize) / 2
		if err := dc.LoadFontFace(fontPath, mid); err != nil {
			break
		}
		lines := dc.WordWrap(text, maxWidth)
		height := mid * 1.15 * float64(len(lines))

		if height > maxHeight {
			maxSize = mid
		} else {
			bestSize = mid
			minSize = mid
		}
	}

	return bestSize
}

func GenerateOG(postTitle string, outPath string, config builder.OGImageConfig) {
	const width = 1200
	const height = 630 // Common OG image size

	dc := gg.NewContext(width, height)
	// Set background color
	dc.SetRGB(244.0/255.0, 239.0/255.0, 229.0/255.0)
	dc.Clear()

	padding := 40.0

	boxWidth := float64(width) - 2*padding
	boxHeight := float64(height) - 200 // leave space for image/footer etc

	// Load a custom font
	fontSize := findMaxFontSize(dc, config.FontPath, postTitle, boxWidth, boxHeight)
	if err := dc.LoadFontFace(config.FontPath, fontSize); err != nil {
		panic(err)
	}

	// The text to be drawn (top of the image)
	text := postTitle
	dc.SetRGB(46/255.0, 74/255.0, 98/255.0)
	wrapped := dc.WordWrap(text, float64(width)-2*padding)
	lineHeight := fontSize * 1.15
	totalTextHeight := lineHeight * float64(len(wrapped))
	footerImageHeight := 80.0
	startY := ((float64(height)-footerImageHeight)-totalTextHeight)/2 + lineHeight/2
	centerX := float64(width) / 2
	for i, line := range wrapped {
		y := startY + float64(i)*lineHeight
		dc.DrawStringAnchored(line, centerX, y, 0.5, 0.5) // center align horizontally
	}

	// Draw the circular headshot image
	img, err := gg.LoadImage(config.IconPath)
	if err != nil {
		panic(err)
	}

	targetX := float64(width) / 2

	imgW := float64(img.Bounds().Dx())
	imgH := float64(img.Bounds().Dy())
	scale := 100.0 / math.Max(imgW, imgH) // fit into 100px size

	targetY := float64(height) - (imgH * scale) - padding

	dc.Push()
	dc.Translate(targetX, targetY)
	dc.Scale(scale, scale)
	dc.Translate(-imgW/2, 0) // center it horizontally
	dc.DrawImage(img, 0, 0)
	dc.Pop()
	// Save the result
	dc.SavePNG(outPath)
}

func main() {
	entrypoint.Start(entrypoint.EasyblogOpts{
		CustomOGGenerator: GenerateOG,
		CustomFuncs: template.FuncMap{
			"RangeEvents":   rangeEvents,
			"GetSoonEvents": getSoonEvents,
			"GetEventURL": func(post builder.Post) string {
				return fmt.Sprintf("https://rva.events/post/%s", post.OGName)
			},
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
			"GetGroup": func(metadata map[string]any) string {
				if group, ok := metadata["Group"].(string); ok {
					return group
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
					layout := "01/02/2006"
					d, err := time.Parse(layout, rawDate)
					if err != nil {
						panic(err)
					}

					loc, err := time.LoadLocation("America/New_York")
					if err != nil {
						panic(err)
					}

					today := time.Now().In(loc)
					d = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
					today = time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)

					days := int(d.Sub(today).Hours() / 24)

					switch {
					case days == 0:
						return "TODAY"
					case days == 1:
						return "TOMORROW"
					case days > 1 && days <= 7:
						return fmt.Sprintf("IN %d DAYS (%s)", days, d.Format("Mon"))
					default:
						return d.Format("Monday, Jan 2, 2006")
					}
				}
				return "Unknown"
			},
		},
	})
}
