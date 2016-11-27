package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/gorilla/sessions"
	"github.com/jaytaylor/html2text"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/deezer"
	"github.com/markbates/goth/providers/instagram"
)

// domains .ONLINE, .STORE, .TECH, .SITE, .SPACE, .WEBSITE, .PRESS & .HOST
var instagramUser goth.User
var m map[string]string
var k map[string]string
var setStopWords []string
var setOfInsta []string

type MainController struct {
	beego.Controller
}

type Likes struct {
	Pagination struct {
	} `json:"pagination"`
	Meta struct {
		Code int `json:"code"`
	} `json:"meta"`
	Data []struct {
		Attribution interface{} `json:"attribution"`
		Tags        []string    `json:"tags"`
		Type        string      `json:"type"`
		Location    interface{} `json:"location"`
		Comments    struct {
			Count int `json:"count"`
		} `json:"comments"`
		Filter      string `json:"filter"`
		CreatedTime string `json:"created_time"`
		Link        string `json:"link"`
		Likes       struct {
			Count int `json:"count"`
		} `json:"likes"`
		Images struct {
			LowResolution struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"low_resolution"`
			Thumbnail struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"thumbnail"`
			StandardResolution struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"standard_resolution"`
		} `json:"images"`
		UsersInPhoto []struct {
			Position struct {
				Y float64 `json:"y"`
				X float64 `json:"x"`
			} `json:"position"`
			User struct {
				Username       string `json:"username"`
				ProfilePicture string `json:"profile_picture"`
				ID             string `json:"id"`
				FullName       string `json:"full_name"`
			} `json:"user"`
		} `json:"users_in_photo"`
		Caption struct {
			CreatedTime string `json:"created_time"`
			Text        string `json:"text"`
			From        struct {
				Username       string `json:"username"`
				ProfilePicture string `json:"profile_picture"`
				ID             string `json:"id"`
				FullName       string `json:"full_name"`
			} `json:"from"`
			ID string `json:"id"`
		} `json:"caption"`
		UserHasLiked bool   `json:"user_has_liked"`
		ID           string `json:"id"`
		User         struct {
			Username       string `json:"username"`
			ProfilePicture string `json:"profile_picture"`
			ID             string `json:"id"`
			FullName       string `json:"full_name"`
		} `json:"user"`
		Videos struct {
			LowResolution struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"low_resolution"`
			StandardResolution struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"standard_resolution"`
			LowBandwidth struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"low_bandwidth"`
		} `json:"videos,omitempty"`
	} `json:"data"`
}

func init() {
	fmt.Println("Loading Zalando brands...")
	res, err := http.Get("https://api.zalando.com/brands?pageSize=2000")

	if err != nil {
		panic("Failed getting zalando brands")
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic("Failed reading zalando brands")
	}

	err = json.Unmarshal(body, &zalandoBrands)
	if err != nil {
		panic("Failed unmarshaling zalando brands")
	}

	for _, brand := range zalandoBrands.Content {
		fmt.Println(brand.Name)
	}

	fmt.Println("Got Zalando brands!")

	zalandoCategories = make(map[string]ZalandoCategories)
	loadZalandoCategories()

	gothic.Store = sessions.NewFilesystemStore(os.TempDir(), []byte("goth-tant"))
	m = make(map[string]string)
	k = make(map[string]string)
	goth.UseProviders(
		instagram.New("xxx", "xxx", "http://tant.store/login/instagram/callback", "public_content"),
		deezer.New("xxx", "xxx", "http://tant.store/login/deezer/callback", "email,listening_history"),
	)
	MapBrands()
}

func (c *MainController) Get() {
	c.TplName = "index.tpl"
}

func (c *MainController) DeezerLogin() {
	c.Ctx.Request.URL.RawQuery += "&provider=deezer"
	fmt.Println(c.Ctx.Request.URL.Query().Get("provider"))
	gothic.BeginAuthHandler(c.Ctx.ResponseWriter, c.Ctx.Request)
}

