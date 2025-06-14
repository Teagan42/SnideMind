package extracttags

import (
	"math"
)

type TagNode struct {
	ID          string
	Name        string
	Description string
	Vector      *[]float64
	ParentID    *string
}

var TagTree = []TagNode{
	{ID: "media", Name: "Media", Description: "Anything related to media like movies, music, or TV"},
	{ID: "media.movies", Name: "Movies", Description: "Films, cinema, or related topics", ParentID: ptr("media")},
	{ID: "media.tv", Name: "TV Shows", Description: "Television series or episodes", ParentID: ptr("media")},
	{ID: "media.new", Name: "New", Description: "Recently added, latest, or newest content", ParentID: ptr("media")},
	{ID: "home", Name: "Home Automation", Description: "Smart home devices and automation"},
	{ID: "home.security", Name: "Security", Description: "Security systems, alerts, cameras", ParentID: ptr("home")},
	{ID: "home.lighting", Name: "Lighting", Description: "Smart lighting systems", ParentID: ptr("home")},
	{ID: "home.climate", Name: "Climate Control", Description: "Thermostats, HVAC systems", ParentID: ptr("home")},
	{ID: "home.fan", Name: "Fan", Description: "Smart fans and ventilation systems", ParentID: ptr("home")},
	{ID: "home.switch", Name: "Switch", Description: "Smart switches and outlets", ParentID: ptr("home")},
	{ID: "home.sensor", Name: "Sensor", Description: "Smart sensors for various applications", ParentID: ptr("home")},
	{ID: "weather", Name: "Weather", Description: "Weather forecasts, current conditions"},
	{ID: "weather.forecast", Name: "Forecast", Description: "Weather forecasts for different regions", ParentID: ptr("weather")},
	{ID: "weather.current", Name: "Current Conditions", Description: "Current weather conditions", ParentID: ptr("weather")},
	{ID: "weather.alerts", Name: "Alerts", Description: "Weather alerts and warnings", ParentID: ptr("weather")},
	{ID: "web", Name: "Web Services", Description: "Web-based services and APIs"},
	{ID: "web.search", Name: "Search", Description: "Web search engines and APIs", ParentID: ptr("web")},
	{ID: "memory", Name: "Memory", Description: "Memory-related topics"},
	{ID: "memory.user", Name: "User Memory", Description: "User related memory", ParentID: ptr("memory")},
	{ID: "memory.projects", Name: "Projects Memory", Description: "Memory related to projects", ParentID: ptr("memory")},
	{ID: "memory.tasks", Name: "Tasks Memory", Description: "Memory related to tasks", ParentID: ptr("memory")},
	{ID: "memory.notes", Name: "Notes Memory", Description: "Memory related to notes", ParentID: ptr("memory")},
	{ID: "memory.ideas", Name: "Ideas Memory", Description: "Memory related to ideas", ParentID: ptr("memory")},
	{ID: "memory.research", Name: "Research Memory", Description: "Memory related to research", ParentID: ptr("memory")},
	{ID: "announcements", Name: "Announcements", Description: "Announcements and updates"},
	{ID: "announcements.updates", Name: "Updates", Description: "Updates to previous announcements", ParentID: ptr("announcements")},
	{ID: "announcements.alerts", Name: "Alert Announcements", Description: "Urgent alert announcements", ParentID: ptr("announcements")},
}

func ptr(s string) *string { return &s }

func CosineSimilarity(a, b []float64) float64 {
	var dot, normA, normB float64
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

type Embedding struct {
	ID        string
	Embedding []float64
}

type ScoredTag struct {
	Tag   TagNode
	Score float64
}
