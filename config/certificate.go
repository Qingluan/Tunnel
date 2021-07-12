package config

import (
	"bufio"
	"bytes"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"strings"
	"time"
)

var (
	host       = "127.0.0.1"
	validFrom  = ""
	validFor   = 5 * 365 * 24 * time.Hour
	isCA       = false
	rsaBits    = 4096  //"Size of RSA key to generate. Ignored if --ecdsa-curve is set")
	ecdsaCurve = ""    // "ECDSA curve to use to generate a key. Valid values are P224, P256 (recommended), P384, P521")
	ed25519Key = false // "Generate an Ed25519 key")
	useRand    = true
	useKey     = ""
)

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	case ed25519.PrivateKey:
		return k.Public().(ed25519.PublicKey)
	default:
		return nil
	}
}

func ConfigCertificate(mhost, mvalidFrom, mecdsaCurve string, useRand, misCa, med25519Key bool, mvalidDateLength time.Duration, mrsaBits int) {
	host = mhost
	validFrom = mvalidFrom
	validFor = mvalidDateLength
	isCA = misCa
	rsaBits = mrsaBits
	ed25519Key = med25519Key
	ecdsaCurve = mecdsaCurve
}

func SetKey(key string) {
	useKey = key
	useRand = false
}

func SetHost(mhost string) {
	host = mhost

}
func SelfCertificate() (pemStringkeyString string) {
	// flag.Parse()

	if len(host) == 0 {
		log.Fatalf("Missing required --host parameter")
	}
	var randReader io.Reader
	if !useRand {
		kk := useKey
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		kk += "Do not go gentle into that good night.Good men, the last wave by, crying how brightTheir frail deeds might have danced in a green bay,Rage, rage against the dying of the light.Wild men who caught and sang the sun in flight,And learn, too late, they grieved it on its way,Do not go gentle into that good night.Grave men, near death, who see with blinding sightBlind eyes could blaze like meteors and be gay,Rage, rage against the dying of the light.And you, my father, there on the sad height,Curse, bless me now with your fierce tears, I pray.Do not go gentle into that good night.Rage, rage against the dying of the light.不要温和地走进那个良夜，老年应当在日暮时燃烧咆哮；怒斥，怒斥光明的消逝。虽然智慧的人临终时懂得黑暗有理，因为他们的话没有迸发出闪电，他们也并不温和地走进那个良夜。善良的人，当最后一浪过去，高呼他们脆弱的善行可能曾会多么光辉地在绿色的海湾里舞蹈，怒斥，怒斥光明的消逝。狂暴的人抓住并歌唱过翱翔的太阳，懂得，但为时太晚，他们使太阳在途中悲伤，也并不温和地走进那个良夜。严肃的人，接近死亡，用炫目的视觉看出失明的眼睛可以像流星一样闪耀欢欣，怒斥，怒斥光明的消逝。您啊，我的父亲．在那悲哀的高处．现在用您的热泪诅咒我，祝福我吧．我求您不要温和地走进那个良夜。怒斥，怒斥光明的消逝。"
		randReader = bufio.NewReader(bytes.NewBuffer([]byte(kk)))
	} else {
		randReader = rand.Reader
	}

	var priv interface{}
	var err error
	switch ecdsaCurve {
	case "":
		if ed25519Key {
			_, priv, err = ed25519.GenerateKey(randReader)
		} else {
			priv, err = rsa.GenerateKey(randReader, rsaBits)
		}
	case "P224":
		priv, err = ecdsa.GenerateKey(elliptic.P224(), randReader)
	case "P256":
		priv, err = ecdsa.GenerateKey(elliptic.P256(), randReader)
	case "P384":
		priv, err = ecdsa.GenerateKey(elliptic.P384(), randReader)
	case "P521":
		priv, err = ecdsa.GenerateKey(elliptic.P521(), randReader)
	default:
		log.Fatalf("Unrecognized elliptic curve: %q", ecdsaCurve)
	}
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}

	// ECDSA, ED25519 and RSA subject keys should have the DigitalSignature
	// KeyUsage bits set in the x509.Certificate template
	keyUsage := x509.KeyUsageDigitalSignature
	// Only RSA subject keys should have the KeyEncipherment KeyUsage bits set. In
	// the context of TLS this KeyUsage is particular to RSA key exchange and
	// authentication.
	if _, isRSA := priv.(*rsa.PrivateKey); isRSA {
		keyUsage |= x509.KeyUsageKeyEncipherment
	}

	var notBefore time.Time
	if len(validFrom) == 0 {
		notBefore = time.Now()
		if !useRand {
			notBefore = time.Date(notBefore.Year(), notBefore.Month(), notBefore.Day(), 0, 0, 0, 0, notBefore.Location())
		}
		log.Println(notBefore)
	} else {
		notBefore, err = time.Parse("Jan 2 15:04:05 2006", validFrom)
		if err != nil {
			log.Fatalf("Failed to parse creation date: %v", err)
		}
	}

	notAfter := notBefore.Add(validFor)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(randReader, serialNumberLimit)

	// if !useRand {

	// 	serialNumber = serialNumberLimit.Mul(big.NewInt(33), big.NewInt(1))
	// }

	if err != nil {
		log.Fatalf("Failed to generate serial number: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              keyUsage,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	hosts := strings.Split(host, ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	if isCA {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
	}

	derBytes, err := x509.CreateCertificate(randReader, &template, &template, publicKey(priv), priv)
	if err != nil {
		log.Fatalf("Failed to create certificate: %v", err)
	}

	certOut := bytes.NewBuffer([]byte{})
	if err != nil {
		log.Fatalf("Failed to open cert.pem for writing: %v", err)
	}
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		log.Fatalf("Failed to write data to cert.pem: %v", err)
	}
	// log.Print("wrote cert.pem\n")
	pemString := certOut.String()

	// keyOut, err := os.OpenFile("key.pem", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	keyOut := bytes.NewBuffer([]byte{})
	if err != nil {
		log.Fatalf("Failed to open key.pem for writing: %v", err)
		return
	}
	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		log.Fatalf("Unable to marshal private key: %v", err)
	}
	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		log.Fatalf("Failed to write data to key.pem: %v", err)
	}
	// if err := keyOut.Close(); err != nil {
	// 	log.Fatalf("Error closing key.pem: %v", err)
	// }
	keyString := keyOut.String()
	// log.Print("wrote key.pem\n")

	pemStringkeyString = base64.StdEncoding.EncodeToString([]byte("tls:" + pemString + "<SEP>" + keyString))
	return
}

