package scraper

import (
	"math/rand"
)

var (
	userAgents = []string{
		"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
		"GalaxyBot/1.0 (http://www.galaxy.com/galaxybot.html)",
		"Googlebot-Image/1.0",
	}
)

func RandomUserAgent() string {
	return userAgents[rand.Intn(len(userAgents))]
}
