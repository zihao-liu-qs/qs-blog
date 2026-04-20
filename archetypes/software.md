---
title: "{{ replace .Name "-" " " | title }}"
date: {{ .Date }}
draft: false
tagline: "一句话描述"
price: "¥0"
featured: false
platform: ["macOS"]
status: "available"   # available | coming-soon
buy_url: "#"          # 后续改为 /api/checkout?product={{ .Name }} 或第三方链接
---

产品描述写在这里。

## 功能

- 功能一
- 功能二

## 系统要求

macOS XX 及以上
