package pipeline

import (
	"encoding/json"
	"go_spider/core/common/com_interfaces"
	"go_spider/core/common/page_items"
	"os"
)

type PipelineFile struct {
	pFile *os.File

	path string
}

func NewPipelineFile(path string) *PipelineFile {
	pFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic("File '" + path + "' in PipelineFile open failed.")
	}
	return &PipelineFile{path: path, pFile: pFile}
}

func (this *PipelineFile) Process(items *page_items.PageItems, t com_interfaces.Task) {
	/* 	this.pFile.WriteString("----------------------------------------------------------------------------------------------\n")
	   	this.pFile.WriteString("Crawled url :\t" + items.GetRequest().GetUrl() + "\n")
		   this.pFile.WriteString("Crawled result : \n")
	*/
	var dArr []interface{}
	for _, value := range items.GetAll() {

		var parsedD interface{}
		c := []byte(value)
		json.Unmarshal(c, &parsedD)
		dArr = append(dArr, parsedD)
		//this.pFile.WriteString(key + "\t:\t" + value + ",\n")
	}

	b, _ := json.Marshal(dArr)
	this.pFile.Write([]byte(b))

}
