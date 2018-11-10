package persist

import (
	"log"

	"gopkg.in/olivere/elastic.v5"
	"coding-180/crawler/engine"
	"coding-180/crawler/persist"
)

type ItemSaverService struct {
	Client *elastic.Client
	Index  string
}

func (s *ItemSaverService) Save(
	item engine.Item, result *string) error {
	err := persist.Save(s.Client, s.Index, item)
	log.Printf("Item %v saved.", item)

	//直接给result赋值就能够把结果返回给client了，贼鸡儿简单...
	if err == nil {
		*result = "ok"
	} else {
		log.Printf("Error saving item %v: %v",
			item, err)
	}
	return err
}