func (c *MainController) InstagramLogin() {
	c.Ctx.Request.URL.RawQuery += "&provider=instagram"
	fmt.Println(c.Ctx.Request.URL.Query().Get("provider"))
	gothic.BeginAuthHandler(c.Ctx.ResponseWriter, c.Ctx.Request)
}

func (c *MainController) InstagramCallback() {
	c.Ctx.Request.URL.RawQuery += "&provider=instagram"

	var err error
	instagramUser, err = gothic.CompleteUserAuth(c.Ctx.ResponseWriter, c.Ctx.Request)
	if err != nil {
		fmt.Fprintln(c.Ctx.ResponseWriter, err)
		return
	}
	c.GetSuggestedBrands()
	c.Data["token"] = instagramUser.AccessToken
	c.Redirect("/login", 302)
}

func (c *MainController) GetLikes() {
	c.TplName = "index3.tpl"
	url := "https://api.instagram.com/v1/users/self/media/liked?access_token=" + instagramUser.AccessToken
	fmt.Print(url)
	fmt.Print(" ")
	resp, err := http.Get(url)
	if err != nil {
		fmt.Print("error get likes")
	}
	defer resp.Body.Close()
	var ll Likes
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &ll)
	fmt.Println()
	for _, d := range ll.Data {
		fmt.Println(d.Link)
	}
	c.Data["likes"] = string(body)
	c.TplName = "index3.tpl"
}

func MapBrands() {
	fmt.Println("start mapping")
	fmt.Print("lenth map")
	fmt.Println(len(m))

	//stops := strings.Replace(stopword, " ", "", -1)
	//sliceStops := strings.Split(stops, ",")
	m = make(map[string]string)
	k = make(map[string]string)
	for _, brand := range zalandoBrands.Content {
		var finalName string
		if brand.BrandFamily.Name == "" {
			finalName = brand.Name
		} else {
			finalName = brand.BrandFamily.Name
		}
		lowerbrand := strings.ToLower(finalName)
		s1 := strings.Replace(lowerbrand, " ", "", -1)
		s2 := strings.Replace(lowerbrand, " ", "_", -1)
		m[s1] = finalName
		m[s2] = finalName
		k[s1] = brand.Key
		k[s2] = brand.Key
		// split := strings.Split(lowerbrand, " ")
		// for _, s := range split {
		// 	if !contains(sliceStops, s) && len(s) > 2 {
		// 		m[s] = brand
		// 	}
		// }
	}
	for k, v := range m {
		fmt.Printf("key[%s] value[%s]\n", k, v)
	}
}

func (c *MainController) GetSuggestedBrands() {
	url := "https://api.instagram.com/v1/users/self/media/liked?access_token=" + instagramUser.AccessToken
	fmt.Print(url)
	fmt.Print(" ")
	resp, err := http.Get(url)
	if err != nil {
		fmt.Print("error get likes")
	}
	defer resp.Body.Close()
	var ll Likes
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &ll)
	fmt.Println()
	fmt.Println("start suggestion")
	fmt.Print("lenth likes")
	fmt.Println(len(ll.Data))
	fmt.Print("lenth map")
	fmt.Println(len(m))
	for _, d := range ll.Data {
		username := strings.ToLower(d.Caption.From.Username)
		full_name := strings.ToLower(d.Caption.From.FullName)
		text := strings.ToLower(d.Caption.Text)
		r, _ := regexp.Compile("@([a-zA-Z0-9_]*?)[^a-zA-Z0-9_]")
		usersMetionated := r.FindAllStringSubmatch(text, -1)
		for k, _ := range m {
			if strings.Contains(username, k) {
				// fmt.Println(d.Link)
				// fmt.Print(k)
				// fmt.Print(" || ")
				// fmt.Println(v)
				// fmt.Print("")
				setOfInsta = append(setOfInsta, k)

			} else if strings.Contains(full_name, k) {
				// fmt.Print(k)
				// fmt.Print(" || ")
				// fmt.Println(v)
				setOfInsta = append(setOfInsta, k)
			} else if strings.Contains(text, k) {
				// fmt.Print(k)
				// fmt.Print(" || ")
				// fmt.Println(v)
				setOfInsta = append(setOfInsta, k)
			}
		}
		for _, mentionedd := range usersMetionated {
			mentioned := mentionedd[1]
			//fmt.Println("Mentioned " + mentioned)
			url = "https://www.instagram.com/" + mentioned + "/?__a=1"
			resp, err := http.Get(url)
			if err != nil {
				fmt.Print("error get likes")
			}
			defer resp.Body.Close()
			var prof InstaProfile
			body2, err := ioutil.ReadAll(resp.Body)
			err = json.Unmarshal(body2, &prof)
			auxprof := strings.ToLower(prof.User.FullName.(string))
			//fmt.Println("############# " + auxprof)
			for k, _ := range m {
				if strings.Contains(auxprof, k) {
					// fmt.Println(d.Link)
					// fmt.Print(k)
					// fmt.Print(" |#| ")
					// fmt.Println(v)
					// fmt.Print("")
					setOfInsta = append(setOfInsta, k)
				}
			}
		}

		for _, etitq := range d.UsersInPhoto {
			//fmt.Print("----------> " + etitq.User.FullName)
			auxName := strings.ToLower(etitq.User.FullName)
			for k, _ := range m {
				if strings.Contains(auxName, k) {
					// fmt.Println(d.Link)
					// fmt.Print(k)
					// fmt.Print(" |>| ")
					// fmt.Println(v)
					// fmt.Print("")
					setOfInsta = append(setOfInsta, k)
				}
			}
		}

		//buscar los de @ en el texto
		//https://api.instagram.com/v1/users//?access_token=4193029779.50199da.a9cda2c19d4c434bb1ee79bd2507276b

	}
}

