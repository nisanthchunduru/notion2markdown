# notion2markdown

`notion2markdown` is a Go package to export Notion pages as markdown

## Usage

Add  `notion2markdown` to your Go project

```
go get github.com/nisanthchunduru/notion2markdown
```

Create a Notion integration, generate a secret and connect that integration to the Notion page you'd like to export https://developers.notion.com/docs/create-a-notion-integration#getting-started

Export a Notion page to markdown

```go
notionToken := "your_notion_secret"
notion2Markdown := notion2markdown.Notion2Markdown{
  NotionToken: notionToken,
}
notionPageId := "a_notion_page_id"
markdown, err := notion2Markdown.exportPage(notionPageId)
```

## Todos

- Support text annotations like italic etc.
