# Go Field Select

[![version](https://img.shields.io/github/v/release/golaxo/gofieldselect)](https://img.shields.io/github/v/release/golaxo/gofieldselect)
[![PR checks](https://github.com/golaxo/gofieldselect/actions/workflows/pr-checks.yml/badge.svg)](https://github.com/golaxo/gofieldselect/actions/workflows/pr-checks.yml)

> [!WARNING]
> GoFieldSelect is under heavy development.

A powerful `select` query parameter implementation based on [Field selection][field-selection] for filtering desired fields in a REST API.

## 🚀 Features

GoFieldSelect provides a way to return only certain fields.

### Root field selection

e.g.

Original JSON output:

```json
{
  "id": 1,
  "name": "John",
  "surname": "Doe",
  "age": 20
}
```

Output for `?fields=id,name`:

```json
{
  "id": 1,
  "name": "John"
}
```

### Nested field selection

e.g.

Original JSON output:

```json
{
  "id": 1,
  "name": "John",
  "address": {
    "street": "Example street",
    "number": "1"
  }
}
```

Output for `?fields=id,name,address(street)`:

```json
{
  "id": 1,
  "name": "John",
  "address": {
    "street": "Example street"
  }
}
```

## ⬇️ Getting Started

To start using it:

```bash
go get github.com/golaxo/gofieldselect@latest
```

And then use it either by 

### Reflection for an existing instance

```go
n, _ := gofieldselect.Parse("name,surname")
src := User{Name: "John", Surname: "Doe", Age: 20}
selected, _ := gofieldselect.GetWithReflection(n, src)
// selected only contains `name` and `surname` set, `age` is set to its default value.
```

Check [examples/getwithreflection](./examples/getwithreflection/main.go) to see it in action.

### Using Get for each field

```go
n, err := gofieldselect.Parse("name,address(street)")
an, _ := fieldSelection.SelectField("address")

dto := examples.User{
    Name:    gofieldselect.Get(n, "name", user.Name),
    Surname: gofieldselect.Get(n, "surname", user.Surname),
    Age:     gofieldselect.Get(n, "age", user.Age),
    Address: &examples.Address{
        Street: gofieldselect.Get(an.Child, "street", user.Address.Street),
        Number: gofieldselect.Get(an.Child, "number", user.Address.Number),
    },
}

// dto only contains `name` and `address.street` set, the rest of the fields contain the default value.
```

Check [examples/get](./examples/get/main.go) to see it in action.

[field-selection]: https://learn.microsoft.com/en-us/azure/data-api-builder/keywords/select-rest
