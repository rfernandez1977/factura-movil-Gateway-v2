# Gateway API Endpoints

## Overview
This document lists all available endpoints for the Gateway API, which integrates with Factura Móvil to manage electronic documents, clients, and products.

## Authentication
All endpoints require an API key passed in the `X-API-Key` header.

## Endpoints

### Document Creation
- **POST /facturas**
  - Creates an invoice.
  - Request Body: JSON object with fields like `date`, `details`, `netTotal`.
  - Response: `200 OK` with `{ "status": "factura created successfully" }`.

- **POST /boletas**
  - Creates a ticket (boleta).
  - Request Body: JSON object with fields like `date`, `details`, `netTotal`.
  - Response: `200 OK` with `{ "status": "boleta created successfully" }`.

- **POST /notas**
  - Creates a credit or debit note.
  - Request Body: JSON object with fields like `date`, `details`, `netTotal`.
  - Response: `200 OK` with `{ "status": "nota created successfully" }`.

- **POST /guias**
  - Creates a dispatch guide.
  - Request Body: JSON object with fields like `date`, `details`, `netTotal`.
  - Response: `200 OK` with `{ "status": "guia created successfully" }`.

### Entity Creation
- **POST /clientes**
  - Creates a client.
  - Request Body: JSON object with fields like `code`, `name`, `address`.
  - Response: `200 OK` with `{ "status": "cliente created successfully" }`.

- **POST /productos**
  - Creates a product.
  - Request Body: JSON object with fields like `code`, `name`, `price`.
  - Response: `200 OK` with `{ "status": "producto created successfully" }`.

### Document Queries
- **GET /documents/:id**
  - Retrieves the status of a document.
  - Response: `200 OK` with the document status in JSON format.

- **GET /documents/:id/pdf**
  - Downloads the PDF of a document.
  - Response: `200 OK` with the PDF content (`application/pdf`).

### Metrics
- **GET /metrics**
  - Exposes Prometheus metrics for monitoring.
  - Response: Prometheus metrics in text format.