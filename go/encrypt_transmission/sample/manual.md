# 使用手册

本工具可以让两个人在不安全的通信信道上安全的传递信息

## 使用方法

> 假设是甲要传递消息给乙（甲 --(msg)--> 乙）
> 假设双方均已持有本工具

1. 乙使用本工具生成一对公私钥对，将公钥发送给甲
    - 启动本工具，输入`g`，将公钥文件(`PUB.KEY`)发送给甲
2. 甲使用本工具将要发送的明文内容加密，将密文发送给乙
    - 启动本工具，输入`i`生成明文文件，将文件中的内容替换成你想要传输的内容（明文）
    - （将乙的公钥文件放到本工具相同目录下）
    - 输入`e`加密，将生成的密文文件(`CIPHER.TXT`)发送给乙
3. 乙使用本工具将密文解密，得到明文消息
    - （将甲的密文文件放到本工具相同目录下）
    - 输入`d`解密，查看解密得到的明文(`message_decrypted.txt`)

```cmd
2025/12/01 15:23:01 > Enter Your Command. ('h' for help)
h
2025/12/01 15:23:03 
> Options:
  - h: this help
  - g: generate public & private key into files ('./priv.key' & './PUB.KEY')
  - i: initialize message file './message.txt'
  - e: encrypt plain text from './message.xxx' and write cipher to './CIPHER.XXX'
  - d: decrypt cipher from './CIPHER.XXX' and write plain text to './message_decrypted.xxx'

note: encrypt/decrypt support automatic recognize 'file extension', in fact, 
when encrypt, we find first file which name matched 'message.[xxx]' and encrypt it into 'CIPHER.[XXX]';
when decrypt, we find first file which name matched 'CIPHER.[XXX]' and decrypt it into 'message_decrypted.[xxx]'
```

## 说明

1. 因为本工具对各文件的命名有要求，所以添加有相关指令支持（`i`命令），正常使用过程中不需要特别关注文件名
2. 使用过程中需要传递两次文件，我们将需要传递的文件（公钥和密文）使用大写字母命名（例如公钥文件：`PUB.KEY`）
3. 乙可以使用其他来源的私钥吗？
   可以，只要使用go官方库`crypto/x509`中的`marshal`函数将私钥转成`[]byte`并写入文件，再将文件名修改为我们的文件名（见帮助）即可。
   至于为什么本工具不支持，主要原因是本工具的设计目的是为了安全传输，如果由本工具操作你的私钥感觉上不太好，
   而且使用本工具生成的一次性密钥就挺好的。
4. 本工具支持文件扩展名识别，会将扩展名随文件一起传递，也就是说，如果你想要加密一张图片，那么只要把图片的文件名改成
   `message.jpg`（或其他图片扩展名，取决于你的文件原始扩展名），在解密之后，乙将直接获得一个图片文件，而不需要甲额外告诉他文件类型
