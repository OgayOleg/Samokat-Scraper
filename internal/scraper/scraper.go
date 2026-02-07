package scraper

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/stealth"

	"samokat-scraper/internal/models"
	"samokat-scraper/internal/utils"
)

type Config struct {
	CategoryURL string
	OutDir      string
	Proxy       string
	APIURL      string
	AuthHeader  string
}

func Run(cfg Config) error {
	if err := os.MkdirAll(cfg.OutDir, 0755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	l := launcher.New().Headless(true).Devtools(false)

	if cfg.Proxy != "" {
		log.Printf("Прокси: %s", cfg.Proxy)
		proxyParts := strings.Split(cfg.Proxy, "@")
		if len(proxyParts) == 2 {
			ipPort := proxyParts[1]
			authParts := strings.Split(strings.TrimPrefix(proxyParts[0], "http://"), ":")
			if len(authParts) == 2 {
				l = l.Proxy(ipPort)
				log.Printf("Аутентификация: %s:%s@%s", authParts[0], authParts[1], ipPort)
			} else {
				l = l.Proxy(cfg.Proxy)
			}
		} else {
			l = l.Proxy(cfg.Proxy)
		}
	}
	u := l.MustLaunch()

	browser := rod.New().ControlURL(u).Timeout(120 * time.Second).MustConnect()
	defer browser.MustClose()

	if strings.Contains(cfg.Proxy, "@") {
		proxyParts := strings.Split(cfg.Proxy, "@")
		if len(proxyParts) == 2 {
			authParts := strings.Split(strings.TrimPrefix(proxyParts[0], "http://"), ":")
			if len(authParts) == 2 {
				go browser.MustHandleAuth(authParts[0], authParts[1])()
				time.Sleep(3 * time.Second)
			}
		}
	}

	page := stealth.MustPage(browser)

	page.MustNavigate(cfg.CategoryURL)
	page.MustWaitLoad()
	time.Sleep(4 * time.Second)

	res := page.MustEval(`() => {
        return fetch("` + cfg.APIURL + `", {
            method: "POST",
            headers: {
                "content-type": "application/json",
                "authorization": "` + cfg.AuthHeader + `",
                "x-application-platform": "web",
            },
            body: JSON.stringify({})
        }).then(r => r.text()).catch(e => e.message);
    }`)

	raw := res.String()
	if len(raw) < 10 || strings.Contains(strings.ToLower(raw), "error") {
		return fmt.Errorf("bad api response: %s", raw)
	}

	var dto models.SamokatResponseDTO
	if err := json.Unmarshal([]byte(raw), &dto); err != nil {
		return fmt.Errorf("unmarshal json: %w", err)
	}

	totalFiles := 0
	totalProducts := 0

	for _, cat := range dto.Categories {
		if len(cat.Products) == 0 {
			continue
		}

		filename := filepath.Join(cfg.OutDir, utils.Slugify(cat.Name)+".txt")
		rows := [][]string{{"name", "price", "url"}}

		for _, product := range cat.Products {
			var price float64
			if product.Prices != nil {
				price = float64(product.Prices.Current) / 100
			}

			rows = append(rows, []string{
				product.Name,
				utils.FormatPrice(price),
				"https://samokat.ru/product/" + product.Slug,
			})
		}

		if err := utils.SaveToTXT(filename, rows); err != nil {
			log.Printf("save txt %s: %v", filename, err)
			continue
		}

		totalFiles++
		totalProducts += len(cat.Products)
	}

	log.Printf("готово: %d файлов, %d товаров в %s/", totalFiles, totalProducts, cfg.OutDir)
	return nil
}
