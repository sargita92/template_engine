package docx

import (
	"archive/zip"
	"os"
	"path/filepath"

	"github.com/sargita92/template_engine/ast"
	wordrenderer "github.com/sargita92/template_engine/render/word"
)

// Write cria um DOCX válido a partir de um AST.Document
func Write(path string, doc ast.Document) error {
	renderer := wordrenderer.New()
	documentXML := renderer.RenderDocument(doc)

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	zipw := zip.NewWriter(f)
	defer zipw.Close()

	if err := writeContentTypes(zipw); err != nil {
		return err
	}

	if err := writeRootRels(zipw); err != nil {
		return err
	}

	if err := writeDocument(zipw, documentXML); err != nil {
		return err
	}

	if err := writeDocumentRels(zipw); err != nil {
		return err
	}

	return nil
}

func writeContentTypes(zipw *zip.Writer) error {
	const contentTypes = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
	<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
	<Default Extension="xml" ContentType="application/xml"/>
	<Override PartName="/word/document.xml"
		ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
</Types>`

	w, err := zipw.Create("[Content_Types].xml")
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(contentTypes))
	return err
}

func writeRootRels(zipw *zip.Writer) error {
	const rels = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
	<Relationship
		Id="rId1"
		Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument"
		Target="word/document.xml"/>
</Relationships>`

	w, err := zipw.Create(filepath.ToSlash("_rels/.rels"))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(rels))
	return err
}

func writeDocument(zipw *zip.Writer, xml string) error {
	w, err := zipw.Create(filepath.ToSlash("word/document.xml"))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(xml))
	return err
}

func writeDocumentRels(zipw *zip.Writer) error {
	const rels = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
</Relationships>`

	w, err := zipw.Create(filepath.ToSlash("word/_rels/document.xml.rels"))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(rels))
	return err
}
