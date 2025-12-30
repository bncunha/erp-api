CREATE TABLE legal_documents (
  id            BIGSERIAL PRIMARY KEY,
  doc_type      TEXT NOT NULL CHECK (doc_type IN ('TERMS', 'PRIVACY')),
  doc_version   TEXT NOT NULL,
  published_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  content_sha256 CHAR(64) NOT NULL,            -- hash do texto publicado (prova)
  is_active     BOOLEAN NOT NULL DEFAULT true,
  UNIQUE (doc_type, doc_version)
);

CREATE TABLE legal_acceptances (
  id              BIGSERIAL PRIMARY KEY,
  user_id         BIGSERIAL NOT NULL REFERENCES users(id),
  tenant_id       BIGSERIAL NOT NULL REFERENCES tenants(id),
  legal_document_id  BIGSERIAL NOT NULL REFERENCES legal_documents(id),
  accepted_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  ip_address      INET NULL,
  user_agent      TEXT NULL,
  accepted        BOOLEAN NOT NULL DEFAULT true,
  UNIQUE (user_id, tenant_id, legal_document_id)
);