var stopword = "a, about, above, across, after, again, against, all, almost, alone, along, already, also, although, always, am, among, an, and, another, any, anybody, anyone, anything, anywhere, are, area, areas, aren't, around, as, ask, asked, asking, asks, at, away, b, back, backed, backing, backs, be, became, because, become, becomes, been, before, began, behind, being, beings, below, best, better, between, big, both, but, by, c, came, can, cannot, can't, case, cases, certain, certainly, clear, clearly, come, could, couldn't, d, did, didn't, differ, different, differently, do, does, doesn't, doing, done, don't, down, downed, downing, downs, during, e, each, early, either, end, ended, ending, ends, enough, even, evenly, ever, every, everybody, everyone, everything, everywhere, f, face, faces, fact, facts, far, felt, few, find, finds, first, for, four, from, full, fully, further, furthered, furthering, furthers, g, gave, general, generally, get, gets, give, given, gives, go, going, good, goods, got, great, greater, greatest, group, grouped, grouping, groups, h, had, hadn't, has, hasn't, have, haven't, having, he, he'd, he'll, her, here, here's, hers, herself, he's, high, higher, highest, him, himself, his, how, however, how's, i, i'd, if, i'll, i'm, important, in, interest, interested, interesting, interests, into, is, isn't, it, its, it's, itself, i've, j, just, k, keep, keeps, kind, knew, know, known, knows, l, large, largely, last, later, latest, least, less, let, lets, let's, like, likely, long, longer, longest, m, made, make, making, man, many, may, me, member, members, men, might, more, most, mostly, mr, mrs, much, must, mustn't, my, myself, n, necessary, need, needed, needing, needs, never, new, newer, newest, next, no, nobody, non, noone, nor, not, nothing, now, nowhere, number, numbers, o, of, off, often, old, older, oldest, on, once, one, only, open, opened, opening, opens, or, order, ordered, ordering, orders, other, others, ought, our, ours, ourselves, out, over, own, p, part, parted, parting, parts, per, perhaps, place, places, point, pointed, pointing, points, possible, present, presented, presenting, presents, problem, problems, put, puts, q, quite, r, rather, really, right, room, rooms, s, said, same, saw, say, says, second, seconds, see, seem, seemed, seeming, seems, sees, several, shall, shan't, she, she'd, she'll, she's, should, shouldn't, show, showed, showing, shows, side, sides, since, small, smaller, smallest, so, some, somebody, someone, something, somewhere, state, states, still, such, sure, t, take, taken, than, that, that's, the, their, theirs, them, themselves, then, there, therefore, there's, these, they, they'd, they'll, they're, they've, thing, things, think, thinks, this, those, though, thought, thoughts, three, through, thus, to, today, together, too, took, toward, turn, turned, turning, turns, two, u, under, until, up, upon, us, use, used, uses, v, very, w, want, wanted, wanting, wants, was, wasn't, way, ways, we, we'd, well, we'll, wells, went, were, we're, weren't, we've, what, what's, when, when's, where, where's, whether, which, while, who, whole, whom, who's, whose, why, why's, will, with, within, without, won't, work, worked, working, works, would, wouldn't, x, y, year, years, yes, yet, you, you'd, you'll, young, younger, youngest, your, you're, yours, yourself, yourselves, you've, z"

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

