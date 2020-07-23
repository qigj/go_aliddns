package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
)

var accessKeyId = flag.String("id", "", "阿里云accesskey")
var accessSecret = flag.String("secret", "", "阿里云accessSecret")
var domain = flag.String("rr", "", "要解析的域名前缀，例如：abc.example.com  那此值为abc")
var basedomain = flag.String("domain", "", "要解析的域名前缀，例如：example.com ")

type mrecord struct {
	Mrequestid string
	Mrr        string
	Mtype      string
	Mvalue     string
}

type ipadd struct {
	Ip string `json:"ip"`
}

func main() {
	flag.Parse()
	if *accessSecret == "" {
		panic("参数不全")
	}
	if *accessKeyId == "" {
		panic("参数不全")

	}
	if *basedomain == "" {
		panic("参数不全")

	}
	if *domain == "" {
		panic("参数不全")
	}
	uprecord := mrecord{}
	client, err := alidns.NewClientWithAccessKey("cn-hangzhou", *accessKeyId, *accessSecret)

	request := alidns.CreateDescribeDomainRecordsRequest()
	request.Scheme = "https"

	request.DomainName = *basedomain

	response, err := client.DescribeDomainRecords(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	// fmt.Printf("%v\n", response.DomainRecords.Record)
	for _, record := range response.DomainRecords.Record {
		if *domain == record.RR {
			// fmt.Printf("%v", record.RecordId)
			uprecord.Mrequestid = record.RecordId
			uprecord.Mrr = record.RR
			uprecord.Mtype = record.Type
			uprecord.Mvalue = record.Value
		}
	}
	// fmt.Printf("%v\n", uprecord)
	// result := getip()
	// ipaddress := ipadd{}
	// err = json.Unmarshal(result, &ipaddress)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("%v", ipaddress.Ip)
	// uprecord.updata(ipaddress.Ip)

	result := getip()
	// fmt.Printf("%s,%s", strings.Split(string(result), "\n")[0], uprecord.Mvalue)
	if uprecord.Mvalue != strings.Split(string(result), "\n")[0] {
		uprecord.updata(string(result))
	} else {
		fmt.Printf("%s", "当前ip地址与域名解析地址相同")
	}

}

func getjsonip() []byte {
	reponse, err := http.Get("https://api.ip.sb/jsonip")
	if err != nil {
		panic(err)
	}
	defer reponse.Body.Close()
	result, _ := ioutil.ReadAll(reponse.Body)
	// fmt.Printf("%s", result)
	return result
}

func getip() []byte {
	reponse, err := http.Get("https://api.ip.sb/ip")
	if err != nil {
		panic(err)
	}
	defer reponse.Body.Close()
	result, _ := ioutil.ReadAll(reponse.Body)
	// fmt.Printf("%s", result)
	return result
}

func (P *mrecord) updata(value string) {
	client, err := alidns.NewClientWithAccessKey("cn-hangzhou", *accessKeyId, *accessSecret)

	request := alidns.CreateUpdateDomainRecordRequest()
	request.Scheme = "https"

	request.RecordId = P.Mrequestid
	request.RR = P.Mrr
	request.Type = P.Mtype
	request.Value = value

	response, err := client.UpdateDomainRecord(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Printf("response is %#v\n", response)
}
