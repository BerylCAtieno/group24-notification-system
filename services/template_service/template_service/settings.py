# -------------------------------
# IMPORTS
# -------------------------------
from kafka import KafkaProducer  # Optional Kafka producer
from pathlib import Path
from pythonjsonlogger import jsonlogger  # JSON logging format
import os
import logging
import dj_database_url  # Optional: parse DATABASE_URL for production
from dotenv import load_dotenv  # Load environment variables from .env
from decouple import config  # Easy configuration with defaults

# -------------------------------
# BASE DIRECTORY
# -------------------------------
BASE_DIR = Path(__file__).resolve().parent.parent

# Load .env file if it exists
load_dotenv(os.path.join(BASE_DIR, ".env"))

# -------------------------------
# LOGGING DIRECTORY
# -------------------------------
LOG_DIR = os.path.join(BASE_DIR, "logs")
os.makedirs(LOG_DIR, exist_ok=True)  # Ensure the directory exists

# -------------------------------
# KAFKA CONFIGURATION
# -------------------------------
USE_KAFKA = os.getenv("USE_KAFKA", "false").lower() == "true"
producer = None

if USE_KAFKA:
    try:
        producer = KafkaProducer(bootstrap_servers="localhost:9092")
    except Exception:
        producer = None  # Safe fallback if Kafka is unavailable

# Kafka topics
KAFKA_BOOTSTRAP_SERVERS = ["kafka:9092"]
KAFKA_TEMPLATE_TOPIC = "template_render_requests"
KAFKA_RESPONSE_TOPIC = "template_render_responses"

# -------------------------------
# SECRET KEY & DEBUG
# -------------------------------
SECRET_KEY = os.getenv("SECRET_KEY", "local-secret-key-for-dev")
DEBUG = os.getenv("DEBUG", "False").lower() in ["true", "1", "yes"]

ALLOWED_HOSTS = ["*"]  # For dev; update in production

# -------------------------------
# INSTALLED APPS
# -------------------------------
INSTALLED_APPS = [
    'django.contrib.admin',
    'django.contrib.auth',
    'django.contrib.contenttypes',
    'django.contrib.sessions',
    'django.contrib.messages',
    'django.contrib.staticfiles',
    "app",
    "rest_framework",
    "django_prometheus",
]

# -------------------------------
# MIDDLEWARE
# -------------------------------
MIDDLEWARE = [
    'django.middleware.security.SecurityMiddleware',
    'whitenoise.middleware.WhiteNoiseMiddleware',  # For static files
    "app.middleware.correlation.CorrelationIdMiddleware",  # Custom middleware
    "django_prometheus.middleware.PrometheusBeforeMiddleware",
    "django_prometheus.middleware.PrometheusAfterMiddleware",
    'django.contrib.sessions.middleware.SessionMiddleware',
    'django.middleware.common.CommonMiddleware',
    'django.middleware.csrf.CsrfViewMiddleware',
    'django.contrib.auth.middleware.AuthenticationMiddleware',
    'django.contrib.messages.middleware.MessageMiddleware',
    'django.middleware.clickjacking.XFrameOptionsMiddleware',
]

# -------------------------------
# STATIC FILES
# -------------------------------
STATICFILES_STORAGE = "whitenoise.storage.CompressedManifestStaticFilesStorage"
STATIC_URL = 'static/'
STATIC_ROOT = os.path.join(BASE_DIR, "staticfiles")

# -------------------------------
# URLS & WSGI
# -------------------------------
ROOT_URLCONF = 'template_service.urls'  # Make sure template_service/urls.py exists
WSGI_APPLICATION = 'template_service.wsgi.application'

# -------------------------------
# TEMPLATES
# -------------------------------
TEMPLATES = [
    {
        'BACKEND': 'django.template.backends.django.DjangoTemplates',
        'DIRS': [],  # Add paths if you have custom templates
        'APP_DIRS': True,
        'OPTIONS': {
            'context_processors': [
                'django.template.context_processors.request',
                'django.contrib.auth.context_processors.auth',
                'django.contrib.messages.context_processors.messages',
            ],
        },
    },
    {
        "BACKEND": "django.template.backends.jinja2.Jinja2",
        "APP_DIRS": True,
        "OPTIONS": {
            "environment": "template_service.jinja_env.environment",  # Custom Jinja2 environment
        },
    },
]

# -------------------------------
# DATABASES
# -------------------------------
# Local: SQLite, Production: can use DATABASE_URL with PostgreSQL
DATABASE_URL = os.getenv("DATABASE_URL")

DATABASES = {
    'default': {
        'ENGINE': config('DATABASE_ENGINE', default='django.db.backends.sqlite3'),
        'NAME': config('DATABASE_NAME', default='db.sqlite3'),
        'USER': config('DATABASE_USER', default=''),
        'PASSWORD': config('DATABASE_PASSWORD', default=''),
        'HOST': config('DATABASE_HOST', default=''),
        'PORT': config('DATABASE_PORT', default=''),
    }
}

# Optional: parse DATABASE_URL if you deploy to Leapcell/PostgreSQL

if DATABASE_URL:
    DATABASES['default'] = dj_database_url.parse(DATABASE_URL, conn_max_age=600)

# -------------------------------
# PASSWORD VALIDATION
# -------------------------------
AUTH_PASSWORD_VALIDATORS = [
    {'NAME': 'django.contrib.auth.password_validation.UserAttributeSimilarityValidator',},
    {'NAME': 'django.contrib.auth.password_validation.MinimumLengthValidator',},
    {'NAME': 'django.contrib.auth.password_validation.CommonPasswordValidator',},
    {'NAME': 'django.contrib.auth.password_validation.NumericPasswordValidator',},
]

# -------------------------------
# INTERNATIONALIZATION
# -------------------------------
LANGUAGE_CODE = 'en-us'
TIME_ZONE = 'UTC'
USE_I18N = True
USE_TZ = True

# -------------------------------
# DEFAULT AUTO FIELD
# -------------------------------
DEFAULT_AUTO_FIELD = 'django.db.models.BigAutoField'

# -------------------------------
# LOGGING CONFIGURATION
# -------------------------------

LOGGING = {
    "version": 1,
    "disable_existing_loggers": False,
    "formatters": {
        "json": {
            "format": (
                '{"time": "%(asctime)s", "level": "%(levelname)s", "message": "%(message)s"}'
            )
        },
    },
    "handlers": {
        "console": {
            "class": "logging.StreamHandler",
            "formatter": "json",
        },
    },
    "root": {
        "handlers": ["console"],
        "level": "INFO",
    },
}

DEFAULT_AUTO_FIELD = "django.db.models.BigAutoField"