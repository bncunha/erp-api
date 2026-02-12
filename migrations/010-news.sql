CREATE TABLE news (
  id BIGSERIAL PRIMARY KEY,
  tenant_id BIGINT NULL REFERENCES companies(id),
  content_html TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE news_roles (
  id BIGSERIAL PRIMARY KEY,
  news_id BIGINT NOT NULL REFERENCES news(id) ON DELETE CASCADE,
  role TEXT NOT NULL CHECK (role IN ('ADMIN', 'RESELLER')),
  UNIQUE (news_id, role)
);

CREATE INDEX idx_news_created_at ON news (created_at DESC);
CREATE INDEX idx_news_tenant_created ON news (tenant_id, created_at DESC);
CREATE INDEX idx_news_roles_role ON news_roles (role, news_id);

