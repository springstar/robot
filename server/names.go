package server

import (
	_ "bytes"
	_ "fmt"
	"github.com/springstar/robot/core"
	"log"
	"io/ioutil"
	"strings"
	_ "unicode/utf8"
	"unicode"
)

type NameManager struct {
	firstNames []string
	female []string
	male []string
}

func newNameManager() *NameManager {
	return &NameManager{}
}

func (m *NameManager) loadNameFiles() {
	m.firstNames = m.loadNames("config/names/firstName.txt")
	m.female = m.loadNames("config/names/firstNameFemale.txt")
	m.male = m.loadNames("config/names/firstNameMale.txt")
}

func (m *NameManager) loadRunes(f string) {
	bin, err := ioutil.ReadFile(f)
	if (err != nil) {
		log.Fatal(err)
	}

	text := string(bin[:])
	core.ScanRunes(text)
}

func (m *NameManager) loadNames(f string) []string {
	bin, err := ioutil.ReadFile(f)
	if (err != nil) {
		log.Fatal(err)
		return nil
	}

	content := string(bin[:])
	strs := strings.Split(content, "\r\n")
	return strs
}

func remove(data []byte) []byte{
	cleanData := strings.Map(func(r rune) rune {
		if unicode.IsControl(r) && !unicode.IsPrint(r) && !unicode.IsSpace(r) {
			return -1
		}
		return r
	}, string(data))

	return []byte(cleanData)
}

func (m *NameManager) randomGenName(sex int) (name string) {
	idx := core.GenRandomInt(len(m.firstNames))	
	firstName := m.firstNames[idx]
	var secondName string

	if sex == ROLE_SEX_MALE {
		idx = core.GenRandomInt(len(m.male))
		secondName = m.male[idx]
	} else if sex == ROLE_SEX_FEMALE {
		idx = core.GenRandomInt(len(m.female))
		secondName = m.female[idx]
	}

	name = strings.Join([]string{firstName, secondName}, "")

	return
}

