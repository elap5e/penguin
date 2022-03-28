// Copyright 2022 Elapse and contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"encoding/json"
	"go/format"
	"html/template"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Face struct {
	ID   string `json:"QSid"`
	Name string `json:"QDes"`
	Safe string
}

type ByID []Face

func (s ByID) Len() int { return len(s) }
func (s ByID) Less(i, j int) bool {
	ii, _ := strconv.Atoi(s[i].ID)
	jj, _ := strconv.Atoi(s[j].ID)
	return ii < jj
}
func (s ByID) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func main() {
	faces := []Face{}
	err := json.Unmarshal([]byte(jsonFaceConfig), &faces)
	if err != nil {
		panic(err)
	}

	for i := range faces {
		faces[i].Name = strings.TrimPrefix(faces[i].Name, "/")
		faces[i].Safe = strings.ToUpper(
			strings.NewReplacer("\"", "", "\\", "").
				Replace(strconv.QuoteToASCII(faces[i].Name)),
		)
	}
	sort.Sort(ByID(faces))

	tmpl, err := template.New("face").Parse(tmplEmoticonFace)
	if err != nil {
		panic(err)
	}
	buf := bytes.Buffer{}
	err = tmpl.Execute(&buf, faces)
	if err != nil {
		panic(err)
	}
	data, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile("daemon/message/face/face.go", os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		panic(err)
	}
}

const tmplEmoticonFace = `// Code generated by message-face-gen, DO NOT EDIT.

package face

import (
	"fmt"
)

type FaceType int16

const (
	FaceTypeUNKNOWN FaceType = -1 // "/未知表情"{{ range . }}
	FaceType{{ .Safe }} FaceType = {{ .ID }} // "/{{ .Name }}"{{ end }}
)

func (t FaceType) String() string {
	switch t {
	case FaceTypeUNKNOWN:
		return "/未知表情"{{ range . }}
	case FaceType{{ .Safe }}:
		return "/{{ .Name }}"{{ end }}
	}
	return "/未知表情"
}

func ParseFaceType(s string) (FaceType, error) {
	switch s {
	case "/未知表情":
		return FaceTypeUNKNOWN, nil{{ range . }}
	case "/{{ .Name }}":
		return FaceType{{ .Safe }}, nil{{ end }}
	}
	return FaceTypeUNKNOWN, fmt.Errorf("Unknown Face String: '%s', defaulting to FaceTypeUNKNOWN", s)
}
`

