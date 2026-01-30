# Photo Listing SaaS

![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)
[![CI](https://github.com/EnockYator/saas-photo-listing-platform/actions/workflows/ci.yml/badge.svg)](https://github.com/EnockYator/saas-photo-listing-platform/actions/workflows/ci.yml)

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org/)
[![Documentation](https://img.shields.io/badge/docs-comprehensive-blue)](https://github.com/EnockYator/saas-photo-listing-platform/docs/overview.md)


## **Introduction**

**Photo Listing SaaS** is a **multi-tenant**, **professional portfolio platform** built exclusively for **photographers**, **studios**, and **creative businesses**. Unlike generic photo-sharing platforms, we provide business-enabling tools with strict **tenant isolation**, subscription-based monetization, and enterprise-grade reliability—all built with **Golang** for **maximum performance** and **maintainability**.

**Core Philosophy**: We don't compete with social networks. We provide the digital business backbone for photography professionals to showcase work, manage client deliveries, and grow their business through a polished, professional platform.

---

## Table of Contents

- [Photo Listing SaaS](#photo-listing-saas)
  - [**Introduction**](#introduction)
  - [Table of Contents](#table-of-contents)
  - [Project Overview](#project-overview)
    - [Key Capabilities](#key-capabilities)
  - [Core Features](#core-features)
  - [Architecture \& Design Principles](#architecture--design-principles)
    - [High-Level Flow](#high-level-flow)
    - [Layer responsibilities](#layer-responsibilities)
  - [Tech Stack](#tech-stack)
  - [Project Structure](#project-structure)
  - [Setup \& Installation](#setup--installation)
    - [Prerequisites](#prerequisites)
    - [Quick Start](#quick-start)
      - [Clone the Repository](#clone-the-repository)
      - [Start all Services](#start-all-services)
      - [Access URLs](#access-urls)
      - [Stopping the Application](#stopping-the-application)
  - [Database \& SQLC](#database--sqlc)
  - [API Documentation](#api-documentation)
  - [Environment Variables](#environment-variables)
    - [Server](#server)
    - [Database](#database)
    - [Object Storage](#object-storage)
  - [Developer Workflow](#developer-workflow)
  - [Contributing](#contributing)
  - [License](#license)
  - [Future Enhancements](#future-enhancements)

---