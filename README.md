# template_engine

## 🇺🇸 English

### Overview

`template_engine` is a Go library whose goal is **not to generate DOCX files from scratch**, but to **apply intelligent replacements to existing DOCX documents**.

In real-world scenarios, documents rarely start from zero.  
Most of the time, there are already **existing templates**, with:

- complex headers and footers
- specific fonts
- predefined text sizes
- spacing, alignment, and corporate styles
- an established visual identity

This library is designed to **reuse those documents exactly as they are**, changing **only what is necessary**, without breaking the original layout.

The project intentionally does **not** aim to be a full DOCX editor.  
Instead, it focuses on safely transforming content inside real-world templates.

---

This library also supports a lightweight markup syntax that can be used **inside replace values** when filling existing DOCX documents.

These markers allow rich formatting while **preserving the original document styles**, such as font family, font size, spacing, headers, footers, and layout configuration.

The markup is intentionally simple and explicit to avoid unexpected side effects.


---

### Supported Markup Tags

| Tag | Description |
|----|------------|
| `[b]...[/b]` | Applies **bold** formatting to the enclosed text |
| `[i]...[/i]` | Applies *italic* formatting |
| `[u]...[/u]` | Applies underline formatting |
| `[text:#RRGGBB]...[/text]` | Changes the **text color** using a hexadecimal color code |
| `[bg:#RRGGBB]...[/bg]` | Applies a **background / highlight color** to the text only |
| `[bgp:#RRGGBB]...[/bgp]` | Applies a **background color to the entire paragraph** |
| `[align:left]...[/align]` | Aligns the paragraph to the left |
| `[align:center]...[/align]` | Aligns the paragraph to the center |
| `[align:right]...[/align]` | Aligns the paragraph to the right |
| `[align:both]...[/align]` | Justifies the paragraph |
| `[ident]...[/ident]` | Applies a default paragraph indentation |
| `[ident:N]...[/ident]` | Applies paragraph indentation with a specific level (`N`) |
| `[tab]` | Inserts a single tab character |
| `[tab:N]` | Inserts `N` tab characters |
| `[br]` | Inserts a line break |

---

### Notes on Behavior

- All formatting is applied **only to the replaced content**
- Existing fonts, font sizes, and document styles are preserved
- `bg` applies background color at the **text (run) level**
- `bgp` applies background color at the **paragraph level**
- The library does not modify headers, footers, margins, or page layout

---

### The problem this library solves

In the Go ecosystem, there is currently **no free, mature, high-quality library** that allows developers to:

- open an existing DOCX file
- preserve its full structure
- safely replace content
- add **bold, italic, colors, alignment, indentation**, etc.
- without losing the original font and size settings

Existing tools often:
- generate documents from scratch
- require full layout control
- or break styles during replacement

This does **not reflect real-world usage**, especially in areas such as:
- legal documents
- contracts
- business reports
- institutional templates

---

### Project goal

The goal of `template_engine` is to:

- work **on top of existing DOCX files**
- replace content using **placeholders**
- allow **rich formatting** in a simple way
- preserve **original fonts, sizes, and styles**
- make document processing predictable and reliable

All of this without recreating the document — only **transforming what already exists**.

---

### Philosophy

- The DOCX file is the source of truth
- The library **does not create documents**, it **transforms them**
- Only `word/document.xml` is modified
- Everything else remains untouched
- The original layout is respected

---

### Key features

- Intelligent placeholders (`{{FIELD}}`, `{{BLOCK:FIELD}}`)
- Simple markup for formatting:
  - bold, italic, underline
  - text and background colors
  - alignment and indentation
- Full preservation of existing styles
- Extensible architecture
- Open-source and free

---

### Placeholder Functions

In addition to simple placeholders like `{{FIELD}}`, the library also supports **placeholder functions**, allowing values to be transformed at replacement time.

The syntax follows this pattern:

{{FUNCTION:FIELD}}

Where:
- **FUNCTION** defines how the value will be transformed
- **FIELD** is the key from the replace map

This approach allows common transformations to be handled directly by the engine, keeping templates clean and avoiding repetitive logic in application code.

