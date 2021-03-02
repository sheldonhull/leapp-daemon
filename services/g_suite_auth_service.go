package services

import (
	"crypto/tls"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/pkg/errors"
	"golang.org/x/net/publicsuffix"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"runtime"
	"strings"
	"time"
)

var httpClient http.Client

func GSuiteAuthFirstStepService(username string, password string) (url.Values, string, string, string, url.Values, string) {
	jar, _ := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})

	httpClient = http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: false},
		},
		Jar: jar,
	}

	firstPageURL := ""

	//=====================
	// BEGIN loadFirstPage
	//=====================
	req, err := http.NewRequest("GET", firstPageURL + "&hl=en&loc=US", nil)
	if err != nil { log.Println(err) }

	req.Header.Set("User-Agent", fmt.Sprintf("leapp-daemon/0.0.1 (%s %s) Noovolari", runtime.GOOS, runtime.GOARCH))

	res, err := httpClient.Do(req)
	if err != nil { log.Println(err) }

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil { log.Println(err) }

	doc.Url, err = url.Parse(firstPageURL + "&hl=en&loc=US")
	if err != nil { log.Println(err) }

	authForm, authURL, err := extractInputsByFormID(doc, "gaia_loginform", "challenge")
	if err != nil { log.Println(err) }

	_, loginPageV1 := authForm["GALX"]

	var postForm url.Values
	// using a field which is known to be in the original login page
	if loginPageV1 {
		// Login page v1
		postForm = url.Values{
			"bgresponse":               []string{"js_disabled"},
			"checkConnection":          []string{""},
			"checkedDomains":           []string{"youtube"},
			"continue":                 []string{authForm.Get("continue")},
			"gxf":                      []string{authForm.Get("gxf")},
			"identifier-captcha-input": []string{""},
			"identifiertoken":          []string{""},
			"identifiertoken_audio":    []string{""},
			"ltmpl":                    []string{"popup"},
			"oauth":                    []string{"1"},
			"Page":                     []string{authForm.Get("Page")},
			"Passwd":                   []string{""},
			"PersistentCookie":         []string{"yes"},
			"ProfileInformation":       []string{""},
			"pstMsg":                   []string{"0"},
			"sarp":                     []string{"1"},
			"scc":                      []string{"1"},
			"SessionState":             []string{authForm.Get("SessionState")},
			"signIn":                   []string{authForm.Get("signIn")},
			"_utf8":                    []string{authForm.Get("_utf8")},
			"GALX":                     []string{authForm.Get("GALX")},
		}
	} else {
		// Login page v2
		postForm = url.Values{
			"challengeId":     []string{"1"},
			"challengeType":   []string{"1"},
			"continue":        []string{authForm.Get("continue")},
			"scc":             []string{"1"},
			"sarp":            []string{"1"},
			"checkeddomains":  []string{"youtube"},
			"checkConnection": []string{"youtube:930:1"},
			"pstMessage":      []string{"1"},
			"oauth":           []string{authForm.Get("oauth")},
			"flowName":        []string{authForm.Get("flowName")},
			"faa":             []string{"1"},
			"Email":           []string{""},
			"Passwd":          []string{""},
			"TrustDevice":     []string{"on"},
			"bgresponse":      []string{"js_disabled"},
		}
		for _, k := range []string{"TL", "gxf"} {
			if v, ok := authForm[k]; ok {
				postForm.Set(k, v[0])
			}
		}
	}
	//===================
	// END loadFirstPage
	//===================

	authForm = postForm
	authForm.Set("Email", username)

	//=====================
	// BEGIN loadLoginPage
	//=====================

	submitURL := authURL+"?hl=en&loc=US"
	referer := firstPageURL+"&hl=en&loc=US"

	req, err = http.NewRequest("POST", submitURL, strings.NewReader(authForm.Encode()))
	if err != nil { log.Println(err) }

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Set("Content-Language", "en-US")
	req.Header.Set("Referer", referer)

	res, err = httpClient.Do(req)
	if err != nil { log.Println(err) }

	doc, err = goquery.NewDocumentFromReader(res.Body)
	if err != nil { log.Println(err) }

	doc.Url, err = url.Parse(submitURL)
	if err != nil { log.Println(err) }

	passwordForm, passwordURL, err := extractInputsByFormID(doc, "gaia_loginform", "challenge")
	if err != nil { log.Println(err) }

	//===================
	// END loadLoginPage
	//===================

	authForm.Set("Passwd", password)
	referingURL := passwordURL

	if _, rawIdPresent := passwordForm["rawidentifier"]; rawIdPresent {
		authForm.Set("rawidentifier", username)
		referingURL = authURL
	}

	if v, tlPresent := passwordForm["TL"]; tlPresent {
		authForm.Set("TL", v[0])
	}
	if v, gxfPresent := passwordForm["gxf"]; gxfPresent {
		authForm.Set("gxf", v[0])
	}

	//=========================
	// BEGIN loadChallengePage
	//=========================

	referer = referingURL

	req, err = http.NewRequest("POST", submitURL + "?hl=en&loc=US", strings.NewReader(authForm.Encode()))
	if err != nil { log.Println(err) }
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Set("Content-Language", "en-US")
	req.Header.Set("Referer", referer)

	res, err = httpClient.Do(req)
	if err != nil { log.Println(err) }

	doc, err = goquery.NewDocumentFromReader(res.Body)
	if err != nil { log.Println(err) }

	doc.Url, err = url.Parse(submitURL)
	if err != nil { log.Println(err) }

	errMsg := mustFindErrorMsg(doc)

	if errMsg != "" {
		log.Println(errMsg)
	}

	secondFactorHeader := "This extra step shows it’s really you trying to sign in"
	secondFactorHeader2 := "This extra step shows that it’s really you trying to sign in"
	secondFactorHeaderJp := "2 段階認証プロセス"

	// have we been asked for 2-Step Verification
	if extractNodeText(doc, "h2", secondFactorHeader) != "" ||
		extractNodeText(doc, "h2", secondFactorHeader2) != "" ||
		extractNodeText(doc, "h1", secondFactorHeaderJp) != "" {

		responseForm, secondActionURL, err := extractInputsByFormID(doc, "challenge")
		if err != nil {
			log.Println(err)
		}

		switch {
		case strings.Contains(secondActionURL, "challenge/totp/"): // handle TOTP challenge

			var token = "" // mfaToken

			responseForm.Set("Pin", token)
			responseForm.Set("TrustDevice", "on") // Don't ask again on this computer

			req, err := http.NewRequest("POST", submitURL, strings.NewReader(responseForm.Encode()))
			if err != nil {
				log.Println(err)
			}

			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Set("Accept-Language", "en")
			req.Header.Set("Content-Language", "en-US")
			req.Header.Set("Referer", submitURL)

			res, err := httpClient.Do(req)
			if err != nil {
				log.Println(err)
			}

			doc, err = goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				log.Println(err)
			}
		}
	}

	//=======================
	// END loadChallengePage
	//=======================

	responseDoc := doc

	captchaInputIds := []string{
		"logincaptcha",
		"identifier-captcha-input",
	}

	var captchaFound *goquery.Selection
	var captchaInputId string

	for _, v := range captchaInputIds {
		captchaFound = responseDoc.Find(fmt.Sprintf("#%s", v))
		if captchaFound != nil && captchaFound.Length() > 0 {
			captchaInputId = v
			break
		}
	}

	var captchaForm url.Values
	var captchaPictureURL string
	var captchaURL string

	for captchaFound != nil && captchaFound.Length() > 0 {
		captchaImgDiv := responseDoc.Find(".captcha-img")
		captchaPictureSrc, found := goquery.NewDocumentFromNode(captchaImgDiv.Children().Nodes[0]).Attr("src")

		if !found {
			log.Println(errors.New("captcha image not found but requested"))
		}

		var err error

		captchaPictureURL, err = generateFullURLIfRelative(captchaPictureSrc, passwordURL)
		if err != nil {
			log.Println(err)
		}

		captchaForm, captchaURL, err = extractInputsByFormID(responseDoc, "gaia_loginform", "challenge")
		if err != nil {
			log.Println(err)
		}

		_, captchaV1 := captchaForm["Passwd"]
		if captchaV1 {
			captchaForm.Set("Passwd", password)
		}
	}

	var loginForm url.Values
	var loginURL string

	if captchaInputId == "" && captchaURL == "" && captchaForm == nil {
		passworddBeingRequested := responseDoc.Find("#password")
		if passworddBeingRequested != nil && passworddBeingRequested.Length() > 0 {

			loginForm, loginURL, err = extractInputsByFormID(responseDoc, "challenge")
			if err != nil {
				log.Println(err)
			}
		}
	}

	return captchaForm, captchaInputId, captchaPictureURL, captchaURL, loginForm, loginURL
}

