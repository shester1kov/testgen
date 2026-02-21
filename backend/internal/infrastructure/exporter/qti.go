package exporter

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/shester1kov/testgen-backend/internal/domain/entity"
)

// QTI (Question and Test Interoperability) format exporter
// QTI is an IMS Global standard for representing assessment content
// This implements QTI 2.1 format used by Canvas, Sakai, and other LMS

// QTIAssessment represents the root element
type QTIAssessment struct {
	XMLName    xml.Name      `xml:"assessmentTest"`
	Xmlns      string        `xml:"xmlns,attr"`
	XmlnsXsi   string        `xml:"xmlns:xsi,attr"`
	Identifier string        `xml:"identifier,attr"`
	Title      string        `xml:"title,attr"`
	TestParts  []QTITestPart `xml:"testPart"`
}

// QTITestPart represents a test part
type QTITestPart struct {
	Identifier         string       `xml:"identifier,attr"`
	NavigationMode     string       `xml:"navigationMode,attr"`
	SubmissionMode     string       `xml:"submissionMode,attr"`
	AssessmentSections []QTISection `xml:"assessmentSection"`
}

// QTISection represents an assessment section
type QTISection struct {
	Identifier string       `xml:"identifier,attr"`
	Title      string       `xml:"title,attr"`
	Visible    string       `xml:"visible,attr"`
	Items      []QTIItemRef `xml:"assessmentItemRef"`
}

// QTIItemRef references an assessment item
type QTIItemRef struct {
	Identifier string `xml:"identifier,attr"`
	Href       string `xml:"href,attr"`
}

// QTIItem represents a single assessment item (question)
type QTIItem struct {
	XMLName       xml.Name               `xml:"assessmentItem"`
	Xmlns         string                 `xml:"xmlns,attr"`
	Identifier    string                 `xml:"identifier,attr"`
	Title         string                 `xml:"title,attr"`
	Adaptive      string                 `xml:"adaptive,attr"`
	TimeDependent string                 `xml:"timeDependent,attr"`
	ResponseDecl  []QTIResponseDecl      `xml:"responseDeclaration"`
	OutcomeDecl   []QTIOutcomeDecl       `xml:"outcomeDeclaration"`
	ItemBody      QTIItemBody            `xml:"itemBody"`
	ResponseProc  *QTIResponseProcessing `xml:"responseProcessing,omitempty"`
}

// QTIResponseDecl declares response variables
type QTIResponseDecl struct {
	Identifier  string          `xml:"identifier,attr"`
	Cardinality string          `xml:"cardinality,attr"`
	BaseType    string          `xml:"baseType,attr"`
	CorrectResp *QTICorrectResp `xml:"correctResponse,omitempty"`
}

// QTICorrectResp contains correct response values
type QTICorrectResp struct {
	Values []QTIValue `xml:"value"`
}

// QTIValue represents a value
type QTIValue struct {
	Value string `xml:",chardata"`
}

// QTIOutcomeDecl declares outcome variables
type QTIOutcomeDecl struct {
	Identifier   string    `xml:"identifier,attr"`
	Cardinality  string    `xml:"cardinality,attr"`
	BaseType     string    `xml:"baseType,attr"`
	DefaultValue *QTIValue `xml:"defaultValue>value,omitempty"`
}

// QTIItemBody contains the question content
type QTIItemBody struct {
	Content string `xml:",innerxml"`
}

// QTIResponseProcessing contains response processing rules
type QTIResponseProcessing struct {
	Template string `xml:"template,attr,omitempty"`
}

// QTIExporter exports tests to QTI format
type QTIExporter struct{}

// NewQTIExporter creates a new QTI exporter
func NewQTIExporter() *QTIExporter {
	return &QTIExporter{}
}

// Format returns the export format identifier
func (e *QTIExporter) Format() ExportFormat {
	return FormatQTI
}

// Name returns human-readable format name
func (e *QTIExporter) Name() string {
	return "QTI 2.1"
}

