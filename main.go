package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

type comment struct {
	small   bool
	content []string
}

type parser struct {
	*bytes.Buffer
	adjacent bool
	comments []comment
	buf      bytes.Buffer
}

func (p *parser) Next() rune {
	r, _, err := p.ReadRune()
	if err != nil {
		return -1
	}
	return r
}

func (p *parser) Is(c ...rune) bool {
	r, _, err := p.ReadRune()
	if err != nil {
		r = -1
	}
	for _, c := range c {
		if r == c {
			return true
		}
	}
	p.UnreadRune()
	return false
}

func (p *parser) parseBlockString(isComment bool) bool {
	i := 0
	for r := p.Next(); r >= 0; r = p.Next() {
		if r != '=' {
			p.UnreadRune()
			break
		}
		i++
	}
	if !p.Is('[') {
		return false
	}
	p.Is('\n') // Skip optional leading newline.
	p.buf.Reset()
loop:
	for r := p.Next(); r >= 0; r = p.Next() {
		switch r {
		case ']':
			j := 0
		parseEnd:
			for r := p.Next(); r >= 0; r = p.Next() {
				switch r {
				case '=':
					j++
				case ']':
					if j != i {
						p.buf.WriteString(strings.Repeat("=", j))
						p.WriteRune(']')
						break parseEnd
					}
					break loop
				default:
					p.buf.WriteString(strings.Repeat("=", j))
					p.WriteRune(r)
				}
			}
		default:
			p.buf.WriteRune(r)
		}
	}
	if isComment {
		p.comments = append(p.comments, comment{
			small:   false,
			content: []string{p.buf.String()},
		})
	}
	return true
}

func (p *parser) parseComment() {
	if p.Is('[') && p.parseBlockString(true) {
		return
	}
	line, _ := p.ReadString('\n')
	line = strings.TrimPrefix(line, " ")
	line = strings.TrimSuffix(line, "\n")
	if i := len(p.comments) - 1; i >= 0 && p.comments[i].small && p.adjacent {
		p.comments[i].content = append(p.comments[i].content, line)
		return
	}
	p.comments = append(p.comments, comment{
		small:   true,
		content: []string{line},
	})
}

func (p *parser) parseString(q rune) {
	for r := p.Next(); r >= 0; r = p.Next() {
		switch r {
		case '\\':
			p.Is('\\', '"', '\'', '\n')
		case q, '\n':
			return
		}
	}
}

func (p *parser) Parse() []string {
	for r := p.Next(); r >= 0; r = p.Next() {
		switch r {
		case '-':
			if p.Is('-') {
				p.parseComment()
				p.adjacent = true
				continue
			}
		case '[':
			p.parseBlockString(false)
		case '"', '\'':
			p.parseString(r)
		}
		p.adjacent = false
	}
	s := make([]string, len(p.comments))
	for i, group := range p.comments {
		s[i] = strings.Join(group.content, "\n")
	}
	return s
}

func varPrefix(s string) string {
	if len(s) == 0 || !('A' <= s[0] && s[0] <= 'Z' || 'a' <= s[0] && s[0] <= 'z' || s[0] == '_') {
		return ""
	}
	for i, c := range s[1:] {
		if !('0' <= c && c <= '9' || 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || c == '_') {
			return s[:i+1]
		}
	}
	return s
}

type Directives map[string]string

func ParseGroup(group string) Directives {
	directives := Directives{}
	var typ string
section:
	if strings.HasPrefix(group, "@") {
		group = group[1:]
		typ = varPrefix(group)
		group = group[len(typ):]
		group = strings.TrimPrefix(group, ":")
		group = strings.TrimSpace(group)
	} else {
		typ = "doc"
		group = strings.TrimSpace(group)
	}
	if i := strings.Index(group, "\n@"); i >= 0 {
		if s, ok := directives[typ]; ok {
			directives[typ] = s + "\n" + group[:i]
		} else {
			directives[typ] = group[:i]
		}
		group = group[i+1:]
		goto section
	}
	if s, ok := directives[typ]; ok {
		directives[typ] = s + "\n" + group
	} else {
		directives[typ] = group
	}
	return directives
}