func GSuiteAuthSecondStepService(captcha string, captchaInputId string, captchaURL string, captchaForm url.Values,
	password string, loginForm url.Values, loginURL string) (bool, url.Values, string) {

	var referer string
	var responseDoc *goquery.Document

	if captcha != "" && captchaInputId != "" && captchaURL != "" && captchaForm != nil {
		captchaForm.Set(captchaInputId, captcha)

		submitURL := captchaURL + "?hl=en&loc=US"
		referer = captchaURL

		//=========================
		// BEGIN loadChallengePage
		//=========================

		req, err := http.NewRequest("POST", submitURL, strings.NewReader(captchaForm.Encode()))
		if err != nil {
			log.Println(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Accept-Language", "en-US")
		req.Header.Set("Content-Language", "en-US")
		req.Header.Set("Referer", referer)

		res, err := httpClient.Do(req)
		if err != nil {
			log.Println(err)
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Println(err)
		}

		doc.Url, err = url.Parse(submitURL)
		if err != nil {
			log.Println(err)
		}

		errMsg := mustFindErrorMsg(doc)

		if errMsg != "" {
			log.Println(errMsg)
		}

		secondFactorHeader := "This extra step shows it’s really you trying to sign in"
		secondFactorHeader2 := "This extra step shows that it’s really you trying to sign in"
		secondFactorHeaderJp := "2 段階認証プロセス"

		// have we been asked for 2-Step Verification
		if extractNodeText(doc, "h2", secondFactorHeader) != "" ||
			extractNodeText(doc, "h2", secondFactorHeader2) != "" ||
			extractNodeText(doc, "h1", secondFactorHeaderJp) != "" {

			responseForm, secondActionURL, err := extractInputsByFormID(doc, "challenge")
			if err != nil {
				log.Println(err)
			}

			switch {
			case strings.Contains(secondActionURL, "challenge/totp/"): // handle TOTP challenge

				var token = "" // mfaToken

				responseForm.Set("Pin", token)
				responseForm.Set("TrustDevice", "on") // Don't ask again on this computer

				req, err := http.NewRequest("POST", submitURL, strings.NewReader(responseForm.Encode()))
				if err != nil {
					log.Println(err)
				}

				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				req.Header.Set("Accept-Language", "en")
				req.Header.Set("Content-Language", "en-US")
				req.Header.Set("Referer", submitURL)

				res, err := httpClient.Do(req)
				if err != nil {
					log.Println(err)
				}

				doc, err = goquery.NewDocumentFromReader(res.Body)
				if err != nil {
					log.Println(err)
				}
			}
		}

		//=======================
		// END loadChallengePage
		//=======================

		responseDoc = doc

		//captchaFound := responseDoc.Find(fmt.Sprintf("#%s", captchaInputId))
	}

	isMfaTokenRequested := false
	var responseForm url.Values
	var submitURL string

	var condition bool

	if loginForm != nil && loginURL != "" {
		condition = true
	} else {
		passworddBeingRequested := responseDoc.Find("#password")
		condition = passworddBeingRequested != nil && passworddBeingRequested.Length() > 0
	}

	if condition {
		var err error

		if loginForm == nil && loginURL == "" {
			loginForm, loginURL, err = extractInputsByFormID(responseDoc, "challenge")
			if err != nil {
				log.Println(err)
			}
		}

		loginForm.Set("Passwd", password)

		submitURL = loginURL+"?hl=en&loc=US"
		referer = loginURL

		req, err := http.NewRequest("POST", submitURL, strings.NewReader(loginForm.Encode()))
		if err != nil {
			log.Println(err)
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Accept-Language", "en-US")
		req.Header.Set("Content-Language", "en-US")
		req.Header.Set("Referer", referer)

		res, err := httpClient.Do(req)
		if err != nil {
			log.Println(err)
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Println(err)
		}

		doc.Url, err = url.Parse(submitURL)
		if err != nil {
			log.Println(err)
		}

		errMsg := mustFindErrorMsg(doc)

		if errMsg != "" {
			log.Println(errMsg)
		}

		secondFactorHeader := "This extra step shows it’s really you trying to sign in"
		secondFactorHeader2 := "This extra step shows that it’s really you trying to sign in"
		secondFactorHeaderJp := "2 段階認証プロセス"

		// have we been asked for 2-Step Verification
		if extractNodeText(doc, "h2", secondFactorHeader) != "" ||
			extractNodeText(doc, "h2", secondFactorHeader2) != "" ||
			extractNodeText(doc, "h1", secondFactorHeaderJp) != "" {

			log.Println("1 - have we been asked for 2-Step Verification")

			var secondActionURL string

			responseForm, secondActionURL, err = extractInputsByFormID(doc, "challenge")
			if err != nil {
				log.Println(err)
			}

			switch {
			case strings.Contains(secondActionURL, "challenge/totp/"): // handle TOTP challenge

				var token = "" // mfaToken

				responseForm.Set("Pin", token)
				responseForm.Set("TrustDevice", "on") // Don't ask again on this computer

				req, err := http.NewRequest("POST", submitURL, strings.NewReader(responseForm.Encode()))
				if err != nil {
					log.Println(err)
				}

				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				req.Header.Set("Accept-Language", "en")
				req.Header.Set("Content-Language", "en-US")
				req.Header.Set("Referer", submitURL)

				res, err := httpClient.Do(req)
				if err != nil {
					log.Println(err)
				}

				doc, err = goquery.NewDocumentFromReader(res.Body)
				if err != nil {
					log.Println(err)
				}

				responseDoc = doc
			}

			skipResponseForm, skipActionURL, err := extractInputsByFormQuery(doc, `[action$="skip"]`)

			if err != nil {
				log.Println(err)
			}

			if skipActionURL == "" {
				log.Println(err)
			}

			referer = submitURL
			submitURL = skipActionURL
			authForm := skipResponseForm

			req, err := http.NewRequest("POST", submitURL, strings.NewReader(authForm.Encode()))
			if err != nil {
				log.Println(err)
			}

			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Set("Accept-Language", "en-US")
			req.Header.Set("Content-Language", "en-US")
			req.Header.Set("Referer", referer)

			res, err := httpClient.Do(req)
			if err != nil {
				log.Println(err)
			}

			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				log.Println(err)
			}

			doc.Url, err = url.Parse(submitURL)
			if err != nil {
				log.Println(err)
			}

			var challengeEntry string

			doc.Find("form[data-challengeentry]").EachWithBreak(func(i int, s *goquery.Selection) bool {
				action, ok := s.Attr("action")
				if !ok {
					return true
				}

				if strings.Contains(action, "challenge/totp/") ||
					strings.Contains(action, "challenge/ipp/") ||
					strings.Contains(action, "challenge/az/") ||
					strings.Contains(action, "challenge/skotp/") {

					challengeEntry, _ = s.Attr("data-challengeentry")
					return false
				}

				return true
			})

			if challengeEntry == "" {
				log.Println(err)
			}

			query := fmt.Sprintf(`[data-challengeentry="%s"]`, challengeEntry)

			var newActionURL string

			responseForm, newActionURL, err = extractInputsByFormQuery(doc, query)
			if err != nil {
				log.Println(err)
			}

			referer = submitURL
			submitURL = newActionURL
			authForm = responseForm

			req, err = http.NewRequest("POST", submitURL, strings.NewReader(authForm.Encode()))

			if err != nil {
				log.Println(err)
			}

			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Set("Accept-Language", "en-US")
			req.Header.Set("Content-Language", "en-US")
			req.Header.Set("Referer", referer)

			res, err = httpClient.Do(req)
			if err != nil {
				log.Println(err)
			}

			doc, err = goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				log.Println(err)
			}

			doc.Url, err = url.Parse(submitURL)
			if err != nil {
				log.Println(err)
			}

			errMsg := mustFindErrorMsg(doc)

			if errMsg != "" {
				log.Println(err)
			}

			secondFactorHeader := "This extra step shows it’s really you trying to sign in"
			secondFactorHeader2 := "This extra step shows that it’s really you trying to sign in"
			secondFactorHeaderJp := "2 段階認証プロセス"

			if extractNodeText(doc, "h2", secondFactorHeader) != "" ||
				extractNodeText(doc, "h2", secondFactorHeader2) != "" ||
				extractNodeText(doc, "h1", secondFactorHeaderJp) != "" {

				log.Println("2 - have we been asked for 2-Step Verification")

				responseForm, secondActionURL, err = extractInputsByFormID(doc, "challenge")
				if err != nil {
					log.Println(err)
				}

				switch {
				case strings.Contains(secondActionURL, "challenge/totp/"): // handle TOTP challenge
					isMfaTokenRequested = true
				}
			}
		}
	}

	return isMfaTokenRequested, responseForm, submitURL
}

func GSuiteAuthThirdStepService(isMfaTokenRequested bool, responseForm url.Values, submitURL string, token string) string {
	var responseDoc *goquery.Document

	if isMfaTokenRequested {
		responseForm.Set("Pin", token)
		responseForm.Set("TrustDevice", "on") // Don't ask again on this computer

		req, err := http.NewRequest("POST", submitURL, strings.NewReader(responseForm.Encode()))
		if err != nil {
			log.Println(err)
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Accept-Language", "en")
		req.Header.Set("Content-Language", "en-US")
		req.Header.Set("Referer", submitURL)

		res, err := httpClient.Do(req)
		if err != nil {
			log.Println(err)
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Println(err)
		}

		responseDoc = doc
	}

	samlAssertion := mustFindInputByName(responseDoc, "SAMLResponse")
	if samlAssertion == "" {
		log.Println("SAML assertion not found")
	}

	var durationSeconds int64 = 3600
	var principalArn = ""
	var roleArn = ""

	request := sts.AssumeRoleWithSAMLInput{
		DurationSeconds: &durationSeconds,
		Policy:          nil,
		PolicyArns:      nil,
		PrincipalArn:    &principalArn,
		RoleArn:         &roleArn,
		SAMLAssertion:   &samlAssertion,
	}

	sess := session.Must(session.NewSession())
	response, _ := sts.New(sess).AssumeRoleWithSAML(&request)

	return response.Credentials.String()
}

func extractInputsByFormID(doc *goquery.Document, formID ...string) (url.Values, string, error) {
	for _, id := range formID {
		formData, actionURL, err := extractInputsByFormQuery(doc, fmt.Sprintf("#%s", id))
		if err != nil && strings.HasPrefix(err.Error(), "could not find form with query ") {
			continue
		}
		return formData, actionURL, err
	}
	return url.Values{}, "", errors.New("could not find any forms matching the provided IDs")
}

func extractInputsByFormQuery(doc *goquery.Document, formQuery string) (url.Values, string, error) {
	formData := url.Values{}
	var actionAttr string

	query := fmt.Sprintf("form%s", formQuery)

	currentURL := doc.Url.String()

	//get action url
	foundForms := doc.Find(query)
	if len(foundForms.Nodes) == 0 {
		return formData, "", fmt.Errorf("could not find form with query %q", query)
	}

	foundForms.Each(func(i int, s *goquery.Selection) {
		action, ok := s.Attr("action")
		if !ok {
			return
		}
		actionAttr = action
	})

	actionURL, err := generateFullURLIfRelative(actionAttr, currentURL)
	if err != nil {
		return formData, "", errors.Wrap(err, "error getting action URL")
	}

	query = fmt.Sprintf("form%s", formQuery)

	// extract form data to passthrough
	doc.Find(query).Find("input").Each(func(i int, s *goquery.Selection) {
		name, ok := s.Attr("name")
		if !ok {
			return
		}
		val, ok := s.Attr("value")
		if !ok {
			return
		}
		formData.Add(name, val)
	})

	return formData, actionURL, nil
}

func generateFullURLIfRelative(destination, currentPageURL string) (string, error) {
	if string(destination[0]) == "/" {
		currentURLParsed, err := url.Parse(currentPageURL)
		if err != nil {
			return "", errors.Wrap(err, "error generating full URL")
		}

		return fmt.Sprintf("%s://%s%s", currentURLParsed.Scheme, currentURLParsed.Host, destination), nil
	} else {
		return destination, nil
	}
}

func mustFindErrorMsg(doc *goquery.Document) string {
	var fieldValue string
	doc.Find(".error-msg").Each(func(i int, s *goquery.Selection) {
		fieldValue = s.Text()

	})
	return fieldValue
}

func extractNodeText(doc *goquery.Document, tag, txt string) string {

	var res string

	doc.Find(tag).Each(func(i int, s *goquery.Selection) {
		if s.Text() == txt {
			res = s.Text()
		}
	})

	return res
}

func mustFindInputByName(doc *goquery.Document, name string) string {

	var fieldValue string

	q := fmt.Sprintf(`input[name="%s"]`, name)

	doc.Find(q).Each(func(i int, s *goquery.Selection) {
		val, ok := s.Attr("value")
		if !ok {
			log.Println(errors.New("mustFindInputByName error"))
		}
		fieldValue = val
	})

	return fieldValue
}
