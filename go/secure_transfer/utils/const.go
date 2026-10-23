package utils

const (
	FileHeaderFlag = "ST"
	DeriveKeyInfo  = "aes-gcm 256 encrypt key"

	PrivateKeyFileName = "priv.key"
	PublicKeyFileName  = "PUB.KEY"

	OriginFileName    = "message"
	DefaultExtension  = "txt"
	EncryptedFileName = "CIPHER"
	DecryptedFileName = "message_decrypted"
)

const (
	DeriveKeySaltLength = 16
	PublicKeyLength     = 32
	AESBaseNonceLength  = 7  // aes-gcm nonce: 12 = 7 + 1 + 4
	FileHeaderAADLength = 32 // sha256 length
	FileHeaderLength    = 1 + 2 + 1 + 4 + DeriveKeySaltLength + AESBaseNonceLength + PublicKeyLength + FileHeaderAADLength

	DeriveKeyLength = 32 // 32 Bytes, match AES-256
	GCMTagLength    = 16 // 加密后的密文比明文多出来的长度

	OnceEncMaxSize = 10 * 1024 * 1024 // 10 MB
	FrameSize      = 1024 * 1024      // 1 MB
)

type EncryptMethod = byte

const (
	EncryptMethod_Once   EncryptMethod = 1 // 一次性读取全部文件到内存中加密
	EncryptMethod_Stream               = 2 // 流式加密
)

func IsEncryptMethod(method byte) bool {
	return method == EncryptMethod_Once || method == EncryptMethod_Stream
}