// Description returns format description
func (e *QTIExporter) Description() string {
	return "IMS QTI 2.1 standard format (Canvas, Sakai, and other LMS)"
}

// SupportedQuestionTypes returns list of supported question types
func (e *QTIExporter) SupportedQuestionTypes() []entity.QuestionType {
	return []entity.QuestionType{
		entity.QuestionTypeSingleChoice,
		entity.QuestionTypeMultipleChoice,
		entity.QuestionTypeTrueFalse,
		entity.QuestionTypeShortAnswer,
	}
}

// Export converts a test with questions and answers to QTI format
// Returns a single XML file containing all items
func (e *QTIExporter) Export(test *entity.Test, questions []*entity.Question, answers map[string][]*entity.Answer) (*ExportResult, error) {
	var builder strings.Builder

	// XML header
	builder.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	builder.WriteString("\n")

	// Create items
	var items []string
	for i, q := range questions {
		qAnswers := answers[q.ID.String()]

		item, err := e.formatItem(q, qAnswers, i+1)
		if err != nil {
			continue
		}
		items = append(items, item)
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("no compatible questions found for QTI export")
	}

	// Wrap in manifest/package structure
	testTitle := "Exported Test"
	if test != nil && test.Title != "" {
		testTitle = test.Title
	}

	builder.WriteString(fmt.Sprintf(`<questestinterop xmlns="http://www.imsglobal.org/xsd/ims_qtiasiv1p2">
  <assessment ident="assessment1" title="%s">
    <section ident="section1">
`, e.escapeXML(testTitle)))

	for _, item := range items {
		builder.WriteString(item)
		builder.WriteString("\n")
	}

	builder.WriteString(`    </section>
  </assessment>
</questestinterop>`)

	return &ExportResult{
		Content:     builder.String(),
		ContentType: "application/xml",
		FileExt:     "xml",
		Format:      FormatQTI,
	}, nil
}

