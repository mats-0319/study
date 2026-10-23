# Windows 开发环境排查笔记

> 记录在 Windows + GoLand 上运行本仓库 CLI 时遇到的"非代码问题"，换机器或复现时可直接对照。
> 本笔记中的所有结论都有日志原文支撑，命令见文末附录。

## 一、现象

在 GoLand 里 Run `secure_transfer_cli` 时：

1. 先弹出一个**新的 cmd 窗口**，标题写着 `Norton CyberCapture`，程序在里面执行一次；
2. 关掉这个窗口后，程序**又在 GoLand 的 Run 控制台里执行一次**；
3. 两次运行打印的 `Work Dir:` **不一样**；
4. 连续 Run 时，后面几次不再弹窗；但 **Run / Debug 交替**时，每次都会弹；
5. 期间 GoLand 偶尔报错：
   `CreateProcess error=4551, 应用程序控制策略已阻止此文件`（连续重试十余次）。

## 二、结论

- **不是 GoLand 的问题，也不是本程序的问题**：全仓库没有任何 `os/exec`、`cmd /c start`、`CREATE_NEW_CONSOLE`，程序自己不会拉新窗口。
- 那个窗口是 **Norton 360 for Gamers（Avast 引擎）的 CyberCapture / 自动沙盒**：它把刚由 `go build` 产出的、无签名无信誉的 exe 拖进沙盒跑一遍；GoLand 自己启的那次才是真实运行。
- 而这个 Norton **并不是手动下载安装的**，是微星主板的 "MSI Driver Utility Installer"（一键驱动工具）**静默捆绑**装上的。

## 三、证据链

### 1. GoLand 侧：直接启动 exe，不经过 cmd

`idea.log` 中的堆栈：

```
GoRunningState.startProcess
GoBuildingRunningState.startProcess
GeneralCommandLine.createProcess -> ProcessBuilder.start
```

即 GoLand 用 Java 的 `ProcessBuilder` **直接** `CreateProcess` 编译产物，工作目录就是运行配置里的 `working_directory`；三个 stdio 全部接管道，因此**正常 Run 不会创建新控制台窗口**（同条件手工实测：无 conhost、无窗口）。

### 2. Norton 侧：沙盒日志

`C:\ProgramData\Norton\Antivirus\log\autosandbox.log`：

```
2026/10/23 14:30:13  Autosandbox candidate: ...\cmd\cli\sample\secure_transfer_cli.exe
    [Source: ]
    [Opened by: ...\cmd\cli\sample\secure_transfer_cli.exe]
    [Reason: 0x00020000]
     --> Result: Sandboxing (no custody)
     --> Instrumentation: Instrumentation inside sandbox requested
```

同一文件里还有另一种结果，对应"后面几次不弹窗"：

```
     --> Result: Not sandboxing (because the file is in the exception list).
```

### 3. 为什么两次工作目录不同

- 沙盒实例的 cwd 由 Norton 沙盒给出，**不是**运行配置的 `working_directory`；
- 沙盒的文件系统是隔离的，所以程序自己的 `log.log` 里**永远只有 `cmd/cli/sample/` 一条 `Work Dir:` 记录**——这也解释了为什么磁盘上找不到第二份 `log.log`。

### 4. 为什么 Run/Debug 交替时每次都弹

Debug 构建的二进制与 Run 不同（`-N -l`、build id 都会变），每次重建都是"新的无信誉文件"，Norton 会重新判定一次。

### 5. Norton 的安装来源

`C:\Program Files (x86)\MSI\MSI_Driver_Utility_Installer\Log\Application_10_15_2026.log`：

```
10-15-2026 08:17:20.369 | APP [25.11.10580.2110] No Install ... Norton 360
10-15-2026 08:17:20.369 | APP [2.19] No Install ... CPU-Z MSI GAMING
10-15-2026 08:31:49.237 | ===== Start process : Norton 360 =====
    ExeSetupFilePath : ...\N360Gamers_25.11.10580.2110\Norton 360\N360Gamers\MSI_MB_GAMER_Norton360_online.exe
    WebPath 0 : https://liveupdate.msi.cn/cd/MSI_WU_DRIVER/MB/Release/Utility/N360Gamers_25.11.10580.2110.zip
```

Norton 安装器自身的自解压日志 `C:\ProgramData\Norton\Icarus\Logs\sfx.log`：

```
[2026-10-15 00:31:51] *** Starting SFX ***
    launched by:'C:\Program Files (x86)\MSI\MSI_Driver_Utility_Installer\MSI_Driver_Utility_Installer.exe'
    SFX started with command line '/silent /JoinEula'
```

安装后的产品标记 `C:\ProgramData\Norton\Icarus\norton-suite\icarus.ini`：

```ini
product-full-name=Norton 360 for Gamers
suite_oem_activation=1
lang-id=2052
```

