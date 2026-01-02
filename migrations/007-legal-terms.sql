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
  tenant_id       BIGSERIAL NOT NULL REFERENCES companies(id),
  legal_document_id  BIGSERIAL NOT NULL REFERENCES legal_documents(id),
  accepted_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
  ip_address      INET NULL,
  user_agent      TEXT NULL,
  accepted        BOOLEAN NOT NULL DEFAULT true,
  UNIQUE (user_id, tenant_id, legal_document_id)
);

INSERT INTO legal_documents (doc_type, doc_version, published_at, content_sha256, is_active)
VALUES ('TERMS', '1.0', now(), '4b38069273fe6bee239b6412ca3357b03fb2ee9fc5caf5d55163e547bac9ea96', true);
INSERT INTO legal_documents (doc_type, doc_version, published_at, content_sha256, is_active)
VALUES ('PRIVACY', '1.0', now(), '189d7aade94cebeb56de4728755c200dfc50a37633b19835dda944c4a1b9acfb', true);