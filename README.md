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
- established visual identity

This library is designed to **reuse those documents exactly as they are**, changing **only what is necessary**, without breaking the original layout.

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

### Sobre o uso de Inteligência Artificial

Este projeto foi desenvolvido com o **apoio de ferramentas de Inteligência Artificial**, utilizadas como auxílio no raciocínio arquitetural, revisão de código e organização de ideias.

Todas as decisões finais, validações técnicas e direcionamentos foram feitos de forma consciente, com foco em qualidade, clareza e aplicabilidade real.

A IA foi usada como **ferramenta**, não como substituto de engenharia.

---

## License

Licensed under the Apache License, Version 2.0.  
See the `LICENSE` file for details.
