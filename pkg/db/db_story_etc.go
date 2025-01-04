package db

import (
	"encoding/json"
	"fmt"
	"strings"
)

func GetAssociateMediaKeysForEditorjsSrc(rawArticle []byte) ([]string, error) {

	var retlist []string

	var editorjsSrc map[string]interface{}

	err := json.Unmarshal(rawArticle, &editorjsSrc)

	if err != nil {

		return nil, fmt.Errorf("failed to unmarshal: %s", err.Error())

	}

	blocks, okay := editorjsSrc["blocks"]

	if !okay {

		return nil, fmt.Errorf("invalid format: %s", "no blocks")
	}

	blocksList := blocks.([]interface{})

	blocksLen := len(blocksList)

	for i := 0; i < blocksLen; i++ {

		blockObj := blocksList[i].(map[string]interface{})

		objType, okay := blockObj["type"]

		if !okay {
			continue
		}

		if objType != "image" {
			continue
		}

		objData, okay := blockObj["data"]

		if !okay {
			continue
		}

		objFields := objData.(map[string]interface{})

		fileField, okay := objFields["file"]

		if !okay {
			continue
		}

		targetProps := fileField.(map[string]interface{})

		urlTarget, okay := targetProps["url"]

		if !okay {
			continue
		}

		target := urlTarget.(string)

		pathList := strings.Split(target, "/")

		key := pathList[len(pathList)-1]

		retlist = append(retlist, key)
	}

	return retlist, nil
}
