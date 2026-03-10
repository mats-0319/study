package scanner

import (
	"regexp"

	"github.com/mats0319/study/go/gocts/token"
)

var requestRE = regexp.MustCompile(`const\s+URI_(\w+)\s*=\s*["']([\w/-]+)["']`)

func matchRequests(fileName string, fileBytes []byte) {
	reMatched := requestRE.FindAllSubmatch(fileBytes, -1)
	for i := range reMatched {
		// 没什么意义，因为匹配函数如果找到了匹配项，这里的长度就必然不小于3,但是因为后续要直接使用下标访问，所以还是判断一下，下同
		if len(reMatched[i]) < 3 {
			continue
		}

		requestName := string(reMatched[i][1])

		token.GeneratorIns.RequestFrom[fileName] = append(token.GeneratorIns.RequestFrom[fileName], requestName)
		token.GeneratorIns.Requests[requestName] = string(reMatched[i][2])
	}
}