type Section struct {
	Path       string
	Order      int
	Heading    string
	Definition string
	Document   string
	Sub        Sections
}

type Sections []*Section

func (s Sections) Len() int {
	return len(s)
}

func (s Sections) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Sections) Less(i, j int) bool {
	if s[i].Order == s[j].Order {
		return s[i].Path < s[j].Path
	}
	return s[i].Order < s[j].Order
}

func parseHeading(sec string) (path string, heading string) {
	sec = strings.TrimSpace(sec)
	i := strings.LastIndex(sec, "{")
	if i >= 0 {
		j := strings.Index(sec[i+1:], "}")
		if j >= 0 {
			path = sec[i+1 : i+1+j]
			sec = sec[:i]
			heading = strings.TrimSpace(sec)
			return path, heading
		}
	}
	path = sec
	heading = sec
	return path, heading
}

func BuildSections(groups []Directives) Sections {
	var sections Sections
	idx := map[string]*Section{}
	for _, dirs := range groups {
		sec := dirs["sec"]
		if sec == "" {
			continue
		}
		if i := strings.Index(sec, "\n"); i >= 0 {
			sec = sec[:i]
		}
		var section Section
		section.Path, section.Heading = parseHeading(sec)
		if len(section.Path) == 0 {
			continue
		}
		if s, ok := dirs["ord"]; ok {
			s = strings.TrimSpace(s)
			if i := strings.Index(s, "\n"); i >= 0 {
				s = s[:i]
			}
			order, err := strconv.ParseInt(s, 10, 64)
			if err == nil {
				section.Order = int(order)
			}
		}
		section.Definition = dirs["def"]
		section.Document = dirs["doc"]
		sections = append(sections, &section)
		idx[section.Path] = &section
	}

	var root Sections
	for _, section := range sections {
		path := strings.Split(section.Path, ".")
		if len(path) <= 1 {
			root = append(root, section)
			continue
		}
		parentPath := path[:len(path)-1]
		if parent, ok := idx[strings.Join(parentPath, ".")]; ok {
			parent.Sub = append(parent.Sub, section)
		} else {
			root = append(root, section)
		}
	}

	for _, section := range sections {
		sort.Stable(section.Sub)
	}
	sort.Stable(root)

	return root
}

func sanitizeAnchorName(text string) string {
	// From github.com/shurcooL/sanitized_anchor_name
	var anchorName []rune
	var futureDash = false
	for _, r := range text {
		switch {
		case unicode.IsLetter(r) || unicode.IsNumber(r):
			if futureDash && len(anchorName) > 0 {
				anchorName = append(anchorName, '-')
			}
			futureDash = false
			anchorName = append(anchorName, unicode.ToLower(r))
		default:
			futureDash = true
		}
	}
	return string(anchorName)
}

func generateHeadingLink(heading string) string {
	heading = strings.ToLower(heading)
	heading = strings.ReplaceAll(heading, " ", "-")
	return "#user-content-" + sanitizeAnchorName(heading)
}

func generateSection(b *strings.Builder, s *Section, level int) {
	b.WriteString(strings.Repeat("#", level))
	b.WriteByte(' ')
	b.WriteString(s.Heading)
	b.WriteByte('\n')
	b.WriteByte('[')
	b.WriteString(s.Path)
	b.WriteString("]: ")
	b.WriteString(generateHeadingLink(s.Heading))
	b.WriteByte('\n')
	if s.Definition != "" {
		b.WriteString("```\n")
		b.WriteString(strings.TrimSpace(s.Definition))
		b.WriteString("\n```\n")
	}
	b.WriteByte('\n')
	if doc := strings.TrimSpace(s.Document); doc != "" {
		b.WriteString(doc)
		b.WriteString("\n\n")
	}
	for _, sub := range s.Sub {
		generateSection(b, sub, level+1)
	}
}

func GenerateMarkdown(sections Sections) []byte {
	var b strings.Builder
	for _, section := range sections {
		generateSection(&b, section, 1)
	}
	return []byte(b.String())
}