Placeholder functions must be declared directly in the DOCX template file using the `{{FUNCTION:FIELD}}` syntax.

They are evaluated during DOCX processing, when the engine reads and transforms the existing `word/document.xml`, and are applied only to the replaced content.

#### Supported Functions

| Function | Description |
|--------|-------------|
| `UPPER` | Converts the value to uppercase |
| `LOWER` | Converts the value to lowercase |
| `TITLE` | Converts the value to title case |
| `MASK_CPF` | Applies Brazilian CPF mask formatting |
| `PHONE_BR` | Formats a Brazilian phone number |
| `DATE_BR` | Formats a date using the Brazilian standard (`DD/MM/YYYY`) |

#### Examples
{{UPPER:name}}
{{LOWER:email}}
{{TITLE:full_name}}
{{DATE_BR:birth_date}}
{{MASK_CPF:cpf}}
{{PHONE_BR:phone}}


#### Execution Order

Placeholder functions are applied **before** markup processing.

This ensures that formatting tags such as `[b]`, `[bg]`, `[bgp]`, and alignment markers operate on the final transformed value.

### About the use of Artificial Intelligence

This project was developed with the **assistance of Artificial Intelligence tools**, used to support architectural reasoning, code review, and idea organization.

All final decisions, technical validations, and directions were made deliberately, with a strong focus on quality and real-world applicability.

AI was used as a **tool**, not as a replacement for engineering.

---

## 🇧🇷 Português

### Visão geral

`template_engine` é uma biblioteca escrita em Go cujo objetivo **não é gerar documentos DOCX do zero**, mas sim **aplicar substituições inteligentes (replaces) em arquivos DOCX já existentes**.

Na vida real, documentos raramente começam do zero.  
Normalmente já existem **modelos prontos**, com:

- cabeçalhos e rodapés complexos
- fontes específicas
- tamanhos de texto definidos
- espaçamentos, alinhamentos e estilos corporativos
- identidade visual já consolidada

O propósito desta biblioteca é **reaproveitar esses arquivos exatamente como eles são**, alterando **apenas o conteúdo necessário**, sem quebrar o layout original.

Este projeto **não tem como objetivo ser um editor completo de DOCX**.  
Ele foca exclusivamente em transformar conteúdo de forma segura dentro de documentos reais já existentes.

---

A biblioteca também suporta uma sintaxe de marcação simples que pode ser utilizada **dentro dos valores de replace** ao preencher documentos DOCX já existentes.

Esses marcadores permitem formatação rica **sem alterar os estilos originais do documento**, como fonte, tamanho de texto, espaçamentos, cabeçalhos, rodapés e configurações de layout.

A marcação foi pensada para ser explícita e previsível, evitando efeitos colaterais inesperados.


---

### Marcadores suportados

| Marcador | Descrição |
|--------|-----------|
| `[b]...[/b]` | Aplica **negrito** ao texto |
| `[i]...[/i]` | Aplica *itálico* |
| `[u]...[/u]` | Aplica sublinhado |
| `[text:#RRGGBB]...[/text]` | Altera a **cor do texto** usando código hexadecimal |
| `[bg:#RRGGBB]...[/bg]` | Aplica **cor de fundo apenas ao texto** |
| `[bgp:#RRGGBB]...[/bgp]` | Aplica **cor de fundo ao parágrafo inteiro** |
| `[align:left]...[/align]` | Alinha o parágrafo à esquerda |
| `[align:center]...[/align]` | Alinha o parágrafo ao centro |
| `[align:right]...[/align]` | Alinha o parágrafo à direita |
| `[align:both]...[/align]` | Justifica o parágrafo |
| `[ident]...[/ident]` | Aplica recuo padrão no parágrafo |
| `[ident:N]...[/ident]` | Aplica recuo no parágrafo com nível específico (`N`) |
| `[tab]` | Insere uma tabulação |
| `[tab:N]` | Insere `N` tabulações |
| `[br]` | Insere uma quebra de linha |

---

### Observações importantes

- A formatação é aplicada **somente ao conteúdo substituído**
- Fontes, tamanhos de texto e estilos do documento original são preservados
- `bg` atua no **nível do texto**
- `bgp` atua no **nível do parágrafo**
- Cabeçalhos, rodapés, margens e layout não são alterados

