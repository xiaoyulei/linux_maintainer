package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"
)

func isDigit(ch byte) bool {
	return strings.IndexByte("0123456789", ch) != -1
}

func isLetter(ch byte) bool {
	return strings.IndexByte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", ch) != -1
}

type stringList []string

func (s stringList) String() string {
	ret := ""
	for _, i := range s {
		if len(ret) != 0 {
			ret = ret + "; "
		}
		ret = ret + i
	}
	return ret
}

type maintainer struct {
	Mails         stringList // M: *Mail* patches to: FullName <address@domain>
	Reviewers     stringList // R: Designated *Reviewer*: FullName <address@domain>. These reviewers should be CCed on patches.
	RelevantMails stringList // L: *Mailing list* that is relevant to this area
	Status        string     // S: *Status*
	Webs          stringList // W: *Web-page* with status/info
	PatchWork     string     // Q: *Patchwork* web based patch tracking system site
	BugsURLs      stringList // B: URI for where to file *bugs*. A web-page with detailed bug filing info, a direct bug tracker link, or a mailto: URI.
	ChatURL       string     // C: URI for *chat* protocol, server and channel where developers usually hang out, for example irc://server/channel.
	Profile       string     // P: Subsystem Profile document for more details submitting patches to the given subsystem. This is either an in-tree file, or a URI. See Documentation/maintainer/maintainer-entry-profile.rst for details.
	SCMTree       stringList // T: *SCM* tree type and location. Type is one of: git, hg, quilt, stgit, topgit
	Files         stringList // F: *Files* and directories wildcard patterns.
	Excluded      stringList // X: *Excluded* files and directories that are NOT maintained, same rules as F:. Files exclusions are tested before file matches.
	FileRegex     stringList // N: Files and directories *Regex* patterns.
	ContentRegex  stringList // K: *Content regex* (perl extended) pattern match in a patch or file.
}

func (m maintainer) CSV() string {
	s := ""
	s = s + m.Mails.String() + ","
	s = s + m.Reviewers.String() + ","
	s = s + m.RelevantMails.String() + ","
	s = s + m.Status + ","
	s = s + m.Webs.String() + ","
	s = s + m.PatchWork + ","
	s = s + m.BugsURLs.String() + ","
	s = s + m.ChatURL + ","
	s = s + m.Profile + ","
	s = s + m.SCMTree.String() + ","
	s = s + m.Files.String() + ","
	s = s + m.Excluded.String() + ","
	s = s + m.FileRegex.String() + ","
	s = s + m.ContentRegex.String()
	return s
}

func (m maintainer) mailList4Excel() string {
	s := ""
	for _, v := range m.Mails {
		s = s + "\t" + v + "\n"
	}
	return s
}

type mudules map[string]*maintainer

func (m mudules) CSV() string {
	s := "Module," + "Mails," + "Reviewers," + "RelevantMails," + "Status," + "Webs," +
		"PatchWork," + "BugsURLs," + "ChatURL," + "Profile," + "SCMTree," + "Files," +
		"Excluded," + "FileRegex," + "ContentRegex" + "\n"
	for k, v := range m {
		s = s + k + "," + v.CSV() + "\n"
	}
	return s
}

func (m mudules) mailList4Excel() string {
	s := ""
	for k, v := range m {
		s = s + k + "\n" + v.mailList4Excel() + "\n"
	}
	return s
}

func (m mudules) mailList4Excel2() string {
	s := ""
	for k, v := range m {
		for _, y := range v.Mails {
			s = s + k + "\t" + y + "\n"
		}
	}
	return s
}

// name->modules
func (m mudules) mailList4Excel3() string {
	nm := map[string]string{}
	for k, v := range m {
		for _, y := range v.Mails {
			mu, ok := nm[y]
			if ok {
				nm[y] = mu + "\n" + k
			} else {
				nm[y] = k
			}
		}
	}

	s := ""
	for k, v := range nm {
		if strings.Contains(v, "\"") {
			v = strings.ReplaceAll(v, "\"", "")
		}
		s = s + k + "\t\"" + v + "\"\n"
	}

	return s
}

