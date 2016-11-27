package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/astaxie/beego"
)

type ZalandoController struct {
	beego.Controller
}

func (c *ZalandoController) Get() {
	c.TplName = "index.tpl"
}

type ZalandoBrands struct {
	Content []struct {
		Key          string `json:"key"`
		Name         string `json:"name"`
		LogoURL      string `json:"logoUrl,omitempty"`
		LogoLargeURL string `json:"logoLargeUrl,omitempty"`
		ShopURL      string `json:"shopUrl"`
		BrandFamily  struct {
			Key     string `json:"key"`
			Name    string `json:"name"`
			ShopURL string `json:"shopUrl"`
		} `json:"brandFamily,omitempty"`
	} `json:"content"`
	TotalElements int `json:"totalElements"`
	TotalPages    int `json:"totalPages"`
	Page          int `json:"page"`
	Size          int `json:"size"`
}

type ZalandoArticles struct {
	Content []struct {
		ID              string        `json:"id"`
		ModelID         string        `json:"modelId"`
		Name            string        `json:"name"`
		ShopURL         string        `json:"shopUrl"`
		Color           string        `json:"color"`
		Available       bool          `json:"available"`
		Season          string        `json:"season"`
		SeasonYear      string        `json:"seasonYear"`
		ActivationDate  time.Time     `json:"activationDate"`
		AdditionalInfos []interface{} `json:"additionalInfos"`
		Tags            []interface{} `json:"tags"`
		Genders         []string      `json:"genders"`
		AgeGroups       []string      `json:"ageGroups"`
		Brand           struct {
			Key          string `json:"key"`
			Name         string `json:"name"`
			LogoURL      string `json:"logoUrl"`
			LogoLargeURL string `json:"logoLargeUrl"`
			BrandFamily  struct {
				Key     string `json:"key"`
				Name    string `json:"name"`
				ShopURL string `json:"shopUrl"`
			} `json:"brandFamily"`
			ShopURL string `json:"shopUrl"`
		} `json:"brand"`
		CategoryKeys []string `json:"categoryKeys"`
		Attributes   []struct {
			Name   string   `json:"name"`
			Values []string `json:"values"`
		} `json:"attributes"`
		Units []struct {
			ID    string `json:"id"`
			Size  string `json:"size"`
			Price struct {
				Currency  string  `json:"currency"`
				Value     float64 `json:"value"`
				Formatted string  `json:"formatted"`
			} `json:"price"`
			OriginalPrice struct {
				Currency  string  `json:"currency"`
				Value     float64 `json:"value"`
				Formatted string  `json:"formatted"`
			} `json:"originalPrice"`
			Available bool `json:"available"`
			Stock     int  `json:"stock"`
		} `json:"units"`
		Media struct {
			Images []struct {
				OrderNumber    int    `json:"orderNumber"`
				Type           string `json:"type"`
				ThumbnailHdURL string `json:"thumbnailHdUrl"`
				SmallURL       string `json:"smallUrl"`
				SmallHdURL     string `json:"smallHdUrl"`
				MediumURL      string `json:"mediumUrl"`
				MediumHdURL    string `json:"mediumHdUrl"`
				LargeURL       string `json:"largeUrl"`
				LargeHdURL     string `json:"largeHdUrl"`
			} `json:"images"`
		} `json:"media"`
	} `json:"content"`
	TotalElements int `json:"totalElements"`
	TotalPages    int `json:"totalPages"`
	Page          int `json:"page"`
	Size          int `json:"size"`
}

type ZalandoCategories struct {
	Content []struct {
		Name             string        `json:"name"`
		Key              string        `json:"key"`
		ParentKey        string        `json:"parentKey"`
		ChildKeys        []interface{} `json:"childKeys"`
		Type             string        `json:"type"`
		Outlet           bool          `json:"outlet"`
		Hidden           bool          `json:"hidden"`
		TargetGroup      string        `json:"targetGroup"`
		SuggestedFilters []string      `json:"suggestedFilters"`
	} `json:"content"`
	TotalElements int `json:"totalElements"`
	TotalPages    int `json:"totalPages"`
	Page          int `json:"page"`
	Size          int `json:"size"`
}

var zalandoBrands ZalandoBrands

var zalandoCategories map[string]ZalandoCategories

var typesOfClothes = []string{"suit", "shirt", "tie", "dress", "blouse", "skirt", "tanktop", "coat", "jacket", "T-shirt", "trousers", "jean", "shorts", "pullover", "sweater", "jumper", "cardigan", "sweatshirt", "hoodie", "glove", "mitten", "romper", "suit", "nappy", "nappies", "swimming trunk", "bikini", "swimming costume", "bathing costume", "pyjamas", "dressing gown", "nightdress", "nightiepyjama", "boxer", "Y-fronts", "bra", "pant", "knicker", "stocking", "tight", "vest", "boot", "wellington", "slipper", "sock", "ladies shoe", "sandal", "beret", "bowler hat", "trilby", "trilbies", "ladies shoe", "knitted hat", "woolen hat ", "bobble hat", "cap", "shoes"}
var colorsList = []string{"black", "brown", "beige", "gray", "white", "blue", "petrol", "turquoise", "green", "olive", "yellow", "orange", "red", "pink", "purple", "gold", "silver"}

func loadZalandoCategories() {
	fmt.Println("Loading Zalando categories...")
	for _, tp := range typesOfClothes {
		fmt.Println(tp)
		fmt.Println("-----------")
		tpurl := url.QueryEscape(tp)
		res, err := http.Get("https://api.zalando.com/categories?name=" + tpurl + "&page=1&pageSize=100")
		defer res.Body.Close()

		if err != nil {
			panic("Failed getting zalando categories")
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic("Failed reading zalando categories")
		}

		var tmp ZalandoCategories
		err = json.Unmarshal(body, &tmp)
		if err != nil {
			panic("Failed unmarshaling zalando categories")
		}

		zalandoCategories[tp] = tmp

		for _, cat := range tmp.Content {
			fmt.Println(cat.Key)
		}
	}
}

func buildApiCall(tp string) string {
	query := "pageSize=200"
	query += "&gender=" + gender

	// ageGroup := "teen"
	// if age > 0 && age <= 3 {
	//     ageGroup = "babies"
	// }

	// if age > 3 && age <= 12 {
	//     ageGroup = "kids"
	// }

	// if age > 12 && age

	for _, db := range deezerBrands {
		query += "&brand=" + db
	}
	for _, col := range deezerColors {
		query += "&color=" + col
	}

	cats := zalandoCategories[tp]
	for _, cat := range cats.Content {
		query += "&category=" + cat.Key
	}

	finalSet := removeDuplicates(setOfInsta)
	for _, nameDavid := range finalSet {
		query += "&brand=" + k[nameDavid]
	}
	return query
}

func getBrands() []string {
	var brands []string
	for _, brand := range zalandoBrands.Content {
		if brand.BrandFamily.Name == "" {
			brands = append(brands, brand.Name)
		} else {
			brands = append(brands, brand.BrandFamily.Name)
		}

	}
	return brands
}
