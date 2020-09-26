package corpus

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"runtime"
)

func New() *Documents {
	return &Documents{
		Docs: make([]*Document, 0, 1),
	}
}

type Documents struct {
	Docs []*Document `xml:"doc"`
	id   int
}

func (d Documents) ToMap() map[int]*Document {
	docsMap := make(map[int]*Document)
	for _, doc := range d.Docs {
		docsMap[doc.ID] = doc
	}
	return docsMap
}

func (d Documents) Len() int {
	return len(d.Docs)
}

func (d Documents) Iterator() func() (int, *Document, bool) {
	index := 0
	return func() (int, *Document, bool) {
		if index >= d.Len() {
			return index, nil, false
		}

		doc := d.Docs[index]
		curr := index
		index++
		return curr, doc, true
	}
}

func (d *Documents) Add(doc *Document) {
	if doc.ID == 0 {
		doc.ID = d.id
		d.id++
	}

	d.Docs = append(d.Docs, doc)
}

type Document struct {
	ID    int
	Title string `xml:"title" json:"title"`
	URL   string `xml:"url" json:"url"`
	Text  string `xml:"abstract" json:"abstract"`
}

func GetAll() (*Documents, error) {
	documents := make([]*Document, 0, 1)

	_, f, _, _ := runtime.Caller(0)

	dir := filepath.Dir(f)
	filePath := filepath.Join(dir, "./data.xml")
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	docs := &Documents{
		Docs: documents,
	}

	if err = xml.NewDecoder(file).Decode(docs); err != nil {
		return nil, err
	}

	for i := range docs.Docs {
		docs.Docs[i].ID = docs.id
		docs.id++
	}

	return docs, nil
}
