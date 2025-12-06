ALTER TABLE companies ADD COLUMN cnpj VARCHAR(18);
ALTER TABLE companies ADD COLUMN cpf VARCHAR(14);
ALTER TABLE companies ADD COLUMN cellphone VARCHAR(20);
ALTER TABLE companies ADD COLUMN legal_name VARCHAR(255);

CREATE UNIQUE INDEX companies_cnpj_unique ON companies (cnpj) WHERE cnpj IS NOT NULL;
CREATE UNIQUE INDEX companies_cpf_unique ON companies (cpf) WHERE cpf IS NOT NULL;

CREATE TABLE addresses (
  id BIGSERIAL PRIMARY KEY,
  street VARCHAR(255) NOT NULL,
  neighborhood VARCHAR(255) NOT NULL,
  number VARCHAR(50) NOT NULL,
  city VARCHAR(255) NOT NULL,
  uf CHAR(2) NOT NULL,
  cep VARCHAR(20) NOT NULL,
  tenant_id BIGINT NOT NULL REFERENCES companies(id),
  created_at TIMESTAMP DEFAULT NOW() NOT NULL
);

CREATE UNIQUE INDEX addresses_tenant_id_unique ON addresses (tenant_id);
