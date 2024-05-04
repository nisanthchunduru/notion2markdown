package notion2markdown

import (
	"context"
	"fmt"
	"strings"

	"github.com/jomei/notionapi"
)

type Notion2Markdown struct {
	NotionToken string
}

func (notion2Markdown Notion2Markdown) PageToMarkdown(pageIdString string) (string, error) {
	jomeiNotionApiClient := notionapi.NewClient(notionapi.Token(notion2Markdown.NotionToken))
	return PageToMarkdown(jomeiNotionApiClient, pageIdString)
}

func PageToMarkdown(jomeiNotionApiClient *notionapi.Client, pageIdString string) (string, error) {
	pagination := notionapi.Pagination{
		PageSize: 100,
	}
	pageId := notionapi.BlockID(pageIdString)
	getChildPageChildrenResponse, err := jomeiNotionApiClient.Block.GetChildren(context.Background(), pageId, &pagination)
	if err != nil {
		return "", err
	}
	childPageBlocks := getChildPageChildrenResponse.Results
	markdown := BlocksToMarkdown(childPageBlocks)
	return markdown, nil
}

func BlocksToMarkdown(blocks []notionapi.Block) string {
	var markdowns []string
	for i, block := range blocks {
		var markdown string

		if block.GetType() == "heading_1" {
			heading1Block := block.(*notionapi.Heading1Block)
			markdown = Heading1ToMarkdown(heading1Block.Heading1)
		} else if block.GetType() == "heading_2" {
			heading2Block := block.(*notionapi.Heading2Block)
			markdown = Heading2ToMarkdown(heading2Block.Heading2)
		} else if block.GetType() == "heading_3" {
			heading3Block := block.(*notionapi.Heading3Block)
			markdown = Heading3ToMarkdown(heading3Block.Heading3)
		} else if block.GetType() == "paragraph" {
			paragraphBlock := block.(*notionapi.ParagraphBlock)
			markdown = ParagraphToMarkdown(paragraphBlock.Paragraph)
		} else if block.GetType() == "bulleted_list_item" {
			bulletedListItemBlock := block.(*notionapi.BulletedListItemBlock)
			markdown = BulletedListItemToMarkdown(bulletedListItemBlock.BulletedListItem)
			if (i + 1) < len(blocks) {
				nextBlock := blocks[i+1]
				if nextBlock.GetType() != "bulleted_list_item" {
					markdown = markdown + "\n"
				}
			}
		} else if block.GetType() == "numbered_list_item" {
			numberedListItemBlock := block.(*notionapi.NumberedListItemBlock)
			markdown = NumberedListItemToMarkdown(numberedListItemBlock.NumberedListItem)
			if (i + 1) < len(blocks) {
				nextBlock := blocks[i+1]
				if nextBlock.GetType() != "numbered_list_item" {
					markdown = markdown + "\n"
				}
			}
		} else if block.GetType() == "code" {
			codeBlock := block.(*notionapi.CodeBlock)
			markdown = CodeToMarkdown(codeBlock.Code)
			if (i + 1) < len(blocks) {
				nextBlock := blocks[i+1]
				if nextBlock.GetType() != "numbered_list_item" {
					markdown = markdown + "\n"
				}
			}
		} else if block.GetType() == "quote" {
			quoteBlock := block.(*notionapi.QuoteBlock)
			markdown = QuoteToMarkdown(quoteBlock.Quote)
		} else if block.GetType() == "image" {
			imageBlock := block.(*notionapi.ImageBlock)
			markdown = ImageToMarkdown(imageBlock.Image)
		}

		if markdown != "" {
			markdowns = append(markdowns, markdown)
		}
	}
	return strings.Join(markdowns, "")
}

func Heading1ToMarkdown(heading1 notionapi.Heading) string {
	markdown := RichTextsToMarkdown(heading1.RichText)
	markdown = "# " + markdown + "\n\n"
	return markdown
}

func Heading2ToMarkdown(heading2 notionapi.Heading) string {
	markdown := RichTextsToMarkdown(heading2.RichText)
	markdown = "## " + markdown + "\n\n"
	return markdown
}

func Heading3ToMarkdown(heading3 notionapi.Heading) string {
	markdown := RichTextsToMarkdown(heading3.RichText)
	markdown = "### " + markdown + "\n\n"
	return markdown
}

func ParagraphToMarkdown(paragraph notionapi.Paragraph) string {
	markdown := RichTextsToMarkdown(paragraph.RichText)
	markdown = markdown + "\n"
	return markdown
}

func BulletedListItemToMarkdown(bulleted_list_item notionapi.ListItem) string {
	markdown := RichTextsToMarkdown(bulleted_list_item.RichText)
	markdown = "- " + markdown + "\n"
	return markdown
}

func NumberedListItemToMarkdown(numbered_list_item notionapi.ListItem) string {
	markdown := RichTextsToMarkdown(numbered_list_item.RichText)
	markdown = "1. " + markdown + "\n"
	return markdown
}

func CodeToMarkdown(code notionapi.Code) string {
	markdown := RichTextsToMarkdown(code.RichText)
	markdown = "```\n" + markdown + "\n```\n"
	return markdown
}

func QuoteToMarkdown(code notionapi.Quote) string {
	markdown := RichTextsToMarkdown(code.RichText)
	markdown = "> " + markdown + "\n\n"
	return markdown
}

func ImageToMarkdown(image notionapi.Image) string {
	var markdown string
	if image.File != nil && image.File.URL != "" {
		markdown = fmt.Sprintf("![Untitled](%s)", image.File.URL)
	}
	return markdown
}

func RichTextsToMarkdown(richTexts []notionapi.RichText) string {
	var markdowns []string
	for _, block := range richTexts {
		var markdown string
		if block.Href != "" {
			markdown = fmt.Sprintf("[%s](%s)", block.PlainText, block.Href)
		} else if block.Annotations.Code {
			markdown = "`" + block.PlainText + "`"
		} else {
			markdown = block.PlainText
		}
		markdowns = append(markdowns, markdown)
	}
	return strings.Join(markdowns, "")
}
