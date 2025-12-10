package parse

import (
	"regexp"

	"github.com/mats9693/study/go/gocts/data"
)

var RequestRE = regexp.MustCompile(`const\s+URI_(\w+)\s*=\s*["']([\w/-]+)["']`)

func matchRequests(filename string, fileBytes []byte) {
	reMatched := RequestRE.FindAllSubmatch(fileBytes, -1)
	for i := range reMatched {
		// 没什么意义，因为匹配函数如果找到了匹配项，这里的长度就必然不小于3,但是因为后续要直接使用下标访问，所以还是判断一下，下同
		if len(reMatched[i]) < 3 {
			continue
		}

		requestName := string(reMatched[i][1])

		data.GeneratorIns.RequestAffiliation[filename] = append(data.GeneratorIns.RequestAffiliation[filename], requestName)
		data.GeneratorIns.Requests[requestName] = string(reMatched[i][2])
	}
}
