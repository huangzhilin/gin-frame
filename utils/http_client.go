package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

//Get 发送get请求方式，Url中直接拼接了请求参数  如：http://xxxx.com?r=1&a=2
func Get(Url string) ([]byte, error) {
	resp, err := http.Get(Url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

//GetWithParams 发送get请求方式，如：http://xxxx.com?r=1添加额外的请求参数map[string]string{"a":"2"}及添加header
func GetWithParams(Url string, params map[string]string, header map[string]string) ([]byte, error) {
	parseUrl, err := url.Parse(Url)
	if err != nil {
		return nil, err
	}

	//返回values类型的字典
	paramValues, err := url.ParseQuery(parseUrl.RawQuery)
	if err != nil {
		return nil, err
	}

	//将额为的params参数加入到values字典中
	if len(params) > 0 {
		for k, v := range params {
			paramValues.Add(k, v)
		}
	}

	parseUrl.RawQuery = paramValues.Encode() //如果参数中有中文，进行urlEncode

	urlPathWithParams := parseUrl.String()
	fmt.Println(urlPathWithParams)

	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, urlPathWithParams, nil)
	if err != nil {
		return nil, err
	}

	//追加请求头
	if len(header) > 0 {
		for k, v := range header {
			req.Header.Add(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

//PostForm 发送post请求方式
func PostForm(Url string, data map[string]string) ([]byte, error) {
	values := url.Values{}
	for k, v := range data {
		values.Add(k, v)
	}

	resp, err := http.PostForm(Url, values)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