type InstaProfile struct {
	User struct {
		Username        string      `json:"username"`
		ConnectedFbPage interface{} `json:"connected_fb_page"`
		Follows         struct {
			Count int `json:"count"`
		} `json:"follows"`
		RequestedByViewer bool `json:"requested_by_viewer"`
		FollowedBy        struct {
			Count int `json:"count"`
		} `json:"followed_by"`
		CountryBlock           interface{} `json:"country_block"`
		HasRequestedViewer     bool        `json:"has_requested_viewer"`
		ExternalURLLinkshimmed interface{} `json:"external_url_linkshimmed"`
		FollowsViewer          bool        `json:"follows_viewer"`
		ProfilePicURL          string      `json:"profile_pic_url"`
		ExternalURL            interface{} `json:"external_url"`
		IsPrivate              bool        `json:"is_private"`
		FullName               interface{} `json:"full_name"`
		Media                  struct {
			Count    int `json:"count"`
			PageInfo struct {
				HasPreviousPage bool   `json:"has_previous_page"`
				StartCursor     string `json:"start_cursor"`
				EndCursor       string `json:"end_cursor"`
				HasNextPage     bool   `json:"has_next_page"`
			} `json:"page_info"`
			Nodes []struct {
				Code       string `json:"code"`
				Dimensions struct {
					Width  int `json:"width"`
					Height int `json:"height"`
				} `json:"dimensions"`
				CommentsDisabled bool `json:"comments_disabled"`
				Owner            struct {
					ID string `json:"id"`
				} `json:"owner"`
				Comments struct {
					Count int `json:"count"`
				} `json:"comments"`
				Caption string `json:"caption,omitempty"`
				Likes   struct {
					Count int `json:"count"`
				} `json:"likes"`
				Date         int    `json:"date"`
				ThumbnailSrc string `json:"thumbnail_src"`
				IsVideo      bool   `json:"is_video"`
				ID           string `json:"id"`
				DisplaySrc   string `json:"display_src"`
			} `json:"nodes"`
		} `json:"media"`
		HasBlockedViewer bool        `json:"has_blocked_viewer"`
		FollowedByViewer bool        `json:"followed_by_viewer"`
		IsVerified       bool        `json:"is_verified"`
		ID               string      `json:"id"`
		Biography        interface{} `json:"biography"`
		BlockedByViewer  bool        `json:"blocked_by_viewer"`
	} `json:"user"`
}

var deezerBrands, deezerColors, deezerTypes []string

