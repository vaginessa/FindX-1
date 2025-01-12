package engines

import (
	"errors"
	"fmt"

	"github.com/bnyro/findx/entities"
	"github.com/bnyro/findx/utilities"
	"github.com/bnyro/findx/web"
)

const resultsPerPage = 25

func FetchImage(query string, page int) ([]entities.Image, error) {
	var images []entities.Image
	var data map[string]interface{}
	offset := (page - 1) * resultsPerPage

	if offset+resultsPerPage >= 250 {
		return images, errors.New("count + offset must be smaller than 250")
	}

	uri := fmt.Sprintf("https://api.qwant.com/v3/search/images?q=%s&offset=%d&locale=en_gb&count=%d", query, offset, resultsPerPage)
	err := web.RequestJson(uri, &data)

	if err != nil {
		return images, err
	}

	results := data["data"].(map[string]interface{})["result"].(map[string]interface{})["items"].([]interface{})
	for _, res := range results {
		result := res.(map[string]interface{})
		image := entities.Image{}
		image.Title = result["title"].(string)
		image.Url = utilities.Redirect(result["url"].(string))
		image.Thumbnail = utilities.RewriteProxied(result["thumbnail"].(string))
		image.Media = result["media"].(string)
		images = append(images, image)
	}

	return images, nil
}
