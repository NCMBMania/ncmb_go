package NCMB

import (
	"fmt"
	"reflect"
	"strings"
	"encoding/json"
	"net/url"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"sort"
	"golang.org/x/exp/slices"
	"time"
)

type UrlOptions struct {
	ObjectId *string
	DefinePath *string
	Queries *map[string]interface{}
}

type Signature struct {
	ncmb *NCMB
	fqdn string
	scriptFqdn string
	signatureMethod string
	signatureVersion string
	version string
	scriptVersion string
	IsScript bool
	baseInfo map[string]string
	time time.Time
}

type PathOptions struct {
	ObjectId *string
	DefinePath *string
}

func (s *Signature) Initialize() {
	s.fqdn = "mbaas.api.nifcloud.com"
	s.scriptFqdn = "script.mbaas.api.nifcloud.com"
	s.signatureMethod = "HmacSHA256"
	s.signatureVersion = "2"
	s.version = "2013-09-01"
	s.scriptVersion = "2015-09-01"
	s.baseInfo = map[string]string{
		"SignatureMethod": s.signatureMethod,
		"SignatureVersion": s.signatureVersion,
		"X-NCMB-Application-Key": s.ncmb.applicationKey,
		"X-NCMB-Timestamp": "",
	}
	s.time = time.Now()
}

func (signature *Signature) Path(className string, options PathOptions) string {
	if signature.IsScript {
		return fmt.Sprintf("/%s/script/%s", signature.scriptVersion, className)
	}
	path := fmt.Sprintf("/%s", signature.version);
  if options.DefinePath != nil {
    return fmt.Sprintf("%s/%s", path, options.DefinePath)
  }
	if slices.Index([]string{"users", "push", "roles", "files", "installations"}, className) > -1 {
		path = fmt.Sprintf("%s/%s", path, className)
	} else {
		path = fmt.Sprintf("%s/classes/%s", path, className)
	}
	if options.ObjectId != nil {
		path = fmt.Sprintf("%s/%s", path, options.ObjectId)
	}
	return path
}


func (signature *Signature) Url(className string, options UrlOptions) (string, error) {
	var queryList []string
	if options.Queries != nil {
		for key, value := range *options.Queries {
			var val string
			if reflect.TypeOf(value) == reflect.TypeOf([]interface{}{}) {
				bytes, err := json.Marshal(value)
				if err != nil {
					return "", err
				}
				val = string(bytes)
			}
			if reflect.TypeOf(value) == reflect.TypeOf(0) {
				val = fmt.Sprintf("%d", value)
			} else {
				val = fmt.Sprintf("%s", value)
			}
			queryList = append(queryList, fmt.Sprintf("%s=%s", key, url.QueryEscape(val)))
		}
	}
	queryString := ""
	if len(queryList) > 0 {
		queryString = fmt.Sprintf("?%s", strings.Join(queryList, "&"))
	}
	params := PathOptions{ObjectId: options.ObjectId, DefinePath: options.DefinePath}
	return fmt.Sprintf("https://%s%s%s", signature.Fqdn(), signature.Path(className, params), queryString), nil
}

func (signature *Signature) Headers(signatureString string) map[string]string {
	baseInfoMap := map[string]string{
		"X-NCMB-Application-Key": signature.ncmb.applicationKey,
		"X-NCMB-Timestamp":       signature.time.Format("2006-01-02T15:04:05.999Z0700"),
		"X-NCMB-Signature":			 signatureString,
	}
	if signature.ncmb.sessionToken != "" {
		baseInfoMap["X-NCMB-Apps-Session-Token"] = signature.ncmb.sessionToken
	}
	return baseInfoMap
}

func (signature *Signature) Generate(method string, className string, options UrlOptions) (string, error) {
	params := PathOptions{ObjectId: options.ObjectId, DefinePath: options.DefinePath}
	path := signature.Path(className, params)
	// baseInfoの定義
	baseInfoMap := map[string]string{
		"X-NCMB-Application-Key": signature.ncmb.applicationKey,
		"SignatureMethod":        signature.signatureMethod,
		"SignatureVersion":       signature.signatureVersion,
		"X-NCMB-Timestamp":       signature.time.Format("2006-01-02T15:04:05.999Z0700"),
	}
	// クエリが存在する場合は、それをbaseInfoMapに追加
	if options.Queries != nil {
		// fmt.Println(options.Queries)
		/*
		for k, v := range options.Queries {
			baseInfoMap[k] = url.QueryEscape(v)
		}
		*/
	}
	// 自然順序でソート
	var keys []string
	for k := range baseInfoMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	// baseInfoの作成
	var baseInfo string
	for _, k := range keys {
		if baseInfo != "" {
			baseInfo += "&"
		}
		baseInfo += fmt.Sprintf("%s=%s", k, baseInfoMap[k])
	}
	// 署名文字列の作成
	signatureString := fmt.Sprintf("%s\n%s\n%s\n%s", method, signature.Fqdn(), path, baseInfo)
	// HMACエンコーディング
	h := hmac.New(sha256.New, []byte(signature.ncmb.clientKey))
	h.Write([]byte(signatureString))
	// Base64エンコーディング
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

func (signature *Signature) Fqdn() string {
	if signature.IsScript {
		return signature.scriptFqdn
	}
	return signature.fqdn
}
