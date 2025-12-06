package helper

import "strings"

func cleanNumeric(value string) string {
    replacer := strings.NewReplacer(".", "", "-", "", "/", "", " ", "", ",", "")
    return replacer.Replace(value)
}

func IsValidCPF(cpf string) bool {
    cpf = cleanNumeric(cpf)
    if len(cpf) != 11 {
        return false
    }

    invalid := []string{"00000000000", "11111111111", "22222222222", "33333333333", "44444444444", "55555555555", "66666666666", "77777777777", "88888888888", "99999999999"}
    for _, inv := range invalid {
        if cpf == inv {
            return false
        }
    }

    sum := 0
    for i := 0; i < 9; i++ {
        sum += int(cpf[i]-'0') * (10 - i)
    }
    firstDigit := (sum * 10) % 11
    if firstDigit == 10 {
        firstDigit = 0
    }
    if firstDigit != int(cpf[9]-'0') {
        return false
    }

    sum = 0
    for i := 0; i < 10; i++ {
        sum += int(cpf[i]-'0') * (11 - i)
    }
    secondDigit := (sum * 10) % 11
    if secondDigit == 10 {
        secondDigit = 0
    }
    return secondDigit == int(cpf[10]-'0')
}

func IsValidCNPJ(cnpj string) bool {
    cnpj = cleanNumeric(cnpj)
    if len(cnpj) != 14 {
        return false
    }

    invalid := []string{"00000000000000", "11111111111111", "22222222222222", "33333333333333", "44444444444444", "55555555555555", "66666666666666", "77777777777777", "88888888888888", "99999999999999"}
    for _, inv := range invalid {
        if cnpj == inv {
            return false
        }
    }

    weights1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
    weights2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}

    calcDigit := func(numbers string, weights []int) int {
        sum := 0
        for i, w := range weights {
            sum += int(numbers[i]-'0') * w
        }
        remainder := sum % 11
        if remainder < 2 {
            return 0
        }
        return 11 - remainder
    }

    firstDigit := calcDigit(cnpj[:12], weights1)
    if firstDigit != int(cnpj[12]-'0') {
        return false
    }
    secondDigit := calcDigit(cnpj[:13], weights2)
    return secondDigit == int(cnpj[13]-'0')
}

func SanitizeDocument(doc string) string {
    return cleanNumeric(doc)
}