type BingResults struct {
	Type     string `json:"_type"`
	WebPages struct {
		WebSearchURL          string `json:"webSearchUrl"`
		TotalEstimatedMatches int    `json:"totalEstimatedMatches"`
		Value                 []struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			URL             string `json:"url"`
			DisplayURL      string `json:"displayUrl"`
			Snippet         string `json:"snippet"`
			DateLastCrawled string `json:"dateLastCrawled"`
			DeepLinks       []struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"deepLinks,omitempty"`
		} `json:"value"`
	} `json:"webPages"`
	RelatedSearches struct {
		ID    string `json:"id"`
		Value []struct {
			Text         string `json:"text"`
			DisplayText  string `json:"displayText"`
			WebSearchURL string `json:"webSearchUrl"`
		} `json:"value"`
	} `json:"relatedSearches"`
	RankingResponse struct {
		Mainline struct {
			Items []struct {
				AnswerType  string `json:"answerType"`
				ResultIndex int    `json:"resultIndex"`
				Value       struct {
					ID string `json:"id"`
				} `json:"value"`
			} `json:"items"`
		} `json:"mainline"`
		Sidebar struct {
			Items []struct {
				AnswerType string `json:"answerType"`
				Value      struct {
					ID string `json:"id"`
				} `json:"value"`
			} `json:"items"`
		} `json:"sidebar"`
	} `json:"rankingResponse"`
}

type DeezerFlow struct {
	Data []struct {
		ID             int    `json:"id"`
		Title          string `json:"title"`
		TitleShort     string `json:"title_short"`
		TitleVersion   string `json:"title_version"`
		Duration       int    `json:"duration"`
		Rank           int    `json:"rank"`
		ExplicitLyrics bool   `json:"explicit_lyrics"`
		Preview        string `json:"preview"`
		Artist         struct {
			ID            int    `json:"id"`
			Name          string `json:"name"`
			Picture       string `json:"picture"`
			PictureSmall  string `json:"picture_small"`
			PictureMedium string `json:"picture_medium"`
			PictureBig    string `json:"picture_big"`
			PictureXl     string `json:"picture_xl"`
			Tracklist     string `json:"tracklist"`
			Type          string `json:"type"`
		} `json:"artist"`
		Album struct {
			ID          int    `json:"id"`
			Title       string `json:"title"`
			Cover       string `json:"cover"`
			CoverSmall  string `json:"cover_small"`
			CoverMedium string `json:"cover_medium"`
			CoverBig    string `json:"cover_big"`
			CoverXl     string `json:"cover_xl"`
			Tracklist   string `json:"tracklist"`
			Type        string `json:"type"`
		} `json:"album"`
		Type string `json:"type"`
	} `json:"data"`
}

var deezerUser goth.User

func (c *MainController) DeezerCallback() {
	c.Ctx.Request.URL.RawQuery += "&provider=deezer"
	defer c.Redirect("/login", 301)

	var err error
	deezerUser, err = gothic.CompleteUserAuth(c.Ctx.ResponseWriter, c.Ctx.Request)
	if err != nil {
		fmt.Fprintln(c.Ctx.ResponseWriter, err)
		return
	}

	// get latest songs
	res, err := http.Get(deezerUser.RawData["tracklist"].(string))
	defer res.Body.Close()
	if err != nil {
		fmt.Println("Error getting flow: " + err.Error())
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading body" + err.Error())
		return
	}

	var flow DeezerFlow
	err = json.Unmarshal(body, &flow)
	if err != nil {
		fmt.Println("Error unmarshaling: " + err.Error())
		return
	}

	// bing search for clothes of latest 25 artists
	cSongs := make(chan int)

	var allBrands, allColors, allTypes []string

	for _, d := range flow.Data {
		go func(name string) {
			defer func() {
				cSongs <- 1
			}()
			//fmt.Println(d.Artist.Name, d.Artist.Type)
			url := fmt.Sprintf("https://api.cognitive.microsoft.com/bing/v5.0/search?q=%s&count=10&offset=0&mkt=en-us&safesearch=Off", url.QueryEscape("what does \""+name+"\" wear"))
			//fmt.Println(url)
			req, _ := http.NewRequest("GET", url, nil)

			req.Header.Add("ocp-apim-subscription-key", "xxx")
			req.Header.Add("cache-control", "no-cache")

			res, _ := http.DefaultClient.Do(req)

			body, err = ioutil.ReadAll(res.Body)
			//fmt.Println(string(body))
			//fmt.Println(res)
			if err != nil {
				fmt.Println("Error reading body " + err.Error())
				return
			}

			var br BingResults
			err = json.Unmarshal(body, &br)

			cWebs := make(chan int, 4)
			added := 0
			for i, result := range br.WebPages.Value {
				if i > 3 {
					break
				}
				added++

				go func(url string, i int) {
					defer func() {
						cWebs <- 1
					}()
					theUrl := url
					if !strings.Contains(theUrl, "http") {
						theUrl = "http://" + theUrl
					}
					res, err = http.Get(theUrl)
					if err != nil {
						fmt.Println("Failed getting URL " + strconv.Itoa(i) + " " + err.Error())
						return
					}

					defer res.Body.Close()

					body, err = ioutil.ReadAll(res.Body)
					if err != nil {
						fmt.Println("Failed reading body URL " + strconv.Itoa(i) + " " + err.Error())
						return
					}

					text, _ := html2text.FromString(string(body))
					brands := extractBrands(text)
					allBrands = append(allBrands, brands...)
					colors := extractColors(text)
					allColors = append(allColors, colors...)
					types := extractTypes(text)
					allTypes = append(allTypes, types...)
					//fmt.Println(name, "\nBrands / ", strings.Join(brands, ","), "\nTypes / ", types, "\nColors / ", colors) /*, "["+theUrl+"]"*/
				}(result.URL, i)
			}

			for i := 0; i < added; i++ {
				<-cWebs
				//fmt.Println("Finished web")
			}
		}(d.Artist.Name)
	}

	for i := 0; i < len(flow.Data); i++ {
		<-cSongs
		//fmt.Println("Finished song", i, len(flow.Data))
	}

	deezerBrands = removeDuplicates(allBrands)
	deezerColors = removeDuplicates(allColors)
	deezerTypes = removeDuplicates(allTypes)
	//fmt.Println("Brands - ", strings.Join(allBrands, ", "))
	//mt.Println("Colors - ", strings.Join(allColors, ", "))
	//fmt.Println("Types - ", strings.Join(allTypes, ", "))

}

func (c *MainController) CallZalando() {
	query := buildApiCall(c.Input().Get("type"))
	res, err := http.Get("https://api.zalando.com/articles?" + query)
	if err != nil {
		fmt.Println("Failed calling zalando api! " + err.Error())
		c.TplName = "error.tpl"
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Failed reading zalando body " + err.Error())
		return
	}

	var articles ZalandoArticles
	err = json.Unmarshal(body, &articles)
	if err != nil {
		fmt.Println("Failed to parse articles JSON " + err.Error())
		return
	}

	//for _, art := range articles.Content {
	//	fmt.Println(art.Name)
	//}

	c.Data["art1"] = articles.Content[rand.Intn(len(articles.Content))]
	c.Data["art2"] = articles.Content[rand.Intn(len(articles.Content))]
	c.Data["art3"] = articles.Content[rand.Intn(len(articles.Content))]
	c.Data["art4"] = articles.Content[rand.Intn(len(articles.Content))]

	c.TplName = "result.tpl"
}

func extractBrands(text string) []string {
	var res []string
	for _, brand := range zalandoBrands.Content {
		//fam := brand.BrandFamily.Name
		//key := brand.BrandFamily.Key
		//if fam == "" {
		fam := brand.Name
		//}

		re := regexp.MustCompile("[^a-zA-Z]" + fam + " [^a-zA-Z]")
		if re.MatchString(text) {
			res = append(res, brand.Key)
		}
	}

	return res
}

func extractColors(text string) []string {
	var res []string
	for _, color := range colorsList {
		re := regexp.MustCompile("[^a-zA-Z]" + color + " [^a-zA-Z]")
		if re.MatchString(text) {
			res = append(res, color)
		}
	}

	return res
}
func extractTypes(text string) []string {
	var res []string
	for _, tp := range typesOfClothes {
		re := regexp.MustCompile("[^a-zA-Z]" + tp + " [^a-zA-Z]")
		if re.MatchString(text) {
			res = append(res, tp)
		}
	}

	return res
}
func removeDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

func (c *MainController) ShowForm() {
	c.TplName = "form.tpl"
}

var gender, age, country string

func (c *MainController) SetInfo() {
	gender = c.GetString("gender")
	age = c.GetString("age")
	country = c.GetString("country")
	c.Redirect("/login", 302)
}

func (c *MainController) Logins() {
	c.Data["instagramUser"] = instagramUser
	c.Data["deezerUser"] = deezerUser
	fmt.Println(deezerUser)
	c.TplName = "logins.tpl"
}

func (c *MainController) Recommendations() {
	c.TplName = "recommendations.tpl"
}
