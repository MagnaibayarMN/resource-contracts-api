package document

import "os/exec"

func CreateDoc(dir string, id string) {
	zip := exec.Command(
		"zip",
		"-r",
		id+".docx",
		"[Content_Types].xml",
		"customXml",
		"docProps",
		"_rels",
		"word",
	)
	zip.Dir = dir + "/" + id
	_, err := zip.Output()
	Check(err)
}
