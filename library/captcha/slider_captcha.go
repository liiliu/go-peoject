package captcha

import (
	"fmt"
	"image"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/wenlng/go-captcha-assets/resources/images_v2"
	"github.com/wenlng/go-captcha-assets/resources/tiles"
	"github.com/wenlng/go-captcha/v2/slide"
)

var slideCapt slide.Captcha

func Initial() {
	builder := slide.NewBuilder(
	// slide.WithGenGraphNumber(2),
	// slide.WithEnableGraphVerticalRandom(true),
	)

	// background images
	imgs, err := images.GetImages()
	if err != nil {
		log.Fatalln(err)
	}

	graphs, err := tiles.GetTiles()
	if err != nil {
		log.Fatalln(err)
	}

	var newGraphs = make([]*slide.GraphImage, 0, len(graphs))
	for i := 0; i < len(graphs); i++ {
		graph := graphs[i]
		newGraphs = append(newGraphs, &slide.GraphImage{
			OverlayImage: graph.OverlayImage,
			MaskImage:    graph.MaskImage,
			ShadowImage:  graph.ShadowImage,
		})
	}

	// set resources
	builder.SetResources(
		slide.WithGraphImages(newGraphs),
		slide.WithBackgrounds(imgs),
	)

	slideCapt = builder.Make()
}

func GenerateCaptcha() (data map[string]interface{}, err error) {
	captData, err := slideCapt.Generate()
	if err != nil {
		return
	}

	dotData := captData.GetData()
	if dotData == nil {
		return
	}

	//dots, _ := json.Marshal(dotData)
	//fmt.Println(">>>>> ", string(dots))

	var mBase64, tBase64 string
	mBase64, err = captData.GetMasterImage().ToBase64()
	if err != nil {
		fmt.Println(err)
	}
	tBase64, err = captData.GetTileImage().ToBase64()
	if err != nil {
		fmt.Println(err)
	}

	data = make(map[string]interface{})
	data["image"] = mBase64
	data["thumb"] = tBase64
	data["thumbX"] = dotData.TileX
	data["thumbY"] = dotData.TileY
	data["thumbWidth"] = dotData.Width
	data["thumbHeight"] = dotData.Height
	data["x"] = dotData.X
	data["y"] = dotData.Y

	return
}

func CheckPoint(srcX, srcY, x, y int64) bool {
	return slide.CheckPoint(srcX, srcY, x, y, 5)
}

// GetImages 获取指定目录下图片资源
func GetImages(dir string) ([]image.Image, error) {
	var imgs []image.Image

	// 支持的图片后缀
	supportedExt := map[string]bool{
		".jpg": true,
		//".jpeg": true,
		//".png":  true,
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && supportedExt[strings.ToLower(filepath.Ext(path))] {
			imgFile, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("打开图片失败: %s, %v", path, err)
			}
			defer imgFile.Close()

			img, _, err := image.Decode(imgFile)
			if err != nil {
				return fmt.Errorf("解码图片失败: %s, %v", path, err)
			}
			imgs = append(imgs, img)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return imgs, nil
}