func parseLine(s string, m *maintainer) {
	h := s[0:2]
	c := strings.TrimSpace(s[2:])
	if strings.Contains(c, "\t") {
		c = strings.ReplaceAll(c, "\t", " ")
	}
	/*
		if strings.Contains(c, "\"") {
			c = strings.ReplaceAll(c, "\"", "\"\"")
		}
		if strings.Contains(c, ",") && !strings.Contains(c, "\"") {
			c = "\"" + c + "\""
		}

			if strings.Contains(c, "VMware, Inc.") {
				c = strings.ReplaceAll(c, "VMware, Inc.", "VMware Inc.")
			}
			if strings.Contains(c, "subscribers-only, for") {
				c = strings.ReplaceAll(c, "subscribers-only, for", "subscribers-only and for")
			}
			if strings.Contains(c, "Lad, Prabhakar") {
				c = strings.ReplaceAll(c, "Lad, Prabhakar", "Lad Prabhakar")
			}
			if strings.Contains(c, "Lee, Chun-Yi") {
				c = strings.ReplaceAll(c, "Lee, Chun-Yi", "Lee Chun-Yi")
			}
	*/
	switch h {
	case "M:":
		if m.Mails == nil {
			m.Mails = []string{}
		}
		m.Mails = append(m.Mails, c)
	case "R:":
		if m.Reviewers == nil {
			m.Reviewers = []string{}
		}
		m.Reviewers = append(m.Reviewers, c)
	case "L:":
		if m.RelevantMails == nil {
			m.RelevantMails = []string{}
		}
		m.RelevantMails = append(m.RelevantMails, c)
	case "S:":
		m.Status = c
	case "W:":
		if m.Webs == nil {
			m.Webs = []string{}
		}
		m.Webs = append(m.Webs, c)
	case "Q:":
		m.PatchWork = c
	case "B:":
		if m.BugsURLs == nil {
			m.BugsURLs = []string{}
		}
		m.BugsURLs = append(m.BugsURLs, c)
	case "C:":
		m.ChatURL = c
	case "P:":
		m.Profile = c
	case "T:":
		if m.SCMTree == nil {
			m.SCMTree = []string{}
		}
		m.SCMTree = append(m.SCMTree, c)
	case "F:":
		if m.Files == nil {
			m.Files = []string{}
		}
		m.Files = append(m.Files, c)
	case "X:":
		if m.Excluded == nil {
			m.Excluded = []string{}
		}
		m.Excluded = append(m.Excluded, c)
	case "N:":
		if m.FileRegex == nil {
			m.FileRegex = []string{}
		}
		m.FileRegex = append(m.FileRegex, c)
	case "K:":
		if m.ContentRegex == nil {
			m.ContentRegex = []string{}
		}
		m.ContentRegex = append(m.ContentRegex, c)
	default:
		log.Fatalf("Unknow flag: %s", s)
	}
}

func write(filename string, s string) {
	f, err := os.Create(filename) //创建文件
	if err != nil {
		log.Fatalf("create file fail. %s", err)
	}

	w := bufio.NewWriter(f) //创建新的 Writer 对象
	_, err = w.WriteString(s)
	if err != nil {
		log.Fatalf("write file fail. %s", err)
	}
	w.Flush()
	f.Close()
}

func main() {
	file, err := os.Open("D:\\code\\MAINTAINERS")
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}

	fileScanner := bufio.NewScanner(file)

	mlist := make(mudules)
	var mtn *maintainer

	flagStart := false

	// read line by line
	for fileScanner.Scan() {
		text := fileScanner.Text()
		if flagStart == false {
			// 把"Maintainers List"前的内容过滤掉
			if strings.Index(text, "Maintainers List") == 0 {
				flagStart = true
			}
			continue
		}

		// 跳过空白行
		if len(text) == 0 || len(strings.TrimSpace(text)) == 0 {
			continue
		}

		// 如果开头的不是字母或数字，也是无用的内容
		ch := text[0]
		if !isLetter(ch) && !isDigit(ch) {
			continue
		}

		// 第二个支付不是冒号的行是模块名
		if text[1] != ':' {
			mtn = new(maintainer)
			/*
				if strings.Contains(text, ",") {
					text = "\"" + text + "\""
				}
			*/
			mlist[text] = mtn
		} else {
			text = strings.TrimSpace(text)
			parseLine(text, mtn)
		}
	}

	// handle first encountered error while reading
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}

	file.Close()

	encode, err := json.MarshalIndent(mlist, "", "\t")
	if err != nil {
		log.Fatalf("Error Marshal json: %s", err)
	}

	write("D:\\code\\mj.json", string(encode))
	//write("D:\\code\\m.csv", mlist.CSV())
	//write("D:\\code\\m.txt", mlist.mailList4Excel2())
	write("D:\\code\\m.txt", mlist.mailList4Excel3())
}