const jsonFaceConfig = `[
{"QSid":"14","QDes":"/微笑","IQLid":"23","AQLid":"23","EMCode":"100"},
{"QSid":"1","QDes":"/撇嘴","IQLid":"40","AQLid":"40","EMCode":"101"},
{"QSid":"2","QDes":"/色","IQLid":"19","AQLid":"19","EMCode":"102"},
{"QSid":"3","QDes":"/发呆","IQLid":"43","AQLid":"43","EMCode":"103"},
{"QSid":"4","QDes":"/得意","IQLid":"21","AQLid":"21","EMCode":"104"},
{"QSid":"6","QDes":"/害羞","IQLid":"20","AQLid":"20","EMCode":"106"},
{"QSid":"7","QDes":"/闭嘴","IQLid":"104","AQLid":"106","EMCode":"107"},
{"QSid":"8","QDes":"/睡","IQLid":"35","AQLid":"35","EMCode":"108"},
{"QSid":"9","QDes":"/大哭","IQLid":"10","AQLid":"10","EMCode":"109"},
{"QSid":"5","QDes":"/流泪","IQLid":"9","AQLid":"9","EMCode":"105","AniStickerType":1,"AniStickerPackId":"1","AniStickerId":"16"},
{"QSid":"10","QDes":"/尴尬","IQLid":"25","AQLid":"25","EMCode":"110"},
{"QSid":"11","QDes":"/发怒","IQLid":"24","AQLid":"24","EMCode":"111"},
{"QSid":"12","QDes":"/调皮","IQLid":"1","AQLid":"1","EMCode":"112"},
{"QSid":"13","QDes":"/呲牙","IQLid":"0","AQLid":"0","EMCode":"113"},
{"QSid":"0","QDes":"/惊讶","IQLid":"33","AQLid":"33","EMCode":"114"},
{"QSid":"15","QDes":"/难过","IQLid":"32","AQLid":"32","EMCode":"115"},
{"QSid":"16","QDes":"/酷","IQLid":"12","AQLid":"12","EMCode":"116"},
{"QSid":"96","QDes":"/冷汗","IQLid":"27","AQLid":"27","EMCode":"117"},
{"QSid":"18","QDes":"/抓狂","IQLid":"13","AQLid":"13","EMCode":"118"},
{"QSid":"19","QDes":"/吐","IQLid":"22","AQLid":"22","EMCode":"119"},
{"QSid":"20","QDes":"/偷笑","IQLid":"3","AQLid":"3","EMCode":"120"},
{"QSid":"21","QDes":"/可爱","IQLid":"18","AQLid":"18","EMCode":"121"},
{"QSid":"22","QDes":"/白眼","IQLid":"30","AQLid":"30","EMCode":"122"},
{"QSid":"23","QDes":"/傲慢","IQLid":"31","AQLid":"31","EMCode":"123"},
{"QSid":"24","QDes":"/饥饿","IQLid":"79","AQLid":"81","EMCode":"124"},
{"QSid":"25","QDes":"/困","IQLid":"80","AQLid":"82","EMCode":"125"},
{"QSid":"26","QDes":"/惊恐","IQLid":"26","AQLid":"26","EMCode":"126"},
{"QSid":"27","QDes":"/流汗","IQLid":"2","AQLid":"2","EMCode":"127"},
{"QSid":"28","QDes":"/憨笑","IQLid":"37","AQLid":"37","EMCode":"128"},
{"QSid":"29","QDes":"/悠闲","IQLid":"50","AQLid":"50","EMCode":"129"},
{"QSid":"30","QDes":"/奋斗","IQLid":"42","AQLid":"42","EMCode":"130"},
{"QSid":"31","QDes":"/咒骂","IQLid":"81","AQLid":"83","EMCode":"131"},
{"QSid":"32","QDes":"/疑问","IQLid":"34","AQLid":"34","EMCode":"132"},
{"QSid":"33","QDes":"/嘘","IQLid":"11","AQLid":"11","EMCode":"133"},
{"QSid":"34","QDes":"/晕","IQLid":"49","AQLid":"49","EMCode":"134"},
{"QSid":"35","QDes":"/折磨","IQLid":"82","AQLid":"84","EMCode":"135"},
{"QSid":"36","QDes":"/衰","isStatic":"1","IQLid":"39","AQLid":"39","EMCode":"136"},
{"QSid":"37","QDes":"/骷髅","isStatic":"1","IQLid":"76","AQLid":"78","EMCode":"137"},
{"QSid":"38","QDes":"/敲打","IQLid":"5","AQLid":"5","EMCode":"138"},
{"QSid":"39","QDes":"/再见","IQLid":"4","AQLid":"4","EMCode":"139"},
{"QSid":"97","QDes":"/擦汗","IQLid":"6","AQLid":"6","EMCode":"140"},
{"QSid":"98","QDes":"/抠鼻","IQLid":"83","AQLid":"85","EMCode":"141"},
{"QSid":"99","QDes":"/鼓掌","IQLid":"84","AQLid":"86","EMCode":"142"},
{"QSid":"100","QDes":"/糗大了","IQLid":"85","AQLid":"87","EMCode":"143"},
{"QSid":"101","QDes":"/坏笑","IQLid":"46","AQLid":"46","EMCode":"144"},
{"QSid":"102","QDes":"/左哼哼","IQLid":"86","AQLid":"88","EMCode":"145"},
{"QSid":"103","QDes":"/右哼哼","IQLid":"44","AQLid":"44","EMCode":"146"},
{"QSid":"104","QDes":"/哈欠","IQLid":"87","AQLid":"89","EMCode":"147"},
{"QSid":"105","QDes":"/鄙视","IQLid":"48","AQLid":"48","EMCode":"148"},
{"QSid":"106","QDes":"/委屈","IQLid":"14","AQLid":"14","EMCode":"149"},
{"QSid":"107","QDes":"/快哭了","IQLid":"88","AQLid":"90","EMCode":"150"},
{"QSid":"108","QDes":"/阴险","IQLid":"41","AQLid":"41","EMCode":"151"},
{"QSid":"305","QDes":"/右亲亲","IQLid":"305","AQLid":"305","EMCode":"10305"},
{"QSid":"109","QDes":"/左亲亲","IQLid":"36","AQLid":"36","EMCode":"152"},
{"QSid":"110","QDes":"/吓","IQLid":"89","AQLid":"91","EMCode":"153"},
{"QSid":"111","QDes":"/可怜","IQLid":"51","AQLid":"51","EMCode":"154"},
{"QSid":"172","QDes":"/眨眼睛","IQLid":"142","AQLid":"164","EMCode":"242"},
{"QSid":"182","QDes":"/笑哭","IQLid":"152","AQLid":"174","EMCode":"252"},
{"QSid":"179","QDes":"/doge","IQLid":"149","AQLid":"171","EMCode":"249"},
{"QSid":"173","QDes":"/泪奔","IQLid":"143","AQLid":"165","EMCode":"243"},
{"QSid":"174","QDes":"/无奈","IQLid":"144","AQLid":"166","EMCode":"244"},
{"QSid":"212","QDes":"/托腮","IQLid":"182","AQLid":"161","EMCode":"282"},
{"QSid":"175","QDes":"/卖萌","IQLid":"145","AQLid":"167","EMCode":"245"},
{"QSid":"178","QDes":"/斜眼笑","IQLid":"148","AQLid":"170","EMCode":"248"},
{"QSid":"177","QDes":"/喷血","IQLid":"147","AQLid":"169","EMCode":"247"},
{"QSid":"180","QDes":"/惊喜","IQLid":"150","AQLid":"172","EMCode":"250"},
{"QSid":"181","QDes":"/骚扰","IQLid":"151","AQLid":"173","EMCode":"251"},
{"QSid":"176","QDes":"/小纠结","IQLid":"146","AQLid":"168","EMCode":"246"},
{"QSid":"183","QDes":"/我最美","IQLid":"153","AQLid":"175","EMCode":"253"},
{"QSid":"245","QDes":"/加油必胜","IQLid":"245","AQLid":"217","QHide":"1","EMCode":"202001"},
{"QSid":"246","QDes":"/加油抱抱","IQLid":"246","AQLid":"218","EMCode":"202002"},
{"QSid":"247","QDes":"/口罩护体","isStatic":"1","IQLid":"247","AQLid":"219","QHide":"1","EMCode":"202003"},
{"QSid":"260","QDes":"/搬砖中","isStatic":"1","IQLid":"260","AQLid":"260","QHide":"1","EMCode":"10260"},
{"QSid":"261","QDes":"/忙到飞起","IQLid":"261","AQLid":"261","QHide":"1","EMCode":"10261"},
{"QSid":"262","QDes":"/脑阔疼","IQLid":"262","AQLid":"262","EMCode":"10262"},
{"QSid":"263","QDes":"/沧桑","IQLid":"263","AQLid":"263","EMCode":"10263"},
{"QSid":"264","QDes":"/捂脸","IQLid":"264","AQLid":"264","EMCode":"10264"},
{"QSid":"265","QDes":"/辣眼睛","IQLid":"265","AQLid":"265","EMCode":"10265"},
{"QSid":"266","QDes":"/哦哟","IQLid":"266","AQLid":"266","EMCode":"10266"},
{"QSid":"267","QDes":"/头秃","IQLid":"267","AQLid":"267","EMCode":"10267"},
{"QSid":"268","QDes":"/问号脸","IQLid":"268","AQLid":"268","EMCode":"10268"},
{"QSid":"269","QDes":"/暗中观察","IQLid":"269","AQLid":"269","EMCode":"10269"},
{"QSid":"270","QDes":"/emm","IQLid":"270","AQLid":"270","EMCode":"10270"},
{"QSid":"271","QDes":"/吃瓜","IQLid":"271","AQLid":"271","EMCode":"10271"},
{"QSid":"272","QDes":"/呵呵哒","IQLid":"272","AQLid":"272","EMCode":"10272"},
{"QSid":"277","QDes":"/汪汪","isStatic":"1","IQLid":"277","AQLid":"277","EMCode":"10277"},
{"QSid":"307","QDes":"/喵喵","isStatic":"1","IQLid":"307","AQLid":"307","EMCode":"10307"},
{"QSid":"306","QDes":"/牛气冲天","isStatic":"1","IQLid":"306","AQLid":"306","EMCode":"10306"},
{"QSid":"281","QDes":"/无眼笑","IQLid":"281","AQLid":"281","EMCode":"10281"},
{"QSid":"282","QDes":"/敬礼","IQLid":"282","AQLid":"282","EMCode":"10282"},
{"QSid":"283","QDes":"/狂笑","IQLid":"283","AQLid":"283","EMCode":"10283"},
{"QSid":"284","QDes":"/面无表情","IQLid":"284","AQLid":"284","EMCode":"10284"},
{"QSid":"285","QDes":"/摸鱼","IQLid":"285","AQLid":"285","EMCode":"10285"},
{"QSid":"293","QDes":"/摸锦鲤","IQLid":"293","AQLid":"293","EMCode":"10293"},
{"QSid":"286","QDes":"/魔鬼笑","IQLid":"286","AQLid":"286","EMCode":"10286"},
{"QSid":"287","QDes":"/哦","IQLid":"287","AQLid":"287","EMCode":"10287"},
{"QSid":"288","QDes":"/请","IQLid":"288","AQLid":"288","EMCode":"10288"},
{"QSid":"289","QDes":"/睁眼","IQLid":"289","AQLid":"289","EMCode":"10289"},
{"QSid":"294","QDes":"/期待","IQLid":"294","AQLid":"294","EMCode":"10294"},
{"QSid":"295","QDes":"/拿到红包","IQLid":"295","AQLid":"295","QHide":"1","EMCode":"10295"},
{"QSid":"296","QDes":"/真好","IQLid":"296","AQLid":"296","QHide":"1","EMCode":"10296"},
{"QSid":"297","QDes":"/拜谢","IQLid":"297","AQLid":"297","EMCode":"10297"},
{"QSid":"298","QDes":"/元宝","IQLid":"298","AQLid":"298","EMCode":"10298"},
{"QSid":"299","QDes":"/牛啊","IQLid":"299","AQLid":"299","EMCode":"10299"},
{"QSid":"300","QDes":"/胖三斤","IQLid":"300","AQLid":"300","EMCode":"10300"},
{"QSid":"301","QDes":"/好闪","IQLid":"301","AQLid":"301","EMCode":"10301"},
{"QSid":"303","QDes":"/右拜年","IQLid":"303","AQLid":"303","QHide":"1","EMCode":"10303"},
{"QSid":"302","QDes":"/左拜年","IQLid":"302","AQLid":"302","QHide":"1","EMCode":"10302"},
{"QSid":"304","QDes":"/红包包","IQLid":"304","AQLid":"304","QHide":"1","EMCode":"10304"},
{"QSid":"322","QDes":"/拒绝","IQLid":"322","AQLid":"322","EMCode":"10322"},
{"QSid":"323","QDes":"/嫌弃","IQLid":"323","AQLid":"323","EMCode":"10323"},
{"QSid":"311","QDes":"/打call","IQLid":"311","AQLid":"311","EMCode":"10311","AniStickerType":1,"AniStickerPackId":"1","AniStickerId":"1"},
{"QSid":"312","QDes":"/变形","IQLid":"312","AQLid":"312","EMCode":"10312","AniStickerType":1,"AniStickerPackId":"1","AniStickerId":"2"},
{"QSid":"313","QDes":"/嗑到了","IQLid":"313","AQLid":"313","QHide":"1","EMCode":"10313","AniStickerType":1,"AniStickerPackId":"1","AniStickerId":"3"},
{"QSid":"314","QDes":"/仔细分析","IQLid":"314","AQLid":"314","EMCode":"10314","AniStickerType":1,"AniStickerPackId":"1","AniStickerId":"4"},
{"QSid":"315","QDes":"/加油","IQLid":"315","AQLid":"315","EMCode":"10315","AniStickerType":1,"AniStickerPackId":"1","AniStickerId":"5"},
{"QSid":"316","QDes":"/我没事","IQLid":"316","AQLid":"316","QHide":"1","EMCode":"10316","AniStickerType":1,"AniStickerPackId":"1","AniStickerId":"6"},
{"QSid":"317","QDes":"/菜汪","IQLid":"317","AQLid":"317","EMCode":"10317","AniStickerType":1,"AniStickerPackId":"1","AniStickerId":"7"},
{"QSid":"318","QDes":"/崇拜","IQLid":"318","AQLid":"318","EMCode":"10318","AniStickerType":1,"AniStickerPackId":"1","AniStickerId":"8"},
{"QSid":"319","QDes":"/比心","IQLid":"319","AQLid":"319","EMCode":"10319","AniStickerType":1,"AniStickerPackId":"1","AniStickerId":"9"},
{"QSid":"320","QDes":"/庆祝","IQLid":"320","AQLid":"320","EMCode":"10320","AniStickerType":1,"AniStickerPackId":"1","AniStickerId":"10"},
{"QSid":"321","QDes":"/老色痞","IQLid":"321","AQLid":"321","QHide":"1","EMCode":"10321","AniStickerType":1,"AniStickerPackId":"1","AniStickerId":"11"},
{"QSid":"324","QDes":"/吃糖","IQLid":"324","AQLid":"324","EMCode":"10324","AniStickerType":1,"AniStickerPackId":"1","AniStickerId":"12"},
{"QSid":"325","QDes":"/惊吓","IQLid":"325","AQLid":"325","EMCode":"10325","AniStickerType":1,"AniStickerPackId":"1","AniStickerId":"14"},
{"QSid":"326","QDes":"/生气","IQLid":"326","AQLid":"326","EMCode":"10326","AniStickerType":1,"AniStickerPackId":"1","AniStickerId":"15"},
{"QSid":"53","QDes":"/蛋糕","IQLid":"59","AQLid":"59","EMCode":"168","AniStickerType":1,"AniStickerPackId":"1","AniStickerId":"17"},
{"QSid":"114","QDes":"/篮球","IQLid":"90","AQLid":"92","EMCode":"158","AniStickerType":2,"AniStickerPackId":"1","AniStickerId":"13"},
{"QSid":"327","QDes":"/加一","QHide":"1","IQLid":"327","AQLid":"327","EMCode":"10327"},
{"QSid":"328","QDes":"/错号","QHide":"1","IQLid":"328","AQLid":"328","EMCode":"10328"},
{"QSid":"329","QDes":"/对号","QHide":"1","IQLid":"329","AQLid":"329","EMCode":"10329"},
{"QSid":"330","QDes":"/完成","QHide":"1","IQLid":"330","AQLid":"330","EMCode":"10330"},
{"QSid":"331","QDes":"/明白","QHide":"1","IQLid":"331","AQLid":"331","EMCode":"10331"},
{"QSid":"49","QDes":"/拥抱","IQLid":"45","AQLid":"45","EMCode":"178"},
{"QSid":"66","QDes":"/爱心","isStatic":"1","IQLid":"28","AQLid":"28","EMCode":"166"},
{"QSid":"63","QDes":"/玫瑰","isStatic":"1","IQLid":"8","AQLid":"8","EMCode":"163"},
{"QSid":"64","QDes":"/凋谢","isStatic":"1","IQLid":"57","AQLid":"57","EMCode":"164"},
{"QSid":"187","QDes":"/幽灵","IQLid":"157","AQLid":"179","EMCode":"257"},
{"QSid":"146","QDes":"/爆筋","isStatic":"1","IQLid":"116","AQLid":"118","EMCode":"121011"},
{"QSid":"116","QDes":"/示爱","IQLid":"29","AQLid":"29","EMCode":"165"},
{"QSid":"67","QDes":"/心碎","IQLid":"72","AQLid":"74","EMCode":"167"},
{"QSid":"60","QDes":"/咖啡","IQLid":"66","AQLid":"66","EMCode":"160"},
{"QSid":"185","QDes":"/羊驼","IQLid":"155","AQLid":"177","EMCode":"255"},
{"QSid":"192","QDes":"/红包","IQLid":"162","AQLid":"184","QHide":"1","EMCode":"262"},
{"QSid":"137","QDes":"/鞭炮","isStatic":"1","IQLid":"107","AQLid":"109","EMCode":"121002"},
{"QSid":"138","QDes":"/灯笼","isStatic":"1","IQLid":"108","AQLid":"110","QHide":"1","EMCode":"121003"},
{"QSid":"136","QDes":"/双喜","isStatic":"1","IQLid":"106","AQLid":"108","QHide":"1","EMCode":"121001"},
{"QSid":"76","QDes":"/赞","IQLid":"52","AQLid":"52","EMCode":"179"},
{"QSid":"124","QDes":"/OK","isStatic":"1","IQLid":"64","AQLid":"64","EMCode":"189"},
{"QSid":"118","QDes":"/抱拳","IQLid":"56","AQLid":"56","EMCode":"183"},
{"QSid":"78","QDes":"/握手","IQLid":"54","AQLid":"54","EMCode":"181"},
{"QSid":"119","QDes":"/勾引","IQLid":"63","AQLid":"63","EMCode":"184"},
{"QSid":"79","QDes":"/胜利","IQLid":"55","AQLid":"55","EMCode":"182"},
{"QSid":"120","QDes":"/拳头","IQLid":"71","AQLid":"73","EMCode":"185"},
{"QSid":"121","QDes":"/差劲","IQLid":"70","AQLid":"72","EMCode":"186"},
{"QSid":"77","QDes":"/踩","IQLid":"53","AQLid":"53","EMCode":"180"},
{"QSid":"122","QDes":"/爱你","IQLid":"65","AQLid":"65","EMCode":"187"},
{"QSid":"123","QDes":"/NO","IQLid":"92","AQLid":"94","EMCode":"188"},
{"QSid":"201","QDes":"/点赞","IQLid":"171","AQLid":"150","EMCode":"271"},
{"QSid":"203","QDes":"/托脸","IQLid":"173","AQLid":"152","EMCode":"273"},
{"QSid":"204","QDes":"/吃","IQLid":"174","AQLid":"153","EMCode":"274"},
{"QSid":"202","QDes":"/无聊","IQLid":"172","AQLid":"151","EMCode":"272"},
{"QSid":"200","QDes":"/拜托","IQLid":"170","AQLid":"149","EMCode":"270"},
{"QSid":"194","QDes":"/不开心","IQLid":"164","AQLid":"143","EMCode":"264"},
{"QSid":"193","QDes":"/大笑","IQLid":"163","AQLid":"185","EMCode":"263"},
{"QSid":"197","QDes":"/冷漠","IQLid":"167","AQLid":"146","QHide":"1","EMCode":"267"},
{"QSid":"211","QDes":"/我不看","IQLid":"181","AQLid":"160","EMCode":"281"},
{"QSid":"210","QDes":"/飙泪","IQLid":"180","AQLid":"159","EMCode":"280"},
{"QSid":"198","QDes":"/呃","IQLid":"168","AQLid":"147","EMCode":"268"},
{"QSid":"199","QDes":"/好棒","IQLid":"169","AQLid":"148","QHide":"1","EMCode":"269"},
{"QSid":"207","QDes":"/花痴","IQLid":"177","AQLid":"156","QHide":"1","EMCode":"277"},
{"QSid":"205","QDes":"/送花","IQLid":"175","AQLid":"154","QHide":"1","EMCode":"275"},
{"QSid":"206","QDes":"/害怕","IQLid":"176","AQLid":"155","EMCode":"276"},
{"QSid":"208","QDes":"/小样儿","IQLid":"178","AQLid":"157","QHide":"1","EMCode":"278"},
{"QSid":"308","QDes":"/求红包","IQLid":"308","isCMEmoji":"1","AQLid":"308","QHide":"1","EMCode":"20243"},
{"QSid":"309","QDes":"/谢红包","IQLid":"309","isCMEmoji":"1","AQLid":"309","QHide":"1","EMCode":"20244"},
{"QSid":"310","QDes":"/新年烟花","IQLid":"310","isCMEmoji":"1","AQLid":"310","QHide":"1","EMCode":"20245"},
{"QSid":"290","QDes":"/敲开心","IQLid":"290","isCMEmoji":"1","AQLid":"290","EMCode":"20240"},
{"QSid":"291","QDes":"/震惊","IQLid":"291","isCMEmoji":"1","AQLid":"291","QHide":"1","EMCode":"20241"},
{"QSid":"292","QDes":"/让我康康","IQLid":"292","isCMEmoji":"1","AQLid":"292","EMCode":"20242"},
{"QSid":"226","QDes":"/拍桌","IQLid":"196","isCMEmoji":"1","AQLid":"198","EMCode":"297"},
{"QSid":"215","QDes":"/糊脸","IQLid":"185","isCMEmoji":"1","AQLid":"187","EMCode":"285"},
{"QSid":"237","QDes":"/偷看","IQLid":"207","isCMEmoji":"1","AQLid":"209","EMCode":"307"},
{"QSid":"214","QDes":"/啵啵","IQLid":"184","isCMEmoji":"1","AQLid":"186","EMCode":"284"},
{"QSid":"235","QDes":"/颤抖","IQLid":"205","isCMEmoji":"1","AQLid":"207","EMCode":"305"},
{"QSid":"222","QDes":"/抱抱","IQLid":"192","isCMEmoji":"1","AQLid":"194","EMCode":"292"},
{"QSid":"217","QDes":"/扯一扯","IQLid":"187","isCMEmoji":"1","AQLid":"189","EMCode":"287"},
{"QSid":"221","QDes":"/顶呱呱","IQLid":"191","isCMEmoji":"1","AQLid":"193","EMCode":"291"},
{"QSid":"225","QDes":"/撩一撩","IQLid":"195","isCMEmoji":"1","AQLid":"197","EMCode":"296"},
{"QSid":"241","QDes":"/生日快乐","IQLid":"211","isCMEmoji":"1","AQLid":"213","EMCode":"311"},
{"QSid":"227","QDes":"/拍手","IQLid":"197","isCMEmoji":"1","AQLid":"199","EMCode":"294"},
{"QSid":"238","QDes":"/扇脸","IQLid":"208","isCMEmoji":"1","AQLid":"210","EMCode":"308"},
{"QSid":"240","QDes":"/喷脸","IQLid":"210","isCMEmoji":"1","AQLid":"212","EMCode":"310"},
{"QSid":"229","QDes":"/干杯","IQLid":"199","isCMEmoji":"1","AQLid":"201","EMCode":"299"},
{"QSid":"216","QDes":"/拍头","IQLid":"186","isCMEmoji":"1","AQLid":"188","EMCode":"286"},
{"QSid":"218","QDes":"/舔一舔","IQLid":"188","isCMEmoji":"1","AQLid":"190","EMCode":"288"},
{"QSid":"233","QDes":"/掐一掐","IQLid":"203","isCMEmoji":"1","AQLid":"205","EMCode":"303"},
{"QSid":"219","QDes":"/蹭一蹭","IQLid":"189","isCMEmoji":"1","AQLid":"191","EMCode":"289"},
{"QSid":"244","QDes":"/扔狗","IQLid":"214","isCMEmoji":"1","AQLid":"216","EMCode":"312"},
{"QSid":"232","QDes":"/佛系","IQLid":"202","isCMEmoji":"1","AQLid":"204","EMCode":"302"},
{"QSid":"243","QDes":"/甩头","IQLid":"213","isCMEmoji":"1","AQLid":"215","EMCode":"313"},
{"QSid":"223","QDes":"/暴击","IQLid":"193","isCMEmoji":"1","AQLid":"195","EMCode":"293"},
{"QSid":"279","QDes":"/打脸","IQLid":"279","isCMEmoji":"1","AQLid":"279","QHide":"1","EMCode":"20238"},
{"QSid":"280","QDes":"/击掌","IQLid":"280","isCMEmoji":"1","AQLid":"280","QHide":"1","EMCode":"20239"},
{"QSid":"231","QDes":"/哼","IQLid":"201","isCMEmoji":"1","AQLid":"203","EMCode":"301"},
{"QSid":"224","QDes":"/开枪","IQLid":"194","isCMEmoji":"1","AQLid":"196","EMCode":"295"},
{"QSid":"278","QDes":"/汗","IQLid":"278","isCMEmoji":"1","AQLid":"278","EMCode":"20237"},
{"QSid":"236","QDes":"/啃头","IQLid":"206","isCMEmoji":"1","AQLid":"208","QHide":"1","EMCode":"306"},
{"QSid":"228","QDes":"/恭喜","IQLid":"198","isCMEmoji":"1","AQLid":"200","QHide":"1","EMCode":"298"},
{"QSid":"220","QDes":"/拽炸天","IQLid":"190","isCMEmoji":"1","AQLid":"192","QHide":"1","EMCode":"290"},
{"QSid":"239","QDes":"/原谅","IQLid":"209","isCMEmoji":"1","AQLid":"211","EMCode":"309"},
{"QSid":"242","QDes":"/头撞击","IQLid":"212","isCMEmoji":"1","AQLid":"214","QHide":"1","EMCode":"314"},
{"QSid":"230","QDes":"/嘲讽","IQLid":"200","isCMEmoji":"1","AQLid":"202","EMCode":"300"},
{"QSid":"234","QDes":"/惊呆","IQLid":"204","isCMEmoji":"1","AQLid":"206","QHide":"1","EMCode":"304"},
{"QSid":"273","QDes":"/我酸了","isStatic":"1","IQLid":"273","AQLid":"273","EMCode":"10273"},
{"QSid":"75","QDes":"/月亮","isStatic":"1","IQLid":"67","AQLid":"68","EMCode":"175"},
{"QSid":"74","QDes":"/太阳","isStatic":"1","IQLid":"73","AQLid":"75","EMCode":"176"},
{"QSid":"46","QDes":"/猪头","isStatic":"1","IQLid":"7","AQLid":"7","EMCode":"162"},
{"QSid":"112","QDes":"/菜刀","IQLid":"17","AQLid":"17","EMCode":"155"},
{"QSid":"56","QDes":"/刀","IQLid":"68","AQLid":"70","EMCode":"171"},
{"QSid":"169","QDes":"/手枪","isStatic":"1","IQLid":"139","AQLid":"141","EMCode":"121034"},
{"QSid":"171","QDes":"/茶","IQLid":"141","AQLid":"163","EMCode":"241"},
{"QSid":"59","QDes":"/便便","IQLid":"15","AQLid":"15","EMCode":"174"},
{"QSid":"144","QDes":"/喝彩","isStatic":"1","IQLid":"114","AQLid":"116","EMCode":"121009"},
{"QSid":"147","QDes":"/棒棒糖","isStatic":"1","IQLid":"117","AQLid":"119","EMCode":"121012"},
{"QSid":"89","QDes":"/西瓜","isStatic":"1","IQLid":"60","AQLid":"60","EMCode":"156"},
{"QSid":"61","QDes":"/饭","isStatic":"1","IQLid":"58","AQLid":"58","QHide":"1","EMCode":"161"},
{"QSid":"148","QDes":"/喝奶","isStatic":"1","IQLid":"118","AQLid":"120","QHide":"1","EMCode":"121013"},
{"QSid":"274","QDes":"/太南了","isStatic":"1","IQLid":"274","AQLid":"274","QHide":"1","EMCode":"10274"},
{"QSid":"113","QDes":"/啤酒","IQLid":"61","AQLid":"61","QHide":"1","EMCode":"157"},
{"QSid":"140","QDes":"/K歌","isStatic":"1","IQLid":"110","AQLid":"112","QHide":"1","EMCode":"121005"},
{"QSid":"188","QDes":"/蛋","IQLid":"158","AQLid":"180","QHide":"1","EMCode":"258"},
{"QSid":"55","QDes":"/炸弹","isStatic":"1","IQLid":"16","AQLid":"16","QHide":"1","EMCode":"170"},
{"QSid":"184","QDes":"/河蟹","IQLid":"154","AQLid":"176","QHide":"1","EMCode":"254"},
{"QSid":"158","QDes":"/钞票","isStatic":"1","IQLid":"128","AQLid":"130","QHide":"1","EMCode":"121023"},
{"QSid":"54","QDes":"/闪电","isStatic":"1","IQLid":"78","AQLid":"80","QHide":"1","EMCode":"169"},
{"QSid":"69","QDes":"/礼物","isStatic":"1","IQLid":"74","AQLid":"76","QHide":"1","EMCode":"177"},
{"QSid":"190","QDes":"/菊花","IQLid":"160","AQLid":"182","QHide":"1","EMCode":"260"},
{"QSid":"151","QDes":"/飞机","isStatic":"1","IQLid":"121","AQLid":"123","QHide":"1","EMCode":"121016"},
{"QSid":"145","QDes":"/祈祷","isStatic":"1","IQLid":"115","AQLid":"117","QHide":"1","EMCode":"121010"},
{"QSid":"117","QDes":"/瓢虫","IQLid":"62","AQLid":"62","QHide":"1","EMCode":"173"},
{"QSid":"168","QDes":"/药","isStatic":"1","IQLid":"138","AQLid":"140","QHide":"1","EMCode":"121033"},
{"QSid":"115","QDes":"/乒乓","IQLid":"91","AQLid":"93","QHide":"1","EMCode":"159"},
{"QSid":"57","QDes":"/足球","IQLid":"75","AQLid":"77","QHide":"1","EMCode":"172"},
{"QSid":"41","QDes":"/发抖","isStatic":"1","IQLid":"69","AQLid":"71","EMCode":"193"},
{"QSid":"125","QDes":"/转圈","IQLid":"95","AQLid":"97","EMCode":"195"},
{"QSid":"42","QDes":"/爱情","IQLid":"38","AQLid":"38","EMCode":"190"},
{"QSid":"43","QDes":"/跳跳","IQLid":"93","AQLid":"95","EMCode":"192"},
{"QSid":"86","QDes":"/怄火","IQLid":"94","AQLid":"96","EMCode":"194"},
{"QSid":"129","QDes":"/挥手","IQLid":"77","AQLid":"79","EMCode":"199"},
{"QSid":"85","QDes":"/飞吻","isStatic":"1","IQLid":"47","AQLid":"47","EMCode":"191"},
{"QSid":"126","QDes":"/磕头","IQLid":"96","AQLid":"98","QHide":"1","EMCode":"196"},
{"QSid":"128","QDes":"/跳绳","IQLid":"98","AQLid":"100","QHide":"1","EMCode":"198"},
{"QSid":"130","QDes":"/激动","IQLid":"99","AQLid":"101","QHide":"1","EMCode":"200"},
{"QSid":"127","QDes":"/回头","IQLid":"97","AQLid":"99","QHide":"1","EMCode":"197"},
{"QSid":"132","QDes":"/献吻","IQLid":"101","AQLid":"103","QHide":"1","EMCode":"202"},
{"QSid":"134","QDes":"/右太极","IQLid":"103","AQLid":"105","QHide":"1","EMCode":"204"},
{"QSid":"133","QDes":"/左太极","IQLid":"102","AQLid":"104","QHide":"1","EMCode":"203"},
{"QSid":"131","QDes":"/街舞","IQLid":"100","AQLid":"102","QHide":"1","EMCode":"201"},
{"QSid":"276","QDes":"/辣椒酱","isStatic":"1","IQLid":"276","AQLid":"276","QHide":"1","EMCode":"10276"}
]`
