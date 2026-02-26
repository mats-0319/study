package generate_avatar

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image/png"
	"os"
)

// GenerateAvatar generate avatar according to 'text'
//
//	@params text: use for calculate hash, distinguish 'A' and 'a'
//	@params size: image size, valid range is [1, 10], image will be '12*size x 12*size' px
func GenerateAvatar(text string, size int) error {
	if !(0 < size && size <= 10) {
		return errors.New("invalid clarity, required: (0, 10]")
	}

	imageImplIns, err := NewImageImpl(text, size)
	if err != nil {
		return err
	}

	// encode to image
	writer := bytes.NewBufferString("")
	{
		b64 := base64.NewEncoder(base64.StdEncoding, writer)
		err := (&png.Encoder{CompressionLevel: png.BestCompression}).Encode(b64, imageImplIns)
		if err != nil {
			return err
		}
		_ = b64.Close()
	}

	// serialize and write file
	imageBytes, err := base64.StdEncoding.DecodeString(writer.String())
	if err != nil {
		return err
	}

	// file name
	filename := text
	if len(filename) > 8 {
		filename = filename[:8]
	}
	filename = fmt.Sprintf("./img/%s_%d.png", filename, size)

	return os.WriteFile(filename, imageBytes, 0777)
}
