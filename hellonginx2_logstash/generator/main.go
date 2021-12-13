package main

import (
	"fmt"
	"strings"
	"time"
	"net/http"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/caarlos0/env/v6"
)

type config struct {
	Rate             float32 `env:"RATE" envDefault:"1"`
	IPv4Percent      int     `env:"IPV4_PERCENT" envDefault:"100"`
	StatusOkPercent  int     `env:"STATUS_OK_PERCENT" envDefault:"80"`
	PathMinLength    int     `env:"PATH_MIN" envDefault:"1"`
	PathMaxLength    int     `env:"PATH_MAX" envDefault:"5"`
	PercentageGet    int     `env:"GET_PERCENT" envDefault:"60"`
	PercentagePost   int     `env:"POST_PERCENT" envDefault:"30"`
	PercentagePut    int     `env:"PUT_PERCENT" envDefault:"0"`
	PercentagePatch  int     `env:"PATCH_PERCENT" envDefault:"0"`
	PercentageDelete int     `env:"DELETE_PERCENT" envDefault:"0"`
}

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}
	checkMinMax(&cfg.PathMinLength, &cfg.PathMaxLength)

	ticker := time.NewTicker(time.Second / time.Duration(500))

	gofakeit.Seed(time.Now().UnixNano())

	var ip, httpMethod, path, httpVersion, referrer, userAgent string
	var statusCode, bodyBytesSent int
	var timeLocal time.Time

	httpVersion = "HTTP/1.1"
	referrer = "-"

	for range ticker.C {
		timeLocal = time.Now()

		ip = weightedIPVersion(cfg.IPv4Percent)
		httpMethod = weightedHTTPMethod(cfg.PercentageGet, cfg.PercentagePost, cfg.PercentagePut, cfg.PercentagePatch, cfg.PercentageDelete)
		path = randomPath(cfg.PathMinLength, cfg.PathMaxLength)
		statusCode = weightedStatusCode(cfg.StatusOkPercent)
		bodyBytesSent = realisticBytesSent(statusCode)
		userAgent = gofakeit.UserAgent()

		if httpMethod == "GET" {
			_, err := http.Get("http://nginx")
			if err != nil {
				// handle err
				println("get error")
				fmt.Println(err)
			}
		} 
		/*

		else if httpMethod == "POST" {
			reqBody := bytes.NewBufferString("Post plain text")
			resp, err := http.Post("nginx", "text/plain", reqBody)
			if err != nil {
				// handle err
				println("post error")
			}
			defer resp.Body.Close()
		}

		 else if httpMethod == "PUT" {
			resp, err := http.Put("127.0.0.1")
			if err != nil {
				// handle err
			}
			defer resp.Body.Close()
		} else if httpMethod == "PATCH" {
			resp, err := http.Patch("127.0.0.1")
			if err != nil {
				// handle err
			}
			defer resp.Body.Close()
		} else if httpMethod == "DELETE" {
			resp, err := http.Delete("127.0.0.1")
			if err != nil {
				// handle err
			}
			defer resp.Body.Close()
		} 
		*/
		
		fmt.Printf("%s - - [%s] \"%s %s %s\" %v %v \"%s\" \"%s\"\n", ip, timeLocal.Format("02/Jan/2006:15:04:05 -0700"), httpMethod, path, httpVersion, statusCode, bodyBytesSent, referrer, userAgent)
	}
}

func realisticBytesSent(statusCode int) int {
	if statusCode != 200 {
		return gofakeit.Number(30, 120)
	}

	return gofakeit.Number(800, 3100)
}

func weightedStatusCode(percentageOk int) int {
	roll := gofakeit.Number(0, 100)
	if roll <= percentageOk {
		return 200
	}

	return gofakeit.HTTPStatusCodeSimple()
}

func weightedHTTPMethod(percentageGet, percentagePost, percentagePut, percentagePatch, percentageDelete int) string {
	if percentageGet+percentagePost >= 100 {
		panic("HTTP method percentages add up to more than 100%")
	}

	roll := gofakeit.Number(0, 100)
	if roll <= percentageGet {
		return "GET"
	} 
/*
else if roll <= percentagePost {
		return "POST"
	} 

else if roll <= percentagePut {
		return "PUT"
	} else if roll <= percentagePatch {
		return "PATCH"
	} else if roll <= percentageDelete {
		return "DELETE"
	}
*/	
	return "GET"
}

func weightedIPVersion(percentageIPv4 int) string {
	roll := gofakeit.Number(0, 100)
	if roll <= percentageIPv4 {
		return gofakeit.IPv4Address()
	} else {
		return gofakeit.IPv6Address()
	}
}

func randomPath(min, max int) string {
	var path strings.Builder
	length := gofakeit.Number(min, max)

	path.WriteString("/")

	for i := 0; i < length; i++ {
		if i > 0 {
			path.WriteString(gofakeit.RandomString([]string{"-", "-", "_", "%20", "/", "/", "/"}))
		}
		path.WriteString(gofakeit.BuzzWord())
	}

	path.WriteString(gofakeit.RandomString([]string{".hmtl", ".php", ".htm", ".jpg", ".png", ".gif", ".svg", ".css", ".js"}))

	result := path.String()
	return strings.Replace(result, " ", "%20", -1)
}

func checkMinMax(min, max *int) {
	if *min < 1 {
		*min = 1
	}
	if *max < 1 {
		*max = 1
	}
	if *min > *max {
		*min, *max = *max, *min
	}
}
