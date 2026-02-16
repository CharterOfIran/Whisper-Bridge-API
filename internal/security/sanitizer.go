package security

import "regexp"

func Clean(msg string) string {
	// ۱. اول محدودیت طول (جلوگیری از پر شدن RAM)
	if len(msg) > 5000 {
		msg = msg[:5000]
	}

	// ۲. فقط اجازه دادن به کاراکترهای مجاز (Regex White-list)
	// این ریجکس حروف فارسی، انگلیسی، اعداد و علائم معمولی رو نگه میداره
	reg := regexp.MustCompile(`[^a-zA-Z0-9آ-ی\s.,!؟?()]+`)
	clean := reg.ReplaceAllString(msg, "")

	return clean
}