---

### O problema que esta biblioteca resolve

No ecossistema Go, **não existe hoje uma biblioteca gratuita, madura e de boa qualidade** que permita:

- abrir um DOCX existente
- preservar toda a estrutura original
- aplicar substituições de texto de forma segura
- adicionar **negrito, itálico, cores, alinhamento, indentação**, etc.
- sem perder o padrão de fonte e tamanho definidos no documento original

Ferramentas existentes costumam:
- gerar documentos do zero
- exigir controle total do layout
- ou quebrar estilos ao substituir conteúdo

Isso **não reflete a realidade de uso**, especialmente em contextos como:
- documentos jurídicos
- contratos
- relatórios empresariais
- modelos institucionais

---

### Objetivo do projeto

O objetivo do `template_engine` é:

- trabalhar **sobre DOCX existentes**
- substituir conteúdo usando **placeholders**
- permitir **formatação rica** de forma simples
- manter **fontes, tamanhos e estilos originais**
- tornar o processo previsível, estável e reutilizável

Tudo isso sem reinventar o documento, apenas **transformando o que já existe**.

---

### Filosofia

- O DOCX é a fonte da verdade
- A biblioteca **não cria documentos**, ela **transforma**
- Apenas o `word/document.xml` é alterado
- Todo o resto do arquivo permanece intacto
- A engine respeita o layout original

---

### Diferenciais

- Replaces inteligentes (`{{FIELD}}`, `{{BLOCK:FIELD}}`)
- Suporte a marcação simples:
  - negrito, itálico, sublinhado
  - cores de texto e fundo
  - alinhamento e indentação
- Preservação total de estilos existentes
- Arquitetura extensível
- Código aberto e gratuito

---

### Funções de Placeholder

Além de placeholders simples como `{{FIELD}}`, a biblioteca também suporta **funções de placeholder**, permitindo transformar valores no momento da substituição.

A sintaxe segue o padrão:

{{FUNCAO:CAMPO}}

Onde:
- **FUNCAO** define como o valor será transformado
- **CAMPO** corresponde à chave no mapa de replaces

Essa abordagem permite tratar transformações comuns diretamente na engine, mantendo os templates limpos e evitando lógica repetitiva no código da aplicação.

As funções de placeholder devem ser declaradas diretamente no arquivo DOCX utilizado como template, utilizando a sintaxe `{{FUNCAO:CAMPO}}`.

Elas são avaliadas durante o processamento do DOCX, no momento em que a engine lê e transforma o `word/document.xml`, sendo aplicadas apenas ao conteúdo substituído.

#### Funções suportadas

| Função | Descrição |
|------|-----------|
| `UPPER` | Converte o valor para maiúsculas |
| `LOWER` | Converte o valor para minúsculas |
| `TITLE` | Converte o valor para formato título |
| `MASK_CPF` | Aplica máscara de CPF |
| `PHONE_BR` | Formata número de telefone brasileiro |
| `DATE_BR` | Formata data no padrão brasileiro (`DD/MM/AAAA`) |

#### Exemplos

{{UPPER:nome}}
{{LOWER:email}}
{{TITLE:nome_completo}}
{{DATE_BR:data_nascimento}}
{{MASK_CPF:cpf}}
{{PHONE_BR:telefone}}

#### Ordem de execução

As funções de placeholder são aplicadas **antes** da marcação visual.

Isso garante que marcadores como `[b]`, `[bg]`, `[bgp]` e alinhamentos atuem sobre o valor final já transformado.

---

### Sobre o uso de Inteligência Artificial

Este projeto foi desenvolvido com o **apoio de ferramentas de Inteligência Artificial**, utilizadas como auxílio no raciocínio arquitetural, revisão de código e organização de ideias.

Todas as decisões finais, validações técnicas e direcionamentos foram feitos de forma consciente, com foco em qualidade, clareza e aplicabilidade real.

A IA foi usada como **ferramenta**, não como substituto de engenharia.

---

## License

Licensed under the Apache License, Version 2.0.  
See the `LICENSE` file for details.
