# vue生态

## pinia

> 学习时间：2026.1
> 文档：https://pinia.vuejs.org/zh/core-concepts/

pinia基于vuex一个升级方向的探索，它们都是用来管理全局状态的。

基础语法：（使用setup store格式）

``` ts
import { ref } from "vue"
import { defineStore } from "pinia"

export let useFlagStore = defineStore("flag", () => {
	let wildScreenFlag = ref<boolean>(true)

	function onScreenWidthChanged(width: number): void {
		wildScreenFlag.value = width > 1280
	}

	return { wildScreenFlag, onScreenWidthChanged }
})
```

解构会破坏响应性，推荐不定义新的变量、直接使用：

```vue
<script setup lang="ts">
const store = useCounterStore()
// ❌ 这将不起作用，因为它破坏了响应性
// 这就和直接解构 `props` 一样
const { name, doubleCount } = store 
// name将始终是 "Eduardo" 
// doubleCount将始终是 0 
setTimeout(() => {
  store.increment()
}, 1000)
// ✅ 这样写是响应式的
// 💡 当然你也可以直接使用 `store.doubleCount`
const doubleValue = computed(() => store.doubleCount)
</script>
```

### 持久化

```code 
const loginStore = useLoginStore();

onMounted(() => {
  window.addEventListener("beforeunload", () => {
    sessionStorage.setItem("login_data", JSON.stringify(loginStore.user))
  })

  let loginData = sessionStorage.getItem("login_data")
  if (!loginData) {
    return
  }

  loginStore.user = deepCopy(JSON.parse(loginData))
  sessionStorage.removeItem("login_data")
})

// deepCopy 简单的deep copy，没有考虑嵌套对象的情况
export function deepCopy<T extends object>(obj: T): T {
	let res: T = {} as T

	for (let key in obj) {
		res[key] = obj[key]
	}

	return res
}
```

## vite配置文档

> 学习时间：2026.1
> 文档：https://cn.vitejs.dev/config/shared-options.html

路径别名：`resolve: { alias: { "@": fileURLToPath(new URL("./src", import.meta.url)) } },`

不清屏：`clearScreen: false`

server.host：如果将此设置为 `0.0.0.0` 或者 `true` 将监听所有地址，包括局域网和公网地址，或者使用下方函数获取本机内网ip

``` ts
import os from "os"
export function getLocalIP(): string {
	const networks = os.networkInterfaces()
	for (let key in networks) {
		// @ts-ignore
		for (let ins of networks[key]) {
			if (ins.family === "IPv4" && !ins.internal) {
				return ins.address
			}
		}
	}

	return "127.0.0.1"
}
```

### 配置允许本机和手机访问

使用本机内网ip：

- server.host使用本机内网ip（获取ip函数见上一节）
- 检查请求目标地址（通常在`.env.development`文件）

## ts编译选项

> 学习时间：2026.1
> ts编译选项文档(tsconfig.json)：https://www.typescriptlang.org/tsconfig/

参考以下json，其他属性默认值即可。

``` json
{
  "compilerOptions": {
    "strict": true,
    "noUnusedLocals": true,
    "noUnusedParameters": true,
    "noEmit": true, // 和下一条一起，允许在import ts文件的时候，添加'.ts'扩展名
    "allowImportingTsExtensions": true,
    "module": "esnext",
    "target": "esnext",
    "moduleResolution": "bundler",
    "paths": {
      "@/*": ["./src/*"]
    },
    "types": [
      "node"
    ],
    "composite": true,
  },
  "include": [
    "src/**/*",
    "src/**/*.vue",
    "vite.config.ts",
    "eslint.config.js"
  ]
}
```

我们注意到，ts和vite的配置都包含路径解析`@`表示`/src/`，其中vite的配置是供程序使用的，而ts的配置是给ide使用的：
如果删除vite的配置，则程序无法正确build；如果删除ts的配置，ide会报错（重启ide）。

## eslint+prettier

> 学习时间：2026.1
> eslint: 9

eslint.config.js

``` js
import { defineConfig, globalIgnores } from "eslint/config"
import vueTs from "@vue/eslint-config-typescript"
import vuePrettier from "@vue/eslint-config-prettier"
import vueParser from "vue-eslint-parser"
import tsParser from "@typescript-eslint/parser"

export default defineConfig([
	globalIgnores(["node_modules/", "dist/"]),
	vueTs(),
	vuePrettier,
	{
		files: ["**/*.{js,ts,vue}"],
		languageOptions: { parser: vueParser, parserOptions: { parser: tsParser, sourceType: "esnext" } },
		rules: {
			"@typescript-eslint/ban-ts-comment": 0,
			"@typescript-eslint/no-array-constructor": 0,
			"prefer-const": 0,
			"@typescript-eslint/no-unused-vars": 0,
			"@typescript-eslint/no-explicit-any": 0
		}
	}
])
```

.prettierrc

```txt
{
  "useTabs": true,
  "tabWidth": 4,
  "printWidth": 120,
  "semi": false,
  "trailingComma": "none",
  "arrowParens": "avoid"
}
```