// formatItem formats a single question as QTI item
func (e *QTIExporter) formatItem(q *entity.Question, answers []*entity.Answer, num int) (string, error) {
	var builder strings.Builder

	identifier := fmt.Sprintf("item%d", num)
	questionText := e.escapeXML(q.QuestionText)

	switch q.QuestionType {
	case entity.QuestionTypeSingleChoice:
		builder.WriteString(fmt.Sprintf(`      <item ident="%s" title="Question %d">
        <presentation>
          <material>
            <mattext texttype="text/plain">%s</mattext>
          </material>
          <response_lid ident="response1" rcardinality="Single">
            <render_choice>
`, identifier, num, questionText))

		correctIdent := ""
		for i, ans := range answers {
			ansIdent := fmt.Sprintf("choice%d", i+1)
			builder.WriteString(fmt.Sprintf(`              <response_label ident="%s">
                <material>
                  <mattext texttype="text/plain">%s</mattext>
                </material>
              </response_label>
`, ansIdent, e.escapeXML(ans.AnswerText)))
			if ans.IsCorrect {
				correctIdent = ansIdent
			}
		}

		builder.WriteString(`            </render_choice>
          </response_lid>
        </presentation>
        <resprocessing>
          <outcomes>
            <decvar varname="SCORE" vartype="Decimal" defaultval="0"/>
          </outcomes>
          <respcondition>
            <conditionvar>
`)
		builder.WriteString(fmt.Sprintf(`              <varequal respident="response1">%s</varequal>
`, correctIdent))
		builder.WriteString(`            </conditionvar>
            <setvar action="Set" varname="SCORE">100</setvar>
          </respcondition>
        </resprocessing>
      </item>`)

	case entity.QuestionTypeMultipleChoice:
		builder.WriteString(fmt.Sprintf(`      <item ident="%s" title="Question %d">
        <presentation>
          <material>
            <mattext texttype="text/plain">%s</mattext>
          </material>
          <response_lid ident="response1" rcardinality="Multiple">
            <render_choice>
`, identifier, num, questionText))

		var correctIdents []string
		for i, ans := range answers {
			ansIdent := fmt.Sprintf("choice%d", i+1)
			builder.WriteString(fmt.Sprintf(`              <response_label ident="%s">
                <material>
                  <mattext texttype="text/plain">%s</mattext>
                </material>
              </response_label>
`, ansIdent, e.escapeXML(ans.AnswerText)))
			if ans.IsCorrect {
				correctIdents = append(correctIdents, ansIdent)
			}
		}

		builder.WriteString(`            </render_choice>
          </response_lid>
        </presentation>
        <resprocessing>
          <outcomes>
            <decvar varname="SCORE" vartype="Decimal" defaultval="0"/>
          </outcomes>
          <respcondition>
            <conditionvar>
              <and>
`)
		for _, ident := range correctIdents {
			builder.WriteString(fmt.Sprintf(`                <varequal respident="response1">%s</varequal>
`, ident))
		}
		builder.WriteString(`              </and>
            </conditionvar>
            <setvar action="Set" varname="SCORE">100</setvar>
          </respcondition>
        </resprocessing>
      </item>`)

	case entity.QuestionTypeTrueFalse:
		isTrue := false
		for _, ans := range answers {
			if ans.IsCorrect {
				if strings.ToLower(strings.TrimSpace(ans.AnswerText)) == "true" {
					isTrue = true
				}
				break
			}
		}

		correctIdent := "false"
		if isTrue {
			correctIdent = "true"
		}

		builder.WriteString(fmt.Sprintf(`      <item ident="%s" title="Question %d">
        <presentation>
          <material>
            <mattext texttype="text/plain">%s</mattext>
          </material>
          <response_lid ident="response1" rcardinality="Single">
            <render_choice>
              <response_label ident="true">
                <material>
                  <mattext texttype="text/plain">True</mattext>
                </material>
              </response_label>
              <response_label ident="false">
                <material>
                  <mattext texttype="text/plain">False</mattext>
                </material>
              </response_label>
            </render_choice>
          </response_lid>
        </presentation>
        <resprocessing>
          <outcomes>
            <decvar varname="SCORE" vartype="Decimal" defaultval="0"/>
          </outcomes>
          <respcondition>
            <conditionvar>
              <varequal respident="response1">%s</varequal>
            </conditionvar>
            <setvar action="Set" varname="SCORE">100</setvar>
          </respcondition>
        </resprocessing>
      </item>`, identifier, num, questionText, correctIdent))

	case entity.QuestionTypeShortAnswer:
		builder.WriteString(fmt.Sprintf(`      <item ident="%s" title="Question %d">
        <presentation>
          <material>
            <mattext texttype="text/plain">%s</mattext>
          </material>
          <response_str ident="response1" rcardinality="Single">
            <render_fib>
              <response_label ident="answer"/>
            </render_fib>
          </response_str>
        </presentation>
        <resprocessing>
          <outcomes>
            <decvar varname="SCORE" vartype="Decimal" defaultval="0"/>
          </outcomes>
`, identifier, num, questionText))

		for _, ans := range answers {
			if ans.IsCorrect {
				builder.WriteString(fmt.Sprintf(`          <respcondition>
            <conditionvar>
              <varequal respident="response1" case="No">%s</varequal>
            </conditionvar>
            <setvar action="Set" varname="SCORE">100</setvar>
          </respcondition>
`, e.escapeXML(ans.AnswerText)))
			}
		}

		builder.WriteString(`        </resprocessing>
      </item>`)

	default:
		return "", fmt.Errorf("unsupported question type: %s", q.QuestionType)
	}

	return builder.String(), nil
}

// escapeXML escapes special XML characters
func (e *QTIExporter) escapeXML(s string) string {
	s = strings.TrimSpace(s)
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
		"'", "&apos;",
	)
	return replacer.Replace(s)
}
