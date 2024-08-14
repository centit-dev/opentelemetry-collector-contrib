package httpbodyprocessor

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"testing"
)

type Person struct {
	Id      *int    `json:"id"`
	Name    string  `json:"name"`
	Age     int     `json:"age"`
	Address Address `json:"address"`
}

type Address struct {
	City    string `json:"city"`
	Country string `json:"country"`
	Room    Room   `json:"room"`
}

type Room struct {
	Unit   string `json:"unit"`
	Number int    `json:"number"`
}

var person = Person{
	Id:   nil,
	Name: "John Doe",
	Age:  30,
	Address: Address{
		City:    "New York",
		Country: "USA",
		Room: Room{
			Unit:   "A",
			Number: 101,
		},
	},
}

func TestFlattenJson(t *testing.T) {
	jsonString, err := json.Marshal(person)
	assert.Nil(t, err)

	var jsonObject map[string]interface{}
	err = json.Unmarshal(jsonString, &jsonObject)
	assert.Nil(t, err)

	output := make(map[string]string)
	flattenJson(jsonObject, "", &output)

	_, ok := output["id"]
	assert.False(t, ok)
	assert.Equal(t, "John Doe", output["name"])
	assert.Equal(t, "30", output["age"])
	assert.Equal(t, "New York", output["address.city"])
	assert.Equal(t, "USA", output["address.country"])
	assert.Equal(t, "A", output["address.room.unit"])
	assert.Equal(t, "101", output["address.room.number"])
}

var xmlString = `
<?xml version= 1.0encoding= uTF-8 ?>
<ProvBOSS>
    <Head>
        <TransactionId>C24572825869</TransactionId>
        <TranactionCode>SrSTEM</TranactionCode>
        <IssueDate>20240613121139</IssueDate>
        <EffectiveDate>20240613121139</EffectiveDate>
    </Head>
    <DATA> 
        <BRAID>O1</BRAID>
        <MSISDI>13910909344</MSISDI>
        <STATUS>00</STATUS>
        <NEWMSISDN></NEWMSISDN>
    </DATA>
</ProvBOSS>
`

func TestFlattenXML(t *testing.T) {
	decoder := xml.NewDecoder(bytes.NewReader([]byte(xmlString)))
	var root Node
	if err := decoder.Decode(&root); err != nil {
		t.Error(err)
	}
	// 平坦化处理
	flatMap := make(map[string]string)
	flattenNode(httpRequestBodyKey, root, flatMap)

	assert.Equal(t, flatMap["http.request.body.ProvBOSS.DATA.BRAID"], "O1")
	assert.Equal(t, flatMap["http.request.body.ProvBOSS.DATA.MSISDI"], "13910909344")
	assert.Equal(t, flatMap["http.request.body.ProvBOSS.DATA.STATUS"], "00")
	assert.Equal(t, flatMap["http.request.body.ProvBOSS.Head.EffectiveDate"], "20240613121139")
	assert.Equal(t, flatMap["http.request.body.ProvBOSS.Head.IssueDate"], "20240613121139")
	assert.Equal(t, flatMap["http.request.body.ProvBOSS.Head.TranactionCode"], "SrSTEM")
	assert.Equal(t, flatMap["http.request.body.ProvBOSS.Head.TransactionId"], "C24572825869")
}

func TestProcessUnknowContentType(t *testing.T) {
	attrs := pcommon.NewMap()
	attrs.PutStr(httpReqContentTypeKey, httpContentTypeXMLKey)
	attrs.PutStr(httpRequestBodyKey, xmlString)

	sli := attrs.PutEmptySlice(httpRespContentTypeKey)
	sliVal := sli.AppendEmpty()
	sliVal.FromRaw(httpContentTypeXMLKey)
	attrs.PutStr(httpResponseBodyKey, "abcdaaksjflaksjfoihow")

	p := NewProcessor(nil)

	p.processResponseBody(&attrs)

	value, ok := attrs.Get(HttpRespContextTypeAttrKey)
	assert.Equal(t, ok, true)
	assert.Equal(t, value.Str(), httpOtherAttrValue)
}

