package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/tidwall/gjson"
)

func run() {
	data, err := os.ReadFile("/home/chengy/go/src/wechat/data.txt")
	if err != nil {
		log.Fatal(err)
	}

	reg := regexp.MustCompile(`APIS.(.*):\s+(\{[\S\s]+?\}),`)
	rs := reg.FindAllStringSubmatch(string(data), -1)
	final := []string{
		"package api",

		`
		type Api struct {
			host    string
			timeout time.Duration
		}
		
		type Config struct{}
		
		func New(host string) *Api {
			return &Api{
				host:    host + "?type=%d",
				timeout: 5 * time.Second,
			}
		}

		`,
	}
	for _, v := range rs {
		apiName := transferName(v[1])
		structTemp := createStruct(apiName, v[2])
		structEmpty := false
		if structTemp == "" {
			structEmpty = true
		}
		final = append(final, structTemp)
		final = append(final, "\n")
		final = append(final, fmt.Sprintf("\ntype %sResult struct{}", apiName))
		final = append(final, createFunc(apiName, v[1], structEmpty))
	}
	final = append(final, `
		type Result struct{}

		func (a *Api) sendApi(t Type, data, v interface{}) error {
			js := []byte(`+"`{}`"+`)
			if data != nil {
				js, _ = json.Marshal(data)
			}
			host := fmt.Sprintf(a.host, t)
			cli := http.DefaultClient
			cli.Timeout = a.timeout
			req, err := http.NewRequest("POST", host, bytes.NewReader(js))
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", "application/json")
			resp, err := cli.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
		
			respData, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			fmt.Println(string(respData))
			if err := json.Unmarshal(respData, v); err != nil {
				return err
			}
			return nil
		}
	`)
	os.WriteFile("/home/chengy/go/src/wechat/api/api.go", []byte(strings.Join(final, "\n")), os.ModePerm)
}

// type MsgStartVoiceHookData struct {
// }

// func MsgStartVoiceHook(data *MsgStartVoiceHookData) {
// 	SendApi(MsgStartVoiceHook, data)
// }

// func SendApi(t Type, data interface{}) {
// 	val, _ := json.Marshal(data)
// }

func createFunc(name, api string, structEmpty bool) string {
	dataFiled := ""
	argsField := "nil"
	resultField := fmt.Sprintf("%sResult", name)
	if !structEmpty {
		dataFiled = "data *" + name + "Data"
		argsField = "data"
	}
	return fmt.Sprintf(`func(a *Api) Request%s(%s)(*%s,error) {
		res := &%s{}
		if err := a.sendApi(%s, %s, res);err != nil{
			return nil,err
		}
		
		return res,nil
	}`, name, dataFiled, resultField, resultField, api, argsField)
}

func createStruct(name, s string) string {
	js := gjson.Parse(s)
	template := "%s %s `json:\"%s\"`"
	var pendding []string
	js.ForEach(func(key, value gjson.Result) bool {
		k := key.String()
		typeValue := "string"
		if value.Type == gjson.Number {
			typeValue = "int"
		}
		pendding = append(pendding, fmt.Sprintf(template,
			transferName(k),
			typeValue,
			k,
		))
		return true
	})
	if len(pendding) == 0 {
		return ""
	}

	rs := fmt.Sprintf("type %sData struct {\n %s \n}", name, strings.Join(pendding, "\n"))
	return rs
}

func transferName(s string) string {
	s = strings.ToLower(s)
	arr := strings.Split(s, "_")
	var data []string
	for _, v := range arr {
		data = append(data, strings.Title(v))
	}
	return strings.Join(data, "")
}