func HandleFile(input, output string) {
	if input == output {
		fmt.Printf("%s: output file is the same as input file", input)
		return
	}
	f, err := os.Open(input)
	if err != nil {
		fmt.Println(err)
		return
	}
	var buf bytes.Buffer
	_, err = io.Copy(&buf, f)
	f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	groups := (&parser{Buffer: &buf}).Parse()
	groupDirectives := make([]Directives, len(groups))
	for i, group := range groups {
		groupDirectives[i] = ParseGroup(group)
	}
	sections := BuildSections(groupDirectives)
	b := GenerateMarkdown(sections)
	if len(bytes.TrimSpace(b)) == 0 {
		return
	}
	os.MkdirAll(filepath.Dir(output), 0755)
	if err := ioutil.WriteFile(output, b, 0666); err != nil {
		fmt.Println(err)
	}
}

func replaceExt(path, ext string) string {
	return path[:len(path)-len(filepath.Ext(path))] + ext
}

const Usage = `qdoc [-ext EXT] [-o DIR] [GLOB...]

Converts documentation within a Lua script into a Markdown file.

Each non-flag argument is a glob matching a number of input files. If no files
are specified, then all files with the .lua extension in the working directory
are matched.

The content of each matched file, assumed to have the Lua syntax, is parsed for
comments. Adjacent line comments are merged into a single group, while block
comments always have their own group.

Groups may contain directives of the form "@name" to control how they are
interpreted and converted. A directive divides the group into chunks; the name
of the directive determines the chunk type, while the rest of the chunk is the
content. An optional ":" may follow the directive to separate it from the
content as needed. The leading and trailing spaces of the content are trimmed.
Multiple directives of the same type can be specified, in which case the content
of each are concatenated together with newlines as separators, to form a single
chunk.

The following directives are processed:

@sec

    Describes a section. Only the first line of the chunk is processed.

    The content determines the text of the section heading. Additionally, it is
    a dot-separated list of names, which determines the level and location of
    the section in the generated Markdown file. For example, section "Foo.Bar"
    will appear under section "Foo". If there is no parent section, then the
    section is placed in the root.

    If the content ends with a path wrapped in curly braces, that is used as the
    path instead, while the text before is used as the heading. This can be used
    to finely control the path, or to make it different from the heading. For
    example:

        Description of Foo.Bar {Foo.Bar}

    Each section defines a footnote that links to the heading, using the path as
    the name, allowing the section to be linked to more easily. For example:

        See [section Foo.Bar][Foo.Bar] for more information.

@ord

    Describes the priority of the section. Only the first line of the chunk is
    processed.

    The content is a signed integer, which determines the order of the section
    among its siblings, sorted by ascending priority. By default, sections have
    a priority of 0. Sections that have the same priority are sorted by path.

@def

    Describes a type definition. The content is enclosed in a fenced code block
    that appears at the top of the section.

    This directive is optional.

@doc

    Describes the content of the section. This is general markdown text.

Other directives are ignored. An entire group is ignored if it does not contain
a valid @sec directive.

An output file is not written if its content is empty or only contains spacing.

The -o flag is a path to a directory to which output files will be written. If
-o is not specified, then an output file is written to the same directory as the
input file.

To determine the location of an output file, the extension of each matched file
is replaced with .md, which can be changed with the -ext flag. An input file is
skipped if the output and input are the same file.

If no input files are specified, then all files with the .lua extension in the
working directory are matched.

Flags:
`

func main() {
	var output string
	var ext string
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), Usage)
		flag.PrintDefaults()
	}
	flag.StringVar(&output, "o", "", "Directory to which files will be written.")
	flag.StringVar(&ext, "ext", ".md", "Extension of output files.")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		args = append(args, "*.lua")
	}
	var files []string
	for _, arg := range args {
		matches, err := filepath.Glob(arg)
		if err != nil {
			fmt.Printf("bad pattern: %s\n", err)
			return
		}
		files = append(files, matches...)
	}

	if output == "" {
		for _, file := range files {
			HandleFile(file, replaceExt(file, ext))
		}
		return
	}
	for _, file := range files {
		HandleFile(file, filepath.Join(output, replaceExt(filepath.Base(file), ext)))
	}
}