**串起来**：装完系统 → 运行微星一键驱动工具 → 它从`liveupdate.msi.cn`下载 `N360Gamers_*.zip` → 用
`/silent /JoinEula`（静默安装 + 自动同意协议）装上 Norton 360 for Gamers，同时顺带装了 CPU-Z MSI GAMING。
**全程无弹窗、无勾选、无需下载**，所以完全不会察觉。

## 四、运行配置本身的坑（顺带修复）

`.run/secure_transfer_cli.run.xml` 里有几个互相矛盾、容易造成"工作目录/启动入口不一致"的字段：

| 字段 | 错误值 | 问题 |
| --- | --- | --- |
| `kind` | `DIRECTORY` | 同时却填写了 `package`、`filePath` |
| `package` | `github.com/mats0319/study/go/secure_transmission` | 模块名应为 `github.com/mats0319/secure_transfer`，且拼写错误（transmission） |
| `filePath` | `$PROJECT_DIR$/go/encrypt_transmission/main.go` | 文件不存在 |
| `output_directory` | `.../cmd/cli/sample` | 把 exe 编译进工作目录，每次 Run 覆盖工作目录里的二进制 |
| `pty_enabled` | `true` | Windows 上 Go 插件已禁用 PTY（日志 `Windows: disable PTY mode`），此开关无意义 |

推荐配置（`working_directory` 必须固定，因为程序按 cwd 读写 `priv.key`、`message.txt`、`log.log`）：

```xml
<component name="ProjectRunConfigurationManager">
  <configuration default="false" name="secure_transfer_cli" type="GoApplicationRunConfiguration" factoryName="Go Application">
    <module name="study" />
    <working_directory value="$PROJECT_DIR$/go/secure_transfer/cmd/cli/sample" />
    <kind value="PACKAGE" />
    <package value="github.com/mats0319/secure_transfer/cmd/cli" />
    <directory value="$PROJECT_DIR$" />
    <method v="2" />
  </configuration>
</component>
```

另外：运行配置应放在**项目根目录**的 `.run/`（即 `study/.run/`），放子目录（如 `go/secure_transfer/.run/`）容易与 `.idea/workspace.xml` 里的自动配置重复、互相覆盖。

## 五、处理步骤

1. 卸载 Norton；重启后用官方清理工具 `NRnR`（Norton Remove and Reinstall）扫尾；
2. 打开「Windows 安全中心 → 病毒和威胁防护」确认 **Defender 已接管**（若提示"无防病毒提供程序"，重启后恢复）；
3. 按上一节修正运行配置；
4. 如果要保留其他杀软，把项目目录、`%LOCALAPPDATA%\go-build`、`%TEMP%` 加入排除项，避免每次 build 都被重新判定；
5. 主板驱动、BIOS 一律去官网按型号手动下载，**不要再运行厂商的"一键驱动/装机工具"**。

## 六、通用清单：新机器检查是否被捆绑安装

- 看 `C:\ProgramData\` 下当天新建的厂商目录（本次是 `C:\ProgramData\Norton`）；
- 看 `C:\Program Files (x86)\` 下不认识的名字；
- 翻厂商工具的安装日志目录（本次 `C:\Program Files (x86)\MSI\MSI_Driver_Utility_Installer\Log\`）；
- 第三方组件的自解压/安装日志常留有 `launched by:` 与命令行（本次 `sfx.log` 记下了 `/silent /JoinEula`）；
- 卸载列表按 `InstallDate` 排序，看是否有一批"同一天、不认识"的软件。

## 附录：排查命令

```powershell
# 是否被 CyberCapture/自动沙盒判定过（看 Result 行）
Get-Content 'C:\ProgramData\Norton\Antivirus\log\autosandbox.log' -Tail 20

# 厂商一键工具装了什么、装件从哪下载
Get-ChildItem 'C:\Program Files (x86)\MSI\MSI_Driver_Utility_Installer\Log' |
  ForEach-Object { Select-String -Path $_.FullName -Pattern 'Start process|WebPath|No Install' }

# 第三方组件是谁调起的、用了什么参数
Select-String -Path 'C:\ProgramData\Norton\Icarus\Logs\sfx.log' -Pattern 'launched by|command line'

# 卸载后残留（服务 / 目录 / 注册表 / 驱动）
Get-Service | Where-Object { $_.Name -match 'Norton|Sym|asw' }
Get-ChildItem 'C:\Program Files\Norton','C:\ProgramData\Norton' -Recurse -ErrorAction SilentlyContinue
Get-ChildItem 'C:\Windows\System32\drivers' | Where-Object Name -match '^asw'

# Defender 是否接管（需管理员）
Get-MpComputerStatus | Select-Object AntivirusEnabled, RealTimeProtectionEnabled
```
