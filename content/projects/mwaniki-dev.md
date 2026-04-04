---
title: mwaniki.dev
tags: Go, HTMX, CSS
description: This site. Server-rendered with Go and HTMX, styled with intention.
icon: "05"
reading_time: 2 min
date: 2025-06-10
---

This portfolio is a complete overhaul built with Go 1.24 standard library, HTMX and pure CSS. No frameworks, no bundlers, no build step.

## Design Principles

- Minimalism — every element earns its place
- Reusability — components that work across pages
- Inclusivity — accessible by default
- Clarity — easy to understand over clever

## Architecture

A single Go binary using `net/http` for routing and `html/template` for rendering. HTMX handles partial page swaps. CSS handles all layout, animation and responsive behaviour.
