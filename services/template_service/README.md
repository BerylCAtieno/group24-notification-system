# 📄 **Template Service**

A robust **Django-based template rendering service** supporting
multi-language templates, dynamic variable substitution, versioning,
preview mode, and optional Kafka event publishing.

## 🚀 Features

-   Create, read, update, delete templates\
-   Multi-language + version support\
-   Single and batch rendering\
-   Jinja2 rendering engine\
-   Optional Kafka integration\
-   Strong validation for variables and types\
-   Health-check endpoint with Kafka status

## 📚 Table of Contents

1.  Installation\
2.  Environment Setup\
3.  Database\
4.  Running the Service\
5.  API Endpoints\
6.  Testing\
7.  Deployment\
8.  Data Model Overview

# 🔧 Installation

``` bash
git clone <repository_url>
cd template_service
python -m venv venv
source venv/bin/activate   # Linux/macOS  
venv\Scripts\activate      # Windows
```

Install dependencies:

``` bash
pip install -r requirements.txt
```

# ⚙️ Environment Setup

Create a `.env` file:

``` env
DEBUG=True
SECRET_KEY=your_secret_key
DJANGO_ALLOWED_HOSTS=localhost 127.0.0.1
DATABASE_URL=postgres://user:password@localhost:5432/template_db
KAFKA_BOOTSTRAP_SERVERS=localhost:9092
```

# 🗄️ Database

### Run migrations

``` bash
python manage.py makemigrations
python manage.py migrate
```

### Database switching

``` python
DATABASES = {
    "default": dj_database_url.config(default="sqlite:///db.sqlite3")
}
```

# ▶️ Running the Service

### Development Server

``` bash
python manage.py runserver
```

### Health Check

    GET /api/v1/health

# 🔌 API Endpoints

## 📍 1. Create Template

### POST /templates/

``` json
{
  "name": "Welcome Email",
  "version": "1.0.0",
  "language": "en",
  "type": "email",
  "subject": "Welcome to {{app_name}}, {{user_name}}!",
  "body": {
    "html": "<h1>Hello {{user_name}}</h1><p>Welcome to {{app_name}}</p>",
    "text": "Hello {{user_name}}, welcome to {{app_name}}"
  },
  "variables": [
    {
      "name": "user_name",
      "type": "string",
      "required": true,
      "description": "User display name"
    },
    {
      "name": "app_name",
      "type": "string",
      "required": true,
      "description": "Application name"
    }
  ],
  "metadata": {
    "created_by": "admin@example.com",
    "tags": ["onboarding", "transactional"]
  }
}
```

## 📍 2. List Templates

### GET /templates/

Filters: `type`, `language`, `tags`, `page`, `per_page`

## 📍 3. Get Template by ID

### GET /templates/`<template_id>`{=html}/

Supports `language` and `version`.

## 📍 4. Render Template

### POST /templates/`<template_id>`{=html}/render/

``` json
{
  "language": "en",
  "version": "latest",
  "variables": {
    "user_name": "John Doe",
    "app_name": "MyApp",
    "verification_url": "https://example.com/verify"
  },
  "preview_mode": false
}
```

## 📍 5. Batch Render Templates

### POST /templates/render/batch/

``` json
{
  "requests": [
    {
      "template_id": "welcome_email",
      "language": "en",
      "variables": {
        "user_name": "John",
        "app_name": "MyApp"
      }
    },
    {
      "template_id": "password_reset",
      "language": "en",
      "variables": {
        "user_name": "Jane",
        "reset_url": "https://example.com/reset"
      }
    }
  ]
}
```

## 📍 6. Delete Template

### DELETE /templates/`<template_id>`{=html}/

Response: 204 No Content

# 🧱 Data Model Overview

``` json
{
  "template_id": "uuid",
  "name": "string",
  "version": "string",
  "language": "string",
  "type": "email | sms | push",
  "subject": "string",
  "body": {
    "html": "string",
    "text": "string"
  },
  "variables": [
    {
      "name": "string",
      "type": "string|number|boolean|date",
      "required": true,
      "description": "string"
    }
  ],
  "metadata": {
    "tags": ["tag1", "tag2"],
    "created_by": "email",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
}
```

# 🧪 Testing

``` bash
python manage.py test
```

# 🚀 Deployment

### Gunicorn

``` bash
gunicorn template_service.wsgi:application --bind 0.0.0.0:8000
```

### Heroku Procfile

    web: gunicorn template_service.wsgi:application --log-file -

# 📝 Notes

-   Kafka optional\
-   Batch render supports up to 50 items\
-   Preview mode skips validation