func CreateCertificate(serverHost string, isCA bool) (configUri, base string) {
	ps := strings.SplitN(serverHost, ":", 2)
	host := ps[0]
	// hostPort := ps[1]
	if len(host) == 0 {
		log.Fatalf("Missing required --host parameter")
	}

	// var priv interface{}
	var err error
	validFor := 365 * 24 * time.Hour
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	pub := &priv.PublicKey
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}

	var notBefore time.Time
	notBefore = time.Now()

	notAfter := notBefore.Add(validFor)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 2048)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("Failed to generate serial number: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization:  []string{"Co"},
			Country:       []string{"JP"},
			Province:      []string{"Tokyo"},
			Locality:      []string{"Tokyo"},
			StreetAddress: []string{"Tasciko"},
			PostalCode:    []string{"10-200-4"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	hosts := strings.Split(host, ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	if isCA {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, pub, priv)
	if err != nil {
		log.Fatalf("Failed to create certificate: %v", err)
	}

	certOut := bytes.NewBuffer([]byte{})
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		log.Fatalf("Failed to write data to cert.pem: %v", err)
	}

	// log.Print("wrote cert.pem\n")

	keyOut := bytes.NewBuffer([]byte{})
	if err != nil {
		log.Fatalf("Failed to open key.pem for writing: %v", err)
		return
	}
	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		log.Fatalf("Unable to marshal private key: %v", err)
	}
	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		log.Fatalf("Failed to write data to key.pem: %v", err)
	}

	// log.Print("wrote key.pem\n")
	password := fmt.Sprintf("%s<SEP>%s", certOut.String(), keyOut.String())
	base = fmt.Sprintf("%s:%s@%s", "tls", password, serverHost)
	ssBody := base64.StdEncoding.EncodeToString([]byte(base))
	configUri = fmt.Sprintf("%s", ssBody)
	return
}
