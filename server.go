package main

import (
	"bufio"
	"errors"
	"html/template"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type ServerConfig struct {
	Title               string
	Mode                string
	TrustedProxies      []string
	Templates           string
	JavaScript          string
	CSS                 string
	Ressources          string
	Address             string
	Port                string
	CleaningInterval    time.Duration
	StatisticInterval   time.Duration
	InactiveDeleteDelay time.Duration
	Colors              map[string]string
	Footer              template.HTML
}

func initServer(router *gin.Engine, config ServerConfig) {
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))
	router.SetTrustedProxies(config.TrustedProxies)
	router.LoadHTMLGlob(strings.TrimRight(config.Templates, "\\/") + "/*")
}

func setStaticRoutes(router *gin.Engine, config ServerConfig) {
	router.Static("/js", config.JavaScript)
	router.Static("/css", config.CSS)
	router.Static("/res", config.Ressources)
	router.StaticFile("/favicon.ico", config.Ressources+"/favicon.png")
}

func setGetRoutes(router *gin.Engine) {
	router.GET("/", viewHome)
	router.GET("/stats", viewStats)
	router.GET("/error", viewError)
	router.GET("/room/:id", viewGame)
	router.GET("/rolls/:id", getAllRolls)
	router.GET("/rolls/:id/:roll_nbr", getRolls)
}

func setPostRoutes(router *gin.Engine) {
	router.POST("/room/:id", rollDice)
	router.POST("/roomSettings", changeRoomSettings)
	router.POST("/", addRoomHandler)
}

func serve(config ServerConfig) {

	gin.SetMode(config.Mode)
	router := gin.New()
	initServer(router, config)
	setStaticRoutes(router, config)
	setGetRoutes(router)
	setPostRoutes(router)
	address := config.Address + ":" + config.Port
	router.Run(address)
}

func (config *ServerConfig) loadConfig(file string) {

	config.setDefaultValues()

	configFile, err := os.Open(file)
	if err == nil {
		defer configFile.Close()
		scanner := bufio.NewScanner(configFile)
		re := regexp.MustCompile(`^\s*([^#]\w*):\s*(([\w ,.-:]|(#\w))*)(#.*)?$`)
		for scanner.Scan() {
			matches := re.FindStringSubmatch(scanner.Text())
			if len(matches) > 2 {
				config.setValue(matches[1], matches[2])
			}
		}
	}
	config.setFooter()
}

func (config *ServerConfig) setDefaultValues() {

	config.Title = getPageTitle("", "")
	config.Address = getAddress("", "")
	config.Port = getPort("", "")
	config.Mode = getMode("", "")

	config.JavaScript = getJS("", "")
	config.CSS = getCSS("", "")
	config.Templates = getTemplates("", "")
	config.Ressources = getRessources("", "")

	config.InactiveDeleteDelay = getInactiveDeleteDelay("", "")
	config.CleaningInterval = getCleaningInterval("", "")
	config.StatisticInterval = getStatisticInterval("", "")
}

func (config *ServerConfig) setValue(key string, value string) bool {

	trimChar := " \t"
	values := strings.Split(value, ",")
	n := len(values)
	if n == 0 {
		return false
	}
	switch key {
	case "trustedProxies":
		config.TrustedProxies = make([]string, n)
		for i := 0; i < n; i++ {
			config.TrustedProxies[i] = strings.Trim(values[i], trimChar)
		}
	case "title":
		config.Title = getPageTitle(values[0], trimChar)
	case "address":
		config.Address = getAddress(values[0], trimChar)
	case "port":
		config.Port = getPort(values[0], trimChar)
	case "mode":
		config.Mode = getMode(values[0], trimChar)

	case "jsDir":
		config.JavaScript = getJS(values[0], trimChar)
	case "cssDir":
		config.CSS = getCSS(values[0], trimChar)
	case "templateDir":
		config.Templates = getTemplates(values[0], trimChar)
	case "ressourceDir":
		config.Ressources = getRessources(values[0], trimChar)

	case "inactiveDeleteDelay":
		config.InactiveDeleteDelay = getInactiveDeleteDelay(values[0], trimChar)
	case "cleaningInterval":
		config.CleaningInterval = getCleaningInterval(values[0], trimChar)
	case "statisticInterval":
		config.StatisticInterval = getStatisticInterval(values[0], trimChar)
	case "color":
		config.addColors(values, trimChar)
	default:
		return false
	}
	return true
}

func getPageTitle(value string, trimChar string) string {
	ret := strings.Trim(value, trimChar)
	if ret != "" {
		return ret
	} else {
		return "Dice Roller"
	}
}

func getAddress(value string, trimChar string) string {
	ret := strings.Trim(value, trimChar)
	if ret != "" {
		return ret
	} else {
		return "localhost"
	}
}

func getPort(value string, trimChar string) string {
	ret := strings.Trim(value, trimChar)
	if ret != "" {
		return ret
	} else {
		return "9000"
	}
}

func getMode(value string, trimChar string) string {
	ret := strings.Trim(value, trimChar)
	if ret != "" {
		return ret
	} else {
		return "release"
	}
}

func getJS(value string, trimChar string) string {
	ret := strings.Trim(value, trimChar)
	if ret != "" {
		return ret
	} else {
		return "./js"
	}
}

func getCSS(value string, trimChar string) string {
	ret := strings.Trim(value, trimChar)
	if ret != "" {
		return ret
	} else {
		return "./css"
	}
}

func getTemplates(value string, trimChar string) string {
	ret := strings.Trim(value, trimChar)
	if ret != "" {
		return ret
	} else {
		return "./templates"
	}
}

func getRessources(value string, trimChar string) string {
	ret := strings.Trim(value, trimChar)
	if ret != "" {
		return ret
	} else {
		return "./res"
	}
}

func (config *ServerConfig) addColors(values []string, trimChar string) {
	if config.Colors == nil {
		config.Colors = make(map[string]string)
	}
	if len(values) >= 2 {
		key := strings.Trim(values[0], trimChar)
		val := strings.Trim(values[1], trimChar)
		config.Colors[key] = val
	} else {
		config.Colors["-"] = "-"
	}
}

func (config *ServerConfig) setFooter() {
	footer, err := ioutil.ReadFile(
		strings.TrimRight(config.Ressources, "\\/") + "/footer.html")
	if err == nil {
		config.Footer = template.HTML(footer)
	}
}

func getInactiveDeleteDelay(value string, trimChar string) time.Duration {
	str := strings.Trim(value, trimChar)
	if str == "" {
		str = "4h"
	}
	ret, err := time.ParseDuration(str)
	if err == nil {
		return ret
	} else {
		panic(errors.New("error loading configuration"))
	}
}

func getCleaningInterval(value string, trimChar string) time.Duration {
	str := strings.Trim(value, trimChar)
	if str == "" {
		str = "15m"
	}
	ret, err := time.ParseDuration(str)
	if err == nil {
		return ret
	} else {
		panic(errors.New("error loading configuration"))
	}
}

func getStatisticInterval(value string, trimChar string) time.Duration {
	str := strings.Trim(value, trimChar)
	if str == "" {
		str = "1m"
	}
	ret, err := time.ParseDuration(str)
	if err == nil {
		return ret
	} else {
		panic(errors.New("error loading configuration"))
	}
}

func cleanup(config ServerConfig) {
	for {
		time.Sleep(config.CleaningInterval)
		deleteOldGames(rooms, playerIds, config.InactiveDeleteDelay)
	}
}

func runStatistics(config ServerConfig) {
	for {
		time.Sleep(config.StatisticInterval)
		stats, _ = updateStatistics(rooms, playerIds)
	}
}
