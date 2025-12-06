package domain

type Address struct {
    Id          int64
    Street      string
    Neighborhood string
    Number      string
    City        string
    UF          string
    Cep         string
    TenantId    int64
}
