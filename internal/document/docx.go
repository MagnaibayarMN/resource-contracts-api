// Package document provides DOCX generation utilities.
package document

import (
	"bufio"
	"bytes"
	"encoding/xml"
)

// CreateTitle creates a Word document heading with Heading1 style.
// Used for contract titles in exported documents.
//
// Parameters:
//   - title: The heading text
//
// Returns:
//   - XML string for Word document heading
func CreateTitle(title string) string {
	return `<w:p><w:pPr><w:pStyle w:val="Heading1" /></w:pPr><w:proofErr w:type="gramStart" /><w:r><w:t>` + title + `</w:t></w:r><w:proofErr w:type="gramEnd" /></w:p>`
}

// CreateParagraph creates a standard paragraph in Word document format.
//
// Parameters:
//   - text: The paragraph text (should be XML-escaped)
//
// Returns:
//   - XML string for Word document paragraph
func CreateParagraph(text string) string {
	return `<w:p><w:r><w:t>` + text + `</w:t></w:r></w:p>`
}

// createHeader generates the XML header for a Word document.
// Includes all necessary namespaces and document structure.
//
// Returns:
//   - XML string for document header
func createHeader() string {
	// return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><w:document xmlns:wpc="http://schemas.microsoft.com/office/word/2010/wordprocessingCanvas" xmlns:mc="http://schemas.openxmlformats.org/markup-compatibility/2006" xmlns:o="urn:schemas-microsoft-com:office:office" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:m="http://schemas.openxmlformats.org/officeDocument/2006/math" xmlns:v="urn:schemas-microsoft-com:vml" xmlns:wp14="http://schemas.microsoft.com/office/word/2010/wordprocessingDrawing" xmlns:wp="http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing" xmlns:w10="urn:schemas-microsoft-com:office:word" xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" xmlns:w14="http://schemas.microsoft.com/office/word/2010/wordml" xmlns:wpg="http://schemas.microsoft.com/office/word/2010/wordprocessingGroup" xmlns:wpi="http://schemas.microsoft.com/office/word/2010/wordprocessingInk" xmlns:wne="http://schemas.microsoft.com/office/word/2006/wordml" xmlns:wps="http://schemas.microsoft.com/office/word/2010/wordprocessingShape" mc:Ignorable="w14 wp14"><w:body><w:sectPr w:rsidR="00F61DEF"><w:headerReference w:type="even" r:id="rId6"/><w:headerReference w:type="default" r:id="rId17"/></w:sectPr>`
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?><w:document xmlns:wpc="http://schemas.microsoft.com/office/word/2010/wordprocessingCanvas" xmlns:cx="http://schemas.microsoft.com/office/drawing/2014/chartex" xmlns:cx1="http://schemas.microsoft.com/office/drawing/2015/9/8/chartex" xmlns:cx2="http://schemas.microsoft.com/office/drawing/2015/10/21/chartex" xmlns:cx3="http://schemas.microsoft.com/office/drawing/2016/5/9/chartex" xmlns:cx4="http://schemas.microsoft.com/office/drawing/2016/5/10/chartex" xmlns:cx5="http://schemas.microsoft.com/office/drawing/2016/5/11/chartex" xmlns:cx6="http://schemas.microsoft.com/office/drawing/2016/5/12/chartex" xmlns:cx7="http://schemas.microsoft.com/office/drawing/2016/5/13/chartex" xmlns:cx8="http://schemas.microsoft.com/office/drawing/2016/5/14/chartex" xmlns:mc="http://schemas.openxmlformats.org/markup-compatibility/2006" xmlns:aink="http://schemas.microsoft.com/office/drawing/2016/ink" xmlns:am3d="http://schemas.microsoft.com/office/drawing/2017/model3d" xmlns:o="urn:schemas-microsoft-com:office:office" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" xmlns:m="http://schemas.openxmlformats.org/officeDocument/2006/math" xmlns:v="urn:schemas-microsoft-com:vml" xmlns:wp14="http://schemas.microsoft.com/office/word/2010/wordprocessingDrawing" xmlns:wp="http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing" xmlns:w10="urn:schemas-microsoft-com:office:word" xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" xmlns:w14="http://schemas.microsoft.com/office/word/2010/wordml" xmlns:w15="http://schemas.microsoft.com/office/word/2012/wordml" xmlns:w16cex="http://schemas.microsoft.com/office/word/2018/wordml/cex" xmlns:w16cid="http://schemas.microsoft.com/office/word/2016/wordml/cid" xmlns:w16="http://schemas.microsoft.com/office/word/2018/wordml" xmlns:w16sdtdh="http://schemas.microsoft.com/office/word/2020/wordml/sdtdatahash" xmlns:w16se="http://schemas.microsoft.com/office/word/2015/wordml/symex" xmlns:wpg="http://schemas.microsoft.com/office/word/2010/wordprocessingGroup" xmlns:wpi="http://schemas.microsoft.com/office/word/2010/wordprocessingInk" xmlns:wne="http://schemas.microsoft.com/office/word/2006/wordml" xmlns:wps="http://schemas.microsoft.com/office/word/2010/wordprocessingShape" mc:Ignorable="w14 w15 w16se w16cid w16 w16cex w16sdtdh wp14"><w:body><w:sdt><w:sdtPr><w:id w:val="1747684861"/><w:docPartObj><w:docPartGallery w:val="Table of Contents"/><w:docPartUnique/></w:docPartObj></w:sdtPr><w:sdtEndPr><w:rPr><w:rFonts w:asciiTheme="minorHAnsi" w:eastAsiaTheme="minorHAnsi" w:hAnsiTheme="minorHAnsi" w:cstheme="minorBidi"/><w:b/><w:bCs/><w:noProof/><w:color w:val="auto"/><w:sz w:val="22"/><w:szCs w:val="22"/></w:rPr></w:sdtEndPr><w:sdtContent><w:p w14:paraId="59CA73FE" w14:textId="11208235" w:rsidR="00B1484F" w:rsidRDefault="00B1484F"><w:pPr><w:pStyle w:val="TOCHeading"/></w:pPr><w:r><w:t>Contents</w:t></w:r></w:p><w:p w14:paraId="46D0C7E4" w14:textId="65D40068" w:rsidR="00B1484F" w:rsidRDefault="00B1484F"><w:fldSimple w:instr=" TOC \o &quot;1-3&quot; \h \z \u "><w:r><w:rPr><w:b/><w:bCs/><w:noProof/></w:rPr><w:t>No table of contents entries found.</w:t></w:r></w:fldSimple></w:p></w:sdtContent></w:sdt><w:p w14:paraId="44DF0BFB" w14:textId="77777777" w:rsidR="00F61DEF" w:rsidRDefault="000448CE"/><w:sectPr w:rsidR="00F61DEF"><w:headerReference w:type="even" r:id="rId7"/><w:headerReference w:type="default" r:id="rId8"/><w:footerReference w:type="even" r:id="rId9"/><w:footerReference w:type="default" r:id="rId10"/><w:headerReference w:type="first" r:id="rId11"/><w:footerReference w:type="first" r:id="rId12"/><w:pgSz w:w="12240" w:h="15840"/><w:pgMar w:top="1440" w:right="1440" w:bottom="1440" w:left="1440" w:header="720" w:footer="720" w:gutter="0"/><w:cols w:space="720"/><w:docGrid w:linePitch="360"/></w:sectPr>`
}

