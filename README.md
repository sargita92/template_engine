# template_engine

Extensible Go library to fill DOCX templates using placeholders and a custom markup language for rich text formatting.

This project provides a powerful yet lightweight template engine focused on **DOCX document generation**, allowing developers to inject dynamic content, styles, and formatting into Word documents programmatically.

---

## ✨ Features

- Placeholder replacement:
    - `{{FIELD}}`
    - `{{BLOCK:FIELD}}`
- Custom lightweight markup language:
    - Text styles: `[b]`, `[i]`, `[u]`
    - Text color: `[text:#RRGGBB]`
    - Background / highlight: `[bg:#RRGGBB]`
    - Alignment: `[align:left|center|right|both]`
    - Indentation: `[ident]`, `[ident:N]`
    - Tabs and line breaks: `[tab]`, `[tab:N]`, `[br]`
- Built-in functions:
    - `UPPER`, `LOWER`, `TITLE`
    - `DATE_BR`
    - `MASK_CPF`
    - `PHONE_BR`
- Automatic promotion of styled placeholders to block rendering
- Native DOCX support (WordprocessingML)
- Designed for extensibility (future renderers: HTML, PDF, etc.)

---

## 📦 Installation

```bash
go get github.com/sargita92/template_engine