func TestProcessXml(t *testing.T) {
	attrs := pcommon.NewMap()
	attrs.PutStr(httpReqContentTypeKey, httpContentTypeXMLKey)
	attrs.PutStr(httpRequestBodyKey, xmlString)

	sli := attrs.PutEmptySlice(httpRespContentTypeKey)
	sliVal := sli.AppendEmpty()
	sliVal.FromRaw(httpContentTypeXMLKey)
	attrs.PutStr(httpResponseBodyKey, xmlString)

	p := NewProcessor(nil)

	p.processHttpRequestBody(&attrs)
	p.processResponseBody(&attrs)

	// assert request
	v, ok := attrs.Get("http.request.body.ProvBOSS.DATA.BRAID")
	assert.Equal(t, ok, true)
	assert.Equal(t, v.Str(), "O1")

	v, ok = attrs.Get("http.request.body.ProvBOSS.DATA.MSISDI")
	assert.Equal(t, ok, true)
	assert.Equal(t, v.Str(), "13910909344")

	v, ok = attrs.Get("http.request.body.ProvBOSS.DATA.STATUS")
	assert.Equal(t, ok, true)
	assert.Equal(t, v.Str(), "00")

	v, ok = attrs.Get("http.request.body.ProvBOSS.Head.EffectiveDate")
	assert.Equal(t, ok, true)
	assert.Equal(t, v.Str(), "20240613121139")

	v, ok = attrs.Get("http.request.body.ProvBOSS.Head.IssueDate")
	assert.Equal(t, ok, true)
	assert.Equal(t, v.Str(), "20240613121139")

	// assert response
	v, ok = attrs.Get("http.response.body.ProvBOSS.DATA.BRAID")
	assert.Equal(t, ok, true)
	assert.Equal(t, v.Str(), "O1")

	v, ok = attrs.Get("http.response.body.ProvBOSS.DATA.MSISDI")
	assert.Equal(t, ok, true)
	assert.Equal(t, v.Str(), "13910909344")

	v, ok = attrs.Get("http.response.body.ProvBOSS.DATA.STATUS")
	assert.Equal(t, ok, true)
	assert.Equal(t, v.Str(), "00")

	v, ok = attrs.Get("http.response.body.ProvBOSS.Head.EffectiveDate")
	assert.Equal(t, ok, true)
	assert.Equal(t, v.Str(), "20240613121139")

	v, ok = attrs.Get("http.response.body.ProvBOSS.Head.IssueDate")
	assert.Equal(t, ok, true)
	assert.Equal(t, v.Str(), "20240613121139")

	v, ok = attrs.Get(HttpRespContextTypeAttrKey)
	assert.Equal(t, ok, true)
	assert.Equal(t, v.Str(), httpXmlAttrValue)

	v, ok = attrs.Get(HttpReqContextTypeAttrKey)
	assert.Equal(t, ok, true)
	assert.Equal(t, v.Str(), httpXmlAttrValue)
}

func TestProcessJson(t *testing.T) {
	jsonString := "{\"id\":\"1\",\"shopId\":\"1\",\"userId\":\"1\",\"avatarUrl\":\"http://www.baidu.com\",\"name\":\"shush\",\"gender\":\"1\",\"mobilePhone\":\"13770321976\",\"level\":null,\"birthday\":null,\"country\":null,\"province\":null,\"city\":null,\"district\":null,\"labels\":null,\"memo\":null}{\"id\":\"1\",\"shopId\":\"1\",\"userId\":\"1\",\"avatarUrl\":\"http://www.baidu.com\",\"name\":\"shush\",\"gender\":\"1\",\"mobilePhone\":\"13770321976\",\"level\":null,\"birthday\":null,\"country\":null,\"province\":null,\"city\":null,\"district\":null,\"labels\":null,\"memo\":null}"
	object := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonString), &object)
	assert.Nil(t, err)
}

var jsonString = `
{
	"val": "hellp"
}
`

func TestParseHttpContentTypeByBody(t *testing.T) {
	jsonVal := pcommon.NewValueStr(jsonString)
	xmlVal := pcommon.NewValueStr(xmlString)
	type args struct {
		value pcommon.Value
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "json",
			args: args{jsonVal},
			want: httpContentTypeJsonKey,
		},
		{
			name: "xml",
			args: args{xmlVal},
			want: httpContentTypeXMLKey,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, ParseHttpContentTypeByBody(tt.args.value), "ParseHttpContentTypeByBody(%v)", tt.args.value)
		})
	}
}