// CreateFooter generates the closing tags for a Word document.
//
// Returns:
//   - XML string for document footer
func CreateFooter() string {
	return `</w:body></w:document>`
}

// CreateTOC : Table of contents
func createTOC() string {
	return `<w:sdt>
		<w:sdtPr>
			<w:id w:val="-875628688" />
			<w:docPartObj>
				<w:docPartGallery w:val="Table of Contents" />
				<w:docPartUnique />
			</w:docPartObj>
		</w:sdtPr>
		<w:sdtEndPr>
			<w:rPr>
				<w:rFonts w:asciiTheme="minorHAnsi" w:eastAsiaTheme="minorHAnsi" w:hAnsiTheme="minorHAnsi" w:cstheme="minorBidi" />
				<w:noProof />
				<w:color w:val="auto" />
				<w:sz w:val="22" />
				<w:szCs w:val="22" />
				<w:lang w:eastAsia="en-US" />
			</w:rPr>
		</w:sdtEndPr>
		<w:sdtContent>
			<w:p w:rsidR="009F46C1" w:rsidRDefault="009F46C1">
				<w:pPr>
					<w:pStyle w:val="TOCHeading" />
				</w:pPr>
				<w:r>
					<w:t>Contents</w:t>
				</w:r>
			</w:p>
			<w:p w:rsidR="009F46C1" w:rsidRDefault="009F46C1">
				<w:fldSimple w:instr=" TOC \o &quot;1-3&quot; \h \z \u ">
					<w:r>
						<w:rPr>
							<w:b />
							<w:bCs />
							<w:noProof />
						</w:rPr>
						<w:t>No table of contents entries found.</w:t>
					</w:r>
				</w:fldSimple>
			</w:p>
		</w:sdtContent>
	</w:sdt>`
}

// XmlEscape escapes special XML characters in text.
// Converts characters like <, >, &, etc. to their XML entities.
//
// Parameters:
//   - value: Raw text string
//
// Returns:
//   - XML-safe escaped string
func XmlEscape(value string) string {
	escaped := &bytes.Buffer{}
	if err := xml.EscapeText(escaped, []byte(value)); err != nil {
		panic(err)
	}
	return escaped.String()
}

// InitializeDocument writes the document header to a buffer.
// Should be called before adding any content to the document.
//
// Parameters:
//   - buffer: Output buffer to write to
func InitializeDocument(buffer *bufio.Writer) {
	buffer.WriteString(createHeader())
	// buffer.WriteString(createTOC())
}
